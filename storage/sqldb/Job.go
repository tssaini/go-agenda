package sqldb

import (
	"time"

	"github.com/tssaini/go-agenda"
)

// Job struct represents the db
type Job struct {
	Name    string
	NextRun time.Time
	LastRun time.Time
	Running bool
	LastErr error
}

// GetAllJobs lists all the job from db
func (db *DB) GetAllJobs() []*agenda.Job {
	// Execute the query
	query := "SELECT * FROM agendaJob"
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	var jobResults []*agenda.Job

	for results.Next() {
		var nextRun, lastRun string
		var jobResult Job
		err = results.Scan(&jobResult.Name, &nextRun, &jobResult.Running, &lastRun)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		jobResult.LastRun, err = time.Parse("2006-01-02 15:04:05", lastRun)
		jobResult.NextRun, err = time.Parse("2006-01-02 15:04:05", nextRun)

		j := agenda.Job{Name: jobResult.Name, NextRun: jobResult.NextRun, LastRun: jobResult.LastRun}

		jobResults = append(jobResults, &j)
	}

	return jobResults
}
