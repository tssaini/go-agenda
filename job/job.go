package job

import (
	"fmt"
	"time"

	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

// Job defines the config of a job
type Job struct {
	Name     string
	JobFunc  func() error
	Schedule cron.Schedule

	NextRun  time.Time
	LastRun  time.Time
	StopChan chan struct{}
	Running  bool
}

// New creates new job object
func New(name string, jobFunc func() error) *Job {
	return &Job{Name: name, JobFunc: jobFunc, StopChan: make(chan struct{}), Running: false}
}

// StartJob launch the job
func (j *Job) StartJob() error {
	if !j.HasSchedule() {
		return fmt.Errorf("Job %v does not have a schedule defined", j.Name)
	}
	if j.Running {
		return fmt.Errorf("Job %v is already running", j.Name)
	}
	j.Running = true
	go j.run()
	return nil
}

func (j *Job) run() {
	for {
		t := j.Schedule.Next(time.Now())
		d := t.Sub(time.Now())

		select {
		case <-j.StopChan:
			log.Infof("Stopping job %v", j.Name)
			return
		case runTime := <-time.After(d):
			//Run job here
			log.Infof("Running job %v", j.Name)
			j.LastRun = runTime
			go j.RunJobFunc()
		}
	}
}

// ScheduleJob set a schedule for the job
func (j *Job) ScheduleJob(parser func(string) (cron.Schedule, error), spec string) error {
	schedule, err := parser(spec)
	if err != nil {
		return err
	}
	j.Schedule = schedule
	return nil
}

// RunJobFunc run the job function
func (j *Job) RunJobFunc() {
	j.JobFunc()
}

// StopJob stop the job
func (j *Job) StopJob() {
	if j.Running {
		j.Running = false
		j.StopChan <- struct{}{}
	}
}

// HasSchedule returns true if job has a schedule
func (j *Job) HasSchedule() bool {
	if j.Schedule != nil {
		return true
	}
	return false
}
