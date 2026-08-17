package api

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"encoding/pem"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"

	"github.com/notepass/vmprov-web/internal/domain"
	"github.com/notepass/vmprov-web/internal/libvirt"
	"github.com/notepass/vmprov-web/internal/repository"
)

type fakeConnRepo struct {
	mu    sync.Mutex
	conns map[int]domain.LibvirtConnection
	next  int
}

var _ repository.LibvirtConnectionRepository = (*fakeConnRepo)(nil)

func newFakeConnRepo() *fakeConnRepo {
	return &fakeConnRepo{conns: map[int]domain.LibvirtConnection{}, next: 1}
}

func (r *fakeConnRepo) Create(_ context.Context, conn domain.LibvirtConnection) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	conn.ID = r.next
	r.next++
	conn.CreatedAt = now
	conn.UpdatedAt = now
	r.conns[conn.ID] = conn
	return conn.ID, nil
}

func (r *fakeConnRepo) GetByID(_ context.Context, id int) (*domain.LibvirtConnection, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	conn, ok := r.conns[id]
	if !ok {
		return nil, nil
	}
	c := conn
	return &c, nil
}

func (r *fakeConnRepo) GetByName(_ context.Context, name string) (*domain.LibvirtConnection, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, conn := range r.conns {
		if conn.Name == name {
			c := conn
			return &c, nil
		}
	}
	return nil, nil
}

func (r *fakeConnRepo) Update(_ context.Context, conn domain.LibvirtConnection) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.conns[conn.ID]; !ok {
		return nil
	}
	existing := r.conns[conn.ID]
	conn.CreatedAt = existing.CreatedAt
	conn.UpdatedAt = time.Now()
	r.conns[conn.ID] = conn
	return nil
}

func (r *fakeConnRepo) Delete(_ context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.conns, id)
	return nil
}

func (r *fakeConnRepo) List(_ context.Context) ([]domain.LibvirtConnection, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	conns := make([]domain.LibvirtConnection, 0, len(r.conns))
	for id := range r.conns {
		conns = append(conns, r.conns[id])
	}
	return conns, nil
}

func (r *fakeConnRepo) get(id int) *domain.LibvirtConnection {
	r.mu.Lock()
	defer r.mu.Unlock()
	conn, ok := r.conns[id]
	if !ok {
		return nil
	}
	c := conn
	return &c
}

type testSetup struct {
	e    *echo.Echo
	repo *fakeConnRepo
	fake *libvirt.FakeClient
}

func newTestSetup(t *testing.T) *testSetup {
	t.Helper()
	repo := newFakeConnRepo()
	fake := &libvirt.FakeClient{}
	h := NewLibvirtConnectionHandler(repo, fake, 5*time.Second, "/dev/null", slog.New(slog.NewTextHandler(io.Discard, nil)))
	e := echo.New()
	h.RegisterRoutes(e)
	return &testSetup{e: e, repo: repo, fake: fake}
}

