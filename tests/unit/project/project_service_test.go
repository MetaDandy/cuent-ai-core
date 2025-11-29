//go:build unit

package project_test

import (
	"errors"
	"testing"

	"github.com/MetaDandy/cuent-ai-core/helper"
	"github.com/MetaDandy/cuent-ai-core/src/model"
	"github.com/MetaDandy/cuent-ai-core/src/modules/project"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Mock Project Repository
type mockProjectRepository struct {
	projects            map[string]*model.Project
	createFunc          func(project *model.Project) error
	updateFunc          func(project *model.Project) error
	findAllFunc         func(opts *helper.FindAllOptions) ([]model.Project, int64, error)
	findByIdFunc        func(id string) (*model.Project, error)
	findByIdUnscopedFunc func(id string) (*model.Project, error)
	softDeleteFunc      func(id string) error
	restoreFunc         func(id string) error
}

func (m *mockProjectRepository) Create(p *model.Project) error {
	if m.createFunc != nil {
		return m.createFunc(p)
	}
	if m.projects == nil {
		m.projects = make(map[string]*model.Project)
	}
	m.projects[p.ID.String()] = p
	return nil
}

func (m *mockProjectRepository) Update(p *model.Project) error {
	if m.updateFunc != nil {
		return m.updateFunc(p)
	}
	if _, ok := m.projects[p.ID.String()]; !ok {
		return gorm.ErrRecordNotFound
	}
	m.projects[p.ID.String()] = p
	return nil
}

func (m *mockProjectRepository) FindAll(opts *helper.FindAllOptions) ([]model.Project, int64, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(opts)
	}
	projects := make([]model.Project, 0, len(m.projects))
	for _, p := range m.projects {
		projects = append(projects, *p)
	}
	return projects, int64(len(projects)), nil
}

