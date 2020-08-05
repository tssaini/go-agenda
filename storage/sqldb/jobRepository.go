package sqldb

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/tssaini/go-agenda/scheduled"
)

// JobRepository sql db impl for job repository
type JobRepository struct {
	*sql.DB
}

// NewJobRepository create new job respoitory given sql connection
func NewJobRepository(db *sql.DB) (*JobRepository, error) {
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &JobRepository{db}, nil
}

// Job struct represents the db
type Job struct {
	Name       string
	NextRun    time.Time
	LastRun    time.Time
	Scheduled  bool
	JobRunning bool
	LastErr    error
}

// FindAllJobs lists all the job from db
func (db *JobRepository) FindAllJobs() ([]*scheduled.Job, error) {
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
		err = results.Scan(&jobResult.Name, &nextRun, &jobResult.JobRunning, &lastRun)
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
func (db *JobRepository) FindJobByName(jobName string) (*scheduled.Job, error) {
	fmt.Println("FindJobByName")
	var jobResult Job
	var lastErr, nextRun, lastRun sql.NullString

	query := "SELECT name, nextRun, lastRun, scheduled, jobRunning, lastErr FROM agendaJob where name = ?"
	// Execute the query
	err := db.QueryRow(query, jobName).Scan(&jobResult.Name, &nextRun, &lastRun, &jobResult.Scheduled, &jobResult.JobRunning, &lastErr)

	if err != nil {
		return nil, err
	}
	if lastRun.Valid {
		jobResult.LastRun, err = time.Parse("2006-01-02 15:04:05", lastRun.String)
	}
	if nextRun.Valid {
		jobResult.NextRun, err = time.Parse("2006-01-02 15:04:05", nextRun.String)
	}
	if lastErr.Valid {
		jobResult.LastErr = errors.New(lastErr.String)
	}
	j := scheduled.Job{Name: jobResult.Name,
		NextRun:     jobResult.NextRun,
		LastRun:     jobResult.LastRun,
		JobMutex:    sync.Mutex{},
		StopChan:    make(chan struct{}),
		StartedChan: make(chan struct{}),
		Scheduled:   jobResult.Scheduled,
		JobRunning:  jobResult.JobRunning,
		AccessMutex: sync.RWMutex{},
		LastErr:     jobResult.LastErr,
	}
	return &j, nil
}

// SaveJob saves the provided job to db
func (db *JobRepository) SaveJob(j *scheduled.Job) error {
	fmt.Println("SaveJob")

	query := "SELECT name FROM agendaJob where name = ?"
	var jobName string
	err := db.QueryRow(query, jobName).Scan(&jobName)

	lastErr := fmt.Sprintf("%v", j.LastErr)

	if err != nil {
		sqlStatement := "UPDATE agendaJob SET nextRun = ?, lastRun = ?, scheduled = ?, jobRunning = ?, lastErr = ? WHERE name = ?"
		_, err := db.Exec(sqlStatement, j.NextRun.Format("2006-01-02 15:04:05"), j.LastRun.Format("2006-01-02 15:04:05"), j.Scheduled, j.JobRunning, lastErr, j.Name)
		if err != nil {
			return err
		}
	} else {
		sqlStatement := "INSERT INTO agendaJob (name, nextRun, lastRun, scheduled, jobRunning, lastErr) VALUES (?, ?, ?, ?, ?, ?)"
		_, err := db.Exec(sqlStatement, j.Name, j.NextRun.Format("2006-01-02 15:04:05"), j.LastRun.Format("2006-01-02 15:04:05"), j.Scheduled, j.JobRunning, lastErr)
		if err != nil {
			return err
		}
	}
	return nil
}
