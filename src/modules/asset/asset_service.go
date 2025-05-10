package asset

import (
	"context"
	"errors"
	"path/filepath"
	"strconv"

	"github.com/MetaDandy/cuent-ai-core/helper"
	"github.com/MetaDandy/cuent-ai-core/src/model"
	generatejob "github.com/MetaDandy/cuent-ai-core/src/modules/generate_job"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
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
	finded, err := s.repo.FindByIdWithGeneratedJobs(id)
	if err != nil {
		return nil, err
	}
	if finded == nil {
		return nil, nil
	}
	dto := AssetToDto(finded)
	return &dto, nil
}

func (s *Service) GenerateOne(id string) (*AssetResponse, error) {
	_, err := s.generate(id)
	if err != nil {
		return nil, err
	}

	reload, _ := s.repo.FindByIdWithGeneratedJobs(id)
	dto := AssetToDto(reload)
	return &dto, nil
}

func (s *Service) GenerateAll(id string) (*[]AssetResponse, error) {
	assets, err := s.repo.FindByScriptID(id)
	if err != nil {
		return nil, err
	}
	if len(assets) == 0 {
		empty := make([]AssetResponse, 0)
		return &empty, nil
	}

	g, ctx := errgroup.WithContext(context.Background())

	for _, a := range assets {
		a := a // captura local requerida
		g.Go(func() error {
			// Abortamos si el contexto se canceló (otro asset falló)
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				_, err := s.generate(a.ID.String())
				return err // si falla, errgroup cancelará a los demás
			}
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err // el primer error aborta el proceso completo
	}

	reloaded, err := s.repo.FindByScriptIDWithGeneratedJobs(id)
	if err != nil {
		return nil, err
	}
	dto := AssetsToListDTO(reloaded)

	return &dto, nil
}

func (s *Service) generate(id string) (*model.Asset, error) {
	asset, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	bucket := "audio"
	dirPath := filepath.Join(asset.ScriptID.String(), asset.Script.UpdatedAt.String())

	if err := s.repo.db.Transaction(func(tx *gorm.DB) error {
		url, historyID, duration, err := helper.AudioOutput(asset.Line, asset.ID.String(), bucket, dirPath)
		if err != nil {
			return err
		}

		chars, err := helper.CharactersUsed(historyID)
		if err != nil {
			chars = 0
		}

		asset.Audio_URL = url
		asset.AudioState = model.StateFinished
		asset.Duration = uint(duration.Seconds())
		if err := tx.Save(asset).Error; err != nil {
			return err
		}

		job := model.GeneratedJob{
			ID:          uuid.New(),
			Provider:    model.ProviderElevenlab,
			Model:       "eleven_monolingual_v1",
			Chars_Used:  uint(chars),
			Token_Spent: strconv.Itoa(chars),
			State:       model.StateFinished,
			Cost:        float64(chars) / 1000 * 0.30, // tarifa por 1 K caracteres
			AssetID:     asset.ID,
		}
		if err := tx.Create(&job).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		asset.AudioState = model.StateError
		asset.Audio_URL = ""
		asset.Duration = 0
		badJob := model.GeneratedJob{
			Error_Message: err.Error(),
			AssetID:       asset.ID,
			State:         model.StateError,
		}

		if e := s.repo.Update(asset); e != nil {
			err = errors.Join(err, e) // Go 1.20+
		}
		if e := s.genRepo.Create(&badJob); e != nil {
			err = errors.Join(err, e)
		}
		return nil, err
	}

	return asset, nil
}
