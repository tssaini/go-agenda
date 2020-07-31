package main

import (
	"errors"
	"sync"

	"github.com/tssaini/go-agenda/job"

	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

// Agenda struct stores all jobs
type Agenda struct {
	jobs       map[string]*job.Job
	jobsMutex  *sync.RWMutex
	running    bool
	stop       chan struct{}
	newJob     chan *job.Job
	cronParser func(string) (cron.Schedule, error)
}

// New creates and returns new agenda object
func New() *Agenda {
	a := Agenda{jobs: make(map[string]*job.Job), jobsMutex: &sync.RWMutex{}, running: false, stop: make(chan struct{}), newJob: make(chan *job.Job), cronParser: cron.ParseStandard}
	return &a
}

// Define instantiate new job
func (a *Agenda) Define(name string, jobFunc func() error) {
	log.Infof("Defining job %v", name)
	j := &job.Job{Name: name, JobFunc: jobFunc, Stop: make(chan struct{}), Running: false}
	//TODO: check if the job already exists
	a.jobsMutex.Lock()
	a.jobs[name] = j
	a.jobsMutex.Unlock()
}

// Now runs the job provided now
func (a *Agenda) Now(name string) error {
	log.Infof("Starting job %v", name)
	a.jobsMutex.RLock()
	job := a.jobs[name]
	a.jobsMutex.RUnlock()
	// err := job.JobFunc()
	// if err != nil {
	// 	log.Errorf("Error while running %v: %v", name, err)
	// 	return err
	// }
	job.RunJob()
	log.Infof("Completed job %v successfully", name)
	return nil
}

// Schedule the next time the job should run
func (a *Agenda) Schedule(name string, time string) error {
	return nil
}

// RepeatEvery when the job should repeat
func (a *Agenda) RepeatEvery(name string, spec string) error {
	schedule, err := a.cronParser(spec)
	if err != nil {
		return err
	}
	a.jobsMutex.RLock()
	j := a.jobs[name]
	a.jobsMutex.RUnlock()
	j.Schedule = schedule
	if a.running {
		a.newJob <- j
	}
	return nil
}

// Start agenda
func (a *Agenda) Start() error {
	if a.running {
		return errors.New("Agenda is already running")
	}
	go a.run()
	a.running = true
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
	a.jobsMutex.RLock()
	for _, j := range a.jobs {
		go j.LaunchJob()
	}
	a.jobsMutex.RUnlock()
	for {
		select {
		case <-a.stop:
			log.Infof("Stopping agenda loop")
			//stop all jobs
			a.jobsMutex.RLock()
			for _, j := range a.jobs {
				if j.Running {
					j.Stop <- struct{}{}
				}
			}
			a.jobsMutex.RUnlock()
			return nil
		case j := <-a.newJob:
			go j.LaunchJob()
		}
	}
}
