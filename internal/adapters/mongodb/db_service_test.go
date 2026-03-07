package mongodb

import (
	"testing"
	"time"

	"github.com/pjover/espigol/internal/domain/model"
	"github.com/pjover/espigol/internal/domain/ports"
)

type MockConfigService struct {
	dbServer string
	dbName   string
}

func (m *MockConfigService) GetString(key string) string {
	switch key {
	case "db.server":
		return m.dbServer
	case "db.name":
		return m.dbName
	}
	return ""
}

func (m *MockConfigService) SetString(key string, value string) error {
	return nil
}

func (m *MockConfigService) GetTime(key string) time.Time {
	return time.Time{}
}

func (m *MockConfigService) SetTime(key string, value time.Time) error {
	return nil
}

func (m *MockConfigService) Init() {
}

func TestNewDbService(t *testing.T) {
	mockConfig := &MockConfigService{
		dbServer: "mongodb://localhost:27017",
		dbName:   "testdb",
	}

	service := NewDbService(mockConfig)

	if service == nil {
		t.Error("NewDbService should return a non-nil service")
	}

	var _ ports.DbService = service
}

func TestDbServiceConfiguration(t *testing.T) {
	mockConfig := &MockConfigService{
		dbServer: "mongodb://mongo:27017",
		dbName:   "espigol",
	}

	dbService := NewDbService(mockConfig).(*dbService)

	if dbService.uri != "mongodb://mongo:27017" {
		t.Errorf("URI mismatch: expected mongodb://mongo:27017, got %s", dbService.uri)
	}
	if dbService.database != "espigol" {
		t.Errorf("Database mismatch: expected espigol, got %s", dbService.database)
	}
}

func TestDbServiceImplementsInterface(t *testing.T) {
	mockConfig := &MockConfigService{
		dbServer: "mongodb://localhost:27017",
		dbName:   "testdb",
	}

	service := NewDbService(mockConfig)

	var _ ports.DbService = service

	if service == nil {
		t.Error("Service should not be nil")
	}
}

func TestDbServiceMethods(t *testing.T) {
	mockConfig := &MockConfigService{
		dbServer: "mongodb://invalid-host-that-does-not-exist:27017/?serverSelectionTimeoutMS=10&connectTimeoutMS=10",
		dbName:   "testdb",
	}

	service := NewDbService(mockConfig)

	partner := model.NewPartner(
		1,
		"John",
		"Doe",
		"VAT123",
		"john@example.com",
		"+34123456789",
		model.Producer,
		42,
		true,
		false,
		time.Now(),
	)

	// This should fail due to connection error (invalid host)
	err := service.UpsertPartner(partner)
	if err == nil {
		t.Error("UpsertPartner should return error when MongoDB is not available")
	}

	// Test FindPartnerByEmail - should return error
	_, err = service.FindPartnerByEmail("test@example.com")
	if err == nil {
		t.Error("FindPartnerByEmail should return error when MongoDB is not available")
	}

	// Test UpsertExpenseForecast - should return error
	forecast := model.NewExpenseForecast(
		1,
		*partner,
		"Test",
		"Test expense",
		100.0,
		time.Now(),
		model.ExpenseSubtypeA1,
		model.ExpenseScopeCommon,
		nil,
		time.Now(),
	)

	err = service.UpsertExpenseForecast(forecast)
	if err == nil {
		t.Error("UpsertExpenseForecast should return error when MongoDB is not available")
	}
}
