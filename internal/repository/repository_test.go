package repository

import (
	"context"
	"testing"

	"github.com/notepass/vmprov-web/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockUserRepository_CRUD(t *testing.T) {
	repo := NewMockUserRepository()

	id, err := repo.Create(context.Background(), domain.User{
		Username: "testuser",
		Email:    "test@example.com",
	})
	require.NoError(t, err)
	assert.Equal(t, 1, id)

	user, err := repo.GetByID(context.Background(), id)
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)

	userByName, err := repo.GetByUsername(context.Background(), "testuser")
	require.NoError(t, err)
	require.NotNil(t, userByName)
	assert.Equal(t, "test@example.com", userByName.Email)

	users, err := repo.List(context.Background())
	require.NoError(t, err)
	assert.Len(t, users, 1)
}

func TestMockUserRepository_NotFound(t *testing.T) {
	repo := NewMockUserRepository()

	user, err := repo.GetByID(context.Background(), 999)
	require.NoError(t, err)
	assert.Nil(t, user)
}

func TestMockTemplateRepository_CRUD(t *testing.T) {
	repo := NewMockTemplateRepository()

	id, err := repo.Create(context.Background(), domain.Template{
		Name:    "ubuntu-base",
		Content: "#cloud-init",
	})
	require.NoError(t, err)
	assert.Equal(t, 1, id)

	tmpl, err := repo.GetByName(context.Background(), "ubuntu-base")
	require.NoError(t, err)
	require.NotNil(t, tmpl)
	assert.Equal(t, "#cloud-init", tmpl.Content)

	err = repo.Update(context.Background(), domain.Template{
		ID:      id,
		Name:    "ubuntu-base",
		Content: "updated content",
	})
	require.NoError(t, err)

	tmpl, err = repo.GetByID(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, "updated content", tmpl.Content)

	err = repo.Delete(context.Background(), id)
	require.NoError(t, err)

	tmpl, err = repo.GetByID(context.Background(), id)
	require.NoError(t, err)
	assert.Nil(t, tmpl)
}

func TestMockLibvirtConnectionRepository_CRUD(t *testing.T) {
	repo := NewMockLibvirtConnectionRepository()

	host := "192.168.1.10"
	username := "libvirt"
	keyPath := "/etc/vmprov/keys/host10"
	id, err := repo.Create(context.Background(), domain.LibvirtConnection{
		Name:                 "host-10",
		Type:                 "ssh",
		Host:                 &host,
		Username:             &username,
		SSHKeyPath:           &keyPath,
		AcceptUnknownHostKey: true,
	})
	require.NoError(t, err)
	assert.Equal(t, 1, id)

	conn, err := repo.GetByID(context.Background(), id)
	require.NoError(t, err)
	require.NotNil(t, conn)
	assert.Equal(t, "host-10", conn.Name)
	assert.Equal(t, "ssh", conn.Type)
	require.NotNil(t, conn.Host)
	assert.Equal(t, "192.168.1.10", *conn.Host)
	assert.True(t, conn.AcceptUnknownHostKey)
	assert.Nil(t, conn.SocketPath)

	socketPath := "/var/run/libvirt/libvirt-sock"
	id, err = repo.Create(context.Background(), domain.LibvirtConnection{
		Name:       "local",
		Type:       "socket",
		SocketPath: &socketPath,
	})
	require.NoError(t, err)
	assert.Equal(t, 2, id)

	conn, err = repo.GetByID(context.Background(), id)
	require.NoError(t, err)
	require.NotNil(t, conn)
	assert.Equal(t, "socket", conn.Type)
	require.NotNil(t, conn.SocketPath)
	assert.Equal(t, "/var/run/libvirt/libvirt-sock", *conn.SocketPath)
	assert.Nil(t, conn.Host)

	conns, err := repo.List(context.Background())
	require.NoError(t, err)
	assert.Len(t, conns, 2)

	err = repo.Delete(context.Background(), id)
	require.NoError(t, err)

	conn, err = repo.GetByID(context.Background(), id)
	require.NoError(t, err)
	assert.Nil(t, conn)

	conns, err = repo.List(context.Background())
	require.NoError(t, err)
	assert.Len(t, conns, 1)
}

func TestMockLibvirtConnectionRepository_Update(t *testing.T) {
	repo := NewMockLibvirtConnectionRepository()

	host := "192.168.1.10"
	username := "libvirt"
	keyPath := "/etc/vmprov/keys/host10"
	id, err := repo.Create(context.Background(), domain.LibvirtConnection{
		Name:       "host-10",
		Type:       "ssh",
		Host:       &host,
		Username:   &username,
		SSHKeyPath: &keyPath,
	})
	require.NoError(t, err)

	newHost := "192.168.1.11"
	newName := "host-11"
	err = repo.Update(context.Background(), domain.LibvirtConnection{
		ID:         id,
		Name:       newName,
		Type:       "ssh",
		Host:       &newHost,
		Username:   &username,
		SSHKeyPath: &keyPath,
	})
	require.NoError(t, err)

	conn, err := repo.GetByID(context.Background(), id)
	require.NoError(t, err)
	require.NotNil(t, conn)
	assert.Equal(t, "host-11", conn.Name)
	require.NotNil(t, conn.Host)
	assert.Equal(t, "192.168.1.11", *conn.Host)

	byName, err := repo.GetByName(context.Background(), "host-11")
	require.NoError(t, err)
	require.NotNil(t, byName)
	assert.Equal(t, id, byName.ID)

	byOldName, err := repo.GetByName(context.Background(), "host-10")
	require.NoError(t, err)
	assert.Nil(t, byOldName)
}

func TestMockLibvirtConnectionRepository_DuplicateName(t *testing.T) {
	repo := NewMockLibvirtConnectionRepository()

	host := "192.168.1.10"
	_, err := repo.Create(context.Background(), domain.LibvirtConnection{
		Name: "host-10",
		Type: "ssh",
		Host: &host,
	})
	require.NoError(t, err)

	_, err = repo.Create(context.Background(), domain.LibvirtConnection{
		Name: "host-10",
		Type: "ssh",
		Host: &host,
	})
	require.Error(t, err)
}

func TestMockLibvirtConnectionRepository_NotFound(t *testing.T) {
	repo := NewMockLibvirtConnectionRepository()

	conn, err := repo.GetByID(context.Background(), 999)
	require.NoError(t, err)
	assert.Nil(t, conn)

	conn, err = repo.GetByName(context.Background(), "missing")
	require.NoError(t, err)
	assert.Nil(t, conn)
}

func TestMockAuditLogRepository_CreateAndList(t *testing.T) {
	repo := NewMockAuditLogRepository()

	userID := 1
	id, err := repo.Create(context.Background(), domain.AuditLog{
		UserID: &userID,
		Action: "login",
	})
	require.NoError(t, err)
	assert.Equal(t, 1, id)

	logs, err := repo.List(context.Background(), 10, 0)
	require.NoError(t, err)
	assert.Len(t, logs, 1)

	userLogs, err := repo.GetByUserID(context.Background(), userID, 10, 0)
	require.NoError(t, err)
	assert.Len(t, userLogs, 1)
}
