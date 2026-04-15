package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Daedalus/project-service/internal/core/domain"
	"github.com/Daedalus/project-service/internal/core/services"
)

// ── Mock Repository ─────────────────────────────────────────────────

type mockRepo struct {
	projects map[string]domain.Project
}

func newMockRepo() *mockRepo {
	return &mockRepo{projects: make(map[string]domain.Project)}
}

func (m *mockRepo) Create(_ context.Context, p domain.Project) (domain.Project, error) {
	m.projects[p.ID] = p
	return p, nil
}

func (m *mockRepo) GetByID(_ context.Context, id string) (*domain.Project, error) {
	p, ok := m.projects[id]
	if !ok {
		return nil, nil
	}
	return &p, nil
}

func (m *mockRepo) ListByStatus(_ context.Context, status string) ([]domain.Project, error) {
	var result []domain.Project
	for _, p := range m.projects {
		switch status {
		case "archived":
			if p.ArchivedAt != nil {
				result = append(result, p)
			}
		case "all":
			result = append(result, p)
		default:
			if p.ArchivedAt == nil {
				result = append(result, p)
			}
		}
	}
	if result == nil {
		result = []domain.Project{}
	}
	return result, nil
}

func (m *mockRepo) Update(_ context.Context, p domain.Project) (domain.Project, error) {
	m.projects[p.ID] = p
	return p, nil
}

func (m *mockRepo) HardDelete(_ context.Context, id string) error {
	delete(m.projects, id)
	return nil
}

// ── Helpers ─────────────────────────────────────────────────────────

func setupHandler() (*ProjectHandler, *http.ServeMux) {
	repo := newMockRepo()
	svc := services.NewProjectService(repo)
	handler := NewProjectHandler(svc)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	return handler, mux
}

func validProjectJSON() []byte {
	data, _ := json.Marshal(map[string]interface{}{
		"name":          "Usine Alpha",
		"industry_type": "Agroalimentaire",
		"location":      "Douala, Cameroun",
		"budget":        1500000.0,
		"floor_width":   120.5,
		"floor_depth":   80.0,
	})
	return data
}

func createProject(t *testing.T, mux *http.ServeMux) map[string]interface{} {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/projects", bytes.NewReader(validProjectJSON()))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create failed: status %d, body: %s", rr.Code, rr.Body.String())
	}
	var result map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &result)
	return result
}

// ── Tests ───────────────────────────────────────────────────────────

func TestHandleCreate_Success(t *testing.T) {
	_, mux := setupHandler()
	result := createProject(t, mux)

	if result["name"] != "Usine Alpha" {
		t.Errorf("expected name 'Usine Alpha', got '%v'", result["name"])
	}
	if result["status"] != "active" {
		t.Errorf("expected status 'active', got '%v'", result["status"])
	}
}

func TestHandleCreate_NoBody(t *testing.T) {
	_, mux := setupHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/projects", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestHandleCreate_ValidationError(t *testing.T) {
	_, mux := setupHandler()
	body, _ := json.Marshal(map[string]interface{}{
		"name":          "Bad",
		"industry_type": "BTP",
		"location":      "Yaoundé",
		"budget":        100.0,
		"floor_width":   0.0,
		"floor_depth":   50.0,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/projects", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", rr.Code)
	}
}

func TestHandleList_Empty(t *testing.T) {
	_, mux := setupHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	var result []interface{}
	json.Unmarshal(rr.Body.Bytes(), &result)
	if len(result) != 0 {
		t.Errorf("expected empty list, got %d items", len(result))
	}
}

func TestHandleList_ReturnsCreated(t *testing.T) {
	_, mux := setupHandler()
	createProject(t, mux)

	req := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	var result []interface{}
	json.Unmarshal(rr.Body.Bytes(), &result)
	if len(result) != 1 {
		t.Errorf("expected 1 project, got %d", len(result))
	}
}

func TestHandleGet_Success(t *testing.T) {
	_, mux := setupHandler()
	created := createProject(t, mux)
	id := created["id"].(string)

	req := httptest.NewRequest(http.MethodGet, "/api/projects/"+id, nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestHandleGet_NotFound(t *testing.T) {
	_, mux := setupHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/projects/nonexistent", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestHandleUpdate_Success(t *testing.T) {
	_, mux := setupHandler()
	created := createProject(t, mux)
	id := created["id"].(string)

	body, _ := json.Marshal(map[string]interface{}{"name": "Usine Beta"})
	req := httptest.NewRequest(http.MethodPut, "/api/projects/"+id, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	var result map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &result)
	if result["name"] != "Usine Beta" {
		t.Errorf("expected 'Usine Beta', got '%v'", result["name"])
	}
}

func TestHandleAutoSave(t *testing.T) {
	_, mux := setupHandler()
	created := createProject(t, mux)
	id := created["id"].(string)

	body, _ := json.Marshal(map[string]interface{}{"budget": 2000000.0})
	req := httptest.NewRequest(http.MethodPatch, "/api/projects/"+id+"/autosave", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	var result map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &result)
	if result["message"] != "Auto-saved" {
		t.Errorf("expected 'Auto-saved', got '%v'", result["message"])
	}
}

func TestHandleArchive(t *testing.T) {
	_, mux := setupHandler()
	created := createProject(t, mux)
	id := created["id"].(string)

	req := httptest.NewRequest(http.MethodPatch, "/api/projects/"+id+"/archive", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	var result map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &result)
	if result["status"] != "archived" {
		t.Errorf("expected 'archived', got '%v'", result["status"])
	}
}

func TestHandleRestore(t *testing.T) {
	_, mux := setupHandler()
	created := createProject(t, mux)
	id := created["id"].(string)

	// Archive first
	archReq := httptest.NewRequest(http.MethodPatch, "/api/projects/"+id+"/archive", nil)
	mux.ServeHTTP(httptest.NewRecorder(), archReq)

	// Restore
	req := httptest.NewRequest(http.MethodPatch, "/api/projects/"+id+"/archive?action=restore", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	var result map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &result)
	if result["status"] != "active" {
		t.Errorf("expected 'active', got '%v'", result["status"])
	}
}

func TestHandleDelete_RequiresConfirm(t *testing.T) {
	_, mux := setupHandler()
	created := createProject(t, mux)
	id := created["id"].(string)

	req := httptest.NewRequest(http.MethodDelete, "/api/projects/"+id, nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestHandleDelete_WithConfirm(t *testing.T) {
	_, mux := setupHandler()
	created := createProject(t, mux)
	id := created["id"].(string)

	req := httptest.NewRequest(http.MethodDelete, "/api/projects/"+id+"?confirm=true", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	// Verify deleted
	getReq := httptest.NewRequest(http.MethodGet, "/api/projects/"+id, nil)
	getRR := httptest.NewRecorder()
	mux.ServeHTTP(getRR, getReq)
	if getRR.Code != http.StatusNotFound {
		t.Errorf("expected 404 after delete, got %d", getRR.Code)
	}
}

func TestCORSMiddleware(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := CORSMiddleware(inner)

	req := httptest.NewRequest(http.MethodOptions, "/api/projects", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("expected 204 for OPTIONS, got %d", rr.Code)
	}
	if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("expected CORS Allow-Origin header")
	}
}
