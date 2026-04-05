package domain

import "time"

// Session is a timestamped record of performing a drill.
type Session struct {
	ID             string
	DrillID        string
	PerformedAt    time.Time
	StepsCompleted int // how many steps completed (supports progressive drills)
	TotalSteps     int
	Notes          string
}

// DueStatus is computed from a Drill and its most recent Session.
// It is never stored — always derived on demand.
type DueStatus struct {
	DrillID       string
	DrillName     string
	LastDoneAt    *time.Time // nil if never performed
	DaysSinceDone int
	DaysOverdue   int // 0 if not overdue
	IsOverdue     bool
}
