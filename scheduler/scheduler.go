package scheduler

import (
	"math/rand"
	"time"
)

// Scheduler manages collection intervals with jitter
type Scheduler struct {
	interval time.Duration
	jitter   time.Duration
}

// NewScheduler creates a new scheduler
func NewScheduler(interval, jitter time.Duration) *Scheduler {
	return &Scheduler{
		interval: interval,
		jitter:   jitter,
	}
}

// Next returns the next collection time with jitter applied
func (s *Scheduler) Next() time.Duration {
	if s.jitter > 0 {
		jitterAmount := time.Duration(rand.Int63n(int64(s.jitter)))
		return s.interval + jitterAmount
	}
	return s.interval
}

// Wait waits for the next scheduled time
func (s *Scheduler) Wait() {
	time.Sleep(s.Next())
}

