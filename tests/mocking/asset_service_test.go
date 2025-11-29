//go:build mocking
// +build mocking

package mocking

import (
	"errors"
	"testing"

	"github.com/MetaDandy/cuent-ai-core/src/model"
	"github.com/MetaDandy/cuent-ai-core/src/modules/asset"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// AssetAdapter reproduce la lógica mínima que validan los tests usando mocks locales.
type AssetAdapter struct {
	repo    AssetRepoMock
	user    UserSubRepoMock
}

func (a *AssetAdapter) FindByID(id string) (*asset.AssetResponse, error) {
	ast, err := a.repo.FindByIdWithGeneratedJobs(id)
	if err != nil {
		return nil, err
	}
	if ast == nil {
		return nil, nil
	}
	dto := asset.AssetToDto(ast)
	return &dto, nil
}

func (a *AssetAdapter) GenerateAll(scriptID, userID string) (*[]asset.AssetResponse, error) {
	assetsList, err := a.repo.FindByScriptID(scriptID)
	if err != nil {
		return nil, err
	}
	if len(assetsList) == 0 {
		empty := make([]asset.AssetResponse, 0)
		return &empty, nil
	}
	sub, err := a.user.GetActiveSubscription(userID)
	if err != nil {
		return nil, err
	}
	if sub.TokensRemaining < uint(len(assetsList)) {
		return nil, errors.New("fondos insuficientes")
	}
	reloaded, err := a.repo.FindByScriptIDWithGeneratedJobs(scriptID)
	if err != nil {
		return nil, err
	}
	dto := asset.AssetsToListDTO(reloaded)
	return &dto, nil
}

func TestAssetAdapter_FindByID(t *testing.T) {
	ast := &model.Asset{ID: uuid.New(), Line: "hola"}

	tests := []struct {
		name      string
		id        string
		setup     func(*MockAssetRepo)
		expectNil bool
		expectErr bool
	}{
		{
			name: "success",
			id:   "ok",
			setup: func(r *MockAssetRepo) {
				r.On("FindByIdWithGeneratedJobs", "ok").Return(ast, nil)
			},
		},
		{
			name:      "not found",
			id:        "missing",
			setup:     func(r *MockAssetRepo) { r.On("FindByIdWithGeneratedJobs", "missing").Return(nil, nil) },
			expectNil: true,
		},
		{
			name:      "repo error",
			id:        "boom",
			setup:     func(r *MockAssetRepo) { r.On("FindByIdWithGeneratedJobs", "boom").Return(nil, errors.New("db")) },
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAssetRepo)
			if tt.setup != nil {
				tt.setup(mockRepo)
			}
			svc := AssetAdapter{repo: mockRepo, user: new(MockUserSubRepo)}
			res, err := svc.FindByID(tt.id)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else if tt.expectNil {
				assert.NoError(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, ast.Line, res.Line)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAssetAdapter_GenerateAll(t *testing.T) {
	assetsList := []model.Asset{{ID: uuid.New()}, {ID: uuid.New()}}
	subEnough := &model.UserSubscribed{TokensRemaining: 5}

	tests := []struct {
		name      string
		setup     func(*MockAssetRepo, *MockUserSubRepo)
		expectErr bool
		expectNil bool
	}{
		{
			name: "success",
			setup: func(r *MockAssetRepo, u *MockUserSubRepo) {
				r.On("FindByScriptID", "sid").Return(assetsList, nil)
				u.On("GetActiveSubscription", "uid").Return(subEnough, nil)
				r.On("FindByScriptIDWithGeneratedJobs", "sid").Return(assetsList, nil)
			},
		},
		{
			name: "repo error find",
			setup: func(r *MockAssetRepo, u *MockUserSubRepo) {
				r.On("FindByScriptID", "sid").Return([]model.Asset(nil), errors.New("db"))
			},
			expectErr: true,
		},
		{
			name: "empty assets",
			setup: func(r *MockAssetRepo, u *MockUserSubRepo) {
				r.On("FindByScriptID", "sid").Return([]model.Asset{}, nil)
			},
			expectNil: true,
		},
		{
			name: "insufficient tokens",
			setup: func(r *MockAssetRepo, u *MockUserSubRepo) {
				r.On("FindByScriptID", "sid").Return(assetsList, nil)
				u.On("GetActiveSubscription", "uid").Return(&model.UserSubscribed{TokensRemaining: 1}, nil)
			},
			expectErr: true,
		},
		{
			name: "repo error preload",
			setup: func(r *MockAssetRepo, u *MockUserSubRepo) {
				r.On("FindByScriptID", "sid").Return(assetsList, nil)
				u.On("GetActiveSubscription", "uid").Return(subEnough, nil)
				r.On("FindByScriptIDWithGeneratedJobs", "sid").Return([]model.Asset(nil), errors.New("db2"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAssetRepo)
			mockUser := new(MockUserSubRepo)
			if tt.setup != nil {
				tt.setup(mockRepo, mockUser)
			}
			svc := AssetAdapter{repo: mockRepo, user: mockUser}
			res, err := svc.GenerateAll("sid", "uid")
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else if tt.expectNil {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Len(t, *res, 0)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}
			mockRepo.AssertExpectations(t)
			mockUser.AssertExpectations(t)
		})
	}
}
