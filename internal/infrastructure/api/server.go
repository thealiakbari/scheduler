package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/thealiakbari/scheduler/internal/application"
)

type Server struct {
	taskService *application.TaskService
	router      *mux.Router
}

func NewServer(taskService *application.TaskService) *Server {
	s := &Server{
		taskService: taskService,
		router:      mux.NewRouter(),
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.HandleFunc("/tasks", s.taskService.HandleCreateTask).Methods("POST")
	s.router.HandleFunc("/tasks/status", s.taskService.HandleGetTaskStatus).Methods("GET")
}

func (s *Server) Router() http.Handler {
	return s.router
}
