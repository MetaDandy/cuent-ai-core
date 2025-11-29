//go:build mocking
// +build mocking

package mocking

import (
	"errors"
	"testing"

	"github.com/MetaDandy/cuent-ai-core/src/modules/project"
	"github.com/stretchr/testify/assert"
)

// ProjectAdapter envuelve la dependencia para probar con mocks sin tocar el servicio real.
type ProjectAdapter struct {
	repo ProjectRepoMock
}

func (p *ProjectAdapter) Get(id string) (*project.ProjectResponse, error) {
	return p.repo.FindById(id)
}

func TestProjectAdapter_Get(t *testing.T) {
	okResp := &project.ProjectResponse{ID: "123", Name: "demo"}

	tests := []struct {
		name      string
		id        string
		setup     func(*MockProjectRepo)
		expectNil bool
		expectErr bool
	}{
		{
			name: "success",
			id:   "123",
			setup: func(r *MockProjectRepo) { r.On("FindById", "123").Return(okResp, nil) },
		},
		{
			name:      "not found",
			id:        "404",
			setup:     func(r *MockProjectRepo) { r.On("FindById", "404").Return(nil, nil) },
			expectNil: true,
		},
		{
			name:      "repo error",
			id:        "500",
			setup:     func(r *MockProjectRepo) { r.On("FindById", "500").Return(nil, errors.New("db")) },
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProjectRepo)
			tt.setup(mockRepo)
			svc := ProjectAdapter{repo: mockRepo}

			res, err := svc.Get(tt.id)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else if tt.expectNil {
				assert.NoError(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, okResp.ID, res.ID)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
