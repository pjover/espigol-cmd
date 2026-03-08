package http_test

import (
	"time"

	"github.com/pjover/espigol/internal/domain/model"
	"github.com/pjover/espigol/internal/domain/ports"
	"github.com/stretchr/testify/mock"
)

// MockConfigService mocks ports.ConfigService
type MockConfigService struct {
	mock.Mock
}

func (m *MockConfigService) Init() {
	m.Called()
}

func (m *MockConfigService) GetString(key string) string {
	args := m.Called(key)
	return args.String(0)
}

func (m *MockConfigService) SetString(key string, value string) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *MockConfigService) GetFloat64(key string) float64 {
	args := m.Called(key)
	return args.Get(0).(float64)
}

func (m *MockConfigService) GetTime(key string) time.Time {
	args := m.Called(key)
	return args.Get(0).(time.Time)
}

func (m *MockConfigService) SetTime(key string, value time.Time) error {
	args := m.Called(key, value)
	return args.Error(0)
}

// Ensure MockConfigService implements ports.ConfigService at compile time.
var _ ports.ConfigService = (*MockConfigService)(nil)

// MockDbService mocks ports.DbService
type MockDbService struct {
	mock.Mock
}

func (m *MockDbService) UpsertPartner(partner *model.Partner) error {
	args := m.Called(partner)
	return args.Error(0)
}

func (m *MockDbService) GetPartnerByID(id int) (*model.Partner, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Partner), args.Error(1)
}

func (m *MockDbService) GetAllPartners() ([]*model.Partner, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Partner), args.Error(1)
}

func (m *MockDbService) DeletePartner(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDbService) FindPartnerByEmail(email string) (*model.Partner, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Partner), args.Error(1)
}

func (m *MockDbService) UpsertExpenseForecast(forecast *model.ExpenseForecast) error {
	args := m.Called(forecast)
	return args.Error(0)
}

func (m *MockDbService) GetExpenseForecastByID(id int) (*model.ExpenseForecast, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ExpenseForecast), args.Error(1)
}

func (m *MockDbService) GetAllExpenseForecasts() ([]*model.ExpenseForecast, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.ExpenseForecast), args.Error(1)
}

func (m *MockDbService) DeleteExpenseForecast(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// Ensure MockDbService implements ports.DbService at compile time.
var _ ports.DbService = (*MockDbService)(nil)
