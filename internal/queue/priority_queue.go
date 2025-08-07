package queue

import (
	"container/heap"
	"sync"

	"github.com/thealiakbari/scheduler/internal/domain"
)

type TaskHeap []*domain.Task

func (h TaskHeap) Len() int { return len(h) }
func (h TaskHeap) Less(i, j int) bool {
	return h[i].Priority < h[j].Priority
}
func (h TaskHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *TaskHeap) Push(x interface{}) {
	*h = append(*h, x.(*domain.Task))
}

func (h *TaskHeap) Pop() interface{} {
	old := *h
	n := len(old)
	task := old[n-1]
	*h = old[0 : n-1]
	return task
}

type PriorityQueue struct {
	tasks TaskHeap
	lock  sync.Mutex
}

func NewPriorityQueue() *PriorityQueue {
	h := &TaskHeap{}
	heap.Init(h)
	return &PriorityQueue{
		tasks: *h,
	}
}

func (pq *PriorityQueue) Push(task *domain.Task) {
	pq.lock.Lock()
	defer pq.lock.Unlock()
	heap.Push(&pq.tasks, task)
}

func (pq *PriorityQueue) Pop() *domain.Task {
	pq.lock.Lock()
	defer pq.lock.Unlock()
	if pq.tasks.Len() == 0 {
		return nil
	}
	return heap.Pop(&pq.tasks).(*domain.Task)
}

func (pq *PriorityQueue) Len() int {
	pq.lock.Lock()
	defer pq.lock.Unlock()
	return pq.tasks.Len()
}
