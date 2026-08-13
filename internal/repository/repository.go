package repository

import (
	"context"

	"github.com/notepass/vmprov-web/internal/domain"
)

// UserRepository defines the interface for user data access.
type UserRepository interface {
	Create(ctx context.Context, user domain.User) (int, error)
	GetByID(ctx context.Context, id int) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	List(ctx context.Context) ([]domain.User, error)
}

// TemplateRepository defines the interface for template data access.
type TemplateRepository interface {
	Create(ctx context.Context, template domain.Template) (int, error)
	GetByID(ctx context.Context, id int) (*domain.Template, error)
	GetByName(ctx context.Context, name string) (*domain.Template, error)
	Update(ctx context.Context, template domain.Template) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]domain.Template, error)
}

// AuditLogRepository defines the interface for audit log data access.
type AuditLogRepository interface {
	Create(ctx context.Context, log domain.AuditLog) (int, error)
	List(ctx context.Context, limit, offset int) ([]domain.AuditLog, error)
	GetByUserID(ctx context.Context, userID int, limit, offset int) ([]domain.AuditLog, error)
}
