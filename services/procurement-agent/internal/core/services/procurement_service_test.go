package services

import (
	"context"
	"errors"
	"testing"

	"github.com/Daedalus/procurement-agent/internal/core/domain"
	"github.com/Daedalus/procurement-agent/internal/core/ports"
)

// ── Mock Repositories ───────────────────────────────────────────────

type mockSearchRepo struct {
	searches map[string]domain.ProcurementSearch
	byCache  map[string]string
}

func newMockSearchRepo() *mockSearchRepo {
	return &mockSearchRepo{
		searches: map[string]domain.ProcurementSearch{},
		byCache:  map[string]string{},
	}
}

func (m *mockSearchRepo) Create(_ context.Context, s domain.ProcurementSearch) (domain.ProcurementSearch, error) {
	m.searches[s.ID] = s
	m.byCache[s.CacheKey] = s.ID
	return s, nil
}
func (m *mockSearchRepo) GetByID(_ context.Context, id string) (*domain.ProcurementSearch, error) {
	s, ok := m.searches[id]
	if !ok {
		return nil, nil
	}
	return &s, nil
}
func (m *mockSearchRepo) GetByCacheKey(_ context.Context, key string) (*domain.ProcurementSearch, error) {
	id, ok := m.byCache[key]
	if !ok {
		return nil, nil
	}
	s := m.searches[id]
	return &s, nil
}
func (m *mockSearchRepo) ListByProject(_ context.Context, projectID string) ([]domain.ProcurementSearch, error) {
	out := []domain.ProcurementSearch{}
	for _, s := range m.searches {
		if projectID == "" || s.ProjectID == projectID {
			out = append(out, s)
		}
	}
	return out, nil
}
func (m *mockSearchRepo) Update(_ context.Context, s domain.ProcurementSearch) (domain.ProcurementSearch, error) {
	m.searches[s.ID] = s
	return s, nil
}

type mockResultRepo struct {
	results map[string]domain.EquipmentResult
}

func newMockResultRepo() *mockResultRepo {
	return &mockResultRepo{results: map[string]domain.EquipmentResult{}}
}

func (m *mockResultRepo) BulkCreate(_ context.Context, items []domain.EquipmentResult) error {
	for _, it := range items {
		m.results[it.ID] = it
	}
	return nil
}
func (m *mockResultRepo) GetByID(_ context.Context, id string) (*domain.EquipmentResult, error) {
	r, ok := m.results[id]
	if !ok {
		return nil, nil
	}
	return &r, nil
}
func (m *mockResultRepo) ListBySearch(_ context.Context, searchID string) ([]domain.EquipmentResult, error) {
	out := []domain.EquipmentResult{}
	for _, r := range m.results {
		if r.SearchID == searchID {
			out = append(out, r)
		}
	}
	return out, nil
}
func (m *mockResultRepo) Update(_ context.Context, r domain.EquipmentResult) (domain.EquipmentResult, error) {
	m.results[r.ID] = r
	return r, nil
}

// ── Mock Supplier ───────────────────────────────────────────────────

type mockSupplier struct {
	name   string
	offers []domain.EquipmentResult
	err    error
}

func (m *mockSupplier) Name() string { return m.name }
func (m *mockSupplier) Search(_ context.Context, _ ports.SupplierQuery) ([]domain.EquipmentResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.offers, nil
}

func makeOffer(name string, priceUSD float64, lead int, rating, spec float64) domain.EquipmentResult {
	r := domain.NewEquipmentResult("")
	r.Name = name
	r.Model = name + "-MDL"
	r.Supplier = name
	r.PriceUSD = priceUSD
	r.LeadTimeDays = lead
	r.SupplierRating = rating
	r.SpecMatch = spec
	return r
}

// ── Helpers ─────────────────────────────────────────────────────────

func newSvc(t *testing.T, supplierList []ports.SupplierCatalog) (*ProcurementService, *mockSearchRepo, *mockResultRepo) {
	t.Helper()
	sr := newMockSearchRepo()
	rr := newMockResultRepo()
	svc := NewProcurementService(sr, rr, supplierList, nil, nil)
	return svc, sr, rr
}

func validSubmit() map[string]interface{} {
	return map[string]interface{}{
		"project_id":     "proj-1",
		"query":          "CNC milling machine 5-axis",
		"category":       "machining",
		"max_budget_usd": float64(50000),
	}
}

// ── Tests ───────────────────────────────────────────────────────────

func TestSubmitSearch_MissingFields(t *testing.T) {
	svc, _, _ := newSvc(t, nil)
	_, err := svc.SubmitSearch(context.Background(), map[string]interface{}{"query": ""})
	if err == nil {
		t.Fatal("expected validation error")
	}
	var ve *domain.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	if _, ok := ve.Errors["project_id"]; !ok {
		t.Error("expected error for project_id")
	}
	if _, ok := ve.Errors["query"]; !ok {
		t.Error("expected error for query")
	}
}

func TestSubmitSearch_RanksAndPersists(t *testing.T) {
	suppliers := []ports.SupplierCatalog{
		&mockSupplier{name: "A", offers: []domain.EquipmentResult{
			makeOffer("A1", 10000, 30, 4.0, 0.7),
		}},
		&mockSupplier{name: "B", offers: []domain.EquipmentResult{
			makeOffer("B1", 5000, 10, 5.0, 0.95),
		}},
		&mockSupplier{name: "C", offers: []domain.EquipmentResult{
			makeOffer("C1", 8000, 20, 3.0, 0.5),
		}},
	}
	svc, _, rr := newSvc(t, suppliers)

	resp, err := svc.SubmitSearch(context.Background(), validSubmit())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Cached {
		t.Error("first submission should not be cached")
	}
	if len(resp.Results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(resp.Results))
	}
	// B1 dominates on every weighted criterion → should rank first.
	if resp.Results[0].Name != "B1" {
		t.Errorf("expected B1 first, got %s (scores: %v)", resp.Results[0].Name, debugScores(resp.Results))
	}
	if resp.Search.Status != domain.SearchStatusCompleted {
		t.Errorf("expected status completed, got %s", resp.Search.Status)
	}
	if len(rr.results) != 3 {
		t.Errorf("expected 3 results persisted, got %d", len(rr.results))
	}
}

