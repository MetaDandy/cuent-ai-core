package script

import (
	"github.com/MetaDandy/cuent-ai-core/helper"
	"github.com/MetaDandy/cuent-ai-core/src/model"
	"github.com/MetaDandy/cuent-ai-core/src/modules/asset"
	"github.com/MetaDandy/cuent-ai-core/src/modules/project"
	"gorm.io/gorm"
)

type Service struct {
	repo        *Repository
	projectRepo *project.Repository
	assetRepo   *asset.Repository
}

func NewService(r *Repository, pr *project.Repository, ar *asset.Repository) *Service {
	return &Service{repo: r, projectRepo: pr, assetRepo: ar}
}

/**
TODO:
- para los mixed, crear endpoints especiales
- Crear un endpoint para traer todos los assets de un script
- Resolver todos los comentarios con !
*/

func (s *Service) Create(input *ScriptCreate) (*ScriptReponse, error) {
	project, err := s.projectRepo.FindById(input.ProjectID)
	if err != nil {
		return nil, err
	}

	var script model.Script
	// ! Contar los tokens para actualizar los tokens del proyecto
	if err := s.repo.db.Transaction(func(tx *gorm.DB) error {
		aiResponse, err := helper.AIFormatter(input.TextEntry)
		if err != nil {
			return err
		}

		script = model.Script{
			Text_Entry:        input.TextEntry,
			ProjectID:         project.ID,
			Prompt_Tokens:     aiResponse.Prompt_Tokens,
			Completion_Tokens: aiResponse.Completion_Tokens,
			Total_Tokens:      aiResponse.Total_Tokens,
			Processed_Text:    aiResponse.Processed_Text,
			State:             model.StateFinished,
			// ! Ver como calcular el total cost
		}
		if err := tx.Create(&script).Error; err != nil {
			return err
		}

		assets := make([]model.Asset, 0, len(aiResponse.Processed_Text_Array))
		for _, line := range aiResponse.Processed_Text_Array {
			assets = append(assets, model.Asset{
				Type:     "LINE",
				Line:     line,
				ScriptID: script.ID,
			})
		}
		if err := tx.Create(&assets).Error; err != nil {
			return err
		}

		// ! Ver como actualizar los cuentokens del proyecto

		return nil
	}); err != nil {
		return nil, err
	}

	// ! Ver si creamos un findby id que busque con relaciones con assets
	reload, _ := s.repo.FindById(script.ID.String())
	dto := ScriptToDTO(reload)
	return &dto, nil
}
