package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/tssaini/go-agenda/scheduled"
)

// DBMock mock the db
type DBMock struct {
	mock.Mock
}

// FindAllJobs mock the FindAllJob
func (db *DBMock) FindAllJobs() ([]*scheduled.Job, error) {
	return nil, nil
}
