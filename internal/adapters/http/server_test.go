package http_test

import (
	"context"
	"testing"
	"time"

	httpAdapter "github.com/pjover/espigol/internal/adapters/http"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	mockConfig := new(MockConfigService)
	mockDb := new(MockDbService)
	mockConfig.On("GetString", "server.port").Return("0")

	srv := httpAdapter.NewHttpServer(mockConfig, mockDb)
	assert.NotNil(t, srv)

	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.Start()
	}()

	time.Sleep(50 * time.Millisecond)

	err := srv.Stop(context.Background())
	assert.NoError(t, err)

	startErr := <-errChan
	assert.NoError(t, startErr)
}