func TestSubmitSearch_CacheHitWithin24h(t *testing.T) {
	suppliers := []ports.SupplierCatalog{
		&mockSupplier{name: "A", offers: []domain.EquipmentResult{makeOffer("A", 1000, 5, 4, 0.9)}},
		&mockSupplier{name: "B", offers: []domain.EquipmentResult{makeOffer("B", 1500, 10, 4, 0.8)}},
		&mockSupplier{name: "C", offers: []domain.EquipmentResult{makeOffer("C", 2000, 15, 4, 0.7)}},
	}
	svc, _, _ := newSvc(t, suppliers)

	first, err := svc.SubmitSearch(context.Background(), validSubmit())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second, err := svc.SubmitSearch(context.Background(), validSubmit())
	if err != nil {
		t.Fatalf("unexpected error on second submit: %v", err)
	}
	if !second.Cached {
		t.Error("expected cache hit on identical second submission (FR-PROC-06)")
	}
	if second.Search.ID != first.Search.ID {
		t.Errorf("expected same search ID, got %s vs %s", second.Search.ID, first.Search.ID)
	}
}

func TestDecideResult_ApproveAndReject(t *testing.T) {
	suppliers := []ports.SupplierCatalog{
		&mockSupplier{name: "A", offers: []domain.EquipmentResult{makeOffer("A", 1000, 5, 4, 0.9)}},
		&mockSupplier{name: "B", offers: []domain.EquipmentResult{makeOffer("B", 1500, 10, 4, 0.8)}},
		&mockSupplier{name: "C", offers: []domain.EquipmentResult{makeOffer("C", 2000, 15, 4, 0.7)}},
	}
	svc, _, _ := newSvc(t, suppliers)

	resp, _ := svc.SubmitSearch(context.Background(), validSubmit())
	target := resp.Results[0].ID

	updated, err := svc.DecideResult(context.Background(), target, domain.DecisionApproved)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Decision != domain.DecisionApproved {
		t.Errorf("expected approved, got %s", updated.Decision)
	}

	updated, err = svc.DecideResult(context.Background(), target, domain.DecisionRejected)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Decision != domain.DecisionRejected {
		t.Errorf("expected rejected, got %s", updated.Decision)
	}
}

func TestDecideResult_InvalidValue(t *testing.T) {
	suppliers := []ports.SupplierCatalog{
		&mockSupplier{name: "A", offers: []domain.EquipmentResult{makeOffer("A", 1000, 5, 4, 0.9)}},
		&mockSupplier{name: "B", offers: []domain.EquipmentResult{makeOffer("B", 1500, 10, 4, 0.8)}},
		&mockSupplier{name: "C", offers: []domain.EquipmentResult{makeOffer("C", 2000, 15, 4, 0.7)}},
	}
	svc, _, _ := newSvc(t, suppliers)
	resp, _ := svc.SubmitSearch(context.Background(), validSubmit())

	_, err := svc.DecideResult(context.Background(), resp.Results[0].ID, "maybe")
	if err == nil {
		t.Fatal("expected InvalidDecisionError")
	}
	var de *domain.InvalidDecisionError
	if !errors.As(err, &de) {
		t.Errorf("expected InvalidDecisionError, got %T", err)
	}
}

func TestDecideResult_NotFound(t *testing.T) {
	svc, _, _ := newSvc(t, nil)
	_, err := svc.DecideResult(context.Background(), "missing-id", domain.DecisionApproved)
	if err == nil {
		t.Fatal("expected ResultNotFoundError")
	}
	var nf *domain.ResultNotFoundError
	if !errors.As(err, &nf) {
		t.Errorf("expected ResultNotFoundError, got %T", err)
	}
}

func TestGetSearch_NotFound(t *testing.T) {
	svc, _, _ := newSvc(t, nil)
	_, err := svc.GetSearch(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected SearchNotFoundError")
	}
	var nf *domain.SearchNotFoundError
	if !errors.As(err, &nf) {
		t.Errorf("expected SearchNotFoundError, got %T", err)
	}
}

func TestRanking_Weights(t *testing.T) {
	results := []domain.EquipmentResult{
		makeOffer("cheap-fast-good", 100, 5, 5.0, 1.0),
		makeOffer("expensive-slow-bad", 1000, 50, 1.0, 0.0),
	}
	ranked := rankResults(results)
	if ranked[0].Name != "cheap-fast-good" {
		t.Errorf("expected cheap-fast-good first, got %s", ranked[0].Name)
	}
	// Best should score 1.0 with full weights normalised.
	if ranked[0].Score < 0.99 {
		t.Errorf("expected top score ≈ 1.0, got %f", ranked[0].Score)
	}
	// Worst still earns 0.15 * (1.0/5) = 0.03 from the non-zero supplier rating.
	if ranked[1].Score > 0.05 {
		t.Errorf("expected bottom score ≤ 0.05, got %f", ranked[1].Score)
	}
}

func debugScores(rs []domain.EquipmentResult) []float64 {
	out := make([]float64, len(rs))
	for i, r := range rs {
		out[i] = r.Score
	}
	return out
}
