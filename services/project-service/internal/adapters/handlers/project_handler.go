package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/Daedalus/project-service/internal/core/domain"
	"github.com/Daedalus/project-service/internal/core/services"
)

// ProjectHandler — HTTP request handlers (Kliops gateway pattern).
type ProjectHandler struct {
	Service *services.ProjectService
}

func NewProjectHandler(svc *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{Service: svc}
}

// RegisterRoutes wires all project routes into the given mux under /api prefix.
func (h *ProjectHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/projects", h.HandleCreate)
	mux.HandleFunc("GET /api/projects", h.HandleList)
	mux.HandleFunc("GET /api/projects/{id}", h.HandleGet)
	mux.HandleFunc("PUT /api/projects/{id}", h.HandleUpdate)
	mux.HandleFunc("PATCH /api/projects/{id}/autosave", h.HandleAutoSave)
	mux.HandleFunc("PATCH /api/projects/{id}/archive", h.HandleArchive)
	mux.HandleFunc("DELETE /api/projects/{id}", h.HandleDelete)
}

// ── CREATE ──────────────────────────────────────────────────────────

func (h *ProjectHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil || len(data) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "No input data provided"})
		return
	}

	project, err := h.Service.CreateProject(r.Context(), data)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toProjectResponse(project))
}

// ── LIST ────────────────────────────────────────────────────────────

func (h *ProjectHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	if status == "" {
		status = "active"
	}

	projects, err := h.Service.ListProjects(r.Context(), status)
	if err != nil {
		h.handleError(w, err)
		return
	}

	result := make([]map[string]interface{}, len(projects))
	for i, p := range projects {
		result[i] = toProjectListItem(p)
	}
	writeJSON(w, http.StatusOK, result)
}

// ── GET ONE ─────────────────────────────────────────────────────────

func (h *ProjectHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	project, err := h.Service.GetProject(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toProjectResponse(project))
}

// ── UPDATE ──────────────────────────────────────────────────────────

func (h *ProjectHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil || len(data) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "No input data provided"})
		return
	}

	project, err := h.Service.UpdateProject(r.Context(), id, data)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toProjectResponse(project))
}

// ── AUTO-SAVE ───────────────────────────────────────────────────────

func (h *ProjectHandler) HandleAutoSave(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		data = make(map[string]interface{})
	}

	result, err := h.Service.AutoSaveProject(r.Context(), id, data)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// ── ARCHIVE / RESTORE ───────────────────────────────────────────────

func (h *ProjectHandler) HandleArchive(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	action := r.URL.Query().Get("action")
	if action == "" {
		action = "archive"
	}

	var project domain.Project
	var err error

	if action == "restore" {
		project, err = h.Service.RestoreProject(r.Context(), id)
	} else {
		project, err = h.Service.ArchiveProject(r.Context(), id)
	}

	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toProjectResponse(project))
}

// ── DELETE ───────────────────────────────────────────────────────────

func (h *ProjectHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	confirm := strings.ToLower(r.URL.Query().Get("confirm")) == "true"

	if err := h.Service.DeleteProject(r.Context(), id, confirm); err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Project permanently deleted"})
}

// ── Error mapping (domain → HTTP) ───────────────────────────────────

func (h *ProjectHandler) handleError(w http.ResponseWriter, err error) {
	var notFound *domain.ProjectNotFoundError
	var validation *domain.ProjectValidationError
	var confirmation *domain.ConfirmationRequiredError

	switch {
	case errors.As(err, &notFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": notFound.Error()})
	case errors.As(err, &validation):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]interface{}{"errors": validation.Errors})
	case errors.As(err, &confirmation):
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": confirmation.Error(),
			"hint":  confirmation.Hint,
		})
	default:
		log.Printf("Internal error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
}

// ── Response serialization ──────────────────────────────────────────

func toProjectResponse(p domain.Project) map[string]interface{} {
	resp := map[string]interface{}{
		"id":              p.ID,
		"name":            p.Name,
		"industry_type":   p.IndustryType,
		"location":        p.Location,
		"budget":          p.Budget,
		"floor_width":     p.FloorWidth,
		"floor_depth":     p.FloorDepth,
		"target_capacity": p.TargetCapacity,
		"status":          p.Status,
		"version":         p.Version,
		"is_archived":     p.IsArchived(),
		"archived_at":     nil,
		"created_at":      p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"updated_at":      p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if p.ArchivedAt != nil {
		resp["archived_at"] = p.ArchivedAt.Format("2006-01-02T15:04:05Z")
	}
	return resp
}

func toProjectListItem(p domain.Project) map[string]interface{} {
	return map[string]interface{}{
		"id":            p.ID,
		"name":          p.Name,
		"industry_type": p.IndustryType,
		"location":      p.Location,
		"budget":        p.Budget,
		"floor_width":   p.FloorWidth,
		"floor_depth":   p.FloorDepth,
		"status":        p.Status,
		"is_archived":   p.IsArchived(),
		"created_at":    p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"updated_at":    p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
