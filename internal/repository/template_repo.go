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

// TemplateRepo is the SQL-backed implementation of TemplateRepository.
type TemplateRepo struct {
	db  *sql.DB
	sb  sq.StatementBuilderType
	log *slog.Logger
}

// NewTemplateRepo creates a new TemplateRepo.
func NewTemplateRepo(db *sql.DB, dialect Dialect, log *slog.Logger) *TemplateRepo {
	sb := sq.StatementBuilder
	if dialect == DialectPostgres {
		sb = sb.PlaceholderFormat(sq.Dollar)
	}
	return &TemplateRepo{db: db, sb: sb, log: log}
}

func (r *TemplateRepo) Create(ctx context.Context, template domain.Template) (int, error) {
	stmt := r.sb.Insert("templates").
		Columns("name", "content", "created_at", "updated_at").
		Values(template.Name, template.Content, time.Now(), time.Now())

	result, err := sq.ExecContextWith(ctx, r.db, stmt)
	if err != nil {
		return 0, fmt.Errorf("failed to insert template: %w", err)
	}
	id, _ := result.LastInsertId()
	return int(id), nil
}

func (r *TemplateRepo) GetByID(ctx context.Context, id int) (*domain.Template, error) {
	query, args, err := r.sb.Select("id", "name", "content", "created_at", "updated_at").
		From("templates").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	t := &domain.Template{}
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&t.ID, &t.Name, &t.Content, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get template: %w", err)
	}
	return t, nil
}

func (r *TemplateRepo) GetByName(ctx context.Context, name string) (*domain.Template, error) {
	query, args, err := r.sb.Select("id", "name", "content", "created_at", "updated_at").
		From("templates").
		Where(sq.Eq{"name": name}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	t := &domain.Template{}
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&t.ID, &t.Name, &t.Content, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get template: %w", err)
	}
	return t, nil
}

func (r *TemplateRepo) Update(ctx context.Context, template domain.Template) error {
	query, args, err := r.sb.Update("templates").
		Set("name", template.Name).
		Set("content", template.Content).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": template.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update template: %w", err)
	}
	return nil
}

func (r *TemplateRepo) Delete(ctx context.Context, id int) error {
	query, args, err := r.sb.Delete("templates").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}
	return nil
}

func (r *TemplateRepo) List(ctx context.Context) ([]domain.Template, error) {
	query, args, err := r.sb.Select("id", "name", "content", "created_at", "updated_at").
		From("templates").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}
	defer rows.Close()

	var templates []domain.Template
	for rows.Next() {
		var t domain.Template
		if err := rows.Scan(&t.ID, &t.Name, &t.Content, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan template: %w", err)
		}
		templates = append(templates, t)
	}
	return templates, nil
}
