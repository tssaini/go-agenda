package db

import (
	"database/sql"
	"fmt"
	"time"

	// needed for mysql
	_ "github.com/go-sql-driver/mysql"
	"github.com/tssaini/go-agenda/job"
)

var db *sql.DB

// InitializeMySQL initialize connection to mysql db
func InitializeMySQL() {
	dBConnection, err := sql.Open("mysql", "")
	if err != nil {
		fmt.Println("Connection Failed!!")
	}
	db = dBConnection
	// dBConnection.SetMaxOpenConns(10)
	// dBConnection.SetMaxIdleConns(5)
	// dBConnection.SetConnMaxLifetime(time.Second * 10)
}

// GetMySQLConnection returns db connection
func GetMySQLConnection() *sql.DB {
	return db
}

// QueryDB query for jobs
func QueryDB(query string) []job.Job {
	// Execute the query
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	jobResults := make([]job.Job, 0)

	for results.Next() {
		var nextRun, lastRun string
		var jobResult job.Job
		err = results.Scan(&jobResult.Name, &nextRun, &jobResult.Locked, &lastRun)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		jobResult.LastRun, err = time.Parse("2006-01-02 15:04:05", lastRun)
		jobResult.NextRun, err = time.Parse("2006-01-02 15:04:05", nextRun)
		jobResults = append(jobResults, jobResult)
	}
	return jobResults
}
