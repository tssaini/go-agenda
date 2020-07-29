package job

import (
	"time"

	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

// Job defines the config of a job
type Job struct {
	Name     string
	JobFunc  func() error
	Schedule cron.Schedule

	NextRun time.Time
	LastRun time.Time
	Stop    chan struct{}
	Running bool
}

// LaunchJob launch the job
func (j *Job) LaunchJob() {
	j.Running = true
	for {
		t := j.Schedule.Next(time.Now())
		d := t.Sub(time.Now())

		select {
		case <-j.Stop:
			log.Infof("Stopping job %v", j.Name)
			j.Running = false
			return
		case runTime := <-time.After(d):
			//Run job here
			log.Infof("Running job %v", j.Name)
			j.LastRun = runTime
			go j.RunJob()
		}
	}
}

// RunJob run the job function
func (j *Job) RunJob() {
	j.JobFunc()
}
