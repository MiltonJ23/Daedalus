package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Daedalus/procurement-agent/internal/core/domain"
	"github.com/Daedalus/procurement-agent/internal/core/ports"
)

// supplierFanOutTimeout caps the per-search supplier fan-out (FR-PROC-02
// / PB-020 ≤30s SLA).
const supplierFanOutTimeout = 30 * time.Second

// minResultsTarget is the PB-020 acceptance criterion floor.
const minResultsTarget = 5

// ProcurementService — use-case orchestrator (Kliops services pattern).
//
// Responsibilities:
//   - validate the incoming request
//   - serve from cache when a recent (≤24h) identical search exists (FR-PROC-06)
//   - otherwise fan-out to all configured supplier adapters concurrently (FR-PROC-02)
//   - merge, convert prices to XAF, rank with the composite score (FR-PROC-03)
//   - persist the search and its ranked results
type ProcurementService struct {
	searches  ports.SearchRepository
	results   ports.ResultRepository
	suppliers []ports.SupplierCatalog
	currency  ports.CurrencyConverter
	extractor ports.SpecExtractor
}

func NewProcurementService(
	searches ports.SearchRepository,
	results ports.ResultRepository,
	suppliers []ports.SupplierCatalog,
	currency ports.CurrencyConverter,
	extractor ports.SpecExtractor,
) *ProcurementService {
	return &ProcurementService{
		searches:  searches,
		results:   results,
		suppliers: suppliers,
		currency:  currency,
		extractor: extractor,
	}
}

// SearchResponse is the bundled view returned to callers.
type SearchResponse struct {
	Search  domain.ProcurementSearch  `json:"search"`
	Results []domain.EquipmentResult  `json:"results"`
	Cached  bool                      `json:"cached"`
}

// SubmitSearch is the main entry point (FR-PROC-01).
func (s *ProcurementService) SubmitSearch(ctx context.Context, data map[string]interface{}) (SearchResponse, error) {
	if err := s.validateSubmit(data); err != nil {
		return SearchResponse{}, err
	}

	projectID := data["project_id"].(string)
	query := strings.TrimSpace(data["query"].(string))
	category, _ := data["category"].(string)
	maxBudget, _ := toFloat64(data["max_budget_usd"])

	cacheKey := buildCacheKey(projectID, query, category, maxBudget)

	// FR-PROC-06: serve from cache when a non-expired identical search exists.
	if cached, err := s.searches.GetByCacheKey(ctx, cacheKey); err == nil && cached != nil && !cached.IsExpired() {
		results, err := s.results.ListBySearch(ctx, cached.ID)
		if err != nil {
			return SearchResponse{}, fmt.Errorf("failed to load cached results: %w", err)
		}
		log.Printf("Procurement cache hit: search=%s (%d results)", cached.ID, len(results))
		return SearchResponse{Search: *cached, Results: results, Cached: true}, nil
	}

	// PB-019: extract a structured spec from the natural-language query.
	spec := s.extractSpec(ctx, query)

	search := domain.NewProcurementSearch(projectID, query, category, cacheKey, maxBudget, spec)
	created, err := s.searches.Create(ctx, search)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("failed to create search: %w", err)
	}

	supplierQuery := ports.SupplierQuery{
		Query:        query,
		Category:     category,
		MaxBudgetUSD: maxBudget,
	}

	// PB-020: cap supplier fan-out at 30s; partial results are still ranked.
	fanCtx, cancel := context.WithTimeout(ctx, supplierFanOutTimeout)
	defer cancel()
	offers := s.fanOutSuppliers(fanCtx, supplierQuery)

	if len(s.suppliers) < 3 {
		// FR-PROC-02 mandates ≥3 supplier sources; log loudly without failing.
		log.Printf("WARNING: only %d supplier(s) configured (FR-PROC-02 requires ≥3)", len(s.suppliers))
	}

	for i := range offers {
		offers[i].SearchID = created.ID
		if offers[i].PriceXAF == 0 && s.currency != nil {
			if xaf, err := s.currency.USDToXAF(ctx, offers[i].PriceUSD); err == nil {
				offers[i].PriceXAF = xaf
			}
		}
	}

	ranked := rankResults(offers)

	if len(ranked) < minResultsTarget {
		log.Printf("WARNING: search %s returned %d offers (PB-020 target ≥%d)",
			created.ID, len(ranked), minResultsTarget)
	}

	if err := s.results.BulkCreate(ctx, ranked); err != nil {
		created.MarkFailed()
		_, _ = s.searches.Update(ctx, created)
		return SearchResponse{}, fmt.Errorf("failed to persist results: %w", err)
	}

	created.MarkCompleted()
	updated, err := s.searches.Update(ctx, created)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("failed to update search status: %w", err)
	}

	log.Printf("Procurement search complete: %s (%d offers from %d suppliers)",
		updated.ID, len(ranked), len(s.suppliers))

	return SearchResponse{Search: updated, Results: ranked, Cached: false}, nil
}

// GetSearch returns a search and its ranked results.
func (s *ProcurementService) GetSearch(ctx context.Context, id string) (SearchResponse, error) {
	search, err := s.searches.GetByID(ctx, id)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("failed to get search: %w", err)
	}
	if search == nil {
		return SearchResponse{}, &domain.SearchNotFoundError{SearchID: id}
	}

	results, err := s.results.ListBySearch(ctx, id)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("failed to load results: %w", err)
	}
	return SearchResponse{Search: *search, Results: results, Cached: false}, nil
}

