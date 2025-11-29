//go:build mocking
// +build mocking

package mocking

import (
	"errors"
	"testing"

	"github.com/MetaDandy/cuent-ai-core/helper"
	"github.com/MetaDandy/cuent-ai-core/src/model"
	generatejob "github.com/MetaDandy/cuent-ai-core/src/modules/generate_job"
	"github.com/stretchr/testify/assert"
)

type GenJobAdapter struct {
	repo GenJobRepoMock
}

func (s *GenJobAdapter) List(opts *helper.FindAllOptions) ([]generatejob.GeneratedJobResponse, error) {
	items, _, err := s.repo.FindAll(opts)
	if err != nil {
		return nil, err
	}
	return generatejob.GeneratedJobToLisDTO(items), nil
}

func (s *GenJobAdapter) Get(id string) (*generatejob.GeneratedJobResponse, error) {
	item, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}
	dto := generatejob.GeneratedJobToDto(item)
	return &dto, nil
}

func TestGenJobAdapter_List(t *testing.T) {
	opts := &helper.FindAllOptions{Limit: 5}
	jobs := []model.GeneratedJob{{Model: "m1"}, {Model: "m2"}}

	tests := []struct {
		name      string
		setup     func(m *MockGenJobRepo)
		expectErr bool
	}{
		{
			name: "success",
			setup: func(m *MockGenJobRepo) {
				m.On("FindAll", opts).Return(jobs, int64(2), nil)
			},
		},
		{
			name: "repo error",
			setup: func(m *MockGenJobRepo) {
				m.On("FindAll", opts).Return([]model.GeneratedJob{}, int64(0), errors.New("db"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockGenJobRepo)
			tt.setup(mockRepo)
			svc := GenJobAdapter{repo: mockRepo}
			resp, err := svc.List(opts)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Len(t, resp, 2)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGenJobAdapter_Get(t *testing.T) {
	okJob := &model.GeneratedJob{Model: "demo"}

	tests := []struct {
		name      string
		id        string
		setup     func(m *MockGenJobRepo)
		expectNil bool
		expectErr bool
	}{
		{
			name: "success",
			id:   "ok",
			setup: func(m *MockGenJobRepo) { m.On("FindById", "ok").Return(okJob, nil) },
		},
		{
			name:      "not found",
			id:        "missing",
			setup:     func(m *MockGenJobRepo) { m.On("FindById", "missing").Return(nil, nil) },
			expectNil: true,
		},
		{
			name:      "repo error",
			id:        "bad",
			setup:     func(m *MockGenJobRepo) { m.On("FindById", "bad").Return(nil, errors.New("db")) },
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockGenJobRepo)
			tt.setup(mockRepo)
			svc := GenJobAdapter{repo: mockRepo}
			res, err := svc.Get(tt.id)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else if tt.expectNil {
				assert.NoError(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "demo", res.Model)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
