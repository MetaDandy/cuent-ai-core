package project

import (
	"github.com/MetaDandy/cuent-ai-core/helper"
	"github.com/MetaDandy/cuent-ai-core/src/core/user"
	"github.com/MetaDandy/cuent-ai-core/src/model"
	"gorm.io/gorm"
)

type Service struct {
	repo     *Repository
	userRepo *user.Repository
}

func NewService(r *Repository, u *user.Repository) *Service {
	return &Service{repo: r, userRepo: u}
}

func (s *Service) FindAll(opts *helper.FindAllOptions) (*helper.PaginatedResponse[ProjectResponse], error) {
	projects, total, err := s.repo.FindAll(opts)
	if err != nil {
		return nil, err
	}
	dtos := ProjectsToListDTO(projects)
	pages := uint((total + int64(opts.Limit) - 1) / int64(opts.Limit))

	return &helper.PaginatedResponse[ProjectResponse]{
		Data:   dtos,
		Total:  total,
		Limit:  opts.Limit,
		Offset: opts.Offset,
		Pages:  pages,
	}, nil
}

func (s *Service) FindByID(id string) (*ProjectResponse, error) {
	project, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, nil
	}
	dto := ProjectToDTO(project)
	return &dto, nil
}

func (s *Service) Create(input *ProjectCreate) (*ProjectResponse, error) {
	user, err := s.userRepo.FindById(input.UserId)
	if err != nil {
		return nil, err
	}

	project := model.Project{
		Name:        input.Name,
		Description: input.Description,
		State:       model.StatePending,
		UserID:      user.ID,
	}

	if err := s.repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&project).Error; err != nil {
			return err
		}
		// ... aquí podrías crear scripts, logs, etc. y devolver err si algo falla
		return nil
	}); err != nil {
		return nil, err
	}

	reload, _ := s.repo.FindById(project.ID.String())
	dto := ProjectToDTO(reload)
	return &dto, nil
}

func (s *Service) Update(id string, input *ProjectUpdate) (*ProjectResponse, error) {
	project, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, nil
	}

	if input.Name != nil {
		project.Name = *input.Name
	}
	if input.Description != nil {
		project.Description = *input.Description
	}

	if err := s.repo.Update(project); err != nil {
		return nil, err
	}

	reloaded, _ := s.repo.FindById(id)
	dto := ProjectToDTO(reloaded)
	return &dto, nil
}

func (s *Service) SoftDelete(id string) (bool, error) {
	project, err := s.repo.FindById(id)
	if err != nil {
		return false, err
	}
	if project == nil {
		return false, nil
	}

	if err := s.repo.SoftDelete(id); err != nil {
		return false, nil
	}

	return true, nil
}

func (s *Service) Restore(id string) (*ProjectResponse, error) {
	project, err := s.repo.FindByIdUnscoped(id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, nil
	}

	if err := s.repo.Restore(id); err != nil {
		return nil, err
	}
	dto := ProjectToDTO(project)
	return &dto, nil
}
