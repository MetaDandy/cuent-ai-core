package script

import (
	"github.com/MetaDandy/cuent-ai-core/helper"
	"github.com/MetaDandy/cuent-ai-core/src/model"
	"github.com/MetaDandy/cuent-ai-core/src/modules/asset"
	"github.com/MetaDandy/cuent-ai-core/src/modules/project"
	"github.com/google/uuid"
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
			ID:                uuid.New(),
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

		// ! Ver como actualizar los cuentokens del proyecto

		return nil
	}); err != nil {
		return nil, err
	}

	// ! Ver si creamos un findby id que busque con relaciones con assets
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
		return nil, nil
	}
	dto := ScriptToDTO(finded)
	return &dto, nil
}