func (m *mockProjectRepository) FindById(id string) (*model.Project, error) {
	if m.findByIdFunc != nil {
		return m.findByIdFunc(id)
	}
	if project, ok := m.projects[id]; ok {
		return project, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockProjectRepository) FindByIdUnscoped(id string) (*model.Project, error) {
	if m.findByIdUnscopedFunc != nil {
		return m.findByIdUnscopedFunc(id)
	}
	if project, ok := m.projects[id]; ok {
		return project, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockProjectRepository) SoftDelete(id string) error {
	if m.softDeleteFunc != nil {
		return m.softDeleteFunc(id)
	}
	if _, ok := m.projects[id]; !ok {
		return gorm.ErrRecordNotFound
	}
	delete(m.projects, id)
	return nil
}

func (m *mockProjectRepository) Restore(id string) error {
	if m.restoreFunc != nil {
		return m.restoreFunc(id)
	}
	return nil
}

// Mock User Repository
type mockUserRepository struct {
	users        map[string]*model.User
	findByIdFunc func(id string) (*model.User, error)
}

func (m *mockUserRepository) FindById(id string) (*model.User, error) {
	if m.findByIdFunc != nil {
		return m.findByIdFunc(id)
	}
	if user, ok := m.users[id]; ok {
		return user, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func TestService_FindByID(t *testing.T) {
	projectID := uuid.New()
	mockProject := &model.Project{
		ID:          projectID,
		Name:        "Test Project",
		Description: "Description",
		State:       model.StatePending,
	}

	tests := []struct {
		name          string
		inputID       string
		mockSetup     func(*mockProjectRepository)
		expectedError bool
		expectedFound bool
	}{
		{
			name:    "Success - Project Found",
			inputID: projectID.String(),
			mockSetup: func(m *mockProjectRepository) {
				m.findByIdFunc = func(id string) (*model.Project, error) {
					if id == projectID.String() {
						return mockProject, nil
					}
					return nil, nil
				}
			},
			expectedError: false,
			expectedFound: true,
		},
		{
			name:    "Success - Project Not Found",
			inputID: "non-existent-id",
			mockSetup: func(m *mockProjectRepository) {
				m.findByIdFunc = func(id string) (*model.Project, error) {
					return nil, nil
				}
			},
			expectedError: false,
			expectedFound: false,
		},
		{
			name:    "Failure - Repository Error",
			inputID: projectID.String(),
			mockSetup: func(m *mockProjectRepository) {
				m.findByIdFunc = func(id string) (*model.Project, error) {
					return nil, errors.New("db error")
				}
			},
			expectedError: true,
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockProjectRepository{}
			tt.mockSetup(mockRepo)

			svc := project.NewService(mockRepo, nil)

			result, err := svc.FindByID(tt.inputID)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				if tt.expectedFound {
					if result == nil {
						t.Errorf("expected result but got nil")
					} else if result.ID != projectID.String() {
						t.Errorf("expected project ID %v, got %v", projectID, result.ID)
					}
				} else {
					if result != nil {
						t.Errorf("expected nil result but got %v", result)
					}
				}
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	projectID := uuid.New()
	mockProject := &model.Project{
		ID:          projectID,
		Name:        "Old Name",
		Description: "Old Description",
		State:       model.StatePending,
	}

	tests := []struct {
		name          string
		inputID       string
		input         *project.ProjectUpdate
		mockSetup     func(*mockProjectRepository)
		expectedError bool
		expectedName  string
	}{
		{
			name:    "Success - Update Name",
			inputID: projectID.String(),
			input: &project.ProjectUpdate{
				Name: stringPtr("New Name"),
			},
			mockSetup: func(m *mockProjectRepository) {
				m.projects = map[string]*model.Project{
					projectID.String(): mockProject,
				}
				m.findByIdFunc = func(id string) (*model.Project, error) {
					if id == projectID.String() {
						return mockProject, nil
					}
					return nil, nil
				}
			},
			expectedError: false,
			expectedName:  "New Name",
		},
		{
			name:    "Success - Update Description",
			inputID: projectID.String(),
			input: &project.ProjectUpdate{
				Description: stringPtr("New Description"),
			},
			mockSetup: func(m *mockProjectRepository) {
				m.projects = map[string]*model.Project{
					projectID.String(): mockProject,
				}
				m.findByIdFunc = func(id string) (*model.Project, error) {
					if id == projectID.String() {
						return mockProject, nil
					}
					return nil, nil
				}
			},
			expectedError: false,
			expectedName:  "Old Name",
		},
		{
			name:    "Success - Update Both",
			inputID: projectID.String(),
			input: &project.ProjectUpdate{
				Name:        stringPtr("New Name"),
				Description: stringPtr("New Description"),
			},
			mockSetup: func(m *mockProjectRepository) {
				m.projects = map[string]*model.Project{
					projectID.String(): mockProject,
				}
				m.findByIdFunc = func(id string) (*model.Project, error) {
					if id == projectID.String() {
						return mockProject, nil
					}
					return nil, nil
				}
			},
			expectedError: false,
			expectedName:  "New Name",
		},
		{
			name:    "Success - Project Not Found",
			inputID: "non-existent-id",
			input: &project.ProjectUpdate{
				Name: stringPtr("New Name"),
			},
			mockSetup: func(m *mockProjectRepository) {
				m.findByIdFunc = func(id string) (*model.Project, error) {
					return nil, nil
				}
			},
			expectedError: false,
			expectedName:  "",
		},
		{
			name:    "Failure - Update Error",
			inputID: projectID.String(),
			input: &project.ProjectUpdate{
				Name: stringPtr("New Name"),
			},
			mockSetup: func(m *mockProjectRepository) {
				m.findByIdFunc = func(id string) (*model.Project, error) {
					return mockProject, nil
				}
				m.updateFunc = func(p *model.Project) error {
					return errors.New("update failed")
				}
			},
			expectedError: true,
			expectedName:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockProjectRepository{
				projects: make(map[string]*model.Project),
			}
			tt.mockSetup(mockRepo)

			svc := project.NewService(mockRepo, nil)
			result, err := svc.Update(tt.inputID, tt.input)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != nil && tt.input.Name != nil {
					if result.Name != tt.expectedName {
						t.Errorf("expected name %s, got %s", tt.expectedName, result.Name)
					}
				}
			}
		})
	}
}

func TestService_SoftDelete(t *testing.T) {
	projectID := uuid.New()
	mockProject := &model.Project{
		ID:    projectID,
		Name:  "Test Project",
		State: model.StatePending,
	}

	tests := []struct {
		name            string
		inputID         string
		mockSetup       func(*mockProjectRepository)
		expectedError   bool
		expectedDeleted bool
	}{
		{
			name:    "Success - Soft Delete",
			inputID: projectID.String(),
			mockSetup: func(m *mockProjectRepository) {
				m.projects = map[string]*model.Project{
					projectID.String(): mockProject,
				}
				m.findByIdFunc = func(id string) (*model.Project, error) {
					if id == projectID.String() {
						return mockProject, nil
					}
					return nil, nil
				}
			},
			expectedError:   false,
			expectedDeleted: true,
		},
		{
			name:    "Success - Project Not Found",
			inputID: "non-existent-id",
			mockSetup: func(m *mockProjectRepository) {
				m.findByIdFunc = func(id string) (*model.Project, error) {
					return nil, nil
				}
			},
			expectedError:   false,
			expectedDeleted: false,
		},
		{
			name:    "Failure - Delete Error",
			inputID: projectID.String(),
			mockSetup: func(m *mockProjectRepository) {
				m.findByIdFunc = func(id string) (*model.Project, error) {
					return mockProject, nil
				}
				m.softDeleteFunc = func(id string) error {
					return errors.New("delete failed")
				}
			},
			expectedError:   false, // El servicio retorna false, nil en caso de error
			expectedDeleted: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockProjectRepository{
				projects: make(map[string]*model.Project),
			}
			tt.mockSetup(mockRepo)

			svc := project.NewService(mockRepo, nil)
			deleted, err := svc.SoftDelete(tt.inputID)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if deleted != tt.expectedDeleted {
					t.Errorf("expected deleted=%v, got %v", tt.expectedDeleted, deleted)
				}
			}
		})
	}
}

func TestService_Restore(t *testing.T) {
	projectID := uuid.New()
	mockProject := &model.Project{
		ID:    projectID,
		Name:  "Test Project",
		State: model.StatePending,
	}

	tests := []struct {
		name          string
		inputID       string
		mockSetup     func(*mockProjectRepository)
		expectedError bool
		expectedFound bool
	}{
		{
			name:    "Success - Restore Project",
			inputID: projectID.String(),
			mockSetup: func(m *mockProjectRepository) {
				m.projects = map[string]*model.Project{
					projectID.String(): mockProject,
				}
				m.findByIdUnscopedFunc = func(id string) (*model.Project, error) {
					if id == projectID.String() {
						return mockProject, nil
					}
					return nil, gorm.ErrRecordNotFound
				}
			},
			expectedError: false,
			expectedFound: true,
		},
		{
			name:    "Failure - Project Not Found",
			inputID: "non-existent-id",
			mockSetup: func(m *mockProjectRepository) {
				m.findByIdUnscopedFunc = func(id string) (*model.Project, error) {
					return nil, gorm.ErrRecordNotFound
				}
			},
			expectedError: true,
			expectedFound: false,
		},
		{
			name:    "Failure - Restore Error",
			inputID: projectID.String(),
			mockSetup: func(m *mockProjectRepository) {
				m.findByIdUnscopedFunc = func(id string) (*model.Project, error) {
					return mockProject, nil
				}
				m.restoreFunc = func(id string) error {
					return errors.New("restore failed")
				}
			},
			expectedError: true,
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockProjectRepository{
				projects: make(map[string]*model.Project),
			}
			tt.mockSetup(mockRepo)

			svc := project.NewService(mockRepo, nil)
			result, err := svc.Restore(tt.inputID)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tt.expectedFound {
					if result == nil {
						t.Errorf("expected result but got nil")
					} else if result.ID != projectID.String() {
						t.Errorf("expected project ID %v, got %v", projectID, result.ID)
					}
				}
			}
		})
	}
}

