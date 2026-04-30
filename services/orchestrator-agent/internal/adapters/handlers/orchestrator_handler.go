package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Daedalus/orchestrator-agent/internal/core/domain"
	"github.com/Daedalus/orchestrator-agent/internal/core/services"
)

type OrchestratorHandler struct {
	Service *services.OrchestratorService
}

func NewOrchestratorHandler(svc *services.OrchestratorService) *OrchestratorHandler {
	return &OrchestratorHandler{Service: svc}
}

func (h *OrchestratorHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/orchestrator/goals", h.HandleSubmit)
	mux.HandleFunc("GET /api/orchestrator/goals", h.HandleList)
	mux.HandleFunc("GET /api/orchestrator/goals/{id}", h.HandleGet)
	mux.HandleFunc("PATCH /api/orchestrator/tasks/{id}", h.HandleTaskUpdate)
}

func (h *OrchestratorHandler) HandleSubmit(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil || len(data) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "No input data provided"})
		return
	}
	graph, err := h.Service.SubmitGoal(r.Context(), data)
	if err != nil {
		OrchestratorOps.WithLabelValues("submit_goal", "error").Inc()
		h.handleError(w, err)
		return
	}
	OrchestratorOps.WithLabelValues("submit_goal", "success").Inc()
	writeJSON(w, http.StatusCreated, graph)
}

func (h *OrchestratorHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	goals, err := h.Service.ListGoals(r.Context(), userID)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, goals)
}

func (h *OrchestratorHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	graph, err := h.Service.GetGoal(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, graph)
}

func (h *OrchestratorHandler) HandleTaskUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}
	status, _ := data["status"].(string)
	errMsg, _ := data["error"].(string)
	updated, err := h.Service.UpdateTaskStatus(r.Context(), id, status, errMsg)
	if err != nil {
		OrchestratorOps.WithLabelValues("task_update", "error").Inc()
		h.handleError(w, err)
		return
	}
	OrchestratorOps.WithLabelValues("task_update", updated.Status).Inc()
	writeJSON(w, http.StatusOK, updated)
}

func (h *OrchestratorHandler) handleError(w http.ResponseWriter, err error) {
	var (
		gnf *domain.GoalNotFoundError
		tnf *domain.TaskNotFoundError
		ve  *domain.ValidationError
		is  *domain.InvalidStatusError
	)
	switch {
	case errors.As(err, &gnf):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": gnf.Error()})
	case errors.As(err, &tnf):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": tnf.Error()})
	case errors.As(err, &ve):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]interface{}{"errors": ve.Errors})
	case errors.As(err, &is):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": is.Error()})
	default:
		log.Printf("internal error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
