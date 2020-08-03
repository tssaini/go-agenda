package mocks

import (
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
