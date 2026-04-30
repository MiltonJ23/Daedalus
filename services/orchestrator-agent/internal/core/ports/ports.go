package ports

import (
	"context"

	"github.com/Daedalus/orchestrator-agent/internal/core/domain"
)

type GoalRepository interface {
	Create(ctx context.Context, g domain.Goal) (domain.Goal, error)
	GetByID(ctx context.Context, id string) (*domain.Goal, error)
	List(ctx context.Context, userID string) ([]domain.Goal, error)
	Update(ctx context.Context, g domain.Goal) (domain.Goal, error)
}

type TaskRepository interface {
	BulkCreate(ctx context.Context, tasks []domain.SubTask) error
	GetByID(ctx context.Context, id string) (*domain.SubTask, error)
	ListByGoal(ctx context.Context, goalID string) ([]domain.SubTask, error)
	Update(ctx context.Context, t domain.SubTask) (domain.SubTask, error)
}

// TaskPublisher dispatches a sub-task to its specialised agent
// (Redis Streams in production; an in-memory stub in tests).
type TaskPublisher interface {
	Publish(ctx context.Context, t domain.SubTask) (streamID string, err error)
}