// ListSearches returns all searches (optionally scoped to a project).
func (s *ProcurementService) ListSearches(ctx context.Context, projectID string) ([]domain.ProcurementSearch, error) {
	return s.searches.ListByProject(ctx, projectID)
}

// DecideResult applies an approve / reject decision (FR-PROC-07).
func (s *ProcurementService) DecideResult(ctx context.Context, resultID, decision string) (domain.EquipmentResult, error) {
	result, err := s.results.GetByID(ctx, resultID)
	if err != nil {
		return domain.EquipmentResult{}, fmt.Errorf("failed to get result: %w", err)
	}
	if result == nil {
		return domain.EquipmentResult{}, &domain.ResultNotFoundError{ResultID: resultID}
	}

	if err := result.ApplyDecision(decision); err != nil {
		return domain.EquipmentResult{}, err
	}

	updated, err := s.results.Update(ctx, *result)
	if err != nil {
		return domain.EquipmentResult{}, fmt.Errorf("failed to persist decision: %w", err)
	}

	log.Printf("Result %s → %s", updated.ID, updated.Decision)
	return updated, nil
}

// FilterOptions is applied in-memory to existing results (PB-023).
type FilterOptions struct {
	Country     string
	MinPriceUSD float64
	MaxPriceUSD float64
	Decision    string
}

// FilterResults loads existing ranked results for a search and filters them
// in-memory — *no* re-search is triggered, satisfying the PB-023 acceptance
// criterion ("Filters apply to existing results without re-running full search").
func (s *ProcurementService) FilterResults(ctx context.Context, searchID string, opts FilterOptions) (SearchResponse, error) {
	search, err := s.searches.GetByID(ctx, searchID)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("failed to get search: %w", err)
	}
	if search == nil {
		return SearchResponse{}, &domain.SearchNotFoundError{SearchID: searchID}
	}

	all, err := s.results.ListBySearch(ctx, searchID)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("failed to load results: %w", err)
	}

	filtered := make([]domain.EquipmentResult, 0, len(all))
	for _, r := range all {
		if opts.Country != "" && !strings.EqualFold(r.Country, opts.Country) {
			continue
		}
		if opts.MinPriceUSD > 0 && r.PriceUSD < opts.MinPriceUSD {
			continue
		}
		if opts.MaxPriceUSD > 0 && r.PriceUSD > opts.MaxPriceUSD {
			continue
		}
		if opts.Decision != "" && !strings.EqualFold(r.Decision, opts.Decision) {
			continue
		}
		filtered = append(filtered, r)
	}

	return SearchResponse{Search: *search, Results: filtered, Cached: true}, nil
}

// PreviewSpec runs the spec extractor without persisting anything (PB-019).
func (s *ProcurementService) PreviewSpec(ctx context.Context, query string) (map[string]interface{}, error) {
	if strings.TrimSpace(query) == "" {
		return nil, &domain.ValidationError{Errors: map[string]string{"query": "query is required"}}
	}
	return s.extractSpec(ctx, query), nil
}

// extractSpec wraps the optional SpecExtractor port with a safe default so
// callers never have to nil-check.
func (s *ProcurementService) extractSpec(ctx context.Context, query string) map[string]interface{} {
	if s.extractor == nil {
		return map[string]interface{}{"raw_query": query}
	}
	spec, err := s.extractor.Extract(ctx, query)
	if err != nil {
		log.Printf("spec extractor failed: %v", err)
		return map[string]interface{}{"raw_query": query}
	}
	return spec
}

// ── internal helpers ────────────────────────────────────────────────

func (s *ProcurementService) fanOutSuppliers(ctx context.Context, q ports.SupplierQuery) []domain.EquipmentResult {
	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		offers  []domain.EquipmentResult
	)

	for _, sup := range s.suppliers {
		wg.Add(1)
		go func(adapter ports.SupplierCatalog) {
			defer wg.Done()
			batch, err := adapter.Search(ctx, q)
			if err != nil {
				log.Printf("supplier %s failed: %v", adapter.Name(), err)
				return
			}
			mu.Lock()
			offers = append(offers, batch...)
			mu.Unlock()
		}(sup)
	}

	wg.Wait()
	return offers
}

func (s *ProcurementService) validateSubmit(data map[string]interface{}) error {
	errs := make(map[string]string)

	required := []string{"project_id", "query"}
	for _, f := range required {
		v, ok := data[f].(string)
		if !ok || strings.TrimSpace(v) == "" {
			errs[f] = fmt.Sprintf("%s is required", f)
		}
	}

	if v, ok := data["max_budget_usd"]; ok {
		if f, ok := toFloat64(v); ok && f < 0 {
			errs["max_budget_usd"] = "max_budget_usd must be zero or positive"
		}
	}

	if len(errs) > 0 {
		return &domain.ValidationError{Errors: errs}
	}
	return nil
}

func buildCacheKey(projectID, query, category string, maxBudget float64) string {
	raw := fmt.Sprintf("%s|%s|%s|%.2f",
		strings.ToLower(projectID),
		strings.ToLower(strings.TrimSpace(query)),
		strings.ToLower(strings.TrimSpace(category)),
		maxBudget,
	)
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func toFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	default:
		return 0, false
	}
}
