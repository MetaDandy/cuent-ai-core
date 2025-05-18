package subscription

import "github.com/MetaDandy/cuent-ai-core/helper"

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) FindAll(opts *helper.FindAllOptions) (*helper.PaginatedResponse[SubscriptionResponse], error) {
	projects, total, err := s.repo.FindAll(opts)
	if err != nil {
		return nil, err
	}
	dtos := SubscriptionToListDTO(projects)
	pages := uint((total + int64(opts.Limit) - 1) / int64(opts.Limit))

	return &helper.PaginatedResponse[SubscriptionResponse]{
		Data:   dtos,
		Total:  total,
		Limit:  opts.Limit,
		Offset: opts.Offset,
		Pages:  pages,
	}, nil
}

func (s *Service) FindByID(id string) (*SubscriptionResponse, error) {
	project, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, nil
	}
	dto := SubscriptionToDTO(project)
	return &dto, nil
}
