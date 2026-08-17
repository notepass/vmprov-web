package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/notepass/vmprov-web/internal/domain"
)

// LibvirtConnectionRepo is the SQL-backed implementation of LibvirtConnectionRepository.
type LibvirtConnectionRepo struct {
	db       *sql.DB
	sb       sq.StatementBuilderType
	log      *slog.Logger
	hostCol  string
	postgres bool
}

// NewLibvirtConnectionRepo creates a new LibvirtConnectionRepo.
func NewLibvirtConnectionRepo(db *sql.DB, dialect Dialect, log *slog.Logger) *LibvirtConnectionRepo {
	sb := sq.StatementBuilder
	hostCol := "host"
	postgres := false
	if dialect == DialectMySQL {
		sb = sb.PlaceholderFormat(sq.Question)
		hostCol = "`host`"
	} else {
		sb = sb.PlaceholderFormat(sq.Dollar)
		postgres = true
	}
	return &LibvirtConnectionRepo{db: db, sb: sb, log: log, hostCol: hostCol, postgres: postgres}
}

func (r *LibvirtConnectionRepo) columns() []string {
	return []string{"id", "name", "type", r.hostCol, "username", "ssh_key_path",
		"accept_unknown_host_key", "socket_path", "description",
		"last_status", "last_checked_at", "created_at", "updated_at"}
}

func (r *LibvirtConnectionRepo) scanConn(row interface {
	Scan(dest ...any) error
}) (*domain.LibvirtConnection, error) {
	c := &domain.LibvirtConnection{}
	if err := row.Scan(&c.ID, &c.Name, &c.Type, &c.Host, &c.Username, &c.SSHKeyPath,
		&c.AcceptUnknownHostKey, &c.SocketPath, &c.Description,
		&c.LastStatus, &c.LastCheckedAt, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return nil, err
	}
	return c, nil
}

func (r *LibvirtConnectionRepo) Create(ctx context.Context, conn domain.LibvirtConnection) (int, error) {
	now := time.Now()
	stmt := r.sb.Insert("libvirt_connections").
		Columns("name", "type", r.hostCol, "username", "ssh_key_path",
			"accept_unknown_host_key", "socket_path", "description",
			"created_at", "updated_at").
		Values(conn.Name, conn.Type, conn.Host, conn.Username, conn.SSHKeyPath,
			conn.AcceptUnknownHostKey, conn.SocketPath, conn.Description,
			now, now)

	if r.postgres {
		// The pgx stdlib driver does not support LastInsertId, so the
		// generated id is returned explicitly.
		stmt = stmt.Suffix("RETURNING id")
	}

	query, args, err := stmt.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %w", err)
	}

	if r.postgres {
		var id int
		if err := r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
			return 0, fmt.Errorf("failed to insert libvirt connection: %w", err)
		}
		return id, nil
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to insert libvirt connection: %w", err)
	}
	id, _ := result.LastInsertId()
	return int(id), nil
}

func (r *LibvirtConnectionRepo) GetByID(ctx context.Context, id int) (*domain.LibvirtConnection, error) {
	query, args, err := r.sb.Select(r.columns()...).
		From("libvirt_connections").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	row := r.db.QueryRowContext(ctx, query, args...)
	c, err := r.scanConn(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get libvirt connection: %w", err)
	}
	return c, nil
}

func (r *LibvirtConnectionRepo) GetByName(ctx context.Context, name string) (*domain.LibvirtConnection, error) {
	query, args, err := r.sb.Select(r.columns()...).
		From("libvirt_connections").
		Where(sq.Eq{"name": name}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	row := r.db.QueryRowContext(ctx, query, args...)
	c, err := r.scanConn(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get libvirt connection: %w", err)
	}
	return c, nil
}

func (r *LibvirtConnectionRepo) Update(ctx context.Context, conn domain.LibvirtConnection) error {
	query, args, err := r.sb.Update("libvirt_connections").
		Set("name", conn.Name).
		Set("type", conn.Type).
		Set(r.hostCol, conn.Host).
		Set("username", conn.Username).
		Set("ssh_key_path", conn.SSHKeyPath).
		Set("accept_unknown_host_key", conn.AcceptUnknownHostKey).
		Set("socket_path", conn.SocketPath).
		Set("description", conn.Description).
		Set("last_status", conn.LastStatus).
		Set("last_checked_at", conn.LastCheckedAt).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": conn.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update libvirt connection: %w", err)
	}
	return nil
}

func (r *LibvirtConnectionRepo) Delete(ctx context.Context, id int) error {
	query, args, err := r.sb.Delete("libvirt_connections").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete libvirt connection: %w", err)
	}
	return nil
}

func (r *LibvirtConnectionRepo) List(ctx context.Context) ([]domain.LibvirtConnection, error) {
	query, args, err := r.sb.Select(r.columns()...).
		From("libvirt_connections").
		OrderBy("id ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list libvirt connections: %w", err)
	}
	defer rows.Close()

	conns := []domain.LibvirtConnection{}
	for rows.Next() {
		c, err := r.scanConn(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan libvirt connection: %w", err)
		}
		conns = append(conns, *c)
	}
	return conns, nil
}
