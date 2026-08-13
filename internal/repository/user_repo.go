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

// Dialect is a database dialect identifier.
type Dialect string

const (
	DialectPostgres Dialect = "postgres"
	DialectMySQL    Dialect = "mysql"
)

// UserRepo is the SQL-backed implementation of UserRepository.
type UserRepo struct {
	db  *sql.DB
	sb  sq.StatementBuilderType
	log *slog.Logger
}

// NewUserRepo creates a new UserRepo.
func NewUserRepo(db *sql.DB, dialect Dialect, log *slog.Logger) *UserRepo {
	sb := sq.StatementBuilder
	if dialect == DialectPostgres {
		sb = sb.PlaceholderFormat(sq.Dollar)
	}
	return &UserRepo{db: db, sb: sb, log: log}
}

func (r *UserRepo) Create(ctx context.Context, user domain.User) (int, error) {
	stmt := r.sb.Insert("users").
		Columns("username", "email", "created_at").
		Values(user.Username, user.Email, time.Now())

	result, err := sq.ExecContextWith(ctx, r.db, stmt)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}
	id, _ := result.LastInsertId()
	return int(id), nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int) (*domain.User, error) {
	query, args, err := r.sb.Select("id", "username", "email", "created_at").
		From("users").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	user := &domain.User{}
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query, args, err := r.sb.Select("id", "username", "email", "created_at").
		From("users").
		Where(sq.Eq{"username": username}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	user := &domain.User{}
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *UserRepo) List(ctx context.Context) ([]domain.User, error) {
	query, args, err := r.sb.Select("id", "username", "email", "created_at").
		From("users").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}
