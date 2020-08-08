package agenda

import (
	"sync"
	"testing"
	"time"

	"github.com/robfig/cron"
	"github.com/tssaini/go-agenda/scheduled"
)

// func TestRepeatEvery(t *testing.T) {

// 	a := Agenda{jobs: make(map[string]*job.Job),
// 		running:    false,
// 		stop:       make(chan struct{}),
// 		newJob:     make(chan *job.Job),
// 		cronParser: cron.ParseStandard}

// }

func TestStart(t *testing.T) {
	t.Run("Starts the agenda loop", func(t *testing.T) {
		a := Agenda{jobs: make(map[string]scheduled.Task),
			jobsMutex:  &sync.RWMutex{},
			running:    false,
			cronParser: cron.ParseStandard,
		}
		a.Start()
		if !a.running {
			t.Errorf("wanted running %v got %v", true, a.running)
		}
	})

	t.Run("Test running agenda twice", func(t *testing.T) {
		a := Agenda{jobs: make(map[string]scheduled.Task),
			jobsMutex:  &sync.RWMutex{},
			running:    false,
			cronParser: cron.ParseStandard,
		}
		a.Start()
		time.Sleep(10 * time.Millisecond)
		err := a.Start()
		if err == nil {
			t.Errorf("wanted err got %v", err)
		}
	})
}

func TestStop(t *testing.T) {
	t.Run("Stop the agenda loop", func(t *testing.T) {
		a := Agenda{jobs: make(map[string]scheduled.Task),
			jobsMutex:  &sync.RWMutex{},
			running:    false,
			cronParser: cron.ParseStandard,
		}
		a.Start()
		time.Sleep(10 * time.Millisecond)
		a.Stop()
		time.Sleep(10 * time.Millisecond)
		if a.running {
			t.Errorf("wanted running %v got %v", false, a.running)
		}
	})

}

func TestDefine(t *testing.T) {
	t.Run("should define a new job", func(t *testing.T) {
		jr := &scheduled.JobRepoMock{}
		a := Agenda{jobs: make(map[string]scheduled.Task),
			jobsMutex:  &sync.RWMutex{},
			running:    false,
			cronParser: cron.ParseStandard,
			jr:         jr,
		}
		jobName := "TestJob"
		jr.On("FindJobByName", jobName).Return(nil, nil)
		jr.On("SaveJob", &scheduled.Job{Name: jobName, Scheduled: false, JobRunning: false}).Return(nil)

		a.Define(jobName, func() error { return nil })
		if _, err := a.getJob("TestJob"); err != nil {
			t.Errorf("wanted TestJob to be defined got %v", err)
		}
	})
}
