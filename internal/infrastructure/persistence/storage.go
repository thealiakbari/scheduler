package persistence

import (
	"errors"
	"sync"

	"github.com/thealiakbari/scheduler/internal/domain"
)

type Storage interface {
	Save(task *domain.Task) error
	Get(id string) (*domain.Task, error)
	Update(task *domain.Task) error
	Close()
}

type InMemoryStorage struct {
	tasks map[string]*domain.Task
	lock  sync.RWMutex
}

func NewStorage() (Storage, error) {
	return &InMemoryStorage{
		tasks: make(map[string]*domain.Task),
	}, nil
}

func (s *InMemoryStorage) Save(task *domain.Task) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.tasks[task.ID] = task
	return nil
}

func (s *InMemoryStorage) Get(id string) (*domain.Task, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if task, ok := s.tasks[id]; ok {
		return task, nil
	}
	return nil, errors.New("task not found")
}

func (s *InMemoryStorage) Update(task *domain.Task) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.tasks[task.ID]; !ok {
		return errors.New("task not found")
	}
	s.tasks[task.ID] = task
	return nil
}

func (s *InMemoryStorage) Close() {
}
