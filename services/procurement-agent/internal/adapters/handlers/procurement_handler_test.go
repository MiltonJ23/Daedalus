package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Daedalus/procurement-agent/internal/core/domain"
	"github.com/Daedalus/procurement-agent/internal/core/ports"
	"github.com/Daedalus/procurement-agent/internal/core/services"
)

// ── Mock repositories & supplier ────────────────────────────────────

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
func (m *mockSearchRepo) GetByCacheKey(_ context.Context, k string) (*domain.ProcurementSearch, error) {
	id, ok := m.byCache[k]
	if !ok {
		return nil, nil
	}
	s := m.searches[id]
	return &s, nil
}
func (m *mockSearchRepo) ListByProject(_ context.Context, p string) ([]domain.ProcurementSearch, error) {
	out := []domain.ProcurementSearch{}
	for _, s := range m.searches {
		if p == "" || s.ProjectID == p {
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
func (m *mockResultRepo) ListBySearch(_ context.Context, sid string) ([]domain.EquipmentResult, error) {
	out := []domain.EquipmentResult{}
	for _, r := range m.results {
		if r.SearchID == sid {
			out = append(out, r)
		}
	}
	return out, nil
}
func (m *mockResultRepo) Update(_ context.Context, r domain.EquipmentResult) (domain.EquipmentResult, error) {
	m.results[r.ID] = r
	return r, nil
}

type stubSupplier struct{ name string }

func (s *stubSupplier) Name() string { return s.name }
func (s *stubSupplier) Search(_ context.Context, _ ports.SupplierQuery) ([]domain.EquipmentResult, error) {
	r := domain.NewEquipmentResult("")
	r.Name = "Stub from " + s.name
	r.Model = s.name + "-X1"
	r.Supplier = s.name
	r.PriceUSD = 1000
	r.PriceXAF = 600000
	r.LeadTimeDays = 14
	r.SupplierRating = 4.0
	r.SpecMatch = 0.8
	r.PowerKW = 5
	return []domain.EquipmentResult{r}, nil
}

// ── Setup helpers ───────────────────────────────────────────────────

func setupHandler() *http.ServeMux {
	svc := services.NewProcurementService(
		newMockSearchRepo(),
		newMockResultRepo(),
		[]ports.SupplierCatalog{
			&stubSupplier{name: "A"},
			&stubSupplier{name: "B"},
			&stubSupplier{name: "C"},
		},
		nil,
		nil,
	)
	h := NewProcurementHandler(svc)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return mux
}

func validBody() []byte {
	b, _ := json.Marshal(map[string]interface{}{
		"project_id":     "proj-1",
		"query":          "industrial mixer 500L",
		"category":       "food-processing",
		"max_budget_usd": 20000.0,
	})
	return b
}

func submitSearch(t *testing.T, mux *http.ServeMux) map[string]interface{} {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/procurement/searches", bytes.NewReader(validBody()))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("submit failed: status %d body=%s", rr.Code, rr.Body.String())
	}
	var out map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &out)
	return out
}

// ── Tests ───────────────────────────────────────────────────────────

func TestHandleSubmit_Success(t *testing.T) {
	mux := setupHandler()
	out := submitSearch(t, mux)

	results, _ := out["results"].([]interface{})
	if len(results) != 3 {
		t.Errorf("expected 3 results from 3 suppliers, got %d", len(results))
	}
	if cached, _ := out["cached"].(bool); cached {
		t.Error("first submit should not be cached")
	}
}

func TestHandleSubmit_NoBody(t *testing.T) {
	mux := setupHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/procurement/searches", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestHandleSubmit_Validation(t *testing.T) {
	mux := setupHandler()
	body, _ := json.Marshal(map[string]interface{}{"category": "x"})
	req := httptest.NewRequest(http.MethodPost, "/api/procurement/searches", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", rr.Code)
	}
}

func TestHandleGet_Success(t *testing.T) {
	mux := setupHandler()
	out := submitSearch(t, mux)
	id := out["search"].(map[string]interface{})["id"].(string)

	req := httptest.NewRequest(http.MethodGet, "/api/procurement/searches/"+id, nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestHandleGet_NotFound(t *testing.T) {
	mux := setupHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/procurement/searches/missing", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestHandleList(t *testing.T) {
	mux := setupHandler()
	submitSearch(t, mux)

	req := httptest.NewRequest(http.MethodGet, "/api/procurement/searches?project_id=proj-1", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	var arr []interface{}
	json.Unmarshal(rr.Body.Bytes(), &arr)
	if len(arr) != 1 {
		t.Errorf("expected 1 search, got %d", len(arr))
	}
}

func TestHandleDecision_Approve(t *testing.T) {
	mux := setupHandler()
	out := submitSearch(t, mux)
	results := out["results"].([]interface{})
	resultID := results[0].(map[string]interface{})["id"].(string)

	body, _ := json.Marshal(map[string]string{"decision": "approved"})
	req := httptest.NewRequest(http.MethodPatch, "/api/procurement/results/"+resultID+"/decision", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rr.Code, rr.Body.String())
	}
	var resp map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["decision"] != "approved" {
		t.Errorf("expected approved, got %v", resp["decision"])
	}
}

func TestHandleDecision_Invalid(t *testing.T) {
	mux := setupHandler()
	out := submitSearch(t, mux)
	results := out["results"].([]interface{})
	resultID := results[0].(map[string]interface{})["id"].(string)

	body, _ := json.Marshal(map[string]string{"decision": "maybe"})
	req := httptest.NewRequest(http.MethodPatch, "/api/procurement/results/"+resultID+"/decision", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestCORSMiddleware(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := CORSMiddleware(inner)

	req := httptest.NewRequest(http.MethodOptions, "/api/procurement/searches", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rr.Code)
	}
	if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("expected CORS header")
	}
}
