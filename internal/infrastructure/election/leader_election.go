package election

import (
	"log"
	"sync"
	"time"
)

type LeaderElector struct {
	isLeader bool
	quit     chan struct{}
	lock     sync.Mutex
}

func NewLeaderElector() *LeaderElector {
	return &LeaderElector{
		isLeader: false,
		quit:     make(chan struct{}),
	}
}

func (le *LeaderElector) Start() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			le.lock.Lock()
			le.isLeader = !le.isLeader
			if le.isLeader {
				log.Println("This node is now the leader")
			} else {
				log.Println("This node is now a follower")
			}
			le.lock.Unlock()
		case <-le.quit:
			ticker.Stop()
			return
		}
	}
}

func (le *LeaderElector) IsLeader() bool {
	le.lock.Lock()
	defer le.lock.Unlock()
	return le.isLeader
}

func (le *LeaderElector) Stop() {
	close(le.quit)
}
