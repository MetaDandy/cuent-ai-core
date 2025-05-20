package script

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"

	"github.com/MetaDandy/cuent-ai-core/helper"
	"github.com/MetaDandy/cuent-ai-core/src/core/user"
	"github.com/MetaDandy/cuent-ai-core/src/model"
	"github.com/MetaDandy/cuent-ai-core/src/modules/asset"
	"github.com/MetaDandy/cuent-ai-core/src/modules/project"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Service struct {
	repo        *Repository
	projectRepo *project.Repository
	assetRepo   *asset.Repository
	userRepo    *user.Repository
}

func NewService(r *Repository, pr *project.Repository, ar *asset.Repository, ur *user.Repository) *Service {
	return &Service{repo: r, projectRepo: pr, assetRepo: ar, userRepo: ur}
}

/**
TODO:
- para los mixed, crear endpoints especiales
*/

func (s *Service) Create(userID string, input *ScriptCreate) (*ScriptReponse, error) {
	project, err := s.projectRepo.FindById(input.ProjectID)
	if err != nil {
		return nil, err
	}

	var script model.Script
	if err := s.repo.db.Transaction(func(tx *gorm.DB) error {
		sub, err := s.userRepo.GetActiveSubscription(userID)
		if err != nil {
			return err
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", sub.ID).
			Take(&sub).Error; err != nil {
			return err
		}

		needed := helper.EstimateCuentokens(input.TextEntry)
		if sub.TokensRemaining < needed {
			return fmt.Errorf(
				"fondos insuficientes: se necesitan aprox. %d cuentokens, tienes %d",
				needed, sub.TokensRemaining,
			)
		}

		aiResponse, err := helper.AIFormatter(input.TextEntry)
		if err != nil {
			return err
		}

		// Cobrar 1 cuentoken por linea
		lines := len(aiResponse.Processed_Text_Array)

		script = model.Script{
			ID:                uuid.New(),
			Text_Entry:        input.TextEntry,
			ProjectID:         project.ID,
			Prompt_Tokens:     aiResponse.Prompt_Tokens,
			Completion_Tokens: aiResponse.Completion_Tokens,
			Total_Tokens:      aiResponse.Total_Tokens,
			Processed_Text:    aiResponse.Processed_Text,
			State:             model.StateFinished,
			Total_Cuentoken:   uint(lines),
		}
		if err := tx.Create(&script).Error; err != nil {
			return err
		}

		if sub.TokensRemaining < uint(lines) {
			sub.TokensRemaining = 0
		} else {
			sub.TokensRemaining = sub.TokensRemaining - uint(lines)
		}

		if err := tx.Save(sub).Error; err != nil {
			return err
		}

		assets := make([]model.Asset, 0, lines)
		for i, line := range aiResponse.Processed_Text_Array {
			assets = append(assets, model.Asset{
				ID:       uuid.New(),
				Type:     "LINE", // ! Hacer que sea un enum entre tts y sfx
				Line:     line,
				ScriptID: script.ID,
				Position: i,
			})
		}
		if err := tx.Create(&assets).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	reload, _ := s.repo.FindByIdWithAssets(script.ID.String())
	dto := ScriptToDTO(reload)
	return &dto, nil
}

func (s *Service) Regenerate(userID, scriptID string) (*ScriptReponse, error) {
	script, err := s.repo.FindById(scriptID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.db.Transaction(func(tx *gorm.DB) error {
		sub, err := s.userRepo.GetActiveSubscription(userID)
		if err != nil {
			return err
		}
		tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&sub, "id = ?", sub.ID)

		needed := 2 * helper.EstimateCuentokens(script.Text_Entry)
		if sub.TokensRemaining < needed {
			return fmt.Errorf(
				"fondos insuficientes: se necesitan aprox. %d cuentokens, tienes %d",
				needed, sub.TokensRemaining,
			)
		}

		aiResponse, err := helper.AIFormatter(script.Text_Entry)
		if err != nil {
			return err
		}

		lines := 2 * len(aiResponse.Processed_Text_Array)

		script.Prompt_Tokens = aiResponse.Prompt_Tokens
		script.Completion_Tokens = aiResponse.Completion_Tokens
		script.Total_Tokens = aiResponse.Total_Tokens
		script.Processed_Text = aiResponse.Processed_Text
		script.Total_Cuentoken += uint(lines)
		if err := tx.Save(&script).Error; err != nil {
			return err
		}

		if sub.TokensRemaining < uint(lines) {
			sub.TokensRemaining = 0
		} else {
			sub.TokensRemaining = sub.TokensRemaining - uint(lines)
		}

		if err := tx.Save(&sub).Error; err != nil {
			return err
		}

		if err := tx.
			Where("script_id = ?", script.ID).
			Delete(&model.Asset{}).Error; err != nil {
			return err
		}

		dirPath := filepath.Join(script.ID.String())
		if err := helper.DeleteFolder(context.TODO(), "audio", dirPath); err != nil {
			log.Printf("error borrando carpeta Supabase: %v", err)
		}

		assets := make([]model.Asset, 0, lines)
		for i, line := range aiResponse.Processed_Text_Array {
			assets = append(assets, model.Asset{
				ID:       uuid.New(),
				Type:     "LINE", // ! Hacer que sea un enum entre tts y sfx
				Line:     line,
				ScriptID: script.ID,
				Position: i,
			})
		}
		if err := tx.Create(&assets).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	reload, _ := s.repo.FindByIdWithAssets(script.ID.String())
	dto := ScriptToDTO(reload)
	return &dto, nil
}

func (s *Service) FindAll(opts *helper.FindAllOptions) (*helper.PaginatedResponse[ScriptReponse], error) {
	finded, total, err := s.repo.FindAll(opts)
	if err != nil {
		return nil, err
	}
	dtos := ScriptToListDTO(finded)
	pages := uint((total + int64(opts.Limit) - 1) / int64(opts.Limit))

	return &helper.PaginatedResponse[ScriptReponse]{
		Data:   dtos,
		Total:  total,
		Limit:  opts.Limit,
		Offset: opts.Offset,
		Pages:  pages,
	}, nil
}

func (s *Service) FindByID(id string) (*ScriptReponse, error) {
	finded, err := s.repo.FindByIdWithAssets(id)
	if err != nil {
		return nil, err
	}
	if finded == nil {
		return nil, errors.New("no se encontró lo solicitado")
	}
	dto := ScriptToDTO(finded)
	return &dto, nil
}

func (s *Service) MixAudio(id, userID string) (*ScriptReponse, error) {
	script, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	assets, err := s.repo.FindByIDWithAssetsPosition(id)
	if err != nil {
		return nil, err
	}

	if err := s.repo.db.Transaction(func(tx *gorm.DB) error {
		sub, err := s.userRepo.GetActiveSubscription(userID)
		if err != nil {
			return err
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", sub.ID).
			Take(&sub).Error; err != nil {
			return err
		}

		needed := uint(len(assets))
		if sub.TokensRemaining < needed {
			return fmt.Errorf(
				"fondos insuficientes: se necesitan aprox. %d cuentokens, tienes %d",
				needed, sub.TokensRemaining,
			)
		}

		url, err := helper.MixAudio(id, assets)
		if err != nil {
			return err
		}

		script.Total_Cuentoken += needed
		script.Mixed_Audio = url
		if err := tx.Save(&script).Error; err != nil {
			return err
		}

		sub.TokensRemaining -= needed
		if err := tx.Save(sub).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	dto := ScriptToDTO(script)
	return &dto, nil
}
