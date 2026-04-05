package ports

import (
	"context"

	"github.com/chwarner-solo/firearms_training/internal/domain"
)

// DrillRepository is the driven port for drill persistence.
type DrillRepository interface {
	Save(ctx context.Context, d domain.Drill) error
	FindByID(ctx context.Context, id string) (domain.Drill, error)
	FindAll(ctx context.Context) ([]domain.Drill, error)
	Delete(ctx context.Context, id string) error
}

// SessionRepository is the driven port for session persistence.
type SessionRepository interface {
	Save(ctx context.Context, s domain.Session) error
	FindByDrillID(ctx context.Context, drillID string) ([]domain.Session, error)
	FindLastByDrillID(ctx context.Context, drillID string) (*domain.Session, error)
}

// ScheduleService is the driving port for due-status queries.
type ScheduleService interface {
	GetDueStatus(ctx context.Context, drillID string) (domain.DueStatus, error)
	GetOverdueDrills(ctx context.Context) ([]domain.DueStatus, error)
}
