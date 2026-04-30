package domain

import (
	"time"

	"github.com/google/uuid"
)

// Decision values for an EquipmentResult (FR-PROC-07).
const (
	DecisionPending  = "pending"
	DecisionApproved = "approved"
	DecisionRejected = "rejected"
)

// SearchStatus values for a ProcurementSearch.
const (
	SearchStatusPending   = "pending"
	SearchStatusCompleted = "completed"
	SearchStatusFailed    = "failed"
	SearchStatusCached    = "cached"
)

// ProcurementSearch — a request submitted by a user to source equipment.
type ProcurementSearch struct {
	ID            string                 `json:"id"`
	ProjectID     string                 `json:"project_id"`
	Query         string                 `json:"query"`
	Category      string                 `json:"category"`
	MaxBudgetUSD  float64                `json:"max_budget_usd"`
	Status        string                 `json:"status"`
	CacheKey      string                 `json:"cache_key"`
	ExtractedSpec map[string]interface{} `json:"extracted_spec"`
	ExpiresAt     time.Time              `json:"expires_at"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// NewProcurementSearch creates a search with sensible defaults.
// Cached results are valid for 24 hours (FR-PROC-06).
func NewProcurementSearch(projectID, query, category, cacheKey string, maxBudgetUSD float64, spec map[string]interface{}) ProcurementSearch {
	now := time.Now().UTC()
	if spec == nil {
		spec = map[string]interface{}{}
	}
	return ProcurementSearch{
		ID:            uuid.New().String(),
		ProjectID:     projectID,
		Query:         query,
		Category:      category,
		MaxBudgetUSD:  maxBudgetUSD,
		Status:        SearchStatusPending,
		CacheKey:      cacheKey,
		ExtractedSpec: spec,
		ExpiresAt:     now.Add(24 * time.Hour),
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// IsExpired reports whether the cached search results have aged out.
func (s *ProcurementSearch) IsExpired() bool {
	return time.Now().UTC().After(s.ExpiresAt)
}

// MarkCompleted marks the search as completed.
func (s *ProcurementSearch) MarkCompleted() {
	s.Status = SearchStatusCompleted
	s.UpdatedAt = time.Now().UTC()
}

// MarkFailed marks the search as failed.
func (s *ProcurementSearch) MarkFailed() {
	s.Status = SearchStatusFailed
	s.UpdatedAt = time.Now().UTC()
}

// Dimensions describes the physical footprint of a piece of equipment in metres.
type Dimensions struct {
	WidthM  float64 `json:"width_m"`
	DepthM  float64 `json:"depth_m"`
	HeightM float64 `json:"height_m"`
}

// EquipmentResult — a ranked equipment offer returned by a supplier.
// Required fields per FR-PROC-04: name, model, supplier, USD+XAF price,
// lead time, specs, dimensions, power.
type EquipmentResult struct {
	ID             string                 `json:"id"`
	SearchID       string                 `json:"search_id"`
	Name           string                 `json:"name"`
	Model          string                 `json:"model"`
	Supplier       string                 `json:"supplier"`
	SupplierRating float64                `json:"supplier_rating"`
	PriceUSD       float64                `json:"price_usd"`
	PriceXAF       float64                `json:"price_xaf"`
	LeadTimeDays   int                    `json:"lead_time_days"`
	SpecMatch      float64                `json:"spec_match"`
	Score          float64                `json:"score"`
	Specifications map[string]interface{} `json:"specifications"`
	Dimensions     Dimensions             `json:"dimensions"`
	PowerKW        float64                `json:"power_kw"`
	Country        string                 `json:"country"`
	Decision       string                 `json:"decision"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// NewEquipmentResult creates a result attached to a search.
func NewEquipmentResult(searchID string) EquipmentResult {
	now := time.Now().UTC()
	return EquipmentResult{
		ID:             uuid.New().String(),
		SearchID:       searchID,
		Specifications: map[string]interface{}{},
		Decision:       DecisionPending,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// ApplyDecision sets the user decision (approve / reject), validating the value.
func (r *EquipmentResult) ApplyDecision(decision string) error {
	switch decision {
	case DecisionApproved, DecisionRejected:
		r.Decision = decision
		r.UpdatedAt = time.Now().UTC()
		return nil
	default:
		return &InvalidDecisionError{Value: decision}
	}
}
