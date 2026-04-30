package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Daedalus/procurement-agent/internal/core/domain"
	"github.com/Daedalus/procurement-agent/internal/core/services"
)

// ProcurementHandler — HTTP request handlers (Kliops gateway pattern).
type ProcurementHandler struct {
	Service *services.ProcurementService
}

func NewProcurementHandler(svc *services.ProcurementService) *ProcurementHandler {
	return &ProcurementHandler{Service: svc}
}

// RegisterRoutes wires all procurement routes into the given mux.
func (h *ProcurementHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/procurement/searches", h.HandleSubmit)
	mux.HandleFunc("GET /api/procurement/searches", h.HandleList)
	mux.HandleFunc("GET /api/procurement/searches/{id}", h.HandleGet)
	mux.HandleFunc("PATCH /api/procurement/results/{id}/decision", h.HandleDecision)
	mux.HandleFunc("POST /api/procurement/spec-preview", h.HandleSpecPreview)
}

// ── SUBMIT ──────────────────────────────────────────────────────────

func (h *ProcurementHandler) HandleSubmit(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil || len(data) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "No input data provided"})
		return
	}

	resp, err := h.Service.SubmitSearch(r.Context(), data)
	if err != nil {
		ProcurementOperations.WithLabelValues("submit", "error").Inc()
		h.handleError(w, err)
		return
	}

	if resp.Cached {
		ProcurementOperations.WithLabelValues("submit", "cache_hit").Inc()
	} else {
		ProcurementOperations.WithLabelValues("submit", "success").Inc()
	}

	status := http.StatusCreated
	if resp.Cached {
		status = http.StatusOK
	}
	writeJSON(w, status, toSearchResponse(resp))
}

// ── LIST ────────────────────────────────────────────────────────────

func (h *ProcurementHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("project_id")

	searches, err := h.Service.ListSearches(r.Context(), projectID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	out := make([]map[string]interface{}, len(searches))
	for i, s := range searches {
		out[i] = toSearchSummary(s)
	}
	writeJSON(w, http.StatusOK, out)
}

// ── GET ONE (with optional in-memory filters PB-023) ────────────────

func (h *ProcurementHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	q := r.URL.Query()
	hasFilters := q.Get("country") != "" || q.Get("min_price_usd") != "" ||
		q.Get("max_price_usd") != "" || q.Get("decision") != ""

	if hasFilters {
		opts := services.FilterOptions{
			Country:  q.Get("country"),
			Decision: q.Get("decision"),
		}
		if v := q.Get("min_price_usd"); v != "" {
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				opts.MinPriceUSD = f
			}
		}
		if v := q.Get("max_price_usd"); v != "" {
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				opts.MaxPriceUSD = f
			}
		}
		resp, err := h.Service.FilterResults(r.Context(), id, opts)
		if err != nil {
			h.handleError(w, err)
			return
		}
		ProcurementOperations.WithLabelValues("filter", "success").Inc()
		writeJSON(w, http.StatusOK, toSearchResponse(resp))
		return
	}

	resp, err := h.Service.GetSearch(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toSearchResponse(resp))
}

// ── SPEC PREVIEW (PB-019) ───────────────────────────────────────────

func (h *ProcurementHandler) HandleSpecPreview(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil || len(data) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "No input data provided"})
		return
	}
	query, _ := data["query"].(string)

	spec, err := h.Service.PreviewSpec(r.Context(), query)
	if err != nil {
		ProcurementOperations.WithLabelValues("spec_preview", "error").Inc()
		h.handleError(w, err)
		return
	}
	ProcurementOperations.WithLabelValues("spec_preview", "success").Inc()
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"query":          query,
		"extracted_spec": spec,
	})
}

// ── DECISION (approve / reject) ─────────────────────────────────────

func (h *ProcurementHandler) HandleDecision(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON body"})
		return
	}

	decision, _ := data["decision"].(string)
	updated, err := h.Service.DecideResult(r.Context(), id, decision)
	if err != nil {
		ProcurementOperations.WithLabelValues("decision", "error").Inc()
		h.handleError(w, err)
		return
	}
	ProcurementOperations.WithLabelValues("decision", updated.Decision).Inc()
	writeJSON(w, http.StatusOK, toResultResponse(updated))
}

// ── Error mapping (domain → HTTP) ───────────────────────────────────

func (h *ProcurementHandler) handleError(w http.ResponseWriter, err error) {
	var (
		searchNotFound *domain.SearchNotFoundError
		resultNotFound *domain.ResultNotFoundError
		validation     *domain.ValidationError
		invalidDec     *domain.InvalidDecisionError
	)

	switch {
	case errors.As(err, &searchNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": searchNotFound.Error()})
	case errors.As(err, &resultNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": resultNotFound.Error()})
	case errors.As(err, &validation):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]interface{}{"errors": validation.Errors})
	case errors.As(err, &invalidDec):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": invalidDec.Error()})
	default:
		log.Printf("Internal error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
}

// ── Response serialization ──────────────────────────────────────────

func toSearchResponse(resp services.SearchResponse) map[string]interface{} {
	results := make([]map[string]interface{}, len(resp.Results))
	for i, r := range resp.Results {
		results[i] = toResultResponse(r)
	}
	return map[string]interface{}{
		"search":  toSearchSummary(resp.Search),
		"results": results,
		"cached":  resp.Cached,
	}
}

func toSearchSummary(s domain.ProcurementSearch) map[string]interface{} {
	return map[string]interface{}{
		"id":             s.ID,
		"project_id":     s.ProjectID,
		"query":          s.Query,
		"category":       s.Category,
		"max_budget_usd": s.MaxBudgetUSD,
		"status":         s.Status,
		"extracted_spec": s.ExtractedSpec,
		"expires_at":     s.ExpiresAt.Format("2006-01-02T15:04:05Z"),
		"created_at":     s.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"updated_at":     s.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toResultResponse(r domain.EquipmentResult) map[string]interface{} {
	return map[string]interface{}{
		"id":              r.ID,
		"search_id":       r.SearchID,
		"name":            r.Name,
		"model":           r.Model,
		"supplier":        r.Supplier,
		"supplier_rating": r.SupplierRating,
		"price_usd":       r.PriceUSD,
		"price_xaf":       r.PriceXAF,
		"lead_time_days":  r.LeadTimeDays,
		"spec_match":      r.SpecMatch,
		"score":           r.Score,
		"specifications":  r.Specifications,
		"dimensions": map[string]interface{}{
			"width_m":  r.Dimensions.WidthM,
			"depth_m":  r.Dimensions.DepthM,
			"height_m": r.Dimensions.HeightM,
		},
		"power_kw":   r.PowerKW,
		"country":    r.Country,
		"decision":   r.Decision,
		"created_at": r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"updated_at": r.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
