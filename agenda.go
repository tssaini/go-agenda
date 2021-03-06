package agenda

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/tssaini/go-agenda/storage/sqldb"

	"github.com/tssaini/go-agenda/scheduled"

	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

// Agenda struct stores all jobs
type Agenda struct {
	jobs       map[string]scheduled.Task
	jobsMutex  *sync.RWMutex
	running    bool
	cronParser func(string) (cron.Schedule, error)
	jr         scheduled.JobRepository
}

// New creates and returns new agenda object
func New(db *sql.DB) (*Agenda, error) {
	jr, err := sqldb.NewJobRepository(db)
	if err != nil {
		return nil, err
	}

	a := Agenda{jobs: make(map[string]scheduled.Task), jobsMutex: &sync.RWMutex{}, running: false, cronParser: cron.ParseStandard, jr: jr}
	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)
	return &a, nil
}

// Define a new job
func (a *Agenda) Define(name string, jobFunc func() error) error {
	j, err := scheduled.NewJob(name, jobFunc, a.jr)
	if err != nil {
		return err
	}
	a.addJob(name, j)
	return nil
}

// Schedule the next time the job should run
// TODO
// func (a *Agenda) Schedule(name string, time time.Time) error {
// 	return nil
// }

// RepeatEvery define when the job should repeat.
// The spec parameter is a cron formated string which specifies when the job is run.
func (a *Agenda) RepeatEvery(name string, spec string) error {
	j, err := a.getJob(name)
	if err != nil {
		return err
	}
	err = j.Schedule(a.cronParser, spec)
	if err != nil {
		return err
	}
	if a.running {
		err := j.Start()
		if err != nil {
			return err
		}
	}
	return nil
}

// Start agenda - schedules all defined jobs
func (a *Agenda) Start() error {
	if a.running {
		return errors.New("Agenda is already running")
	}
	a.running = true
	// run all jobs
	a.jobsMutex.RLock()
	for _, j := range a.jobs {
		if j.HasSchedule() {
			if err := j.Start(); err != nil {
				return err
			}
		}
	}
	a.jobsMutex.RUnlock()

	return nil
}

// Stop all jobs defined
func (a *Agenda) Stop() error {
	//stop all jobs
	a.jobsMutex.RLock()
	for _, j := range a.jobs {
		j.Stop()
	}
	a.jobsMutex.RUnlock()
	a.running = false
	return nil
}

func (a *Agenda) getJob(name string) (scheduled.Task, error) {
	a.jobsMutex.RLock()
	j, ok := a.jobs[name]
	a.jobsMutex.RUnlock()
	if ok {
		return j, nil
	}
	return j, fmt.Errorf("Job %v does not exist", name)
}

func (a *Agenda) addJob(name string, j scheduled.Task) {
	//TODO: check if the job already exists
	a.jobsMutex.Lock()
	a.jobs[name] = j
	a.jobsMutex.Unlock()
}

// Now runs the provided job
func (a *Agenda) Now(name string) error {
	job, err := a.getJob(name)
	if err != nil {
		return err
	}

	job.ScheduleNow()

	return nil
}
