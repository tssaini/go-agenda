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

	AccessMutex sync.RWMutex
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
		log.Errorf("%v", err)
	}
	//job exists in the db
	if existingJob != nil {
		log.Infof("Job %v found in db", name)
		existingJob.JobFunc = jobFunc
		existingJob.Repository = r
		existingJob.Scheduled = false
		return existingJob, nil
	}
	log.Infof("Job %v not found in db", name)
	j := &Job{Name: name, JobFunc: jobFunc, StopChan: make(chan struct{}), StartedChan: make(chan struct{}), Scheduled: false, JobRunning: false, JobMutex: sync.Mutex{}, AccessMutex: sync.RWMutex{}, Repository: r}
	err = r.SaveJob(j)
	if err != nil {
		log.Errorf("%v", err)
	}
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
	err := j.Repository.SaveJob(j)
	if err != nil {
		log.Errorf("%v", err)
	}
	return nil
}

func (j *Job) run() {
	j.setScheduled(true)
	defer j.setScheduled(false)
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
			err := j.Repository.SaveJob(j)
			if err != nil {
				log.Errorf("%v", err)
			}
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
		log.Errorf("%v", err)
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
	err := j.Repository.SaveJob(j)
	if err != nil {
		log.Errorf("%v", err)
	}

	j.LastErr = j.JobFunc()
	j.JobRunning = false
	err = j.Repository.SaveJob(j)
	if err != nil {
		log.Errorf("%v", err)
	}
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
	j.AccessMutex.RLock()
	defer j.AccessMutex.RUnlock()
	return j.Scheduled
}

func (j *Job) setScheduled(running bool) {
	j.AccessMutex.Lock()
	j.Scheduled = running
	j.AccessMutex.Unlock()
}
