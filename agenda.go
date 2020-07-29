package main

import (
	"time"

	"github.com/tssaini/go-agenda/job"

	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

// Agenda struct stores all jobs
type Agenda struct {
	jobs map[string]*job.Job
}

// New creates and returns new agenda object
func New() *Agenda {
	a := Agenda{jobs: make(map[string]*job.Job)}
	return &a
}

// Define instantiate new job
func (a *Agenda) Define(name string, jobFunc func() error) {
	log.Infof("Defining job %v", name)
	a.jobs[name] = &job.Job{Name: name, JobFunc: jobFunc}
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
	scheduler, err := cron.ParseStandard(spec)
	if err != nil {
		return err
	}
	nextRun := scheduler.Next(time.Now())
	a.jobs[name].NextRun = nextRun
	// log.Infof("Scheduling %v for %v", name, nextRun)
	return nil
}
