package asset

import (
	"errors"
	"path/filepath"
	"strconv"

	"github.com/MetaDandy/cuent-ai-core/helper"
	"github.com/MetaDandy/cuent-ai-core/src/model"
	generatejob "github.com/MetaDandy/cuent-ai-core/src/modules/generate_job"
	"github.com/google/uuid"
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

func (s *Service) GenerateAll(id string, regenerate bool) (*[]AssetResponse, error) {
	assets, err := s.repo.FindByScriptID(id)
	if err != nil {
		return nil, err
	}
	if len(assets) == 0 {
		empty := make([]AssetResponse, 0)
		return &empty, nil
	}

	var errs []error

	// * No usar go rutine por el tema de los rate limits estrictos de eleven labs
	for _, a := range assets {
		assetID := a.ID.String()
		var err error
		if regenerate {
			_, err = s.generate(assetID)
		} else {
			if a.AudioState == model.StatePending || a.AudioState == model.StateError {
				_, err = s.generate(assetID)
			}
		}
		if err != nil {
			errs = append(errs, err)
		}
	}

	reloaded, err := s.repo.FindByScriptIDWithGeneratedJobs(id)
	if err != nil {
		return nil, err
	}
	dto := AssetsToListDTO(reloaded)

	if len(errs) > 0 {
		return &dto, errors.Join(errs...)
	}

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
			ID:            uuid.New(),
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
