package libvirt

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"errors"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

func writeTestKey(t *testing.T, dir string, name string, passphrase []byte) string {
	t.Helper()
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	var block *pem.Block
	if passphrase == nil {
		block, err = ssh.MarshalPrivateKey(priv, "test")
	} else {
		block, err = ssh.MarshalPrivateKeyWithPassphrase(priv, "test", passphrase)
	}
	require.NoError(t, err)
	keyBytes := pem.EncodeToMemory(block)

	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, keyBytes, 0600))
	return path
}

func TestValidateSSHKey_Valid(t *testing.T) {
	dir := t.TempDir()
	path := writeTestKey(t, dir, "id_ed25519", nil)

	assert.NoError(t, ValidateSSHKey(path))
}

func TestValidateSSHKey_Missing(t *testing.T) {
	err := ValidateSSHKey(filepath.Join(t.TempDir(), "does-not-exist"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot read ssh key file")
}

func TestValidateSSHKey_Unreadable(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("running as root, file permission checks are bypassed")
	}

	dir := t.TempDir()
	path := writeTestKey(t, dir, "id_ed25519", nil)
	require.NoError(t, os.Chmod(path, 0000))
	t.Cleanup(func() { os.Chmod(path, 0600) })

	err := ValidateSSHKey(path)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot read ssh key file")
}

func TestValidateSSHKey_PassphraseProtected(t *testing.T) {
	dir := t.TempDir()
	path := writeTestKey(t, dir, "id_ed25519", []byte("secret"))

	err := ValidateSSHKey(path)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrKeyPassphraseProtected), "expected passphrase error, got: %v", err)
}

func TestValidateSSHKey_NotAKey(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "garbage")
	require.NoError(t, os.WriteFile(path, []byte("not a key"), 0600))

	err := ValidateSSHKey(path)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not a valid private key")
}

func TestValidateSSHKey_EmptyPath(t *testing.T) {
	err := ValidateSSHKey("")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "empty")
}

func TestTestConnection_SocketDialFailure(t *testing.T) {
	c := New()
	socketPath := filepath.Join(t.TempDir(), "libvirt-sock")

	_, err := c.TestConnection(context.Background(), Connection{
		Type:       TypeSocket,
		SocketPath: socketPath,
		Timeout:    time.Second,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect to libvirt")
}

func TestTestConnection_InvalidType(t *testing.T) {
	c := New()

	_, err := c.TestConnection(context.Background(), Connection{
		Type:    "tcp",
		Timeout: time.Second,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported connection type")
}

func TestTestConnection_SocketMissingPathField(t *testing.T) {
	c := New()

	_, err := c.TestConnection(context.Background(), Connection{
		Type:    TypeSocket,
		Timeout: time.Second,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "socket_path is required")
}

func TestTestConnection_SSHMissingHost(t *testing.T) {
	c := New()

	_, err := c.TestConnection(context.Background(), Connection{
		Type:    TypeSSH,
		Timeout: time.Second,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "host is required")
}

func TestTestConnection_Timeout(t *testing.T) {
	dir := t.TempDir()
	socketPath := filepath.Join(dir, "libvirt-sock")

	ln, err := net.Listen("unix", socketPath)
	require.NoError(t, err)
	defer ln.Close()

	accepted := make(chan net.Conn, 1)
	go func() {
		conn, err := ln.Accept()
		if err == nil {
			accepted <- conn
		}
	}()

	c := New()
	start := time.Now()
	_, err = c.TestConnection(context.Background(), Connection{
		Type:       TypeSocket,
		SocketPath: socketPath,
		Timeout:    500 * time.Millisecond,
	})
	elapsed := time.Since(start)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "timed out")
	assert.Less(t, elapsed, 4*time.Second, "should return at the configured timeout, not hang")

	select {
	case conn := <-accepted:
		conn.Close()
	default:
	}
}

func TestFormatLibVersion(t *testing.T) {
	assert.Equal(t, "10.12.0", formatLibVersion(10*1000000+12*1000+0))
	assert.Equal(t, "9.0.1", formatLibVersion(9*1000000+0*1000+1))
	assert.Equal(t, "0.0.0", formatLibVersion(0))
}
