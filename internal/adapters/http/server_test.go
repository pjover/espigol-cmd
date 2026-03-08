package http_test

import (
	"context"
	"testing"
	"time"

	httpAdapter "github.com/pjover/espigol/internal/adapters/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func (m *MockConfigService) GetTime(key string) time.Time {
	args := m.Called(key)
	return args.Get(0).(time.Time)
}

func (m *MockConfigService) SetTime(key string, value time.Time) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func TestNewServer(t *testing.T) {
	mockConfig := new(MockConfigService)
	// We'll use port 0 so the OS assigns a random available port, preventing conflicts during tests
	mockConfig.On("GetString", "server.port").Return("0")

	srv := httpAdapter.NewServer(mockConfig)
	assert.NotNil(t, srv)

	// Test Start and Stop asynchronously
	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.Start()
	}()

	// Give it a moment to start
	time.Sleep(50 * time.Millisecond)

	// Verify stop
	err := srv.Stop(context.Background())
	assert.NoError(t, err)

	startErr := <-errChan
	assert.NoError(t, startErr) // http.ErrServerClosed shouldn't be returned as error from our Start wrapper
}
