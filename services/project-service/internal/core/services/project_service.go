package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Daedalus/project-service/internal/core/domain"
	"github.com/Daedalus/project-service/internal/core/ports"
)

// ProjectService — use-case orchestrator (Kliops services pattern).
type ProjectService struct {
	repo ports.ProjectRepository
}

func NewProjectService(repo ports.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

// CreateProject validates input and persists a new project.
func (s *ProjectService) CreateProject(ctx context.Context, data map[string]interface{}) (domain.Project, error) {
	if err := s.validateCreate(data); err != nil {
		return domain.Project{}, err
	}

	var targetCap *string
	if v, ok := data["target_capacity"].(string); ok {
		targetCap = &v
	}

	budget, _ := toFloat64(data["budget"])
	width, _ := toFloat64(data["floor_width"])
	depth, _ := toFloat64(data["floor_depth"])

	project := domain.NewProject(
		data["name"].(string),
		data["industry_type"].(string),
		data["location"].(string),
		budget, width, depth,
		targetCap,
	)

	created, err := s.repo.Create(ctx, project)
	if err != nil {
		return domain.Project{}, fmt.Errorf("failed to create project: %w", err)
	}

	log.Printf("Project created: %s (%s)", created.Name, created.ID)
	return created, nil
}

// ListProjects returns projects filtered by status.
func (s *ProjectService) ListProjects(ctx context.Context, status string) ([]domain.Project, error) {
	return s.repo.ListByStatus(ctx, status)
}

// GetProject retrieves a single project by ID.
func (s *ProjectService) GetProject(ctx context.Context, id string) (domain.Project, error) {
	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Project{}, fmt.Errorf("failed to get project: %w", err)
	}
	if project == nil {
		return domain.Project{}, &domain.ProjectNotFoundError{ProjectID: id}
	}
	return *project, nil
}

// UpdateProject applies a partial update to an existing project.
func (s *ProjectService) UpdateProject(ctx context.Context, id string, data map[string]interface{}) (domain.Project, error) {
	project, err := s.GetProject(ctx, id)
	if err != nil {
		return domain.Project{}, err
	}

	if err := s.validateUpdate(data); err != nil {
		return domain.Project{}, err
	}

	project.ApplyUpdate(data)

	updated, err := s.repo.Update(ctx, project)
	if err != nil {
		return domain.Project{}, fmt.Errorf("failed to update project %s: %w", id, err)
	}

	log.Printf("Project updated: %s (v%d)", id, updated.Version)
	return updated, nil
}

// AutoSaveResponse is the return type for autosave operations.
type AutoSaveResponse struct {
	Message   string    `json:"message"`
	Version   int       `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AutoSaveProject handles incremental auto-save (optimized PATCH).
func (s *ProjectService) AutoSaveProject(ctx context.Context, id string, data map[string]interface{}) (AutoSaveResponse, error) {
	project, err := s.GetProject(ctx, id)
	if err != nil {
		return AutoSaveResponse{}, err
	}

	if len(data) == 0 {
		return AutoSaveResponse{
			Message:   "Nothing to save",
			Version:   project.Version,
			UpdatedAt: project.UpdatedAt,
		}, nil
	}

	if err := s.validateUpdate(data); err != nil {
		return AutoSaveResponse{}, err
	}

	project.ApplyUpdate(data)

	updated, err := s.repo.Update(ctx, project)
	if err != nil {
		return AutoSaveResponse{}, fmt.Errorf("failed to autosave project %s: %w", id, err)
	}

	return AutoSaveResponse{
		Message:   "Auto-saved",
		Version:   updated.Version,
		UpdatedAt: updated.UpdatedAt,
	}, nil
}

// ArchiveProject soft-deletes a project.
func (s *ProjectService) ArchiveProject(ctx context.Context, id string) (domain.Project, error) {
	project, err := s.GetProject(ctx, id)
	if err != nil {
		return domain.Project{}, err
	}

	project.SoftDelete()

	updated, err := s.repo.Update(ctx, project)
	if err != nil {
		return domain.Project{}, fmt.Errorf("failed to archive project %s: %w", id, err)
	}
	return updated, nil
}

// RestoreProject un-archives a project.
func (s *ProjectService) RestoreProject(ctx context.Context, id string) (domain.Project, error) {
	project, err := s.GetProject(ctx, id)
	if err != nil {
		return domain.Project{}, err
	}

	project.Restore()

	updated, err := s.repo.Update(ctx, project)
	if err != nil {
		return domain.Project{}, fmt.Errorf("failed to restore project %s: %w", id, err)
	}
	return updated, nil
}

// DeleteProject permanently removes a project (requires explicit confirmation).
func (s *ProjectService) DeleteProject(ctx context.Context, id string, confirm bool) error {
	if !confirm {
		return domain.NewConfirmationRequired()
	}

	_, err := s.GetProject(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.HardDelete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete project %s: %w", id, err)
	}

	log.Printf("Project permanently deleted: %s", id)
	return nil
}

// ── Validation helpers ──────────────────────────────────────────────

func (s *ProjectService) validateCreate(data map[string]interface{}) error {
	errors := make(map[string]string)

	required := []string{"name", "industry_type", "location", "budget", "floor_width", "floor_depth"}
	for _, field := range required {
		if _, ok := data[field]; !ok {
			errors[field] = fmt.Sprintf("%s is required", field)
		}
	}

	s.validateDimensions(data, errors)

	if len(errors) > 0 {
		return &domain.ProjectValidationError{Errors: errors}
	}
	return nil
}

func (s *ProjectService) validateUpdate(data map[string]interface{}) error {
	errors := make(map[string]string)
	s.validateDimensions(data, errors)
	if len(errors) > 0 {
		return &domain.ProjectValidationError{Errors: errors}
	}
	return nil
}

func (s *ProjectService) validateDimensions(data map[string]interface{}, errors map[string]string) {
	if v, ok := data["floor_width"]; ok {
		if f, ok := toFloat64(v); ok && f <= 0 {
			errors["floor_width"] = "Floor width must be a positive non-zero value."
		}
	}
	if v, ok := data["floor_depth"]; ok {
		if f, ok := toFloat64(v); ok && f <= 0 {
			errors["floor_depth"] = "Floor depth must be a positive non-zero value."
		}
	}
	if v, ok := data["budget"]; ok {
		if f, ok := toFloat64(v); ok && f < 0 {
			errors["budget"] = "Budget must be zero or positive."
		}
	}
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
