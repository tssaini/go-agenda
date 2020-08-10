# go-agenda
go-agenda is a Go library that allows you to schedule jobs and persist them in a SQL database. It is inspired by the [Agenda](https://github.com/agenda/agenda) Node.js library.

[![Build Status](https://travis-ci.org/tssaini/go-agenda.svg?branch=master)](https://travis-ci.org/tssaini/go-agenda)
[![Go Report Card](https://goreportcard.com/badge/github.com/tssaini/go-agenda)](https://goreportcard.com/report/github.com/tssaini/go-agenda)
[![GoDoc](https://godoc.org/github.com/tssaini/go-agenda?status.svg)](https://godoc.org/github.com/tssaini/go-agenda)

## Features
- Schedules job using a cron expression
- Stores the job metadata such as next run time, last run time and any errors produced in the SQL database

## Usage

```go
db, _ := sql.Open("mysql", "user:test@/goagenda")

agenda, _ := agenda.New(db)
agenda.Define("print hello", func() error {
    fmt.Println("Hello world")
    return nil
})
// Will run the job at the beginning of every minute
agenda.RepeatEvery("print hello", "* * * * *")
agenda.Start()
```

See cmd/main.go for a working example.

