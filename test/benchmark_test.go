package test

import (
	"github.com/thealiakbari/scheduler/internal/domain"
	"github.com/thealiakbari/scheduler/internal/queue"
	"testing"
)

func BenchmarkPriorityQueuePushPop(b *testing.B) {
	pq := queue.NewPriorityQueue()
	task := domain.NewTask(domain.Medium, []byte("payload"))
	for n := 0; n < b.N; n++ {
		pq.Push(task)
		_ = pq.Pop()
	}
}
