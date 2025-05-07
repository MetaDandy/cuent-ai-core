package user

import "github.com/MetaDandy/cuent-ai-core/helper"

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) FindAll(opts *helper.FindAllOptions) (*helper.PaginatedResponse[UserResponse], error) {
	users, total, err := s.repo.FindAll(opts)
	if err != nil {
		return nil, err
	}
	dtos := UsersToListDTO(users)
	pages := uint((total + int64(opts.Limit) - 1) / int64(opts.Limit))

	return &helper.PaginatedResponse[UserResponse]{
		Data:   dtos,
		Total:  total,
		Limit:  opts.Limit,
		Offset: opts.Offset,
		Pages:  pages,
	}, nil
}
