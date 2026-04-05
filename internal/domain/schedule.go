package domain

import "time"

// ComputeDueStatus derives the due status for a drill given its last session.
// Pass nil for lastSession if the drill has never been performed.
// now is injected so this function is pure and testable without time.Now().
func ComputeDueStatus(drill Drill, lastSession *Session, now time.Time) DueStatus {
	status := DueStatus{
		DrillID:   drill.ID,
		DrillName: drill.Name,
	}

	// Unscheduled drills are never overdue.
	if drill.ScheduleDays == 0 {
		return status
	}

	// Never performed — always overdue.
	if lastSession == nil {
		status.IsOverdue = true
		return status
	}

	last := lastSession.PerformedAt.UTC()
	status.LastDoneAt = &last

	daysSince := int(now.UTC().Sub(last).Hours() / 24)
	status.DaysSinceDone = daysSince

	daysOverdue := daysSince - drill.ScheduleDays
	if daysOverdue > 0 {
		status.IsOverdue = true
		status.DaysOverdue = daysOverdue
	}

	return status
}
