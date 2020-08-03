package agenda

import (
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

// ScheduledTask interface
type ScheduledTask interface {
	Start() error
	Schedule(parser func(string) (cron.Schedule, error), spec string) error
	Stop()
	ScheduleNow()
	HasSchedule() bool
}

// Job defines the config of a job
type Job struct {
	Name         string
	JobFunc      func() error
	Scheduler    cron.Schedule
	JobMutex     sync.Mutex
	NextRun      time.Time
	LastRun      time.Time
	StopChan     chan struct{}
	StartedChan  chan struct{}
	Running      bool
	RunningMutex sync.RWMutex
	LastErr      error
}

// NewJob creates new job object
func NewJob(name string, jobFunc func() error) *Job {
	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)
	return &Job{Name: name, JobFunc: jobFunc, StopChan: make(chan struct{}), StartedChan: make(chan struct{}), Running: false, JobMutex: sync.Mutex{}, RunningMutex: sync.RWMutex{}}
}

// Start launch the job to be run on schedule
func (j *Job) Start() error {
	if !j.HasSchedule() {
		return fmt.Errorf("Job %v does not have a schedule defined", j.Name)
	}
	if j.isRunning() {
		return fmt.Errorf("Job %v is already running", j.Name)
	}

	go j.run()
	<-j.StartedChan
	return nil
}

func (j *Job) run() {
	j.setRunning(true)
	j.StartedChan <- struct{}{}
	for {
		t := j.Scheduler.Next(time.Now())
		d := t.Sub(time.Now())

		select {
		case <-j.StopChan:
			return
		case <-time.After(d):
			//Run job here
			go j.runJobFunc()
		}
	}
}

// Schedule set a schedule for the job
func (j *Job) Schedule(parser func(string) (cron.Schedule, error), spec string) error {
	schedule, err := parser(spec)
	if err != nil {
		return err
	}
	j.Scheduler = schedule
	return nil
}

// runJobFunc runs the job function
func (j *Job) runJobFunc() {
	log.Infof("Running job %v", j.Name)
	j.JobMutex.Lock()
	j.LastRun = time.Now()
	err := j.JobFunc()
	j.LastErr = err
	j.JobMutex.Unlock()
}

// ScheduleNow runs the job now
func (j *Job) ScheduleNow() {
	go j.runJobFunc()
}

// Stop the job
func (j *Job) Stop() {
	log.Infof("Stopping job %v", j.Name)
	if j.isRunning() {
		j.StopChan <- struct{}{}
		j.setRunning(false)
	}
}

// HasSchedule returns true if job has a schedule
func (j *Job) HasSchedule() bool {
	if j.Scheduler != nil {
		return true
	}
	return false
}

func (j *Job) isRunning() bool {
	j.RunningMutex.RLock()
	defer j.RunningMutex.RUnlock()
	return j.Running
}

func (j *Job) setRunning(running bool) {
	j.RunningMutex.Lock()
	j.Running = running
	j.RunningMutex.Unlock()
}
