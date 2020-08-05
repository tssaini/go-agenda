package scheduled

import (
	"time"

	"github.com/robfig/cron"
	"github.com/stretchr/testify/mock"
)

// JobMock mocks the Agenda.Job struct
type JobMock struct {
	mock.Mock
}

//Start implements Start() of ScheduleTask Interface
func (j *JobMock) Start() error {
	return nil
}

//Schedule implements Schedule() of ScheduleTask Interface
func (j *JobMock) Schedule(parser func(string) (cron.Schedule, error), spec string) error {
	return nil
}

//Stop implements Stop() of ScheduleTask Interface
func (j *JobMock) Stop() {

}

//ScheduleNow implements ScheduleNow() of ScheduleTask Interface
func (j *JobMock) ScheduleNow() {

}

//HasSchedule implements HasSchedule() of ScheduleTask Interface
func (j *JobMock) HasSchedule() bool {
	return false
}



// JobRepoMock mock the db
type JobRepoMock struct {
	mock.Mock
}

// FindAllJobs mock the FindAllJob
func (jr *JobRepoMock) FindAllJobs() ([]*Job, error) {
	return nil, nil
}

// FindJobByName mock the FindJobByName
func (jr *JobRepoMock) FindJobByName(jobName string) (*Job, error) {
	return nil, nil
}

// SaveJob mock the SaveJob
func (jr *JobRepoMock) SaveJob(j *Job) error {
	return nil
}



type scheduleMock struct {
}

func (s scheduleMock) Next(ti time.Time) time.Time {
	return time.Now().Add(10 * time.Millisecond)
}

var parserCalls []string

func parserMock(spec string) (cron.Schedule, error) {
	parserCalls = append(parserCalls, spec)
	return scheduleMock{}, nil
}
