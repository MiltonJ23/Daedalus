package services

import (
	"context"
	"testing"

	"github.com/Daedalus/project-service/internal/core/domain"
)

// ── Mock Repository (Kliops pattern) ────────────────────────────────

type mockProjectRepo struct {
	projects map[string]domain.Project
}

func newMockRepo() *mockProjectRepo {
	return &mockProjectRepo{projects: make(map[string]domain.Project)}
}

func (m *mockProjectRepo) Create(_ context.Context, p domain.Project) (domain.Project, error) {
	m.projects[p.ID] = p
	return p, nil
}

func (m *mockProjectRepo) GetByID(_ context.Context, id string) (*domain.Project, error) {
	p, ok := m.projects[id]
	if !ok {
		return nil, nil
	}
	return &p, nil
}

func (m *mockProjectRepo) ListByStatus(_ context.Context, status string) ([]domain.Project, error) {
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
	return result, nil
}

func (m *mockProjectRepo) Update(_ context.Context, p domain.Project) (domain.Project, error) {
	m.projects[p.ID] = p
	return p, nil
}

func (m *mockProjectRepo) HardDelete(_ context.Context, id string) error {
	delete(m.projects, id)
	return nil
}

// ── Helper ──────────────────────────────────────────────────────────

func validCreateData() map[string]interface{} {
	return map[string]interface{}{
		"name":          "Usine Alpha",
		"industry_type": "Agroalimentaire",
		"location":      "Douala, Cameroun",
		"budget":        float64(1500000),
		"floor_width":   float64(120.5),
		"floor_depth":   float64(80.0),
	}
}

// ── Tests ───────────────────────────────────────────────────────────

func TestCreateProject_Success(t *testing.T) {
	svc := NewProjectService(newMockRepo())
	p, err := svc.CreateProject(context.Background(), validCreateData())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "Usine Alpha" {
		t.Errorf("expected name 'Usine Alpha', got '%s'", p.Name)
	}
	if p.Status != "active" {
		t.Errorf("expected status 'active', got '%s'", p.Status)
	}
	if p.Version != 1 {
		t.Errorf("expected version 1, got %d", p.Version)
	}
}

func TestCreateProject_MissingField(t *testing.T) {
	svc := NewProjectService(newMockRepo())
	data := map[string]interface{}{"name": "Incomplete"}
	_, err := svc.CreateProject(context.Background(), data)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	ve, ok := err.(*domain.ProjectValidationError)
	if !ok {
		t.Fatalf("expected ProjectValidationError, got %T", err)
	}
	if _, exists := ve.Errors["industry_type"]; !exists {
		t.Error("expected error for 'industry_type'")
	}
}

func TestCreateProject_NegativeWidth(t *testing.T) {
	svc := NewProjectService(newMockRepo())
	data := validCreateData()
	data["floor_width"] = float64(-10)
	_, err := svc.CreateProject(context.Background(), data)
	if err == nil {
		t.Fatal("expected validation error")
	}
	ve, ok := err.(*domain.ProjectValidationError)
	if !ok {
		t.Fatalf("expected ProjectValidationError, got %T", err)
	}
	if _, exists := ve.Errors["floor_width"]; !exists {
		t.Error("expected error for 'floor_width'")
	}
}

func TestCreateProject_NegativeBudget(t *testing.T) {
	svc := NewProjectService(newMockRepo())
	data := validCreateData()
	data["budget"] = float64(-100)
	_, err := svc.CreateProject(context.Background(), data)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestGetProject_Found(t *testing.T) {
	repo := newMockRepo()
	svc := NewProjectService(repo)
	created, _ := svc.CreateProject(context.Background(), validCreateData())

	got, err := svc.GetProject(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("expected ID %s, got %s", created.ID, got.ID)
	}
}

func TestGetProject_NotFound(t *testing.T) {
	svc := NewProjectService(newMockRepo())
	_, err := svc.GetProject(context.Background(), "nonexistent-id")
	if err == nil {
		t.Fatal("expected error")
	}
	if _, ok := err.(*domain.ProjectNotFoundError); !ok {
		t.Errorf("expected ProjectNotFoundError, got %T", err)
	}
}

func TestUpdateProject_Success(t *testing.T) {
	repo := newMockRepo()
	svc := NewProjectService(repo)
	created, _ := svc.CreateProject(context.Background(), validCreateData())

	updated, err := svc.UpdateProject(context.Background(), created.ID, map[string]interface{}{
		"name": "Usine Beta",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Name != "Usine Beta" {
		t.Errorf("expected name 'Usine Beta', got '%s'", updated.Name)
	}
	if updated.Version != 2 {
		t.Errorf("expected version 2, got %d", updated.Version)
	}
}

func TestUpdateProject_InvalidDimension(t *testing.T) {
	repo := newMockRepo()
	svc := NewProjectService(repo)
	created, _ := svc.CreateProject(context.Background(), validCreateData())

	_, err := svc.UpdateProject(context.Background(), created.ID, map[string]interface{}{
		"floor_depth": float64(0),
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestAutoSave_WithData(t *testing.T) {
	repo := newMockRepo()
	svc := NewProjectService(repo)
	created, _ := svc.CreateProject(context.Background(), validCreateData())

	result, err := svc.AutoSaveProject(context.Background(), created.ID, map[string]interface{}{
		"budget": float64(2000000),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Message != "Auto-saved" {
		t.Errorf("expected 'Auto-saved', got '%s'", result.Message)
	}
	if result.Version != 2 {
		t.Errorf("expected version 2, got %d", result.Version)
	}
}

func TestAutoSave_Empty(t *testing.T) {
	repo := newMockRepo()
	svc := NewProjectService(repo)
	created, _ := svc.CreateProject(context.Background(), validCreateData())

	result, err := svc.AutoSaveProject(context.Background(), created.ID, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Message != "Nothing to save" {
		t.Errorf("expected 'Nothing to save', got '%s'", result.Message)
	}
}

func TestArchiveRestore(t *testing.T) {
	repo := newMockRepo()
	svc := NewProjectService(repo)
	created, _ := svc.CreateProject(context.Background(), validCreateData())

	archived, err := svc.ArchiveProject(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if archived.Status != "archived" {
		t.Errorf("expected 'archived', got '%s'", archived.Status)
	}
	if !archived.IsArchived() {
		t.Error("expected IsArchived() to return true")
	}

	restored, err := svc.RestoreProject(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if restored.Status != "active" {
		t.Errorf("expected 'active', got '%s'", restored.Status)
	}
}

func TestDeleteProject_WithoutConfirm(t *testing.T) {
	repo := newMockRepo()
	svc := NewProjectService(repo)
	created, _ := svc.CreateProject(context.Background(), validCreateData())

	err := svc.DeleteProject(context.Background(), created.ID, false)
	if err == nil {
		t.Fatal("expected ConfirmationRequiredError")
	}
	if _, ok := err.(*domain.ConfirmationRequiredError); !ok {
		t.Errorf("expected ConfirmationRequiredError, got %T", err)
	}
}

func TestDeleteProject_WithConfirm(t *testing.T) {
	repo := newMockRepo()
	svc := NewProjectService(repo)
	created, _ := svc.CreateProject(context.Background(), validCreateData())

	err := svc.DeleteProject(context.Background(), created.ID, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = svc.GetProject(context.Background(), created.ID)
	if err == nil {
		t.Fatal("expected ProjectNotFoundError after deletion")
	}
}
