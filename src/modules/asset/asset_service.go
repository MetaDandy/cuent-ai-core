package asset

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/MetaDandy/cuent-ai-core/helper"
	"github.com/MetaDandy/cuent-ai-core/src/core/user"
	"github.com/MetaDandy/cuent-ai-core/src/model"
	generatejob "github.com/MetaDandy/cuent-ai-core/src/modules/generate_job"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Service struct {
	repo     *Repository
	genRepo  *generatejob.Repository
	userRepo *user.Repository
}

func NewService(r *Repository, gnr *generatejob.Repository, ur *user.Repository) *Service {
	return &Service{repo: r, genRepo: gnr, userRepo: ur}
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

func (s *Service) FindByScriptID(id string) (*[]AssetResponse, error) {
	assets, err := s.repo.FindByScriptID(id)
	if err != nil {
		return nil, err
	}

	dto := AssetsToListDTO(assets)

	return &dto, nil
}

func (s *Service) GenerateOne(id, userID string) (*AssetResponse, error) {
	_, err := s.generate(id, userID)
	if err != nil {
		return nil, err
	}

	reload, _ := s.repo.FindByIdWithGeneratedJobs(id)
	dto := AssetToDto(reload)
	return &dto, nil
}

func (s *Service) GenerateAll(id, userID string, regenerate bool) (*[]AssetResponse, error) {
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
			_, err = s.generate(assetID, userID)
		} else {
			if a.AudioState == model.StatePending || a.AudioState == model.StateError {
				_, err = s.generate(assetID, userID)
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

func (s *Service) generate(id, userID string) (*model.Asset, error) {
	asset, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	sub, err := s.userRepo.GetActiveSubscription(userID)
	if err != nil {
		return nil, err
	}

	var tokens uint
	if strings.HasPrefix(asset.Line, "*") {
		tokens = 40
	} else {
		tokens = uint(utf8.RuneCountInString(asset.Line))
	}

	if sub.TokensRemaining < tokens {
		return nil, fmt.Errorf(
			"fondos insuficientes: se necesitan aprox. %d cuentokens, tienes %d",
			tokens, sub.TokensRemaining,
		)
	}

	bucket := "audio"
	dirPath := filepath.Join(asset.ScriptID.String())

	if err := s.repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", sub.ID).
			Take(&sub).Error; err != nil {
			return err
		}

		url, historyID, _, duration, err := helper.AudioOutput(asset.Line, asset.ID.String(), bucket, dirPath)
		if err != nil {
			return err
		}

		chars, err := helper.CharactersUsed(historyID)
		if err != nil {
			chars = 0
		}

		asset.Audio_URL = url
		asset.AudioState = model.StateFinished
		asset.Duration = duration.Seconds()
		if err := tx.Save(asset).Error; err != nil {
			return err
		}

		job := model.GeneratedJob{
			ID:              uuid.New(),
			Provider:        model.ProviderElevenlab,
			Model:           "eleven_monolingual_v1",
			Chars_Used:      uint(chars),
			Cuentoken_Spent: tokens,
			Token_Spent:     strconv.Itoa(chars),
			State:           model.StateFinished,
			Cost:            float64(chars) / 1000 * 0.30, // tarifa por 1 K caracteres
			AssetID:         asset.ID,
		}
		if err := tx.Create(&job).Error; err != nil {
			return err
		}

		sub.TokensRemaining -= tokens
		if err := tx.Save(sub).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		// ! Ver si es factible cobrar la mitad si ocurre un error
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

func (s *Service) GenerateVideo(id, userID string, key_words GenerateVideo) (*model.Asset, error) {
	asset, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	if asset.AudioState == model.StateError || asset.AudioState == model.StatePending {
		return nil, errors.New("para generar un video, primero debe haber generado el audio")
	}

	sub, err := s.userRepo.GetActiveSubscription(userID)
	if err != nil {
		return nil, err
	}

	tokens := uint(asset.Duration) * 50

	if sub.TokensRemaining < tokens {
		return nil, fmt.Errorf(
			"fondos insuficientes: se necesitan aprox. %d cuentokens, tienes %d",
			tokens, sub.TokensRemaining,
		)
	}

	bucket := "video"
	dirPath := filepath.Join(asset.ScriptID.String())

	if err := s.repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", sub.ID).
			Take(&sub).Error; err != nil {
			return err
		}

		images, err := helper.SearchImage(key_words.KeyWords)
		if err != nil {
			return err
		}

		rawVideo, err := helper.GenerateVideo(images, asset.Audio_URL, asset.Duration)
		if err != nil {
			return err
		}

		video := bytes.NewReader(rawVideo)
		fileName := asset.ID.String() + ".mp4"

		url, err := helper.Upload(context.TODO(), bucket, dirPath, fileName, video, "video/mp4", false)
		if err != nil {
			return err
		}

		asset.Video_URL = url
		asset.VideoState = model.StateFinished
		if err := tx.Save(asset).Error; err != nil {
			return err
		}

		job := model.GeneratedJob{
			ID:              uuid.New(),
			Provider:        model.ProviderGemini,
			Model:           "veo_2",
			Cuentoken_Spent: tokens,
			State:           model.StateFinished,
			AssetID:         asset.ID,
		}
		if err := tx.Create(&job).Error; err != nil {
			return err
		}

		sub.TokensRemaining -= tokens
		if err := tx.Save(sub).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		asset.VideoState = model.StateError
		asset.Video_URL = ""
		badJob := model.GeneratedJob{
			ID:            uuid.New(),
			Error_Message: err.Error(),
			AssetID:       asset.ID,
			State:         model.StateError,
		}

		if e := s.repo.Update(asset); e != nil {
			err = errors.Join(err, e)
		}
		if e := s.genRepo.Create(&badJob); e != nil {
			err = errors.Join(err, e)
		}
		return nil, err
	}

	return asset, nil
}
