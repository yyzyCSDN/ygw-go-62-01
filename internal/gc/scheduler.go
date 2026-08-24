package gc

import (
	"sync"
	"time"
)

// Scheduler runs reclamation passes on an interval.
type Scheduler struct {
	interval time.Duration
	run      func(now time.Time) (*Report, error)
	mu       sync.Mutex
	last     *Report
}

// NewScheduler wraps a run function with a fixed interval.
func NewScheduler(interval time.Duration, run func(now time.Time) (*Report, error)) *Scheduler {
	return &Scheduler{interval: interval, run: run}
}

// Tick executes one pass when the interval has elapsed since the last run.
func (s *Scheduler) Tick(now time.Time) (*Report, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.last != nil && now.Sub(s.last.FinishedAt) < s.interval {
		return s.last, nil
	}
	report, err := s.run(now)
	if err != nil {
		return nil, err
	}
	s.last = report
	return report, nil
}
