package sqldb

import (
	"fmt"
	"time"

	"github.com/tssaini/go-agenda/scheduled"
)

// Job struct represents the db
type Job struct {
	Name    string
	NextRun time.Time
	LastRun time.Time
	Running bool
	LastErr error
}

// FindAllJobs lists all the job from db
func (db *DB) FindAllJobs() ([]*scheduled.Job, error) {
	// Execute the query
	query := "SELECT * FROM agendaJob"
	results, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	var jobResults []*scheduled.Job

	for results.Next() {
		var nextRun, lastRun string
		var jobResult Job
		err = results.Scan(&jobResult.Name, &nextRun, &jobResult.Running, &lastRun)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		jobResult.LastRun, err = time.Parse("2006-01-02 15:04:05", lastRun)
		jobResult.NextRun, err = time.Parse("2006-01-02 15:04:05", nextRun)

		j := scheduled.Job{Name: jobResult.Name, NextRun: jobResult.NextRun, LastRun: jobResult.LastRun}

		jobResults = append(jobResults, &j)
	}

	return jobResults, nil
}

// FindJobByName returns the job given the name
func (db *DB) FindJobByName(jobName string) (*scheduled.Job, error) {
	fmt.Println("FindJobByName")
	return nil, nil
}

// SaveJob saves the provided job to db
func (db *DB) SaveJob(j *scheduled.Job) error {
	fmt.Println("SaveJob")
	return nil
}