func TestService_FindAll(t *testing.T) {
	mockProject1 := &model.Project{
		ID:    uuid.New(),
		Name:  "Project 1",
		State: model.StatePending,
	}
	mockProject2 := &model.Project{
		ID:    uuid.New(),
		Name:  "Project 2",
		State: model.StateActive,
	}

	tests := []struct {
		name          string
		opts          *helper.FindAllOptions
		mockSetup     func(*mockProjectRepository)
		expectedCount int
		expectedError bool
	}{
		{
			name: "Success - Find All Projects",
			opts: &helper.FindAllOptions{
				Limit:  10,
				Offset: 0,
			},
			mockSetup: func(m *mockProjectRepository) {
				m.projects = map[string]*model.Project{
					mockProject1.ID.String(): mockProject1,
					mockProject2.ID.String(): mockProject2,
				}
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name: "Success - Empty Result",
			opts: &helper.FindAllOptions{
				Limit:  10,
				Offset: 0,
			},
			mockSetup: func(m *mockProjectRepository) {
				m.projects = make(map[string]*model.Project)
			},
			expectedCount: 0,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockProjectRepository{
				projects: make(map[string]*model.Project),
			}
			tt.mockSetup(mockRepo)

			svc := project.NewService(mockRepo, nil)
			result, err := svc.FindAll(tt.opts)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != nil {
					if len(result.Data) != tt.expectedCount {
						t.Errorf("expected %d projects, got %d", tt.expectedCount, len(result.Data))
					}
				}
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