func doRequest(t *testing.T, e *echo.Echo, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var reader *bytes.Reader
	if body != nil {
		b, err := json.Marshal(body)
		require.NoError(t, err)
		reader = bytes.NewReader(b)
	} else {
		reader = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

func sshRequest(t *testing.T, name string) LibvirtConnectionRequest {
	t.Helper()
	dir := t.TempDir()
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	block, err := ssh.MarshalPrivateKey(priv, "test")
	require.NoError(t, err)
	keyPath := filepath.Join(dir, "id_ed25519")
	require.NoError(t, os.WriteFile(keyPath, pem.EncodeToMemory(block), 0o600))

	return LibvirtConnectionRequest{
		Name:       name,
		Type:       libvirt.TypeSSH,
		Host:       ptr("host.example.com"),
		Username:   ptr("admin"),
		SSHKeyPath: ptr(keyPath),
	}
}

func ptr[T any](v T) *T { return &v }

func TestList_Empty(t *testing.T) {
	s := newTestSetup(t)

	rec := doRequest(t, s.e, http.MethodGet, "/api/v1/remotes/libvirt/connections", nil)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp []LibvirtConnectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Empty(t, resp)
}

func TestList_ReturnsConnections(t *testing.T) {
	s := newTestSetup(t)
	now := time.Now()
	s.repo.conns[1] = domain.LibvirtConnection{
		ID: 1, Name: "prod", Type: libvirt.TypeSSH,
		Host: ptr("host"), Username: ptr("u"), SSHKeyPath: ptr("/k"),
		CreatedAt: now, UpdatedAt: now,
	}

	rec := doRequest(t, s.e, http.MethodGet, "/api/v1/remotes/libvirt/connections", nil)

	require.Equal(t, http.StatusOK, rec.Code)
	var resp []LibvirtConnectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp, 1)
	assert.Equal(t, 1, resp[0].ID)
	assert.Equal(t, "prod", resp[0].Name)
	assert.Equal(t, "ssh", resp[0].Type)
}

func TestCreate_SSH(t *testing.T) {
	s := newTestSetup(t)
	desc := "production host"
	req := sshRequest(t, "prod")
	req.Description = &desc

	rec := doRequest(t, s.e, http.MethodPost, "/api/v1/remotes/libvirt/connections", req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var resp LibvirtConnectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, 1, resp.ID)
	assert.Equal(t, "prod", resp.Name)
	assert.Equal(t, "ssh", resp.Type)
	assert.Equal(t, "host.example.com", *resp.Host)
	assert.Equal(t, "admin", *resp.Username)
	assert.Equal(t, *req.SSHKeyPath, *resp.SSHKeyPath)
	assert.False(t, resp.AcceptUnknownHostKey)

	stored := s.repo.get(1)
	require.NotNil(t, stored)
	assert.Equal(t, "prod", stored.Name)
}

func TestCreate_Socket(t *testing.T) {
	s := newTestSetup(t)
	req := LibvirtConnectionRequest{
		Name:       "local",
		Type:       libvirt.TypeSocket,
		SocketPath: ptr("/var/run/libvirt/libvirt-sock"),
	}

	rec := doRequest(t, s.e, http.MethodPost, "/api/v1/remotes/libvirt/connections", req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var resp LibvirtConnectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "socket", resp.Type)
	assert.Equal(t, "/var/run/libvirt/libvirt-sock", *resp.SocketPath)
}

func TestCreate_AcceptUnknownHostKey(t *testing.T) {
	s := newTestSetup(t)
	req := sshRequest(t, "prod")
	req.AcceptUnknownHostKey = true

	rec := doRequest(t, s.e, http.MethodPost, "/api/v1/remotes/libvirt/connections", req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var resp LibvirtConnectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.True(t, resp.AcceptUnknownHostKey)
}

func TestCreate_DuplicateName(t *testing.T) {
	s := newTestSetup(t)
	doRequest(t, s.e, http.MethodPost, "/api/v1/remotes/libvirt/connections", sshRequest(t, "prod"))

	rec := doRequest(t, s.e, http.MethodPost, "/api/v1/remotes/libvirt/connections", sshRequest(t, "prod"))

	require.Equal(t, http.StatusConflict, rec.Code)
	var errResp ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &errResp))
	assert.Contains(t, errResp.Error, "already exists")
}

func TestCreate_MissingName(t *testing.T) {
	s := newTestSetup(t)
	req := LibvirtConnectionRequest{Type: libvirt.TypeSocket, SocketPath: ptr("/sock")}

	rec := doRequest(t, s.e, http.MethodPost, "/api/v1/remotes/libvirt/connections", req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	var errResp ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &errResp))
	assert.Contains(t, errResp.Error, "name is required")
}

func TestCreate_InvalidType(t *testing.T) {
	s := newTestSetup(t)
	req := LibvirtConnectionRequest{Name: "bad", Type: "tcp"}

	rec := doRequest(t, s.e, http.MethodPost, "/api/v1/remotes/libvirt/connections", req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	var errResp ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &errResp))
	assert.Contains(t, errResp.Error, "type must be")
}

func TestCreate_SSH_MissingHost(t *testing.T) {
	s := newTestSetup(t)
	req := sshRequest(t, "prod")
	req.Host = nil

	rec := doRequest(t, s.e, http.MethodPost, "/api/v1/remotes/libvirt/connections", req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	var errResp ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &errResp))
	assert.Contains(t, errResp.Error, "host is required")
}

func TestCreate_SSH_MissingUsername(t *testing.T) {
	s := newTestSetup(t)
	req := sshRequest(t, "prod")
	req.Username = nil

	rec := doRequest(t, s.e, http.MethodPost, "/api/v1/remotes/libvirt/connections", req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	var errResp ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &errResp))
	assert.Contains(t, errResp.Error, "username is required")
}

func TestCreate_SSH_MissingKeyPath(t *testing.T) {
	s := newTestSetup(t)
	req := sshRequest(t, "prod")
	req.SSHKeyPath = nil

	rec := doRequest(t, s.e, http.MethodPost, "/api/v1/remotes/libvirt/connections", req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	var errResp ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &errResp))
	assert.Contains(t, errResp.Error, "ssh_key_path is required")
}

