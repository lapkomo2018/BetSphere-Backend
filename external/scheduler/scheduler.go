package scheduler

import (
	"sync"
	"time"
)

type (
	Job struct {
		Interval time.Duration
		Task     func()
		ticker   *time.Ticker
		stop     chan struct{}
	}

	Scheduler struct {
		jobs []*Job
		wg   sync.WaitGroup
	}
)

func New() *Scheduler {
	return &Scheduler{
		jobs: make([]*Job, 0),
	}
}

func (s *Scheduler) AddJob(interval time.Duration, task func()) *Scheduler {
	job := &Job{
		Interval: interval,
		Task:     task,
		ticker:   time.NewTicker(interval),
		stop:     make(chan struct{}),
	}
	s.jobs = append(s.jobs, job)

	s.wg.Add(1)
	go func(j *Job) {
		defer s.wg.Done()
		for {
			select {
			case <-j.ticker.C:
				j.Task()
			case <-j.stop:
				j.ticker.Stop()
				return
			}
		}
	}(job)

	return s
}

func (s *Scheduler) StopAll() {
	for _, job := range s.jobs {
		close(job.stop)
	}
	s.wg.Wait()
}
