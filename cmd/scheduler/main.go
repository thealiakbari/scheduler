package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/thealiakbari/scheduler/internal/application"
	"github.com/thealiakbari/scheduler/internal/infrastructure/api"
	"github.com/thealiakbari/scheduler/internal/infrastructure/election"
	"github.com/thealiakbari/scheduler/internal/infrastructure/metrics"
	"github.com/thealiakbari/scheduler/internal/infrastructure/persistence"
	"github.com/thealiakbari/scheduler/internal/queue"
	"github.com/thealiakbari/scheduler/internal/worker"
)

func main() {
	storage, err := persistence.NewStorage()
	if err != nil {
		log.Fatalf("failed to init storage: %v", err)
	}
	defer storage.Close()

	pq := queue.NewPriorityQueue()

	svc := application.NewTaskService(storage, pq)

	leaderElector := election.NewLeaderElector()
	go leaderElector.Start()

	pool := worker.NewPool(5, svc)
	pool.Start()

	srv := api.NewServer(svc)
	go func() {
		log.Println("API server is running on :8080")
		if err := http.ListenAndServe(":8080", srv.Router()); err != nil {
			log.Fatalf("failed to start API server: %v", err)
		}
	}()

	go func() {
		log.Println("Metrics server running on :9090")
		if err := http.ListenAndServe(":9090", metrics.Handler()); err != nil {
			log.Fatalf("failed to start metrics server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gracefully...")
	pool.Stop()
	leaderElector.Stop()
	time.Sleep(2 * time.Second)
}
