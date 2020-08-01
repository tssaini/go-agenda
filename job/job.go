package job

import (
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

// Job defines the config of a job
type Job struct {
	Name         string
	jobFunc      func() error
	schedule     cron.Schedule
	jobMutex     sync.Mutex
	nextRun      time.Time
	lastRun      time.Time
	stopChan     chan struct{}
	startedChan  chan struct{}
	running      bool
	runningMutex sync.RWMutex
	lastErr      error
}

// New creates new job object
func New(name string, jobFunc func() error) *Job {
	return &Job{Name: name, jobFunc: jobFunc, stopChan: make(chan struct{}), startedChan: make(chan struct{}), running: false, jobMutex: sync.Mutex{}, runningMutex: sync.RWMutex{}}
}

// StartJob launch the job to be run on schedule
func (j *Job) StartJob() error {
	if !j.HasSchedule() {
		return fmt.Errorf("Job %v does not have a schedule defined", j.Name)
	}
	if j.isRunning() {
		return fmt.Errorf("Job %v is already running", j.Name)
	}

	go j.run()
	<-j.startedChan
	return nil
}

func (j *Job) run() {
	j.setRunning(true)
	j.startedChan <- struct{}{}
	for {
		t := j.schedule.Next(time.Now())
		d := t.Sub(time.Now())

		select {
		case <-j.stopChan:
			return
		case <-time.After(d):
			//Run job here
			go j.runJobFunc()
		}
	}
}

// ScheduleJob set a schedule for the job
func (j *Job) ScheduleJob(parser func(string) (cron.Schedule, error), spec string) error {
	schedule, err := parser(spec)
	if err != nil {
		return err
	}
	j.schedule = schedule
	return nil
}

// runJobFunc run the job function
func (j *Job) runJobFunc() {
	log.Infof("Running job %v", j.Name)
	j.jobMutex.Lock()
	j.lastRun = time.Now()
	err := j.jobFunc()
	j.lastErr = err
	j.jobMutex.Unlock()
}

// ScheduleJobNow run the job now
func (j *Job) ScheduleJobNow() {
	go j.runJobFunc()
}

// StopJob stop the job
func (j *Job) StopJob() {
	log.Infof("Stopping job %v", j.Name)
	if j.isRunning() {
		j.stopChan <- struct{}{}
		j.setRunning(false)
	}
}

// HasSchedule returns true if job has a schedule
func (j *Job) HasSchedule() bool {
	if j.schedule != nil {
		return true
	}
	return false
}

func (j *Job) isRunning() bool {
	j.runningMutex.RLock()
	defer j.runningMutex.RUnlock()
	return j.running
}

func (j *Job) setRunning(running bool) {
	j.runningMutex.Lock()
	j.running = running
	j.runningMutex.Unlock()
}
