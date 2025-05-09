package scheduler

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type (
	Scheduler struct {
		jobs map[string]*Job

		mu sync.RWMutex
		wg sync.WaitGroup
	}
)

var (
	ErrJobNotFound = errors.New("job not found")
)

func New() *Scheduler {
	return &Scheduler{
		jobs: make(map[string]*Job),
	}
}

func (s *Scheduler) Jobs() []*Job {
	s.mu.RLock()
	jobs := make([]*Job, 0, len(s.jobs))
	for _, job := range s.jobs {
		jobs = append(jobs, job)
	}
	s.mu.RUnlock()

	return jobs
}

func (s *Scheduler) Job(id string) *Job {
	s.mu.RLock()
	defer s.mu.RUnlock()

	job, ok := s.jobs[id]
	if ok {
		return job
	}

	return nil
}

func (s *Scheduler) Add(name string, interval time.Duration, task TaskFunc) *Job {
	j := newJob(name, interval, task)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[j.ID()] = j

	return j
}

func (s *Scheduler) Start(jobID string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	job, ok := s.jobs[jobID]
	if !ok {
		return ErrJobNotFound
	}

	if err := job.start(func() { s.wg.Add(1) }, s.wg.Done); err != nil {
		return err
	}

	return nil
}

func (s *Scheduler) Stop(jobID string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	job, ok := s.jobs[jobID]
	if !ok {
		return ErrJobNotFound
	}

	return job.stop()
}

func (s *Scheduler) StartAll() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, job := range s.jobs {
		if err := job.start(func() { s.wg.Add(1) }, s.wg.Done); err != nil {
			return err
		}
	}

	return nil
}

func (s *Scheduler) StopAll() []error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var errs []error
	for _, j := range s.jobs {
		if err := j.stop(); err != nil {
			err = errors.Join(err, fmt.Errorf("`%s job %s stop failed", j.ID(), j.Name()))
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs
	}

	s.wg.Wait()
	return nil
}
