package libvirt

import (
	"context"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultIntegrationSocketPath = "/var/run/libvirt/libvirt-sock"
	defaultIntegrationSSHHost    = "127.0.0.1"
	defaultIntegrationTimeout    = 15 * time.Second
)

var integrationVersionRe = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

// requireIntegration skips the test unless LIBVIRT_INTEGRATION is set to a
// non-empty value, so default and -short runs stay green on machines without
// a libvirt daemon.
func requireIntegration(t *testing.T) {
	t.Helper()
	if os.Getenv("LIBVIRT_INTEGRATION") == "" {
		t.Skip("LIBVIRT_INTEGRATION not set; skipping live libvirt integration test")
	}
}

func integrationSocketPath() string {
	if p := os.Getenv("LIBVIRT_SOCKET_PATH"); p != "" {
		return p
	}
	return defaultIntegrationSocketPath
}

func integrationSSHHost() string {
	if h := os.Getenv("LIBVIRT_SSH_HOST"); h != "" {
		return h
	}
	return defaultIntegrationSSHHost
}

func integrationSSHUser(t *testing.T) string {
	t.Helper()
	if u := os.Getenv("LIBVIRT_SSH_USER"); u != "" {
		return u
	}
	if cu, err := user.Current(); err == nil {
		return cu.Username
	}
	t.Skip("cannot determine the current user for LIBVIRT_SSH_USER")
	return ""
}

// integrationSSHKey returns the LIBVIRT_SSH_KEY path, skipping the test when
// it is unset (e.g. local runs without SSH fixtures).
func integrationSSHKey(t *testing.T) string {
	t.Helper()
	key := os.Getenv("LIBVIRT_SSH_KEY")
	if key == "" {
		t.Skip("LIBVIRT_SSH_KEY not set; skipping SSH integration test")
	}
	return key
}

// integrationKnownHosts returns the LIBVIRT_KNOWN_HOSTS path, skipping the
// test when it is unset.
func integrationKnownHosts(t *testing.T) string {
	t.Helper()
	kh := os.Getenv("LIBVIRT_KNOWN_HOSTS")
	if kh == "" {
		t.Skip("LIBVIRT_KNOWN_HOSTS not set; skipping strict SSH integration test")
	}
	return kh
}

func integrationTimeout() time.Duration {
	if v := os.Getenv("LIBVIRT_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return defaultIntegrationTimeout
}

func assertQEMUHypervisor(t *testing.T, result *TestResult) {
	t.Helper()
	assert.True(t, strings.EqualFold(result.HypervisorType, "QEMU"),
		"expected QEMU hypervisor type, got %q", result.HypervisorType)
}

func TestIntegration_Socket(t *testing.T) {
	requireIntegration(t)

	c := New()
	result, err := c.TestConnection(context.Background(), Connection{
		Type:       TypeSocket,
		SocketPath: integrationSocketPath(),
		Timeout:    integrationTimeout(),
	})
	require.NoError(t, err)
	assert.Regexp(t, integrationVersionRe, result.LibvirtVersion)
	assertQEMUHypervisor(t, result)
	assert.GreaterOrEqual(t, result.TotalDomains, result.ActiveDomains)
	assert.GreaterOrEqual(t, result.ActiveDomains, 0)
}

func TestIntegration_SSH_StrictKnownHost(t *testing.T) {
	requireIntegration(t)
	keyPath := integrationSSHKey(t)
	knownHosts := integrationKnownHosts(t)

	c := New()
	result, err := c.TestConnection(context.Background(), sshIntegrationConnection(t, keyPath, knownHosts, false))
	require.NoError(t, err)
	assert.Regexp(t, integrationVersionRe, result.LibvirtVersion)
	assertQEMUHypervisor(t, result)
}

func TestIntegration_SSH_UnknownHostKeyRejected(t *testing.T) {
	requireIntegration(t)
	keyPath := integrationSSHKey(t)

	knownHosts := filepath.Join(t.TempDir(), "known_hosts")
	require.NoError(t, os.WriteFile(knownHosts, nil, 0600))

	c := New()
	_, err := c.TestConnection(context.Background(), sshIntegrationConnection(t, keyPath, knownHosts, false))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "key is unknown",
		"expected a host key verification error, got: %v", err)
}

func TestIntegration_SSH_AcceptUnknownHostKey(t *testing.T) {
	requireIntegration(t)
	keyPath := integrationSSHKey(t)

	knownHosts := filepath.Join(t.TempDir(), "known_hosts")
	require.NoError(t, os.WriteFile(knownHosts, nil, 0600))

	c := New()
	result, err := c.TestConnection(context.Background(), sshIntegrationConnection(t, keyPath, knownHosts, true))
	require.NoError(t, err)
	assert.Regexp(t, integrationVersionRe, result.LibvirtVersion)

	data, err := os.ReadFile(knownHosts)
	require.NoError(t, err)
	assert.Contains(t, string(data), integrationSSHHost(),
		"expected the host key for %q to be appended to known_hosts", integrationSSHHost())
}

func sshIntegrationConnection(t *testing.T, keyPath, knownHosts string, acceptUnknown bool) Connection {
	t.Helper()
	return Connection{
		Type:                 TypeSSH,
		Host:                 integrationSSHHost(),
		Username:             integrationSSHUser(t),
		SSHKeyPath:           keyPath,
		KnownHostsFile:       knownHosts,
		AcceptUnknownHostKey: acceptUnknown,
		Timeout:              integrationTimeout(),
	}
}
