package agenda

import (
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
	stop       chan struct{}
	newJob     chan scheduled.Task
	cronParser func(string) (cron.Schedule, error)
	db         *sqldb.DB
}

// New creates and returns new agenda object
func New() (*Agenda, error) {
	d, err := sqldb.NewDB("mysql", "root:googlechrome@/goagenda")
	if err != nil {
		return nil, err
	}

	a := Agenda{jobs: make(map[string]scheduled.Task), jobsMutex: &sync.RWMutex{}, running: false, stop: make(chan struct{}), newJob: make(chan scheduled.Task), cronParser: cron.ParseStandard, db: d}
	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)
	return &a, nil
}

// Define new job
func (a *Agenda) Define(name string, jobFunc func() error) {
	log.Infof("Defining job %v", name)
	j, err := scheduled.NewJob(name, jobFunc, a.db)
	if err != nil {
		log.Errorf("Unable to define job %v: %v", name, err)
		return
	}
	a.addJob(name, j)
}

// Schedule the next time the job should run
// TODO
func (a *Agenda) Schedule(name string, time string) error {
	return nil
}

// RepeatEvery define when the job should repeat
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
		if j.HasSchedule() {
			if err := j.Start(); err != nil {
				return err
			}
		}
	}
	a.jobsMutex.RUnlock()
	for {
		select {
		case <-a.stop:
			log.Infof("Stopping agenda loop")
			//stop all jobs
			a.jobsMutex.RLock()
			for _, j := range a.jobs {
				j.Stop()
			}
			a.jobsMutex.RUnlock()
			return nil
		case j := <-a.newJob:
			if err := j.Start(); err != nil {
				return err
			}
		}
	}
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

// Now runs the job provided now
func (a *Agenda) Now(name string) error {
	log.Infof("Starting job %v", name)
	job, err := a.getJob(name)
	if err != nil {
		return err
	}
	// err := job.JobFunc()
	// if err != nil {
	// 	log.Errorf("Error while running %v: %v", name, err)
	// 	return err
	// }
	job.ScheduleNow()
	// log.Infof("Completed job %v successfully", name)
	return nil
}
