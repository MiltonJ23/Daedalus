package ports

import (
	"context"

	"github.com/Daedalus/procurement-agent/internal/core/domain"
)

// SupplierQuery is what we pass down to each supplier adapter.
type SupplierQuery struct {
	Query        string
	Category     string
	MaxBudgetUSD float64
}

// SupplierCatalog — port interface for an external supplier source (FR-PROC-02).
// Each adapter must implement this interface; the service queries at least 3
// of them concurrently and merges the offers before ranking.
type SupplierCatalog interface {
	Name() string
	Search(ctx context.Context, q SupplierQuery) ([]domain.EquipmentResult, error)
}

// CurrencyConverter converts amounts between currencies.
// USD → XAF is required by FR-PROC-04 (every result carries both prices).
type CurrencyConverter interface {
	USDToXAF(ctx context.Context, amountUSD float64) (float64, error)
}