func TestCreate_SSH_InvalidKey(t *testing.T) {
	s := newTestSetup(t)
	req := sshRequest(t, "prod")
	garbage := filepath.Join(t.TempDir(), "garbage")
	require.NoError(t, os.WriteFile(garbage, []byte("not a key"), 0o600))
	req.SSHKeyPath = &garbage

	rec := doRequest(t, s.e, http.MethodPost, "/api/v1/remotes/libvirt/connections", req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	var errResp ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &errResp))
	assert.Contains(t, errResp.Error, "not a valid private key")
}

func TestCreate_SSH_MissingKeyFile(t *testing.T) {
	s := newTestSetup(t)
	req := sshRequest(t, "prod")
	missing := filepath.Join(t.TempDir(), "does-not-exist")
	req.SSHKeyPath = &missing

	rec := doRequest(t, s.e, http.MethodPost, "/api/v1/remotes/libvirt/connections", req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	var errResp ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &errResp))
	assert.Contains(t, errResp.Error, "cannot read ssh key file")
}

func TestCreate_Socket_MissingSocketPath(t *testing.T) {
	s := newTestSetup(t)
	req := LibvirtConnectionRequest{Name: "local", Type: libvirt.TypeSocket}

	rec := doRequest(t, s.e, http.MethodPost, "/api/v1/remotes/libvirt/connections", req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	var errResp ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &errResp))
	assert.Contains(t, errResp.Error, "socket_path is required")
}

func TestGet_Found(t *testing.T) {
	s := newTestSetup(t)
	now := time.Now()
	status := "ok"
	checked := now
	s.repo.conns[7] = domain.LibvirtConnection{
		ID: 7, Name: "prod", Type: libvirt.TypeSSH,
		Host: ptr("h"), Username: ptr("u"), SSHKeyPath: ptr("/k"),
		LastStatus: &status, LastCheckedAt: &checked,
		CreatedAt: now, UpdatedAt: now,
	}

	rec := doRequest(t, s.e, http.MethodGet, "/api/v1/remotes/libvirt/connections/7", nil)

	require.Equal(t, http.StatusOK, rec.Code)
	var resp LibvirtConnectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, 7, resp.ID)
	assert.Equal(t, "ok", *resp.LastStatus)
	require.NotNil(t, resp.LastCheckedAt)
	assert.Equal(t, checked.UTC(), resp.LastCheckedAt.UTC())
}

func TestGet_NotFound(t *testing.T) {
	s := newTestSetup(t)

	rec := doRequest(t, s.e, http.MethodGet, "/api/v1/remotes/libvirt/connections/999", nil)

	require.Equal(t, http.StatusNotFound, rec.Code)
	var errResp ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &errResp))
	assert.Contains(t, errResp.Error, "not found")
}

func TestGet_InvalidID(t *testing.T) {
	s := newTestSetup(t)

	rec := doRequest(t, s.e, http.MethodGet, "/api/v1/remotes/libvirt/connections/abc", nil)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	var errResp ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &errResp))
	assert.Contains(t, errResp.Error, "invalid connection id")
}

func TestUpdate_Found(t *testing.T) {
	s := newTestSetup(t)
	doRequest(t, s.e, http.MethodPost, "/api/v1/remotes/libvirt/connections", sshRequest(t, "old-name"))

	req := sshRequest(t, "new-name")
	desc := "updated"
	req.Description = &desc

	rec := doRequest(t, s.e, http.MethodPut, "/api/v1/remotes/libvirt/connections/1", req)

	require.Equal(t, http.StatusOK, rec.Code)
	var resp LibvirtConnectionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, 1, resp.ID)
	assert.Equal(t, "new-name", resp.Name)
	assert.Equal(t, "updated", *resp.Description)
}

func TestUpdate_PreservesCheckStatus(t *testing.T) {
	s := newTestSetup(t)
	now := time.Now()
	status := "error"
	checked := now
	s.repo.conns[1] = domain.LibvirtConnection{
		ID: 1, Name: "prod", Type: libvirt.TypeSSH,
		Host: ptr("h"), Username: ptr("u"), SSHKeyPath: ptr("/k"),
		LastStatus: &status, LastCheckedAt: &checked,
		CreatedAt: now, UpdatedAt: now,
	}

	req := sshRequest(t, "prod")
	rec := doRequest(t, s.e, http.MethodPut, "/api/v1/remotes/libvirt/connections/1", req)
	require.Equal(t, http.StatusOK, rec.Code)

	stored := s.repo.get(1)
	require.NotNil(t, stored)
	require.NotNil(t, stored.LastStatus)
	assert.Equal(t, "error", *stored.LastStatus)
	require.NotNil(t, stored.LastCheckedAt)
	assert.Equal(t, checked, *stored.LastCheckedAt)
}

