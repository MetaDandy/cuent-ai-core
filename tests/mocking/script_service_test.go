//go:build mocking
// +build mocking

package mocking

import (
	"errors"
	"testing"

	"github.com/MetaDandy/cuent-ai-core/src/model"
	"github.com/MetaDandy/cuent-ai-core/src/modules/script"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type ScriptAdapter struct {
	repo  ScriptRepoMock
	asset AssetRepoMock
	user  UserSubRepoMock
}

func (s *ScriptAdapter) FindByID(id string) (*script.ScriptReponse, error) {
	sc, err := s.repo.FindByIdWithAssets(id)
	if err != nil {
		return nil, err
	}
	if sc == nil {
		return nil, errors.New("no se encontró lo solicitado")
	}
	dto := script.ScriptToDTO(sc)
	return &dto, nil
}

func (s *ScriptAdapter) MixAudio(id, userID string) (*script.ScriptReponse, error) {
	sc, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	assets, err := s.asset.FindByScriptID(id)
	if err != nil {
		return nil, err
	}
	sub, err := s.user.GetActiveSubscription(userID)
	if err != nil {
		return nil, err
	}
	if sub.TokensRemaining < uint(len(assets)) {
		return nil, errors.New("fondos insuficientes")
	}
	dto := script.ScriptToDTO(sc)
	return &dto, nil
}

func TestScriptAdapter_FindByID(t *testing.T) {
	sc := &model.Script{ID: uuid.New(), Text_Entry: "hola"}

	tests := []struct {
		name      string
		id        string
		setup     func(*MockScriptRepo)
		expectErr bool
	}{
		{
			name: "success",
			id:   "ok",
			setup: func(r *MockScriptRepo) { r.On("FindByIdWithAssets", "ok").Return(sc, nil) },
		},
		{
			name:      "not found",
			id:        "missing",
			setup:     func(r *MockScriptRepo) { r.On("FindByIdWithAssets", "missing").Return(nil, nil) },
			expectErr: true,
		},
		{
			name:      "repo error",
			id:        "bad",
			setup:     func(r *MockScriptRepo) { r.On("FindByIdWithAssets", "bad").Return(nil, errors.New("db")) },
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockScriptRepo)
			if tt.setup != nil {
				tt.setup(mockRepo)
			}
			svc := ScriptAdapter{repo: mockRepo, asset: new(MockAssetRepo), user: new(MockUserSubRepo)}
			res, err := svc.FindByID(tt.id)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, sc.Text_Entry, res.Text_Entry)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestScriptAdapter_MixAudio(t *testing.T) {
	mockRepo := new(MockScriptRepo)
	mockAsset := new(MockAssetRepo)
	mockUser := new(MockUserSubRepo)

	mockRepo.On("FindById", "ok").Return(&model.Script{ID: uuid.New()}, nil)
	mockAsset.On("FindByScriptID", "ok").Return([]model.Asset{{}, {}}, nil)
	mockUser.On("GetActiveSubscription", "uid").Return(&model.UserSubscribed{TokensRemaining: 5}, nil)

	svc := ScriptAdapter{repo: mockRepo, asset: mockAsset, user: mockUser}

	resp, err := svc.MixAudio("ok", "uid")
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	mockRepo.AssertExpectations(t)
	mockAsset.AssertExpectations(t)
	mockUser.AssertExpectations(t)
}

func TestScriptAdapter_MixAudio_Errors(t *testing.T) {
	mockRepo := new(MockScriptRepo)
	mockAsset := new(MockAssetRepo)
	mockUser := new(MockUserSubRepo)

	mockRepo.On("FindById", "bad").Return(nil, errors.New("db"))
	mockRepo.On("FindById", "insufficient").Return(&model.Script{}, nil)
	mockAsset.On("FindByScriptID", "insufficient").Return([]model.Asset{{}, {}}, nil)
	mockUser.On("GetActiveSubscription", "uid").Return(&model.UserSubscribed{TokensRemaining: 1}, nil)
	mockAsset.On("FindByScriptID", "badAssets").Return(nil, errors.New("db2"))
	mockRepo.On("FindById", "badAssets").Return(&model.Script{}, nil)

	svc := ScriptAdapter{repo: mockRepo, asset: mockAsset, user: mockUser}

	_, err := svc.MixAudio("bad", "uid")
	assert.Error(t, err)

	_, err = svc.MixAudio("badAssets", "uid")
	assert.Error(t, err)

	_, err = svc.MixAudio("insufficient", "uid")
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
	mockAsset.AssertExpectations(t)
	mockUser.AssertExpectations(t)
}
