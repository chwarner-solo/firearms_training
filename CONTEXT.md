# Project Context — Drill Tracker

## What We're Building Right Now

A personal firearms training tool. Track drill definitions (live and dry fire), schedule how often each drill should be done, log when you last performed it, and surface what's overdue. Includes a full-screen drill runner mode for use on an iPad or large display during training sessions.

---

## Platform

- **Backend**: Go
- **Storage**: SQLite (local, single file, no server required)
- **Frontend v1**: PWA (Progressive Web App — Go serves static files + API, browser renders)
- **Frontend v2 (planned)**: Flutter — backend API unchanged
- **Deployment**: Single binary, runs locally

---

## Architecture

Ports & Adapters (Hexagonal). Domain logic has zero knowledge of HTTP, SQLite, or any delivery mechanism.

```
┌─────────────────────────────────────────────┐
│                  Domain                      │
│  Drill | DrillStep | Session | Schedule      │
│  Pure Go structs, no external dependencies  │
└────────────────┬────────────────────────────┘
                 │ ports (interfaces)
     ┌───────────┼───────────┐
     │           │           │
 HTTP API    SQLite Repo   (future)
  adapter     adapter
     │
  PWA / Flutter
```

---

## Domain Model (settled)

### Drill
```go
type Drill struct {
    ID           string
    Name         string
    Type         DrillType   // Live | Dry
    Steps        []DrillStep
    ScheduleDays int         // target interval in days; 0 = unscheduled
    Notes        string
}

type DrillType string
const (
    DrillTypeLive DrillType = "live"
    DrillTypeDry  DrillType = "dry"
)
```

### DrillStep
```go
type DrillStep struct {
    ID           string
    Order        int
    Description  string  // Markdown — rendered full-screen during drill runner
    RepsMin      int
    RepsMax      int     // same as RepsMin if exact
    PassCriteria string  // Markdown — what "passing" this step looks like
    Notes        string  // optional coaching cues
}
```

### Session
```go
type Session struct {
    ID             string
    DrillID        string
    PerformedAt    time.Time
    StepsCompleted int
    TotalSteps     int
    Notes          string
}
```

### DueStatus (computed, not stored)
```go
type DueStatus struct {
    DrillID       string
    DrillName     string
    LastDoneAt    *time.Time  // nil if never performed
    DaysSinceDone int
    DaysOverdue   int         // 0 if not overdue
    IsOverdue     bool
}
```

---

## Ports (Interfaces)

```go
type DrillRepository interface {
    Save(ctx context.Context, d domain.Drill) error
    FindByID(ctx context.Context, id string) (domain.Drill, error)
    FindAll(ctx context.Context) ([]domain.Drill, error)
    Delete(ctx context.Context, id string) error
}

type SessionRepository interface {
    Save(ctx context.Context, s domain.Session) error
    FindByDrillID(ctx context.Context, drillID string) ([]domain.Session, error)
    FindLastByDrillID(ctx context.Context, drillID string) (*domain.Session, error)
}

type ScheduleService interface {
    GetDueStatus(ctx context.Context, drillID string) (domain.DueStatus, error)
    GetOverdueDrills(ctx context.Context) ([]domain.DueStatus, error)
}
```

---

## What Exists

```
internal/
  domain/
    drill.go          — Drill, DrillStep types
    session.go        — Session, DueStatus types
    schedule.go       — ComputeDueStatus() — pure function, no dependencies
    schedule_test.go  — 5 tests, all passing (never done, today, overdue,
                        not yet overdue, unscheduled)
  ports/
    ports.go          — DrillRepository, SessionRepository, ScheduleService interfaces
  adapters/
    sqlite/           — empty, not yet implemented
    http/             — empty, not yet implemented
cmd/
  drill-tracker/      — empty, not yet implemented
go.mod                — module: github.com/chrisarmstrong/drill-tracker
```

**Test status**: `go test ./internal/domain/...` → 5/5 PASS

---

## What We're Working On

> SQLite adapters — implementing DrillRepository and SessionRepository against a local SQLite file.

---

## What's Next

1. ~~Repo scaffold — Go module, folder structure~~ ✓
2. ~~Domain structs with unit tests~~ ✓
3. ~~`ComputeDueStatus` implementation (pure domain logic)~~ ✓
4. SQLite adapter for `DrillRepository`
5. SQLite adapter for `SessionRepository`
6. `ScheduleService` implementation wiring repos + domain logic
7. HTTP API (basic CRUD + due status endpoint)
8. PWA drill runner UI (step-through mode, markdown rendering, pass/fail)

---

## Key Decisions (Do Not Revisit Without Discussion)

| Decision | Rationale |
|---|---|
| DrillStep.Description is Markdown | Rendered full-screen on iPad during drill runner |
| RepsMin/RepsMax (not exact) | Dry fire uses ranges (10-15), live fire may be exact |
| PassCriteria per step | Enables progressive drills ("all 5 in the 1" circle before advancing") |
| Schedule is interval-based (days) | Directly answers "am I overdue?" |
| DueStatus is computed, not stored | Derived from Drill + last Session — no sync problems |
| SQLite local file | No server, no setup, single binary |
| PWA first, Flutter later | API is the port — frontend is swappable |
| ComputeDueStatus takes `now time.Time` | Injected so function is pure and testable without time.Now() |
| Module path: github.com/chwarner-solo/firearms_training | Matches actual GitHub repo: github.com/chwarner-solo/firearms_training |

---

## Conventions

- TDD — test first, always
- Table-driven tests in Go
- No framework in the domain layer — pure Go
- IDs are UUIDs (string)
- All times are UTC
- Errors are wrapped with context (`fmt.Errorf("finding drill: %w", err)`)

---

## How To Use This File

Paste this at the start of every session. Update **What Exists** and **What We're Working On** together at the end of each session. Add settled decisions to the decisions table. Never delete rows — mark superseded if they change.
