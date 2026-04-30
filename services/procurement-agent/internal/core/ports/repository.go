package ports

import (
	"context"

	"github.com/Daedalus/procurement-agent/internal/core/domain"
)

// SearchRepository — port interface for procurement search persistence.
type SearchRepository interface {
	Create(ctx context.Context, search domain.ProcurementSearch) (domain.ProcurementSearch, error)
	GetByID(ctx context.Context, id string) (*domain.ProcurementSearch, error)
	GetByCacheKey(ctx context.Context, cacheKey string) (*domain.ProcurementSearch, error)
	ListByProject(ctx context.Context, projectID string) ([]domain.ProcurementSearch, error)
	Update(ctx context.Context, search domain.ProcurementSearch) (domain.ProcurementSearch, error)
}

// ResultRepository — port interface for equipment result persistence.
type ResultRepository interface {
	BulkCreate(ctx context.Context, results []domain.EquipmentResult) error
	GetByID(ctx context.Context, id string) (*domain.EquipmentResult, error)
	ListBySearch(ctx context.Context, searchID string) ([]domain.EquipmentResult, error)
	Update(ctx context.Context, result domain.EquipmentResult) (domain.EquipmentResult, error)
}
