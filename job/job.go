package job

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Job defines the config of a job
type Job struct {
	Name    string
	JobFunc func() error
	NextRun time.Time
	Locked  bool
	LastRun time.Time
}

var jobsMap = make(map[string]Job)

// Define instantiate new job
func Define(name string, jobFunc func() error) {
	log.Infof("Creating job %v", name)
	jobsMap[name] = Job{Name: name, JobFunc: jobFunc}
}

// Now runs the job provided now
func Now(name string) {
	log.Infof("Starting job %v", name)
	job := jobsMap[name]
	err := job.JobFunc()
	if err != nil {
		log.Errorf("Error while running %v: %v", name, err)
	} else {
		log.Infof("Completed job %v successfully", name)
	}
}
