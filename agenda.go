package main

import (
	"errors"

	"github.com/tssaini/go-agenda/job"

	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

// Agenda struct stores all jobs
type Agenda struct {
	jobs    map[string]*job.Job
	running bool
	stop    chan struct{}
}

// New creates and returns new agenda object
func New() *Agenda {
	a := Agenda{jobs: make(map[string]*job.Job), running: false, stop: make(chan struct{})}
	return &a
}

// Define instantiate new job
func (a *Agenda) Define(name string, jobFunc func() error) {
	log.Infof("Defining job %v", name)
	a.jobs[name] = &job.Job{Name: name, JobFunc: jobFunc, Stop: make(chan struct{}), Running: false}
}

// Now runs the job provided now
func (a *Agenda) Now(name string) error {
	log.Infof("Starting job %v", name)
	job := a.jobs[name]
	err := job.JobFunc()
	if err != nil {
		log.Errorf("Error while running %v: %v", name, err)
		return err
	}
	log.Infof("Completed job %v successfully", name)
	return nil
}

// Schedule the next time the job should run
func (a *Agenda) Schedule(name string, time string) error {
	return nil
}

// RepeatEvery when the job should repeat
func (a *Agenda) RepeatEvery(name string, spec string) error {
	schedule, err := cron.ParseStandard(spec)
	if err != nil {
		return err
	}
	a.jobs[name].Schedule = schedule
	return nil
}

// Start agenda
func (a *Agenda) Start() error {
	if a.running {
		return errors.New("Agenda is already running")
	}
	a.running = true
	go a.run()
	return nil
}

// Stop agenda
func (a *Agenda) Stop() error {
	a.stop <- struct{}{}
	a.running = false
	return nil
}

func (a *Agenda) run() error {
	log.Infof("Running agenda loop")
	// run all jobs
	for _, j := range a.jobs {
		go j.LaunchJob()
	}
	for {
		select {
		case <-a.stop:
			//stop all jobs
			for _, j := range a.jobs {
				if j.Running {
					j.Stop <- struct{}{}
				}
			}
			return nil
		}
	}

}
