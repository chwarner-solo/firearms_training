package domain_test

import (
	"testing"
	"time"

	"github.com/chwarner-solo/firearms_training/internal/domain"
)

func TestDueStatus_NeverPerformed(t *testing.T) {
	drill := domain.Drill{
		ID:           "1",
		Name:         "1-inch Circle",
		ScheduleDays: 7,
	}

	status := domain.ComputeDueStatus(drill, nil, time.Now())

	if !status.IsOverdue {
		t.Error("expected drill never performed to be overdue")
	}
	if status.LastDoneAt != nil {
		t.Error("expected LastDoneAt to be nil for never-performed drill")
	}
}

func TestDueStatus_PerformedToday(t *testing.T) {
	now := time.Now().UTC()
	drill := domain.Drill{ID: "1", Name: "Dot Drill", ScheduleDays: 3}
	last := &domain.Session{PerformedAt: now}

	status := domain.ComputeDueStatus(drill, last, now)

	if status.IsOverdue {
		t.Errorf("expected drill performed today to not be overdue, got DaysOverdue=%d", status.DaysOverdue)
	}
	if status.DaysSinceDone != 0 {
		t.Errorf("expected DaysSinceDone=0, got %d", status.DaysSinceDone)
	}
}

func TestDueStatus_Overdue(t *testing.T) {
	now := time.Now().UTC()
	drill := domain.Drill{ID: "1", Name: "Dry Fire Fundamentals", ScheduleDays: 3}
	last := &domain.Session{PerformedAt: now.Add(-5 * 24 * time.Hour)}

	status := domain.ComputeDueStatus(drill, last, now)

	if !status.IsOverdue {
		t.Error("expected drill last done 5 days ago with 3-day schedule to be overdue")
	}
	if status.DaysOverdue != 2 {
		t.Errorf("expected DaysOverdue=2, got %d", status.DaysOverdue)
	}
}

func TestDueStatus_NotYetOverdue(t *testing.T) {
	now := time.Now().UTC()
	drill := domain.Drill{ID: "1", Name: "Lena Miculek", ScheduleDays: 7}
	last := &domain.Session{PerformedAt: now.Add(-3 * 24 * time.Hour)}

	status := domain.ComputeDueStatus(drill, last, now)

	if status.IsOverdue {
		t.Error("expected drill last done 3 days ago with 7-day schedule to not be overdue")
	}
	if status.DaysOverdue != 0 {
		t.Errorf("expected DaysOverdue=0, got %d", status.DaysOverdue)
	}
}

func TestDueStatus_Unscheduled(t *testing.T) {
	drill := domain.Drill{ID: "1", Name: "Free Drill", ScheduleDays: 0}

	status := domain.ComputeDueStatus(drill, nil, time.Now())

	if status.IsOverdue {
		t.Error("expected unscheduled drill to never be overdue")
	}
}
