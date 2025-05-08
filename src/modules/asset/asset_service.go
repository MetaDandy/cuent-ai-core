package asset

import (
	"github.com/MetaDandy/cuent-ai-core/helper"
	generatejob "github.com/MetaDandy/cuent-ai-core/src/modules/generate_job"
)

type Service struct {
	repo    *Repository
	genRepo *generatejob.Repository
}

func NewService(r *Repository, gnr *generatejob.Repository) *Service {
	return &Service{repo: r, genRepo: gnr}
}

func (s *Service) FindAll(opts *helper.FindAllOptions) (*helper.PaginatedResponse[AssetResponse], error) {
	finded, total, err := s.repo.FindAll(opts)
	if err != nil {
		return nil, err
	}
	dtos := AssetsToListDTO(finded)
	pages := uint((total + int64(opts.Limit) - 1) / int64(opts.Limit))

	return &helper.PaginatedResponse[AssetResponse]{
		Data:   dtos,
		Total:  total,
		Limit:  opts.Limit,
		Offset: opts.Offset,
		Pages:  pages,
	}, nil
}

func (s *Service) FindByID(id string) (*AssetResponse, error) {
	finded, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	if finded == nil {
		return nil, nil
	}
	dto := AssetToDto(finded)
	return &dto, nil
}

func (s *Service) Generate(id string) (*AssetResponse, error) {
	asset, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	// ! Lógica de la generación

	dto := AssetToDto(asset)
	return &dto, nil
}
