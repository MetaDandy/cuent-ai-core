//go:build mocking
// +build mocking

package mocking

import (
	"errors"
	"testing"
	"time"

	"github.com/MetaDandy/cuent-ai-core/helper"
	"github.com/MetaDandy/cuent-ai-core/src/core/subscription"
	"github.com/MetaDandy/cuent-ai-core/src/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type SubscriptionAdapter struct {
	repo SubRepoMock
}

func (s *SubscriptionAdapter) FindAll(opts *helper.FindAllOptions) (*helper.PaginatedResponse[subscription.SubscriptionResponse], error) {
	items, total, err := s.repo.FindAll(opts)
	if err != nil {
		return nil, err
	}
	dto := subscription.SubscriptionToListDTO(items)
	pages := uint((total + int64(opts.Limit) - 1) / int64(opts.Limit))
	return &helper.PaginatedResponse[subscription.SubscriptionResponse]{
		Data:   dto,
		Total:  total,
		Limit:  opts.Limit,
		Offset: opts.Offset,
		Pages:  pages,
	}, nil
}

func (s *SubscriptionAdapter) FindByID(id string) (*subscription.SubscriptionResponse, error) {
	sub, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return nil, nil
	}
	dto := subscription.SubscriptionToDTO(sub)
	return &dto, nil
}

func TestSubscriptionAdapter_FindAll(t *testing.T) {
	repo := new(MockSubRepo)
	opts := &helper.FindAllOptions{Limit: 2, Offset: 0}

	now := time.Now()
	sub1 := model.Subscription{ID: uuid.New(), Name: "Free", Cuentokens: 50, Duration: now}
	sub2 := model.Subscription{ID: uuid.New(), Name: "Pro", Cuentokens: 200, Duration: now}

	tests := []struct {
		name      string
		setup     func()
		expectErr bool
	}{
		{
			name: "success",
			setup: func() {
				repo.On("FindAll", opts).Return([]model.Subscription{sub1, sub2}, int64(5), nil)
			},
		},
		{
			name: "repo error",
			setup: func() {
				repo.On("FindAll", opts).Return([]model.Subscription{}, int64(0), errors.New("db"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.ExpectedCalls = nil
			repo.Calls = nil
			tt.setup()
			svc := SubscriptionAdapter{repo: repo}
			resp, err := svc.FindAll(opts)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, uint(3), resp.Pages)
				assert.Len(t, resp.Data, 2)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestSubscriptionAdapter_FindByID(t *testing.T) {
	repo := new(MockSubRepo)
	target := &model.Subscription{
		ID:         uuid.New(),
		Name:       "Pro",
		Cuentokens: 120,
		Duration:   time.Now(),
	}

	tests := []struct {
		name      string
		id        string
		setup     func()
		expectNil bool
		expectErr bool
	}{
		{
			name: "success",
			id:   target.ID.String(),
			setup: func() { repo.On("FindById", target.ID.String()).Return(target, nil) },
		},
		{
			name:      "not found",
			id:        "missing",
			setup:     func() { repo.On("FindById", "missing").Return(nil, nil) },
			expectNil: true,
		},
		{
			name:      "error",
			id:        "boom",
			setup:     func() { repo.On("FindById", "boom").Return(nil, errors.New("db")) },
			expectErr: true,
		},
	}

	svc := SubscriptionAdapter{repo: repo}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.ExpectedCalls = nil
			repo.Calls = nil
			tt.setup()
			res, err := svc.FindByID(tt.id)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else if tt.expectNil {
				assert.NoError(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, target.ID.String(), res.ID)
			}
			repo.AssertExpectations(t)
		})
	}
}
