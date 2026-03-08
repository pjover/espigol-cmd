package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/pjover/espigol/internal/domain/model"
	"github.com/pjover/espigol/internal/domain/ports"
)

// PartnerResponse is the JSON representation of a Partner.
type PartnerResponse struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Surname          string `json:"surname"`
	VATCode          string `json:"vatCode"`
	Email            string `json:"email"`
	Mobile           string `json:"mobile"`
	PartnerType      string `json:"partnerType"`
	RiaNumber        int    `json:"riaNumber"`
	OliveSection     bool   `json:"oliveSection"`
	LivestockSection bool   `json:"livestockSection"`
	AddedOn          string `json:"addedOn"`
}

// PartnerRequest is the JSON payload for creating or updating a Partner.
type PartnerRequest struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Surname          string `json:"surname"`
	VATCode          string `json:"vatCode"`
	Email            string `json:"email"`
	Mobile           string `json:"mobile"`
	PartnerType      string `json:"partnerType"`
	RiaNumber        int    `json:"riaNumber"`
	OliveSection     bool   `json:"oliveSection"`
	LivestockSection bool   `json:"livestockSection"`
	AddedOn          string `json:"addedOn"`
}

// PartnerHandler handles HTTP requests for Partner resources.
type PartnerHandler struct {
	db ports.DbService
}

// NewPartnerHandler creates a new PartnerHandler.
func NewPartnerHandler(db ports.DbService) *PartnerHandler {
	return &PartnerHandler{db: db}
}

// RegisterRoutes registers all partner routes on the given ServeMux.
func (h *PartnerHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /partners", h.listPartners)
	mux.HandleFunc("POST /partners", h.createPartner)
	mux.HandleFunc("GET /partners/{id}", h.getPartner)
	mux.HandleFunc("PUT /partners/{id}", h.updatePartner)
	mux.HandleFunc("DELETE /partners/{id}", h.deletePartner)
}

// @Summary  List all partners
// @Tags     partners
// @Produce  json
// @Success  200 {array}  PartnerResponse
// @Router   /partners [get]
func (h *PartnerHandler) listPartners(w http.ResponseWriter, r *http.Request) {
	partners, err := h.db.GetAllPartners()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]PartnerResponse, 0, len(partners))
	for _, p := range partners {
		resp = append(resp, partnerToResponse(p))
	}
	writeJSON(w, http.StatusOK, resp)
}

// @Summary  Get a partner by ID
// @Tags     partners
// @Produce  json
// @Param    id   path     int  true  "Partner ID"
// @Success  200  {object} PartnerResponse
// @Failure  400  {object} errorResponse
// @Failure  404  {object} errorResponse
// @Router   /partners/{id} [get]
func (h *PartnerHandler) getPartner(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	partner, err := h.db.GetPartnerByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, partnerToResponse(partner))
}

// @Summary  Create a new partner
// @Tags     partners
// @Accept   json
// @Produce  json
// @Param    partner  body     PartnerRequest  true  "Partner payload"
// @Success  201      {object} PartnerResponse
// @Failure  400      {object} errorResponse
// @Failure  500      {object} errorResponse
// @Router   /partners [post]
func (h *PartnerHandler) createPartner(w http.ResponseWriter, r *http.Request) {
	var req PartnerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	partner, err := requestToPartner(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.db.UpsertPartner(partner); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, partnerToResponse(partner))
}

// @Summary  Update an existing partner
// @Tags     partners
// @Accept   json
// @Produce  json
// @Param    id       path     int             true  "Partner ID"
// @Param    partner  body     PartnerRequest  true  "Partner payload"
// @Success  200      {object} PartnerResponse
// @Failure  400      {object} errorResponse
// @Failure  500      {object} errorResponse
// @Router   /partners/{id} [put]
func (h *PartnerHandler) updatePartner(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var req PartnerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}
	req.ID = id

	partner, err := requestToPartner(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.db.UpsertPartner(partner); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, partnerToResponse(partner))
}

// @Summary  Delete a partner by ID
// @Tags     partners
// @Param    id  path  int  true  "Partner ID"
// @Success  204
// @Failure  400  {object} errorResponse
// @Failure  404  {object} errorResponse
// @Router   /partners/{id} [delete]
func (h *PartnerHandler) deletePartner(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	if err := h.db.DeletePartner(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- helpers ---

type errorResponse struct {
	Error string `json:"error"`
}

func parseID(w http.ResponseWriter, r *http.Request) (int, bool) {
	raw := r.PathValue("id")
	id, err := strconv.Atoi(raw)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id: "+raw)
		return 0, false
	}
	return id, true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}

func partnerToResponse(p *model.Partner) PartnerResponse {
	return PartnerResponse{
		ID:               p.ID(),
		Name:             p.Name(),
		Surname:          p.Surname(),
		VATCode:          p.VATCode(),
		Email:            p.Email(),
		Mobile:           p.Mobile(),
		PartnerType:      p.PartnerType().String(),
		RiaNumber:        p.RiaNumber(),
		OliveSection:     p.OliveSection(),
		LivestockSection: p.LivestockSection(),
		AddedOn:          p.AddedOn().Format("2006-01-02"),
	}
}

func requestToPartner(req PartnerRequest) (*model.Partner, error) {
	addedOn, err := time.Parse("2006-01-02", req.AddedOn)
	if err != nil {
		addedOn = time.Now()
	}
	return model.NewPartner(
		req.ID,
		req.Name,
		req.Surname,
		req.VATCode,
		req.Email,
		req.Mobile,
		model.PartnerType(req.PartnerType),
		req.RiaNumber,
		req.OliveSection,
		req.LivestockSection,
		addedOn,
	), nil
}
