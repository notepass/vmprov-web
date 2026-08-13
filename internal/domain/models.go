package domain

import "time"

// User represents an application user.
type User struct {
	ID        int       `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
}

// Template represents a cloud-init template.
type Template struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// AuditLog represents an audit log entry.
type AuditLog struct {
	ID        int       `db:"id"`
	UserID    *int      `db:"user_id"`
	Action    string    `db:"action"`
	Details   *string   `db:"details"`
	CreatedAt time.Time `db:"created_at"`
}
