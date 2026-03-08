package http_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	httpAdapter "github.com/pjover/espigol/internal/adapters/http"
	"github.com/pjover/espigol/internal/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// partnerFixture returns a sample Partner for use in tests.
func partnerFixture(id int) *model.Partner {
	return model.NewPartner(
		id,
		"Joan",
		"Bosch",
		"43120518A",
		"joan@example.com",
		"+34600000001",
		model.Producer,
		12345,
		true,
		false,
		time.Date(2023, 4, 21, 0, 0, 0, 0, time.UTC),
	)
}

func newTestMux(db *MockDbService) *http.ServeMux {
	mux := http.NewServeMux()
	httpAdapter.NewPartnerHandler(db).RegisterRoutes(mux)
	return mux
}

func partnerRequestBody() httpAdapter.PartnerRequest {
	return httpAdapter.PartnerRequest{
		Name:             "Joan",
		Surname:          "Bosch",
		VATCode:          "43120518A",
		Email:            "joan@example.com",
		Mobile:           "+34600000001",
		PartnerType:      "Productor",
		RiaNumber:        12345,
		OliveSection:     true,
		LivestockSection: false,
		AddedOn:          "2023-04-21",
	}
}

// --- listPartners ---

func TestListPartners_200(t *testing.T) {
	mockDb := new(MockDbService)
	mockDb.On("GetAllPartners").Return([]*model.Partner{partnerFixture(1), partnerFixture(2)}, nil)

	req := httptest.NewRequest(http.MethodGet, "/partners", nil)
	rec := httptest.NewRecorder()
	newTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp []map[string]interface{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Len(t, resp, 2)
	mockDb.AssertExpectations(t)
}

func TestListPartners_500(t *testing.T) {
	mockDb := new(MockDbService)
	mockDb.On("GetAllPartners").Return(nil, fmt.Errorf("db error"))

	req := httptest.NewRequest(http.MethodGet, "/partners", nil)
	rec := httptest.NewRecorder()
	newTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockDb.AssertExpectations(t)
}

// --- getPartner ---

func TestGetPartner_200(t *testing.T) {
	mockDb := new(MockDbService)
	mockDb.On("GetPartnerByID", 1).Return(partnerFixture(1), nil)

	req := httptest.NewRequest(http.MethodGet, "/partners/1", nil)
	rec := httptest.NewRecorder()
	newTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp map[string]interface{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, float64(1), resp["id"])
	assert.Equal(t, "Joan", resp["name"])
	mockDb.AssertExpectations(t)
}

func TestGetPartner_404(t *testing.T) {
	mockDb := new(MockDbService)
	mockDb.On("GetPartnerByID", 99).Return(nil, fmt.Errorf("partner with ID 99 not found"))

	req := httptest.NewRequest(http.MethodGet, "/partners/99", nil)
	rec := httptest.NewRecorder()
	newTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	mockDb.AssertExpectations(t)
}

func TestGetPartner_400_InvalidID(t *testing.T) {
	mockDb := new(MockDbService)

	req := httptest.NewRequest(http.MethodGet, "/partners/abc", nil)
	rec := httptest.NewRecorder()
	newTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// --- createPartner ---

func TestCreatePartner_201(t *testing.T) {
	mockDb := new(MockDbService)
	mockDb.On("UpsertPartner", mock.MatchedBy(func(p *model.Partner) bool {
		return p.ID() == 0 && p.Name() == "Joan" && p.Email() == "joan@example.com"
	})).Return(nil)

	payload, _ := json.Marshal(partnerRequestBody())
	req := httptest.NewRequest(http.MethodPost, "/partners", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	newTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var resp map[string]interface{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "Joan", resp["name"])
	mockDb.AssertExpectations(t)
}

func TestCreatePartner_400_BadJSON(t *testing.T) {
	mockDb := new(MockDbService)

	req := httptest.NewRequest(http.MethodPost, "/partners", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	newTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreatePartner_500(t *testing.T) {
	mockDb := new(MockDbService)
	mockDb.On("UpsertPartner", mock.MatchedBy(func(p *model.Partner) bool {
		return p.Name() == "Joan"
	})).Return(fmt.Errorf("db error"))

	payload, _ := json.Marshal(partnerRequestBody())
	req := httptest.NewRequest(http.MethodPost, "/partners", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	newTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockDb.AssertExpectations(t)
}

// --- updatePartner ---

func TestUpdatePartner_200(t *testing.T) {
	mockDb := new(MockDbService)
	mockDb.On("UpsertPartner", mock.MatchedBy(func(p *model.Partner) bool {
		return p.ID() == 5 && p.Name() == "Joan"
	})).Return(nil)

	payload, _ := json.Marshal(partnerRequestBody())
	req := httptest.NewRequest(http.MethodPut, "/partners/5", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	newTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp map[string]interface{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, float64(5), resp["id"])
	mockDb.AssertExpectations(t)
}

func TestUpdatePartner_400_InvalidID(t *testing.T) {
	mockDb := new(MockDbService)

	payload, _ := json.Marshal(partnerRequestBody())
	req := httptest.NewRequest(http.MethodPut, "/partners/abc", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	newTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdatePartner_400_BadJSON(t *testing.T) {
	mockDb := new(MockDbService)

	req := httptest.NewRequest(http.MethodPut, "/partners/1", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	newTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// --- deletePartner ---

func TestDeletePartner_204(t *testing.T) {
	mockDb := new(MockDbService)
	mockDb.On("DeletePartner", 3).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/partners/3", nil)
	rec := httptest.NewRecorder()
	newTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockDb.AssertExpectations(t)
}

func TestDeletePartner_404(t *testing.T) {
	mockDb := new(MockDbService)
	mockDb.On("DeletePartner", 99).Return(fmt.Errorf("partner with ID 99 not found"))

	req := httptest.NewRequest(http.MethodDelete, "/partners/99", nil)
	rec := httptest.NewRecorder()
	newTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	mockDb.AssertExpectations(t)
}

func TestDeletePartner_400_InvalidID(t *testing.T) {
	mockDb := new(MockDbService)

	req := httptest.NewRequest(http.MethodDelete, "/partners/xyz", nil)
	rec := httptest.NewRecorder()
	newTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
