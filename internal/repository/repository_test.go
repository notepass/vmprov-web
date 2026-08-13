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
