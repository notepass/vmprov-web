package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Defaults(t *testing.T) {
	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, DefaultPort, cfg.ServerPort)
	assert.Equal(t, DefaultLogLevel, cfg.LogLevel)
	assert.Equal(t, DefaultMaxOpenConns, cfg.DBMaxOpenConns)
	assert.Equal(t, DefaultMaxIdleConns, cfg.DBMaxIdleConns)
	assert.Equal(t, DefaultConnMaxLifetime, cfg.DBConnMaxLifetime)
	assert.Equal(t, DefaultLibvirtConnectTimeout, cfg.LibvirtConnectTimeout)

	home, err := os.UserHomeDir()
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(home, ".ssh", "known_hosts"), cfg.LibvirtKnownHostsFile)

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
	t.Setenv("DB_MAX_OPEN_CONNS", "25")
	t.Setenv("DB_MAX_IDLE_CONNS", "10")
	t.Setenv("DB_CONN_MAX_LIFETIME", "300")
	t.Setenv("LIBVIRT_CONNECT_TIMEOUT", "30")
	t.Setenv("LIBVIRT_KNOWN_HOSTS_FILE", "/etc/vmprov/known_hosts")

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, 3000, cfg.ServerPort)
	assert.Equal(t, "DEBUG", cfg.LogLevel)
	assert.Equal(t, 25, cfg.DBMaxOpenConns)
	assert.Equal(t, 10, cfg.DBMaxIdleConns)
	assert.Equal(t, 300, cfg.DBConnMaxLifetime)
	assert.Equal(t, 30, cfg.LibvirtConnectTimeout)
	assert.Equal(t, "/etc/vmprov/known_hosts", cfg.LibvirtKnownHostsFile)
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
