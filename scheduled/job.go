package scheduled

import (
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

// Task interface
type Task interface {
	Start() error
	Schedule(parser func(string) (cron.Schedule, error), spec string) error
	Stop()
	ScheduleNow()
	HasSchedule() bool
}

// Job defines the config of a job
type Job struct {
	Name        string
	JobFunc     func() error
	Scheduler   cron.Schedule
	JobMutex    sync.Mutex
	NextRun     time.Time
	LastRun     time.Time
	StopChan    chan struct{}
	StartedChan chan struct{}

	Scheduled  bool
	JobRunning bool

	accessMutex sync.RWMutex
	LastErr     error
	Repository  JobRepository
}

// JobRepository provides access to storage
type JobRepository interface {
	// FindAllJobs() ([]*Job, error)
	FindJobByName(jobName string) (*Job, error)
	SaveJob(j *Job) error
}

// NewJob creates new job object
func NewJob(name string, jobFunc func() error, r JobRepository) (*Job, error) {
	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)

	existingJob, err := r.FindJobByName(name)
	if err != nil {
		return nil, err
	}
	//job exists in the db
	if existingJob != nil {
		existingJob.JobFunc = jobFunc
		existingJob.Repository = r
		return existingJob, nil
	}

	j := &Job{Name: name, JobFunc: jobFunc, StopChan: make(chan struct{}), StartedChan: make(chan struct{}), Scheduled: false, JobRunning: false, JobMutex: sync.Mutex{}, accessMutex: sync.RWMutex{}, Repository: r}
	r.SaveJob(j)
	return j, nil
}

// Start launch the job to be run on schedule
func (j *Job) Start() error {
	if !j.HasSchedule() {
		return fmt.Errorf("Job %v does not have a schedule defined", j.Name)
	}
	if j.isScheduled() {
		return fmt.Errorf("Job %v is already scheduled to run", j.Name)
	}

	go j.run()
	<-j.StartedChan
	return nil
}

func (j *Job) run() {
	j.setScheduled(true)
	j.StartedChan <- struct{}{}

	if j.NextRun != (time.Time{}) && j.NextRun.Before(time.Now()) {
		j.runJobFunc()
	}

	for {
		if j.NextRun == (time.Time{}) {
			j.calcNextRun()
		}

		d := j.NextRun.Sub(time.Now())

		select {
		case <-j.StopChan:
			return
		case <-time.After(d):
			//Run job here
			j.runJobFunc()
			j.calcNextRun()
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
	j.calcNextRun()
	err = j.Repository.SaveJob(j)
	if err != nil {
		return err
	}
	return nil
}

func (j *Job) calcNextRun() {
	j.NextRun = j.Scheduler.Next(time.Now())
}

// runJobFunc runs the job function
func (j *Job) runJobFunc() {
	log.Infof("Running job %v", j.Name)
	j.JobMutex.Lock()
	defer j.JobMutex.Unlock()

	j.JobRunning = true
	j.LastRun = time.Now()
	j.Repository.SaveJob(j)

	j.LastErr = j.JobFunc()
	j.JobRunning = false
	j.Repository.SaveJob(j)
}

// ScheduleNow runs the job now
func (j *Job) ScheduleNow() {
	go j.runJobFunc()
}

// Stop the job
func (j *Job) Stop() {
	log.Infof("Stopping job %v", j.Name)
	if j.isScheduled() {
		j.StopChan <- struct{}{}
		j.setScheduled(false)
	}
}

// HasSchedule returns true if job has a schedule
func (j *Job) HasSchedule() bool {
	if j.Scheduler != nil {
		return true
	}
	return false
}

func (j *Job) isScheduled() bool {
	j.accessMutex.RLock()
	defer j.accessMutex.RUnlock()
	return j.Scheduled
}

func (j *Job) setScheduled(running bool) {
	j.accessMutex.Lock()
	j.Scheduled = running
	j.accessMutex.Unlock()
}
