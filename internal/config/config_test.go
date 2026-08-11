package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Defaults(t *testing.T) {
	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, DefaultPort, cfg.ServerPort)
	assert.Equal(t, DefaultLogLevel, cfg.LogLevel)
	assert.Empty(t, cfg.DBConnString)
	assert.Empty(t, cfg.DBUsername)
	assert.Empty(t, cfg.DBPassword)
}

func TestLoad_EnvOverride(t *testing.T) {
	t.Setenv("SERVER_PORT", "3000")
	t.Setenv("LOG_LEVEL", "DEBUG")
	t.Setenv("DB_CONN_STRING", "postgres://localhost/testdb")
	t.Setenv("DB_USERNAME", "testuser")
	t.Setenv("DB_PASSWORD", "testpass")

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, 3000, cfg.ServerPort)
	assert.Equal(t, "DEBUG", cfg.LogLevel)
	assert.Equal(t, "postgres://localhost/testdb", cfg.DBConnString)
	assert.Equal(t, "testuser", cfg.DBUsername)
	assert.Equal(t, "testpass", cfg.DBPassword)
}

func TestLoad_ConfigFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "config-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	configPath := tmpDir + "/config.yaml"
	err = os.WriteFile(configPath, []byte("server_port: 9090\nlog_level: WARN\n"), 0644)
	require.NoError(t, err)

	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, 9090, cfg.ServerPort)
	assert.Equal(t, "WARN", cfg.LogLevel)
}
