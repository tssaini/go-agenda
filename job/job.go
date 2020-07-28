package job

import (
	"time"
)

// Job defines the config of a job
type Job struct {
	Name    string
	NextRun time.Time
	Locked  bool
	LastRun time.Time
}
