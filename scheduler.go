package smash

import (
	"errors"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

var (
	errTaskNotFound = errors.New("Task not found")
)

type Scheduler struct {
	tasks    map[string]Task
	stopChan chan struct{}
}

func (s *Scheduler) RunAfter(duration time.Duration, function TaskFunction) (string, error) {
	return s.RunAt(time.Now().Add(duration), function)
}

func (s *Scheduler) RunAt(time time.Time, function TaskFunction) (string, error) {
	id, err := gonanoid.New()

	if err != nil {
		return "", err
	}

	s.tasks[id] = Task{
		Schedule: Schedule{
			NextRun: time,
		},
		Func: function,
	}

	return id, nil
}

func (s *Scheduler) RunEvery(duration time.Duration, function TaskFunction) (string, error) {
	id, err := gonanoid.New()

	if err != nil {
		return "", err
	}

	s.tasks[id] = Task{
		Schedule: Schedule{
			IsRecurring: true,
			Duration:    duration,
			NextRun:     time.Now().Add(duration),
		},
		Func: function,
	}

	return id, nil
}

func (s *Scheduler) Cancel(id string) error {
	if _, ok := s.tasks[id]; !ok {
		return errTaskNotFound
	}

	delete(s.tasks, id)

	return nil
}

func (s *Scheduler) Start() error {
	go func() {
		ticker := time.NewTicker(1 * time.Millisecond)

		for {
			select {
			case <-ticker.C:
				s.runPending()
			case <-s.stopChan:
				close(s.stopChan)
			}
		}
	}()

	return nil
}

func (s *Scheduler) Stop() {
	s.stopChan <- struct{}{}
}

func (s *Scheduler) Wait() {
	<-s.stopChan
}

func (s *Scheduler) runPending() {
	for id, task := range s.tasks {
		if task.IsDue() {
			go task.Run()

			if !task.IsRecurring {
				delete(s.tasks, id)
			}
		}
	}
}

func NewScheduler() (*Scheduler, error) {
	return &Scheduler{
		tasks:    map[string]Task{},
		stopChan: make(chan struct{}),
	}, nil
}

type TaskFunction func()

type Schedule struct {
	IsRecurring bool
	LastRun     time.Time
	NextRun     time.Time
	Duration    time.Duration
}

type Task struct {
	Schedule
	Func TaskFunction
}

func (t *Task) Run() {
	t.scheduleNextRun()
	t.Func()
}

func (t *Task) IsDue() bool {
	timeNow := time.Now()
	return timeNow == t.NextRun || timeNow.After(t.NextRun)
}

func (task *Task) scheduleNextRun() {
	if !task.IsRecurring {
		return
	}

	task.LastRun = task.NextRun
	task.NextRun = task.NextRun.Add(task.Duration)
}
