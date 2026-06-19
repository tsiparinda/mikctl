package models

import (
	"time"
)

type ScriptResult struct {
	Router Router
	Script string
	Tag    string

	StartedAt  time.Time
	FinishedAt time.Time

	Success bool
	Output  string
	Error   string
}
