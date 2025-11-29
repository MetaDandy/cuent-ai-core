//go:build mocking
// +build mocking

package mocking

import (
	"errors"
	"testing"

	"github.com/MetaDandy/cuent-ai-core/helper"
	"github.com/MetaDandy/cuent-ai-core/src/core/user"
	"github.com/MetaDandy/cuent-ai-core/src/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// UserServiceAdapter implementa las rutas de negocio necesarias para testear con mocks locales.
type UserServiceAdapter struct {
	repo UserRepoMock
}

func (s *UserServiceAdapter) Signin(in *user.Signin) (*user.UserResponse, string, error) {
	u, err := s.repo.FindByEmail(in.Email)
	if err != nil {
		return nil, "", err
	}
	if !helper.CheckPasswordHash(in.Password, u.Password) {
		return nil, "", errors.New("la contraseña no coincide")
	}
	token, err := helper.GenerateJwt(u.ID.String(), u.Email)
	if err != nil {
		return nil, "", err
	}
	dto := user.UserToDTO(u)
	return &dto, token, nil
}

func (s *UserServiceAdapter) ChangePassword(id string, in *user.ChangePassoword) (*user.UserResponse, error) {
	u, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	if !helper.CheckPasswordHash(in.Old_Password, u.Password) {
		return nil, errors.New("la contraseña no coincide")
	}
	if in.New_Password != in.Confirm_Password {
		return nil, errors.New("la nueva contraseña no coincide con la confirmación")
	}
	hash, err := helper.HashPassword(in.New_Password)
	if err != nil {
		return nil, err
	}
	u.Password = hash
	if err := s.repo.Update(u); err != nil {
		return nil, err
	}
	dto := user.UserToDTO(u)
	return &dto, nil
}

func TestUserServiceAdapter_Signin(t *testing.T) {
	t.Setenv("JWT_SECRET", "secret")
	hashed, _ := helper.HashPassword("pass123")
	validUser := &model.User{ID: uuid.New(), Email: "a@test.com", Password: hashed}

	tests := []struct {
		name       string
		input      user.Signin
		setup      func(m *MockUserRepo)
		expectErr  bool
		expectUser bool
	}{
		{
			name:  "success",
			input: user.Signin{Email: validUser.Email, Password: "pass123"},
			setup: func(m *MockUserRepo) {
				m.On("FindByEmail", validUser.Email).Return(validUser, nil)
			},
			expectUser: true,
		},
		{
			name:  "wrong password",
			input: user.Signin{Email: validUser.Email, Password: "bad"},
			setup: func(m *MockUserRepo) {
				m.On("FindByEmail", validUser.Email).Return(validUser, nil)
			},
			expectErr: true,
		},
		{
			name:  "repo error",
			input: user.Signin{Email: "x@test.com", Password: "any"},
			setup: func(m *MockUserRepo) {
				m.On("FindByEmail", "x@test.com").Return(nil, errors.New("db"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			if tt.setup != nil {
				tt.setup(mockRepo)
			}
			svc := UserServiceAdapter{repo: mockRepo}
			res, token, err := svc.Signin(&tt.input)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, res)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.NotEmpty(t, token)
				assert.Equal(t, tt.input.Email, res.Email)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserServiceAdapter_ChangePassword(t *testing.T) {
	t.Setenv("JWT_SECRET", "secret")
	oldHash, _ := helper.HashPassword("old")
	userID := uuid.New().String()

	tests := []struct {
		name      string
		input     user.ChangePassoword
		setup     func(m *MockUserRepo, u *model.User)
		expectErr bool
	}{
		{
			name: "success",
			input: user.ChangePassoword{
				Old_Password:     "old",
				New_Password:     "new12345",
				Confirm_Password: "new12345",
			},
			setup: func(m *MockUserRepo, u *model.User) {
				m.On("FindById", userID).Return(u, nil)
				m.On("Update", mock.AnythingOfType("*model.User")).Return(nil)
			},
		},
		{
			name: "old mismatch",
			input: user.ChangePassoword{
				Old_Password:     "bad",
				New_Password:     "new12345",
				Confirm_Password: "new12345",
			},
			setup: func(m *MockUserRepo, u *model.User) {
				m.On("FindById", userID).Return(u, nil)
			},
			expectErr: true,
		},
		{
			name: "confirm mismatch",
			input: user.ChangePassoword{
				Old_Password:     "old",
				New_Password:     "new12345",
				Confirm_Password: "zzz",
			},
			setup: func(m *MockUserRepo, u *model.User) {
				m.On("FindById", userID).Return(u, nil)
			},
			expectErr: true,
		},
		{
			name: "update error",
			input: user.ChangePassoword{
				Old_Password:     "old",
				New_Password:     "new12345",
				Confirm_Password: "new12345",
			},
			setup: func(m *MockUserRepo, u *model.User) {
				m.On("FindById", userID).Return(u, nil)
				m.On("Update", mock.AnythingOfType("*model.User")).Return(errors.New("save"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			mockUser := &model.User{ID: uuid.MustParse(userID), Password: oldHash}
			if tt.setup != nil {
				tt.setup(mockRepo, mockUser)
			}

			svc := UserServiceAdapter{repo: mockRepo}
			res, err := svc.ChangePassword(userID, &tt.input)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.True(t, helper.CheckPasswordHash(tt.input.New_Password, mockUser.Password))
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
