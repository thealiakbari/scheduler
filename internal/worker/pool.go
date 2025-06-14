package worker

import (
	"log"
	"sync"
	"time"

	"github.com/thealiakbari/scheduler/internal/application"
)

type Pool struct {
	size      int
	taskSvc   *application.TaskService
	quit      chan struct{}
	waitGroup sync.WaitGroup
}

func NewPool(size int, taskSvc *application.TaskService) *Pool {
	return &Pool{
		size:    size,
		taskSvc: taskSvc,
		quit:    make(chan struct{}),
	}
}

func (p *Pool) Start() {
	for i := 0; i < p.size; i++ {
		p.waitGroup.Add(1)
		go p.worker(i)
	}
}

func (p *Pool) worker(id int) {
	defer p.waitGroup.Done()
	for {
		select {
		case <-p.quit:
			log.Printf("Worker %d shutting down", id)
			return
		default:
			if err := p.taskSvc.ProcessNextTask(); err != nil {
				// در صورتی که تسکی نباشد کمی استراحت می‌کنیم.
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func (p *Pool) Stop() {
	close(p.quit)
	p.waitGroup.Wait()
}
