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

// AuditLogRepo is the SQL-backed implementation of AuditLogRepository.
type AuditLogRepo struct {
	db  *sql.DB
	sb  sq.StatementBuilderType
	log *slog.Logger
}

// NewAuditLogRepo creates a new AuditLogRepo.
func NewAuditLogRepo(db *sql.DB, dialect Dialect, log *slog.Logger) *AuditLogRepo {
	sb := sq.StatementBuilder
	if dialect == DialectPostgres {
		sb = sb.PlaceholderFormat(sq.Dollar)
	}
	return &AuditLogRepo{db: db, sb: sb, log: log}
}

func (r *AuditLogRepo) Create(ctx context.Context, logEntry domain.AuditLog) (int, error) {
	stmt := r.sb.Insert("audit_logs").
		Columns("user_id", "action", "details", "created_at").
		Values(logEntry.UserID, logEntry.Action, logEntry.Details, time.Now())

	result, err := sq.ExecContextWith(ctx, r.db, stmt)
	if err != nil {
		return 0, fmt.Errorf("failed to insert audit log: %w", err)
	}
	id, _ := result.LastInsertId()
	return int(id), nil
}

func (r *AuditLogRepo) List(ctx context.Context, limit, offset int) ([]domain.AuditLog, error) {
	stmt := r.sb.Select("id", "user_id", "action", "details", "created_at").
		From("audit_logs").
		OrderBy("created_at DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	query, args, err := stmt.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list audit logs: %w", err)
	}
	defer rows.Close()

	var logs []domain.AuditLog
	for rows.Next() {
		var l domain.AuditLog
		if err := rows.Scan(&l.ID, &l.UserID, &l.Action, &l.Details, &l.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}
		logs = append(logs, l)
	}
	return logs, nil
}

func (r *AuditLogRepo) GetByUserID(ctx context.Context, userID int, limit, offset int) ([]domain.AuditLog, error) {
	stmt := r.sb.Select("id", "user_id", "action", "details", "created_at").
		From("audit_logs").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("created_at DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	query, args, err := stmt.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs by user: %w", err)
	}
	defer rows.Close()

	var logs []domain.AuditLog
	for rows.Next() {
		var l domain.AuditLog
		if err := rows.Scan(&l.ID, &l.UserID, &l.Action, &l.Details, &l.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}
		logs = append(logs, l)
	}
	return logs, nil
}
