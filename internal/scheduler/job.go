package scheduler

import (
	"errors"
	"sync"
	"time"

	"stavki/internal/log"

	"github.com/google/uuid"
)

type (
	Job struct {
		id       string
		name     string
		interval time.Duration
		lastRun  time.Time
		task     TaskFunc
		ticker   *time.Ticker
		stopCh   chan struct{}

		logger log.FieldLogger

		mu      sync.RWMutex
		started bool
	}

	TaskFunc func(logger log.FieldLogger) error
)

var (
	ErrJobStarted    = errors.New("job already started")
	ErrJobNotStarted = errors.New("job not started")
)

func newJob(logger log.FieldLogger, name string, interval time.Duration, task TaskFunc) *Job {
	id := uuid.New().String()

	if logger == nil {
		logger = log.New()
	}

	logger = logger.WithFields(log.Fields{
		"job_id":   id,
		"name":     name,
		"interval": interval,
	})

	return &Job{
		id:       id,
		name:     name,
		interval: interval,
		task:     task,
		ticker:   time.NewTicker(interval),
		stopCh:   make(chan struct{}),
		logger:   logger,
	}
}

func (j *Job) ID() string {
	return j.id
}

func (j *Job) Name() string {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.name
}

func (j *Job) ChangeName(name string) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.name = name
	j.logger = j.logger.WithField("name", name)
}

func (j *Job) Interval() time.Duration {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.interval
}

func (j *Job) ChangeInterval(interval time.Duration) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.interval = interval
	if j.ticker != nil {
		j.ticker.Stop()
	}
	j.ticker = time.NewTicker(interval)
	j.logger.WithField("interval", j.interval)
}

func (j *Job) LastRun() time.Time {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.lastRun
}

func (j *Job) Started() bool {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.started
}

func (j *Job) start(onStart func(), onStop func()) error {
	j.mu.Lock()
	if j.started {
		j.mu.Unlock()
		return ErrJobStarted
	}
	j.started = true
	j.mu.Unlock()

	j.ticker = time.NewTicker(j.interval)
	j.stopCh = make(chan struct{})
	onStart()
	go func() {
		defer onStop()

		for {
			select {
			case <-j.ticker.C:
				j.wrapped()
			case <-j.stopCh:
				j.ticker.Stop()
				return
			}
		}
	}()

	return nil
}

func (j *Job) stop() error {
	j.mu.Lock()
	if !j.started {
		j.mu.Unlock()
		return ErrJobNotStarted
	}
	j.started = false
	j.mu.Unlock()

	close(j.stopCh)
	return nil
}

func (j *Job) wrapped() {
	start := time.Now()
	j.mu.Lock()
	logger := j.logger
	j.lastRun = start
	j.mu.Unlock()
	logger.Info("Job started")

	var err error
	defer func() {
		logger = logger.WithField("duration", time.Since(start).String())

		if r := recover(); r != nil {
			logger.WithField("panic", r).Error("Job panicked")
			return
		}

		if err != nil {
			logger.WithError(err).Error("Job failed")
			return
		}

		logger.Info("Job finished")
	}()

	err = j.task(logger)
}
