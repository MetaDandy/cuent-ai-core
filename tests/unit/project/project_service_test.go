//go:build unit

package project_test

import (
	"errors"
	"testing"

	"github.com/MetaDandy/cuent-ai-core/helper"
	"github.com/MetaDandy/cuent-ai-core/src/model"
	"github.com/MetaDandy/cuent-ai-core/src/modules/project"
	"github.com/google/uuid"
)

// --- Mock Repository ---

type mockRepository struct {
	createFunc           func(project *model.Project) error
	updateFunc           func(project *model.Project) error
	findAllFunc          func(opts *helper.FindAllOptions) ([]model.Project, int64, error)
	findByIdFunc         func(id string) (*model.Project, error)
	findByIdUnscopedFunc func(id string) (*model.Project, error)
	softDeleteFunc       func(id string) error
	restoreFunc          func(id string) error
}

func (m *mockRepository) Create(p *model.Project) error {
	if m.createFunc != nil {
		return m.createFunc(p)
	}
	return nil
}
func (m *mockRepository) Update(p *model.Project) error {
	if m.updateFunc != nil {
		return m.updateFunc(p)
	}
	return nil
}
func (m *mockRepository) FindAll(opts *helper.FindAllOptions) ([]model.Project, int64, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(opts)
	}
	return nil, 0, nil
}
func (m *mockRepository) FindById(id string) (*model.Project, error) {
	if m.findByIdFunc != nil {
		return m.findByIdFunc(id)
	}
	return nil, nil
}
func (m *mockRepository) FindByIdUnscoped(id string) (*model.Project, error) {
	if m.findByIdUnscopedFunc != nil {
		return m.findByIdUnscopedFunc(id)
	}
	return nil, nil
}
func (m *mockRepository) SoftDelete(id string) error {
	if m.softDeleteFunc != nil {
		return m.softDeleteFunc(id)
	}
	return nil
}
func (m *mockRepository) Restore(id string) error {
	if m.restoreFunc != nil {
		return m.restoreFunc(id)
	}
	return nil
}

// --- Tests ---

func TestService_FindByID(t *testing.T) {
	// Datos de prueba
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
		mockBehavior  func(m *mockRepository)
		expectedError bool
		expectedFound bool
	}{
		{
			name:    "Success - Project Found",
			inputID: projectID.String(),
			mockBehavior: func(m *mockRepository) {
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
			mockBehavior: func(m *mockRepository) {
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
			mockBehavior: func(m *mockRepository) {
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
			// Setup
			mockRepo := &mockRepository{}
			tt.mockBehavior(mockRepo)

			// Inject mock (pass nil for userRepo as it's not used in FindByID)
			svc := project.NewService(mockRepo, nil)

			// Execute
			result, err := svc.FindByID(tt.inputID)

			// Verify
			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
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
		})
	}
}
