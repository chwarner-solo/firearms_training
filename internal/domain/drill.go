package domain

// DrillType distinguishes live fire from dry fire drills.
type DrillType string

const (
	DrillTypeLive DrillType = "live"
	DrillTypeDry  DrillType = "dry"
)

// Drill defines a named training drill and its schedule.
type Drill struct {
	ID           string
	Name         string
	Type         DrillType
	Steps        []DrillStep
	ScheduleDays int // target interval in days; 0 means unscheduled
	Notes        string
}

// DrillStep is one step in a drill's sequence.
// Description and PassCriteria are Markdown — rendered full-screen during the drill runner.
type DrillStep struct {
	ID           string
	Order        int
	Description  string // Markdown
	RepsMin      int
	RepsMax      int    // same as RepsMin for exact counts
	PassCriteria string // Markdown — what passing this step looks like
	Notes        string // optional coaching cues
}
