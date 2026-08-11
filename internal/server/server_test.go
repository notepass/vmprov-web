package server

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/notepass/vmprov-web/internal/config"
)

func TestNew_ReturnsHTTPServer(t *testing.T) {
	cfg := &config.Config{
		ServerPort: 8080,
		LogLevel:   "INFO",
	}
	logger := slog.Default()

	srv := New(cfg, logger)
	require.NotNil(t, srv)
	assert.NotNil(t, srv.ReadTimeout)
}
