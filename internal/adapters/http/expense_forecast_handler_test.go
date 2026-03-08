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

func expenseForecastFixture(id int) *model.ExpenseForecast {
	partner := partnerFixture(1)
	return model.NewExpenseForecast(
		id,
		*partner,
		"Fertilitzants",
		"Compra de fertilitzants",
		1500.50,
		time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
		model.ExpenseSubtypeA6,
		model.ExpenseScopeOliveSection,
		[]string{"invoice.pdf"},
		time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
	)
}

func newForecastTestMux(db *MockDbService) *http.ServeMux {
	mux := http.NewServeMux()
	httpAdapter.NewExpenseForecastHandler(db).RegisterRoutes(mux)
	return mux
}

func forecastRequestBody() httpAdapter.ExpenseForecastRequest {
	return httpAdapter.ExpenseForecastRequest{
		PartnerID:      1,
		Concept:        "Fertilitzants",
		Description:    "Compra de fertilitzants",
		GrossAmount:    1500.50,
		PlannedDate:    "2024-06-15",
		ExpenseSubtype: "[a6] Despeses de fertilitzants, productes d'alimentació animal i ormejos",
		Scope:          "Secció d'oliva",
		Attachments:    []string{"invoice.pdf"},
		AddedOn:        "2024-01-10",
	}
}

func TestListExpenseForecasts_200(t *testing.T) {
	mockDb := new(MockDbService)
	mockDb.On("GetAllExpenseForecasts").Return(
		[]*model.ExpenseForecast{expenseForecastFixture(1), expenseForecastFixture(2)}, nil,
	)

	req := httptest.NewRequest(http.MethodGet, "/expense-forecasts", nil)
	rec := httptest.NewRecorder()
	newForecastTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp []map[string]interface{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Len(t, resp, 2)
	mockDb.AssertExpectations(t)
}

func TestListExpenseForecasts_500(t *testing.T) {
	mockDb := new(MockDbService)
	mockDb.On("GetAllExpenseForecasts").Return(nil, fmt.Errorf("db error"))

	req := httptest.NewRequest(http.MethodGet, "/expense-forecasts", nil)
	rec := httptest.NewRecorder()
	newForecastTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockDb.AssertExpectations(t)
}

func TestGetExpenseForecast_200(t *testing.T) {
	mockDb := new(MockDbService)
	mockDb.On("GetExpenseForecastByID", 1).Return(expenseForecastFixture(1), nil)

	req := httptest.NewRequest(http.MethodGet, "/expense-forecasts/1", nil)
	rec := httptest.NewRecorder()
	newForecastTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp map[string]interface{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, float64(1), resp["id"])
	assert.Equal(t, "Fertilitzants", resp["concept"])
	mockDb.AssertExpectations(t)
}

func TestGetExpenseForecast_404(t *testing.T) {
	mockDb := new(MockDbService)
	mockDb.On("GetExpenseForecastByID", 99).Return(nil, fmt.Errorf("expense forecast with ID 99 not found"))

	req := httptest.NewRequest(http.MethodGet, "/expense-forecasts/99", nil)
	rec := httptest.NewRecorder()
	newForecastTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	mockDb.AssertExpectations(t)
}

func TestGetExpenseForecast_400_InvalidID(t *testing.T) {
	mockDb := new(MockDbService)

	req := httptest.NewRequest(http.MethodGet, "/expense-forecasts/abc", nil)
	rec := httptest.NewRecorder()
	newForecastTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateExpenseForecast_201(t *testing.T) {
	mockDb := new(MockDbService)
	mockDb.On("GetPartnerByID", 1).Return(partnerFixture(1), nil)
	mockDb.On("UpsertExpenseForecast", mock.MatchedBy(func(f *model.ExpenseForecast) bool {
		return f.ID() == 0 && f.Concept() == "Fertilitzants" && f.Partner().ID() == 1
	})).Return(nil)

	payload, _ := json.Marshal(forecastRequestBody())
	req := httptest.NewRequest(http.MethodPost, "/expense-forecasts", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	newForecastTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var resp map[string]interface{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "Fertilitzants", resp["concept"])
	mockDb.AssertExpectations(t)
}

func TestCreateExpenseForecast_400_BadJSON(t *testing.T) {
	mockDb := new(MockDbService)

	req := httptest.NewRequest(http.MethodPost, "/expense-forecasts", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	newForecastTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateExpenseForecast_400_PartnerNotFound(t *testing.T) {
	mockDb := new(MockDbService)
	mockDb.On("GetPartnerByID", 1).Return(nil, fmt.Errorf("partner with ID 1 not found"))

	payload, _ := json.Marshal(forecastRequestBody())
	req := httptest.NewRequest(http.MethodPost, "/expense-forecasts", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	newForecastTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockDb.AssertExpectations(t)
}

func TestCreateExpenseForecast_500(t *testing.T) {
	mockDb := new(MockDbService)
	mockDb.On("GetPartnerByID", 1).Return(partnerFixture(1), nil)
	mockDb.On("UpsertExpenseForecast", mock.MatchedBy(func(f *model.ExpenseForecast) bool {
		return f.Concept() == "Fertilitzants"
	})).Return(fmt.Errorf("db error"))

	payload, _ := json.Marshal(forecastRequestBody())
	req := httptest.NewRequest(http.MethodPost, "/expense-forecasts", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	newForecastTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockDb.AssertExpectations(t)
}

func TestUpdateExpenseForecast_200(t *testing.T) {
	mockDb := new(MockDbService)
	mockDb.On("GetPartnerByID", 1).Return(partnerFixture(1), nil)
	mockDb.On("UpsertExpenseForecast", mock.MatchedBy(func(f *model.ExpenseForecast) bool {
		return f.ID() == 7 && f.Concept() == "Fertilitzants"
	})).Return(nil)

	payload, _ := json.Marshal(forecastRequestBody())
	req := httptest.NewRequest(http.MethodPut, "/expense-forecasts/7", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	newForecastTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp map[string]interface{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, float64(7), resp["id"])
	mockDb.AssertExpectations(t)
}

func TestUpdateExpenseForecast_400_InvalidID(t *testing.T) {
	mockDb := new(MockDbService)

	payload, _ := json.Marshal(forecastRequestBody())
	req := httptest.NewRequest(http.MethodPut, "/expense-forecasts/abc", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	newForecastTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateExpenseForecast_400_BadJSON(t *testing.T) {
	mockDb := new(MockDbService)

	req := httptest.NewRequest(http.MethodPut, "/expense-forecasts/1", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	newForecastTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteExpenseForecast_204(t *testing.T) {
	mockDb := new(MockDbService)
	mockDb.On("DeleteExpenseForecast", 5).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/expense-forecasts/5", nil)
	rec := httptest.NewRecorder()
	newForecastTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockDb.AssertExpectations(t)
}

func TestDeleteExpenseForecast_404(t *testing.T) {
	mockDb := new(MockDbService)
	mockDb.On("DeleteExpenseForecast", 99).Return(fmt.Errorf("expense forecast with ID 99 not found"))

	req := httptest.NewRequest(http.MethodDelete, "/expense-forecasts/99", nil)
	rec := httptest.NewRecorder()
	newForecastTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	mockDb.AssertExpectations(t)
}

func TestDeleteExpenseForecast_400_InvalidID(t *testing.T) {
	mockDb := new(MockDbService)

	req := httptest.NewRequest(http.MethodDelete, "/expense-forecasts/xyz", nil)
	rec := httptest.NewRecorder()
	newForecastTestMux(mockDb).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
