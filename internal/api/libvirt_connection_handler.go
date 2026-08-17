// Package api contains the HTTP handlers for the vmprov-web REST API.
package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/notepass/vmprov-web/internal/domain"
	"github.com/notepass/vmprov-web/internal/libvirt"
	"github.com/notepass/vmprov-web/internal/repository"
)

// LibvirtConnectionHandler handles the libvirt connection endpoints.
type LibvirtConnectionHandler struct {
	repo       repository.LibvirtConnectionRepository
	client     libvirt.Client
	timeout    time.Duration
	knownHosts string
	log        *slog.Logger
}

// NewLibvirtConnectionHandler creates a new LibvirtConnectionHandler.
func NewLibvirtConnectionHandler(repo repository.LibvirtConnectionRepository, client libvirt.Client, timeout time.Duration, knownHostsFile string, log *slog.Logger) *LibvirtConnectionHandler {
	return &LibvirtConnectionHandler{
		repo:       repo,
		client:     client,
		timeout:    timeout,
		knownHosts: knownHostsFile,
		log:        log,
	}
}

// LibvirtConnectionRequest is the JSON body for create and update.
type LibvirtConnectionRequest struct {
	Name                 string  `json:"name"`
	Type                 string  `json:"type"`
	Host                 *string `json:"host"`
	Username             *string `json:"username"`
	SSHKeyPath           *string `json:"ssh_key_path"`
	AcceptUnknownHostKey bool    `json:"accept_unknown_host_key"`
	SocketPath           *string `json:"socket_path"`
	Description          *string `json:"description"`
}

