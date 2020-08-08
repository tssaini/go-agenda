package scheduled

import (
	"sync"
	"testing"
	"time"

	"github.com/robfig/cron"
)

func TestNewJob(t *testing.T) {

	t.Run("Should use existing job from db", func(t *testing.T) {
		jr := new(JobRepoMock)
		jobName := "testJob"
		jobFunc := func() error {
			return nil
		}

		jr.On("FindJobByName", jobName).Return(&Job{Name: jobName}, nil)

		j, err := NewJob(jobName, jobFunc, jr)

		if j.Name != jobName {
			t.Errorf("wanted job name to be %v but got %v", jobName, j.Name)
		}
		if j.JobFunc == nil {
			t.Error("wanted job func to be assigned")
		}

		if err != nil {
			t.Errorf("wanted err to be nil but got %v", err)
		}
		jr.AssertExpectations(t)
	})

	t.Run("Should not find job in db", func(t *testing.T) {
		jr := new(JobRepoMock)
		jobName := "testJob"
		jobFunc := func() error {
			return nil
		}

		jr.On("FindJobByName", jobName).Return(nil, nil)
		jr.On("SaveJob", &Job{Name: jobName, Scheduled: false, JobRunning: false}).Return(nil)

		j, err := NewJob(jobName, jobFunc, jr)

		if j.Name != jobName {
			t.Errorf("wanted job name to be %v but got %v", jobName, j.Name)
		}
		if j.JobFunc == nil {
			t.Error("wanted job func to be assigned")
		}

		if err != nil {
			t.Errorf("wanted err to be nil but got %v", err)
		}
		jr.AssertExpectations(t)
	})
}

func TestStartJob(t *testing.T) {

	t.Run("Should fail to start job without schedule", func(t *testing.T) {
		jobName := "testJob"
		j := &Job{Name: jobName, Scheduled: false}
		err := j.Start()
		if err == nil {
			t.Errorf("wanted err but got %v", err)
		}
	})

	t.Run("Should fail to start an already running job", func(t *testing.T) {
		jr := &JobRepoMock{}
		jobName := "testJob"

		sm := &scheduleMock{}
		tNow := time.Now().Add(1 * time.Minute)
		sm.On("Next").Return(tNow)
		parserMock := func(spec string) (cron.Schedule, error) {
			return sm, nil
		}

		j := &Job{Name: jobName, Repository: jr, Scheduled: false, StopChan: make(chan struct{}), StartedChan: make(chan struct{}), JobMutex: sync.Mutex{}, AccessMutex: sync.RWMutex{}}
		jr.On("SaveJob", &Job{Name: jobName, Scheduled: true, JobRunning: false}).Return(nil)
		jr.On("SaveJob", &Job{Name: jobName, Scheduled: true, JobRunning: false, NextRun: tNow}).Return(nil)
		jr.On("SaveJob", &Job{Name: jobName, Scheduled: false, JobRunning: false, NextRun: tNow}).Return(nil)

		j.Schedule(parserMock, "* * * * *")
		err := j.Start()
		if err != nil {
			t.Errorf("failed to start job %v", err)
		}
		time.Sleep(10 * time.Millisecond) //Wait for the run go routine to start
		err = j.Start()
		if err == nil {
			t.Errorf("wanted err but got %v", err)
		}
		j.Stop()
		jr.AssertExpectations(t)
	})
}

func TestStopJob(t *testing.T) {
	t.Run("Should stop the job", func(t *testing.T) {

		jr := &JobRepoMock{}
		jobName := "testJob"

		sm := &scheduleMock{}
		tNow := time.Now().Add(1 * time.Minute)
		sm.On("Next").Return(tNow)
		parserMock := func(spec string) (cron.Schedule, error) {
			return sm, nil
		}

		j := &Job{Name: jobName, Repository: jr, Scheduled: false, StopChan: make(chan struct{}), StartedChan: make(chan struct{}), JobMutex: sync.Mutex{}, AccessMutex: sync.RWMutex{}}
		jr.On("SaveJob", &Job{Name: jobName, Scheduled: true, JobRunning: false}).Return(nil)
		jr.On("SaveJob", &Job{Name: jobName, Scheduled: true, JobRunning: false, NextRun: tNow}).Return(nil)
		jr.On("SaveJob", &Job{Name: jobName, Scheduled: false, JobRunning: false, NextRun: tNow}).Return(nil)

		j.Schedule(parserMock, "* * * * *")
		j.Start()
		if !j.isScheduled() {
			t.Error("testJob is not scheduled")
		}
		time.Sleep(10 * time.Millisecond) //Wait for the run go routine to start
		j.Stop()
		if j.isScheduled() {
			t.Error("wanted testJob to be stopped")
		}
		jr.AssertExpectations(t)
	})
}

func TestRunJobFunc(t *testing.T) {
	t.Run("Should run the job", func(t *testing.T) {

		jr := &JobRepoMock{}
		jobName := "testJob"
		wasRun := false
		jobFunc := func() error {
			wasRun = true
			return nil
		}

		j := &Job{Name: jobName, JobFunc: jobFunc, Repository: jr, Scheduled: false, StopChan: make(chan struct{}), StartedChan: make(chan struct{}), JobMutex: sync.Mutex{}, AccessMutex: sync.RWMutex{}}
		jr.On("SaveJob", &Job{Name: jobName, Scheduled: false, JobRunning: false}).Return(nil)
		jr.On("SaveJob", &Job{Name: jobName, Scheduled: false, JobRunning: true}).Return(nil)

		j.runJobFunc()
		if !wasRun {
			t.Error("Expected job function to have been run")
		}
		jr.AssertExpectations(t)
	})
}
