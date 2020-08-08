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
	if j.IsScheduled() {
		return fmt.Errorf("Job %v is already scheduled to run", j.Name)
	}

	go j.run()
	<-j.StartedChan
	return nil
}

func (j *Job) run() {
	j.setScheduled(true)
	err := j.Repository.SaveJob(j)
	if err != nil {
		log.Errorf("%v", err)
	}
	j.StartedChan <- struct{}{}

	if j.GetNextRun() != (time.Time{}) && j.GetNextRun().Before(time.Now().UTC()) {
		j.runJobFunc()
	}

	for {
		if j.GetNextRun() == (time.Time{}) || j.GetNextRun().Before(time.Now().UTC()) {
			j.calcNextRun()
		}

		d := j.GetNextRun().Sub(time.Now().UTC())

		select {
		case <-j.StopChan:
			return
		case <-time.After(d):
			//Run job here
			j.runJobFunc()
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
	if err != nil {
		log.Errorf("%v", err)
	}
	return nil
}

func (j *Job) calcNextRun() {
	j.setNextRun(j.Scheduler.Next(time.Now().UTC()))
	err := j.Repository.SaveJob(j)
	if err != nil {
		log.Errorf("%v", err)
	}
}

// runJobFunc runs the job function
func (j *Job) runJobFunc() {
	log.Infof("Running job %v", j.Name)
	j.JobMutex.Lock()
	defer j.JobMutex.Unlock()

	j.setRunning(true)
	j.setLastRun(time.Now().UTC())
	err := j.Repository.SaveJob(j)
	if err != nil {
		log.Errorf("%v", err)
	}

	j.setLastErr(j.JobFunc())
	j.setRunning(false)
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
	if j.IsScheduled() {
		j.StopChan <- struct{}{}
		j.setScheduled(false)
		err := j.Repository.SaveJob(j)
		if err != nil {
			log.Errorf("%v", err)
		}
	}
}

// HasSchedule returns true if job has a schedule
func (j *Job) HasSchedule() bool {
	if j.Scheduler != nil {
		return true
	}
	return false
}

// IsScheduled return true if scheduled
func (j *Job) IsScheduled() bool {
	j.AccessMutex.RLock()
	defer j.AccessMutex.RUnlock()
	return j.Scheduled
}

func (j *Job) setScheduled(scheduled bool) {
	j.AccessMutex.Lock()
	j.Scheduled = scheduled
	j.AccessMutex.Unlock()
	// err := j.Repository.SaveJob(j)
	// if err != nil {
	// 	log.Errorf("%v", err)
	// }
}

func (j *Job) setRunning(running bool) {
	j.AccessMutex.Lock()
	j.JobRunning = running
	j.AccessMutex.Unlock()
}

// IsRunning returns true if job is running
func (j *Job) IsRunning() bool {
	j.AccessMutex.RLock()
	defer j.AccessMutex.RUnlock()
	return j.JobRunning
}

func (j *Job) setLastErr(e error) {
	j.AccessMutex.Lock()
	j.LastErr = e
	j.AccessMutex.Unlock()
}

//GetLastErr returns the last err from job
func (j *Job) GetLastErr() error {
	j.AccessMutex.RLock()
	defer j.AccessMutex.RUnlock()
	return j.LastErr
}

func (j *Job) setLastRun(t time.Time) {
	j.AccessMutex.Lock()
	j.LastRun = t
	j.AccessMutex.Unlock()
}

//GetLastRun returns the time of the last run
func (j *Job) GetLastRun() time.Time {
	j.AccessMutex.RLock()
	defer j.AccessMutex.RUnlock()
	return j.LastRun
}

func (j *Job) setNextRun(t time.Time) {
	j.AccessMutex.Lock()
	j.NextRun = t
	j.AccessMutex.Unlock()
}

//GetNextRun returns the next time the job will run
func (j *Job) GetNextRun() time.Time {
	j.AccessMutex.RLock()
	defer j.AccessMutex.RUnlock()
	return j.NextRun
}