// LibvirtConnectionResponse is the JSON representation of a connection.
type LibvirtConnectionResponse struct {
	ID                   int        `json:"id"`
	Name                 string     `json:"name"`
	Type                 string     `json:"type"`
	Host                 *string    `json:"host,omitempty"`
	Username             *string    `json:"username,omitempty"`
	SSHKeyPath           *string    `json:"ssh_key_path,omitempty"`
	AcceptUnknownHostKey bool       `json:"accept_unknown_host_key"`
	SocketPath           *string    `json:"socket_path,omitempty"`
	Description          *string    `json:"description,omitempty"`
	LastStatus           *string    `json:"last_status,omitempty"`
	LastCheckedAt        *time.Time `json:"last_checked_at,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

// LibvirtConnectionTestResponse is the JSON response of a connection test.
type LibvirtConnectionTestResponse struct {
	Status         string `json:"status"`
	Message        string `json:"message,omitempty"`
	LibvirtVersion string `json:"libvirt_version,omitempty"`
	HypervisorType string `json:"hypervisor_type,omitempty"`
	TotalDomains   int    `json:"total_domains,omitempty"`
	ActiveDomains  int    `json:"active_domains,omitempty"`
}

// ErrorResponse is the JSON body for error responses.
type ErrorResponse struct {
	Error string `json:"error"`
}

// RegisterRoutes registers the libvirt connection routes on the Echo instance.
func (h *LibvirtConnectionHandler) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/api/v1/remotes/libvirt/connections")
	g.GET("", h.List)
	g.POST("", h.Create)
	g.GET("/:id", h.Get)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
	g.POST("/:id/test", h.Test)
}

func (h *LibvirtConnectionHandler) List(c echo.Context) error {
	ctx := c.Request().Context()

	conns, err := h.repo.List(ctx)
	if err != nil {
		return h.internalError(c, err)
	}

	resp := make([]LibvirtConnectionResponse, 0, len(conns))
	for _, conn := range conns {
		resp = append(resp, toResponse(conn))
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *LibvirtConnectionHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()

	var req LibvirtConnectionRequest
	if err := c.Bind(&req); err != nil {
		return h.badRequest(c, fmt.Sprintf("invalid request body: %v", err))
	}

	conn, err := h.validateRequest(&req)
	if err != nil {
		return h.badRequest(c, err.Error())
	}

	if existing, err := h.repo.GetByName(ctx, conn.Name); err != nil {
		return h.internalError(c, err)
	} else if existing != nil {
		return h.conflict(c, fmt.Sprintf("connection name %q already exists", conn.Name))
	}

	id, err := h.repo.Create(ctx, *conn)
	if err != nil {
		return h.internalError(c, err)
	}

	created, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return h.internalError(c, err)
	}
	if created == nil {
		return h.internalError(c, fmt.Errorf("created connection %d not found", id))
	}

	return c.JSON(http.StatusCreated, toResponse(*created))
}

func (h *LibvirtConnectionHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := parseID(c)
	if err != nil {
		return h.badRequest(c, err.Error())
	}

	conn, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return h.internalError(c, err)
	}
	if conn == nil {
		return h.notFound(c, "connection not found")
	}

	return c.JSON(http.StatusOK, toResponse(*conn))
}

func (h *LibvirtConnectionHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := parseID(c)
	if err != nil {
		return h.badRequest(c, err.Error())
	}

	existing, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return h.internalError(c, err)
	}
	if existing == nil {
		return h.notFound(c, "connection not found")
	}

	var req LibvirtConnectionRequest
	if err := c.Bind(&req); err != nil {
		return h.badRequest(c, fmt.Sprintf("invalid request body: %v", err))
	}

	conn, err := h.validateRequest(&req)
	if err != nil {
		return h.badRequest(c, err.Error())
	}

	if req.Name != existing.Name {
		if other, err := h.repo.GetByName(ctx, req.Name); err != nil {
			return h.internalError(c, err)
		} else if other != nil && other.ID != id {
			return h.conflict(c, fmt.Sprintf("connection name %q already exists", req.Name))
		}
	}

	conn.ID = id
	conn.LastStatus = existing.LastStatus
	conn.LastCheckedAt = existing.LastCheckedAt

	if err := h.repo.Update(ctx, *conn); err != nil {
		return h.internalError(c, err)
	}

	updated, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return h.internalError(c, err)
	}
	if updated == nil {
		return h.internalError(c, fmt.Errorf("updated connection %d not found", id))
	}

	return c.JSON(http.StatusOK, toResponse(*updated))
}

func (h *LibvirtConnectionHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := parseID(c)
	if err != nil {
		return h.badRequest(c, err.Error())
	}

	existing, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return h.internalError(c, err)
	}
	if existing == nil {
		return h.notFound(c, "connection not found")
	}

	if err := h.repo.Delete(ctx, id); err != nil {
		return h.internalError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *LibvirtConnectionHandler) Test(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := parseID(c)
	if err != nil {
		return h.badRequest(c, err.Error())
	}

	conn, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return h.internalError(c, err)
	}
	if conn == nil {
		return h.notFound(c, "connection not found")
	}

	lvConn := h.toLibvirtConn(conn)

	result, testErr := h.client.TestConnection(ctx, lvConn)

	now := time.Now()
	status := "ok"
	if testErr != nil {
		status = "error"
	}
	conn.LastStatus = &status
	conn.LastCheckedAt = &now
	if updateErr := h.repo.Update(ctx, *conn); updateErr != nil {
		h.log.Warn("failed to persist last check status", "id", id, "error", updateErr)
	}

	if testErr != nil {
		h.log.Warn("libvirt connection test failed", "id", id, "name", conn.Name, "error", testErr)
		return c.JSON(http.StatusBadGateway, LibvirtConnectionTestResponse{
			Status:  "error",
			Message: testErr.Error(),
		})
	}

	return c.JSON(http.StatusOK, LibvirtConnectionTestResponse{
		Status:         "ok",
		LibvirtVersion: result.LibvirtVersion,
		HypervisorType: result.HypervisorType,
		TotalDomains:   result.TotalDomains,
		ActiveDomains:  result.ActiveDomains,
	})
}

// validateRequest checks per-type required fields and SSH key validity.
func (h *LibvirtConnectionHandler) validateRequest(req *LibvirtConnectionRequest) (*domain.LibvirtConnection, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	switch req.Type {
	case libvirt.TypeSSH, libvirt.TypeSocket:
	default:
		return nil, fmt.Errorf("type must be %q or %q", libvirt.TypeSSH, libvirt.TypeSocket)
	}

	conn := &domain.LibvirtConnection{
		Name:                 req.Name,
		Type:                 req.Type,
		Host:                 req.Host,
		Username:             req.Username,
		SSHKeyPath:           req.SSHKeyPath,
		AcceptUnknownHostKey: req.AcceptUnknownHostKey,
		SocketPath:           req.SocketPath,
		Description:          req.Description,
	}

	switch req.Type {
	case libvirt.TypeSSH:
		if req.Host == nil || *req.Host == "" {
			return nil, fmt.Errorf("host is required for ssh connections")
		}
		if req.Username == nil || *req.Username == "" {
			return nil, fmt.Errorf("username is required for ssh connections")
		}
		if req.SSHKeyPath == nil || *req.SSHKeyPath == "" {
			return nil, fmt.Errorf("ssh_key_path is required for ssh connections")
		}
		if err := libvirt.ValidateSSHKey(*req.SSHKeyPath); err != nil {
			return nil, err
		}
	case libvirt.TypeSocket:
		if req.SocketPath == nil || *req.SocketPath == "" {
			return nil, fmt.Errorf("socket_path is required for socket connections")
		}
	}

	return conn, nil
}

func (h *LibvirtConnectionHandler) toLibvirtConn(conn *domain.LibvirtConnection) libvirt.Connection {
	lvConn := libvirt.Connection{
		Type:                 conn.Type,
		AcceptUnknownHostKey: conn.AcceptUnknownHostKey,
		KnownHostsFile:       h.knownHosts,
		Timeout:              h.timeout,
	}
	if conn.Host != nil {
		lvConn.Host = *conn.Host
	}
	if conn.Username != nil {
		lvConn.Username = *conn.Username
	}
	if conn.SSHKeyPath != nil {
		lvConn.SSHKeyPath = *conn.SSHKeyPath
	}
	if conn.SocketPath != nil {
		lvConn.SocketPath = *conn.SocketPath
	}
	return lvConn
}

func parseID(c echo.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, fmt.Errorf("invalid connection id %q: must be a number", c.Param("id"))
	}
	return id, nil
}

func toResponse(conn domain.LibvirtConnection) LibvirtConnectionResponse {
	return LibvirtConnectionResponse{
		ID:                   conn.ID,
		Name:                 conn.Name,
		Type:                 conn.Type,
		Host:                 conn.Host,
		Username:             conn.Username,
		SSHKeyPath:           conn.SSHKeyPath,
		AcceptUnknownHostKey: conn.AcceptUnknownHostKey,
		SocketPath:           conn.SocketPath,
		Description:          conn.Description,
		LastStatus:           conn.LastStatus,
		LastCheckedAt:        conn.LastCheckedAt,
		CreatedAt:            conn.CreatedAt,
		UpdatedAt:            conn.UpdatedAt,
	}
}

func (h *LibvirtConnectionHandler) badRequest(c echo.Context, msg string) error {
	return c.JSON(http.StatusBadRequest, ErrorResponse{Error: msg})
}

func (h *LibvirtConnectionHandler) notFound(c echo.Context, msg string) error {
	return c.JSON(http.StatusNotFound, ErrorResponse{Error: msg})
}

func (h *LibvirtConnectionHandler) conflict(c echo.Context, msg string) error {
	return c.JSON(http.StatusConflict, ErrorResponse{Error: msg})
}

func (h *LibvirtConnectionHandler) internalError(c echo.Context, err error) error {
	h.log.Error("internal error", "error", err)
	return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
}
