package agenda

import (
	"fmt"
	"testing"
	"time"

	"github.com/robfig/cron"
)

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

func TestStartJob(t *testing.T) {
	j := NewJob("testJob", func() error {
		return nil
	})

	t.Run("Should fail to start job without schedule", func(t *testing.T) {
		err := j.Start()
		if err == nil {
			t.Errorf("wanted err but got %v", err)
		}
	})

	t.Run("Should fail to start an already running job", func(t *testing.T) {
		j.Schedule(cron.ParseStandard, "* * * * *")
		err := j.Start()
		if err != nil {
			t.Errorf("failed to start job %v", err)
		}
		err = j.Start()
		if err == nil {
			t.Errorf("wanted err but got %v", err)
		}
		j.Stop()
	})

	t.Run("Should run the job", func(t *testing.T) {
		parserCalls = nil
		funcChan := make(chan struct{})
		j := NewJob("testJob", func() error {
			funcChan <- struct{}{}
			return nil
		})
		j.Schedule(parserMock, "* * * * *")
		err := j.Start()
		if err != nil {
			t.Errorf("failed to start job %v", err)
		}
		if !j.isRunning() {
			t.Error("wanted testJob to be running")
		}
		time.Sleep(20 * time.Millisecond) // Allow the test to run
		j.Stop()

		select {
		case <-funcChan:
			fmt.Println("RUn")
		case <-time.After(10 * time.Millisecond):
			t.Error("wanted testJobFunc to have been run")
		}
	})
}

func TestScheduleJob(t *testing.T) {
	j := NewJob("testJob", func() error {
		return nil
	})

	t.Run("Should schedule job", func(t *testing.T) {
		parserCalls = nil
		spec := "* * * * *"
		err := j.Schedule(parserMock, spec)
		if err != nil {
			t.Errorf("wanted err but got %v", err)
		}
		if parserCalls[0] != spec {
			t.Errorf("wanted first call to be %v got %v", spec, parserCalls[0])
		}
	})
}

func TestStopJob(t *testing.T) {
	t.Run("Should stop the job", func(t *testing.T) {
		parserCalls = nil
		j := NewJob("testJob", func() error {
			return nil
		})
		j.Schedule(parserMock, "* * * * *")
		j.Start()
		if !j.isRunning() {
			t.Error("testJob is not running")
		}
		j.Stop()
		if j.isRunning() {
			t.Error("wanted testJob to be stopped")
		}
	})
}
