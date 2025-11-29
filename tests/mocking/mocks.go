//go:build mocking
// +build mocking

package mocking

import (
	"github.com/MetaDandy/cuent-ai-core/helper"
	"github.com/MetaDandy/cuent-ai-core/src/model"
	"github.com/MetaDandy/cuent-ai-core/src/modules/project"
	"github.com/stretchr/testify/mock"
)

// Repos compartidos para tests de mocking.

type AssetRepoMock interface {
	FindByIdWithGeneratedJobs(id string) (*model.Asset, error)
	FindByScriptID(scriptID string) ([]model.Asset, error)
	FindByScriptIDWithGeneratedJobs(scriptID string) ([]model.Asset, error)
}

type MockAssetRepo struct{ mock.Mock }

func (m *MockAssetRepo) FindByIdWithGeneratedJobs(id string) (*model.Asset, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Asset), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockAssetRepo) FindByScriptID(scriptID string) ([]model.Asset, error) {
	args := m.Called(scriptID)
	if args.Get(0) != nil {
		return args.Get(0).([]model.Asset), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockAssetRepo) FindByScriptIDWithGeneratedJobs(scriptID string) ([]model.Asset, error) {
	args := m.Called(scriptID)
	if args.Get(0) != nil {
		return args.Get(0).([]model.Asset), args.Error(1)
	}
	return nil, args.Error(1)
}

type UserSubRepoMock interface {
	GetActiveSubscription(userID string) (*model.UserSubscribed, error)
}

type MockUserSubRepo struct{ mock.Mock }

func (m *MockUserSubRepo) GetActiveSubscription(userID string) (*model.UserSubscribed, error) {
	args := m.Called(userID)
	if args.Get(0) != nil {
		return args.Get(0).(*model.UserSubscribed), args.Error(1)
	}
	return nil, args.Error(1)
}

type GenJobRepoMock interface {
	FindAll(opts *helper.FindAllOptions) ([]model.GeneratedJob, int64, error)
	FindById(id string) (*model.GeneratedJob, error)
}

type MockGenJobRepo struct{ mock.Mock }

func (m *MockGenJobRepo) FindAll(opts *helper.FindAllOptions) ([]model.GeneratedJob, int64, error) {
	args := m.Called(opts)
	return args.Get(0).([]model.GeneratedJob), args.Get(1).(int64), args.Error(2)
}

func (m *MockGenJobRepo) FindById(id string) (*model.GeneratedJob, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.GeneratedJob), args.Error(1)
	}
	return nil, args.Error(1)
}

type ProjectRepoMock interface {
	FindById(id string) (*project.ProjectResponse, error)
}

type MockProjectRepo struct{ mock.Mock }

func (m *MockProjectRepo) FindById(id string) (*project.ProjectResponse, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*project.ProjectResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

type ScriptRepoMock interface {
	FindByIdWithAssets(id string) (*model.Script, error)
	FindAll(opts *helper.FindAllOptions) ([]model.Script, int64, error)
	FindById(id string) (*model.Script, error)
}

type MockScriptRepo struct{ mock.Mock }

func (m *MockScriptRepo) FindByIdWithAssets(id string) (*model.Script, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Script), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockScriptRepo) FindAll(opts *helper.FindAllOptions) ([]model.Script, int64, error) {
	args := m.Called(opts)
	return args.Get(0).([]model.Script), args.Get(1).(int64), args.Error(2)
}
func (m *MockScriptRepo) FindById(id string) (*model.Script, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Script), args.Error(1)
	}
	return nil, args.Error(1)
}

type UserRepoMock interface {
	FindByEmail(email string) (*model.User, error)
	FindById(id string) (*model.User, error)
	Update(u *model.User) error
}

type MockUserRepo struct{ mock.Mock }

func (m *MockUserRepo) FindByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUserRepo) FindById(id string) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUserRepo) Update(u *model.User) error {
	return m.Called(u).Error(0)
}

type SubRepoMock interface {
	FindAll(opts *helper.FindAllOptions) ([]model.Subscription, int64, error)
	FindById(id string) (*model.Subscription, error)
}

type MockSubRepo struct{ mock.Mock }

func (m *MockSubRepo) FindAll(opts *helper.FindAllOptions) ([]model.Subscription, int64, error) {
	args := m.Called(opts)
	return args.Get(0).([]model.Subscription), args.Get(1).(int64), args.Error(2)
}
func (m *MockSubRepo) FindById(id string) (*model.Subscription, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Subscription), args.Error(1)
	}
	return nil, args.Error(1)
}