func TestUpdate_NotFound(t *testing.T) {
	s := newTestSetup(t)
	req := sshRequest(t, "prod")

	rec := doRequest(t, s.e, http.MethodPut, "/api/v1/remotes/libvirt/connections/999", req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdate_DuplicateName(t *testing.T) {
	s := newTestSetup(t)
	doRequest(t, s.e, http.MethodPost, "/api/v1/remotes/libvirt/connections", sshRequest(t, "a"))
	doRequest(t, s.e, http.MethodPost, "/api/v1/remotes/libvirt/connections", sshRequest(t, "b"))

	rec := doRequest(t, s.e, http.MethodPut, "/api/v1/remotes/libvirt/connections/1", sshRequest(t, "b"))

	require.Equal(t, http.StatusConflict, rec.Code)
	var errResp ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &errResp))
	assert.Contains(t, errResp.Error, "already exists")
}

func TestUpdate_InvalidBody(t *testing.T) {
	s := newTestSetup(t)
	doRequest(t, s.e, http.MethodPost, "/api/v1/remotes/libvirt/connections", sshRequest(t, "prod"))

	rec := doRequest(t, s.e, http.MethodPut, "/api/v1/remotes/libvirt/connections/1", LibvirtConnectionRequest{Name: "", Type: "bad"})

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDelete_Found(t *testing.T) {
	s := newTestSetup(t)
	doRequest(t, s.e, http.MethodPost, "/api/v1/remotes/libvirt/connections", sshRequest(t, "prod"))

	rec := doRequest(t, s.e, http.MethodDelete, "/api/v1/remotes/libvirt/connections/1", nil)

	require.Equal(t, http.StatusNoContent, rec.Code)
	assert.Nil(t, s.repo.get(1))
}

func TestDelete_NotFound(t *testing.T) {
	s := newTestSetup(t)

	rec := doRequest(t, s.e, http.MethodDelete, "/api/v1/remotes/libvirt/connections/999", nil)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestTest_Success(t *testing.T) {
	s := newTestSetup(t)
	doRequest(t, s.e, http.MethodPost, "/api/v1/remotes/libvirt/connections", sshRequest(t, "prod"))
	s.fake.Result = &libvirt.TestResult{
		LibvirtVersion: "10.6.0",
		HypervisorType: "kvm",
		TotalDomains:   4,
		ActiveDomains:  2,
	}

	rec := doRequest(t, s.e, http.MethodPost, "/api/v1/remotes/libvirt/connections/1/test", nil)

	require.Equal(t, http.StatusOK, rec.Code)
	var resp LibvirtConnectionTestResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "ok", resp.Status)
	assert.Equal(t, "10.6.0", resp.LibvirtVersion)
	assert.Equal(t, "kvm", resp.HypervisorType)
	assert.Equal(t, 4, resp.TotalDomains)
	assert.Equal(t, 2, resp.ActiveDomains)

	require.Len(t, s.fake.Calls, 1)
	assert.Equal(t, "host.example.com", s.fake.Calls[0].Host)
	assert.Equal(t, "/dev/null", s.fake.Calls[0].KnownHostsFile)

	stored := s.repo.get(1)
	require.NotNil(t, stored)
	require.NotNil(t, stored.LastStatus)
	assert.Equal(t, "ok", *stored.LastStatus)
	require.NotNil(t, stored.LastCheckedAt)
}

func TestTest_Failure(t *testing.T) {
	s := newTestSetup(t)
	doRequest(t, s.e, http.MethodPost, "/api/v1/remotes/libvirt/connections", sshRequest(t, "prod"))
	s.fake.Err = assert.AnError

	rec := doRequest(t, s.e, http.MethodPost, "/api/v1/remotes/libvirt/connections/1/test", nil)

	require.Equal(t, http.StatusBadGateway, rec.Code)
	var resp LibvirtConnectionTestResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "error", resp.Status)
	assert.NotEmpty(t, resp.Message)

	stored := s.repo.get(1)
	require.NotNil(t, stored)
	require.NotNil(t, stored.LastStatus)
	assert.Equal(t, "error", *stored.LastStatus)
}

func TestTest_NotFound(t *testing.T) {
	s := newTestSetup(t)

	rec := doRequest(t, s.e, http.MethodPost, "/api/v1/remotes/libvirt/connections/999/test", nil)

	require.Equal(t, http.StatusNotFound, rec.Code)
}
