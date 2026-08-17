// Package libvirt wraps the digitalocean/go-libvirt dialers behind an
// injectable Client interface for connection testing.
package libvirt

import (
	"context"
	"errors"
	"fmt"
	"time"

	goLibvirt "github.com/digitalocean/go-libvirt"
	"github.com/digitalocean/go-libvirt/socket"
	"github.com/digitalocean/go-libvirt/socket/dialers"
)

// Connection types.
const (
	TypeSSH    = "ssh"
	TypeSocket = "socket"
)

// DefaultTimeout is used when a connection has no timeout set.
const DefaultTimeout = 10 * time.Second

// Connection describes a typed libvirt endpoint to test.
type Connection struct {
	Type                 string
	Host                 string
	Username             string
	SSHKeyPath           string
	AcceptUnknownHostKey bool
	SocketPath           string
	KnownHostsFile       string
	Timeout              time.Duration
}

// TestResult is the outcome of a successful connection test.
type TestResult struct {
	LibvirtVersion string
	HypervisorType string
	TotalDomains   int
	ActiveDomains  int
}

// Client abstracts libvirt connection testing so handlers can be tested
// without a live libvirt host.
type Client interface {
	TestConnection(ctx context.Context, conn Connection) (*TestResult, error)
}

type client struct{}

// New returns a Client backed by the go-libvirt dialers.
func New() Client {
	return &client{}
}

func (c *client) TestConnection(ctx context.Context, conn Connection) (*TestResult, error) {
	if conn.Timeout <= 0 {
		conn.Timeout = DefaultTimeout
	}
	ctx, cancel := context.WithTimeout(ctx, conn.Timeout)
	defer cancel()

	type outcome struct {
		result *TestResult
		err    error
	}
	ch := make(chan outcome, 1)
	go func() {
		result, err := c.testConn(conn)
		ch <- outcome{result: result, err: err}
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("libvirt connection test timed out after %s: %w", conn.Timeout, ctx.Err())
	case o := <-ch:
		return o.result, o.err
	}
}

func (c *client) buildDialer(conn Connection) (socket.Dialer, error) {
	switch conn.Type {
	case TypeSSH:
		if conn.Host == "" {
			return nil, errors.New("host is required for ssh connections")
		}
		opts := []dialers.SSHOption{
			dialers.UseSSHUsername(conn.Username),
			dialers.UseKeyFile(conn.SSHKeyPath),
			dialers.WithSSHAuthMethods((&dialers.SSHAuthMethods{}).PrivKey()),
			dialers.UseKnownHostsFile(conn.KnownHostsFile),
		}
		if conn.AcceptUnknownHostKey {
			opts = append(opts, dialers.WithAcceptUnknownHostKey())
		}
		return dialers.NewSSH(conn.Host, opts...), nil
	case TypeSocket:
		if conn.SocketPath == "" {
			return nil, errors.New("socket_path is required for socket connections")
		}
		return dialers.NewLocal(
			dialers.WithSocket(conn.SocketPath),
			dialers.WithLocalTimeout(conn.Timeout),
		), nil
	default:
		return nil, fmt.Errorf("unsupported connection type %q", conn.Type)
	}
}

func (c *client) testConn(conn Connection) (*TestResult, error) {
	dialer, err := c.buildDialer(conn)
	if err != nil {
		return nil, err
	}

	lv := goLibvirt.NewWithDialer(dialer)
	defer lv.Disconnect()

	if err := lv.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to libvirt: %w", err)
	}

	version, err := lv.ConnectGetLibVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get libvirt version: %w", err)
	}

	hypervisorType, err := lv.ConnectGetType()
	if err != nil {
		return nil, fmt.Errorf("failed to get hypervisor type: %w", err)
	}

	allDomains, _, err := lv.ConnectListAllDomains(0,
		goLibvirt.ConnectListDomainsActive|goLibvirt.ConnectListDomainsInactive)
	if err != nil {
		return nil, fmt.Errorf("failed to list domains: %w", err)
	}

	activeDomains, _, err := lv.ConnectListAllDomains(0, goLibvirt.ConnectListDomainsActive)
	if err != nil {
		return nil, fmt.Errorf("failed to list active domains: %w", err)
	}

	return &TestResult{
		LibvirtVersion: formatLibVersion(version),
		HypervisorType: hypervisorType,
		TotalDomains:   len(allDomains),
		ActiveDomains:  len(activeDomains),
	}, nil
}

// formatLibVersion converts the libvirt version integer into a dotted string.
// The integer encodes major*1,000,000 + minor*1,000 + micro.
func formatLibVersion(v uint64) string {
	major := v / 1000000
	v %= 1000000
	minor := v / 1000
	v %= 1000
	micro := v
	return fmt.Sprintf("%d.%d.%d", major, minor, micro)
}
