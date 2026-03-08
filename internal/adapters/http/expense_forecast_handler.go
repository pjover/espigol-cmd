package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pjover/espigol/internal/domain/model"
	"github.com/pjover/espigol/internal/domain/ports"
)

// ExpenseForecastResponse is the JSON representation of an ExpenseForecast.
type ExpenseForecastResponse struct {
	ID             int      `json:"id"`
	PartnerID      int      `json:"partnerId"`
	Concept        string   `json:"concept"`
	Description    string   `json:"description"`
	GrossAmount    float64  `json:"grossAmount"`
	PlannedDate    string   `json:"plannedDate"`
	ExpenseSubtype string   `json:"expenseSubtype"`
	Scope          string   `json:"scope"`
	Attachments    []string `json:"attachments"`
	AddedOn        string   `json:"addedOn"`
}

// ExpenseForecastRequest is the JSON payload for creating or updating an ExpenseForecast.
type ExpenseForecastRequest struct {
	ID             int      `json:"id"`
	PartnerID      int      `json:"partnerId"`
	Concept        string   `json:"concept"`
	Description    string   `json:"description"`
	GrossAmount    float64  `json:"grossAmount"`
	PlannedDate    string   `json:"plannedDate"`
	ExpenseSubtype string   `json:"expenseSubtype"`
	Scope          string   `json:"scope"`
	Attachments    []string `json:"attachments"`
	AddedOn        string   `json:"addedOn"`
}

// ExpenseForecastHandler handles HTTP requests for ExpenseForecast resources.
type ExpenseForecastHandler struct {
	db ports.DbService
}

// NewExpenseForecastHandler creates a new ExpenseForecastHandler.
func NewExpenseForecastHandler(db ports.DbService) *ExpenseForecastHandler {
	return &ExpenseForecastHandler{db: db}
}

// RegisterRoutes registers all expense forecast routes on the given ServeMux.
func (h *ExpenseForecastHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /expense-forecasts", h.listExpenseForecasts)
	mux.HandleFunc("POST /expense-forecasts", h.createExpenseForecast)
	mux.HandleFunc("GET /expense-forecasts/{id}", h.getExpenseForecast)
	mux.HandleFunc("PUT /expense-forecasts/{id}", h.updateExpenseForecast)
	mux.HandleFunc("DELETE /expense-forecasts/{id}", h.deleteExpenseForecast)
}

// @Summary  List all expense forecasts
// @Tags     expense-forecasts
// @Produce  json
// @Success  200 {array}  ExpenseForecastResponse
// @Failure  500 {object} errorResponse
// @Router   /expense-forecasts [get]
func (h *ExpenseForecastHandler) listExpenseForecasts(w http.ResponseWriter, r *http.Request) {
	forecasts, err := h.db.GetAllExpenseForecasts()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]ExpenseForecastResponse, 0, len(forecasts))
	for _, f := range forecasts {
		resp = append(resp, forecastToResponse(f))
	}
	writeJSON(w, http.StatusOK, resp)
}

// @Summary  Get an expense forecast by ID
// @Tags     expense-forecasts
// @Produce  json
// @Param    id   path     int  true  "ExpenseForecast ID"
// @Success  200  {object} ExpenseForecastResponse
// @Failure  400  {object} errorResponse
// @Failure  404  {object} errorResponse
// @Router   /expense-forecasts/{id} [get]
func (h *ExpenseForecastHandler) getExpenseForecast(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	forecast, err := h.db.GetExpenseForecastByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, forecastToResponse(forecast))
}

// @Summary  Create a new expense forecast
// @Tags     expense-forecasts
// @Accept   json
// @Produce  json
// @Param    forecast  body     ExpenseForecastRequest  true  "ExpenseForecast payload"
// @Success  201       {object} ExpenseForecastResponse
// @Failure  400       {object} errorResponse
// @Failure  500       {object} errorResponse
// @Router   /expense-forecasts [post]
func (h *ExpenseForecastHandler) createExpenseForecast(w http.ResponseWriter, r *http.Request) {
	var req ExpenseForecastRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	forecast, err := h.requestToForecast(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.db.UpsertExpenseForecast(forecast); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, forecastToResponse(forecast))
}

// @Summary  Update an existing expense forecast
// @Tags     expense-forecasts
// @Accept   json
// @Produce  json
// @Param    id        path     int                     true  "ExpenseForecast ID"
// @Param    forecast  body     ExpenseForecastRequest  true  "ExpenseForecast payload"
// @Success  200       {object} ExpenseForecastResponse
// @Failure  400       {object} errorResponse
// @Failure  500       {object} errorResponse
// @Router   /expense-forecasts/{id} [put]
func (h *ExpenseForecastHandler) updateExpenseForecast(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var req ExpenseForecastRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}
	req.ID = id

	forecast, err := h.requestToForecast(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.db.UpsertExpenseForecast(forecast); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, forecastToResponse(forecast))
}

// @Summary  Delete an expense forecast by ID
// @Tags     expense-forecasts
// @Param    id  path  int  true  "ExpenseForecast ID"
// @Success  204
// @Failure  400  {object} errorResponse
// @Failure  404  {object} errorResponse
// @Router   /expense-forecasts/{id} [delete]
func (h *ExpenseForecastHandler) deleteExpenseForecast(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	if err := h.db.DeleteExpenseForecast(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- helpers ---

func forecastToResponse(f *model.ExpenseForecast) ExpenseForecastResponse {
	attachments := f.Attachments()
	if attachments == nil {
		attachments = []string{}
	}
	return ExpenseForecastResponse{
		ID:             f.ID(),
		PartnerID:      f.Partner().ID(),
		Concept:        f.Concept(),
		Description:    f.Description(),
		GrossAmount:    f.GrossAmount(),
		PlannedDate:    f.PlannedDate().Format("2006-01-02"),
		ExpenseSubtype: string(f.ExpenseSubtype()),
		Scope:          f.Scope().String(),
		Attachments:    attachments,
		AddedOn:        f.AddedOn().Format("2006-01-02"),
	}
}

func (h *ExpenseForecastHandler) requestToForecast(req ExpenseForecastRequest) (*model.ExpenseForecast, error) {
	partner, err := h.db.GetPartnerByID(req.PartnerID)
	if err != nil {
		return nil, fmt.Errorf("partner with id %d not found: %w", req.PartnerID, err)
	}

	plannedDate, err := time.Parse("2006-01-02", req.PlannedDate)
	if err != nil {
		plannedDate = time.Now()
	}

	addedOn, err := time.Parse("2006-01-02", req.AddedOn)
	if err != nil {
		addedOn = time.Now()
	}

	attachments := req.Attachments
	if attachments == nil {
		attachments = []string{}
	}

	return model.NewExpenseForecast(
		req.ID,
		*partner,
		req.Concept,
		req.Description,
		req.GrossAmount,
		plannedDate,
		model.ExpenseSubtype(req.ExpenseSubtype),
		model.ExpenseScope(req.Scope),
		attachments,
		addedOn,
	), nil
}
