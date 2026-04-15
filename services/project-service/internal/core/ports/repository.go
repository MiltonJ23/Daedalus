package ports

import (
	"context"

	"github.com/Daedalus/project-service/internal/core/domain"
)

// ProjectRepository — port interface for project persistence (Kliops pattern).
type ProjectRepository interface {
	Create(ctx context.Context, project domain.Project) (domain.Project, error)
	GetByID(ctx context.Context, id string) (*domain.Project, error)
	ListByStatus(ctx context.Context, status string) ([]domain.Project, error)
	Update(ctx context.Context, project domain.Project) (domain.Project, error)
	HardDelete(ctx context.Context, id string) error
}
