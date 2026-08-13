package repository

import (
	"context"
	"sync"

	"github.com/notepass/vmprov-web/internal/domain"
)

// MockUserRepository is a thread-safe mock implementation of UserRepository.
type MockUserRepository struct {
	mu      sync.RWMutex
	users   map[int]domain.User
	nextID  int
	byName  map[string]int
}

// NewMockUserRepository creates a new mock user repository.
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:  make(map[int]domain.User),
		byName: make(map[string]int),
	}
}

func (m *MockUserRepository) Create(_ context.Context, user domain.User) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nextID++
	id := m.nextID
	user.ID = id
	m.users[id] = user
	m.byName[user.Username] = id
	return id, nil
}

func (m *MockUserRepository) GetByID(_ context.Context, id int) (*domain.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	u, ok := m.users[id]
	if !ok {
		return nil, nil
	}
	return &u, nil
}

func (m *MockUserRepository) GetByUsername(_ context.Context, username string) (*domain.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	id, ok := m.byName[username]
	if !ok {
		return nil, nil
	}
	u := m.users[id]
	return &u, nil
}

func (m *MockUserRepository) List(_ context.Context) ([]domain.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	users := make([]domain.User, 0, len(m.users))
	for _, u := range m.users {
		users = append(users, u)
	}
	return users, nil
}

// MockTemplateRepository is a thread-safe mock implementation of TemplateRepository.
type MockTemplateRepository struct {
	mu        sync.RWMutex
	templates map[int]domain.Template
	nextID    int
	byName    map[string]int
}

// NewMockTemplateRepository creates a new mock template repository.
func NewMockTemplateRepository() *MockTemplateRepository {
	return &MockTemplateRepository{
		templates: make(map[int]domain.Template),
		byName:    make(map[string]int),
	}
}

func (m *MockTemplateRepository) Create(_ context.Context, t domain.Template) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nextID++
	id := m.nextID
	t.ID = id
	m.templates[id] = t
	m.byName[t.Name] = id
	return id, nil
}

func (m *MockTemplateRepository) GetByID(_ context.Context, id int) (*domain.Template, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.templates[id]
	if !ok {
		return nil, nil
	}
	return &t, nil
}

func (m *MockTemplateRepository) GetByName(_ context.Context, name string) (*domain.Template, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	id, ok := m.byName[name]
	if !ok {
		return nil, nil
	}
	t := m.templates[id]
	return &t, nil
}

func (m *MockTemplateRepository) Update(_ context.Context, t domain.Template) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.templates[t.ID]; !ok {
		return nil
	}
	m.templates[t.ID] = t
	return nil
}

func (m *MockTemplateRepository) Delete(_ context.Context, id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.templates, id)
	return nil
}

func (m *MockTemplateRepository) List(_ context.Context) ([]domain.Template, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	templates := make([]domain.Template, 0, len(m.templates))
	for _, t := range m.templates {
		templates = append(templates, t)
	}
	return templates, nil
}

// MockAuditLogRepository is a thread-safe mock implementation of AuditLogRepository.
type MockAuditLogRepository struct {
	mu   sync.RWMutex
	logs []domain.AuditLog
}

// NewMockAuditLogRepository creates a new mock audit log repository.
func NewMockAuditLogRepository() *MockAuditLogRepository {
	return &MockAuditLogRepository{
		logs: make([]domain.AuditLog, 0),
	}
}

func (m *MockAuditLogRepository) Create(_ context.Context, logEntry domain.AuditLog) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	id := len(m.logs) + 1
	logEntry.ID = id
	m.logs = append(m.logs, logEntry)
	return id, nil
}

func (m *MockAuditLogRepository) List(_ context.Context, limit, offset int) ([]domain.AuditLog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	end := offset + limit
	if end > len(m.logs) {
		end = len(m.logs)
	}
	return m.logs[offset:end], nil
}

func (m *MockAuditLogRepository) GetByUserID(_ context.Context, userID int, limit, offset int) ([]domain.AuditLog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var filtered []domain.AuditLog
	for _, l := range m.logs {
		if l.UserID != nil && *l.UserID == userID {
			filtered = append(filtered, l)
		}
	}
	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[offset:end], nil
}
