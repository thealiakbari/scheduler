package application

import (
	"errors"
	"sync"

	"github.com/thealiakbari/scheduler/internal/domain"
	"github.com/thealiakbari/scheduler/internal/infrastructure/persistence"
	"github.com/thealiakbari/scheduler/internal/queue"
)

type TaskService struct {
	storage persistence.Storage
	queue   *queue.PriorityQueue
	mu      sync.Mutex
}

func NewTaskService(storage persistence.Storage, queue *queue.PriorityQueue) *TaskService {
	return &TaskService{
		storage: storage,
		queue:   queue,
	}
}

func (s *TaskService) CreateTask(priority domain.Priority, payload []byte) (*domain.Task, error) {
	task := domain.NewTask(priority, payload)
	if err := s.storage.Save(task); err != nil {
		return nil, err
	}
	s.queue.Push(task)
	return task, nil
}

func (s *TaskService) GetTaskStatus(id string) (domain.Status, error) {
	task, err := s.storage.Get(id)
	if err != nil {
		return "", err
	}
	return task.Status, nil
}

func (s *TaskService) ProcessNextTask() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	task := s.queue.Pop()
	if task == nil {
		return errors.New("no task available")
	}
	task.Status = domain.StatusRunning
	if err := s.storage.Update(task); err != nil {
		task.Status = domain.StatusFailed
		return err
	}
	task.Status = domain.StatusCompleted
	s.storage.Update(task)
	return nil
}
