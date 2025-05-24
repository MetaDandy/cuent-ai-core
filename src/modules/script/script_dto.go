package script

import (
	"time"

	"github.com/MetaDandy/cuent-ai-core/src/model"
	"github.com/MetaDandy/cuent-ai-core/src/modules/asset"
)

type ScriptCreate struct {
	TextEntry string `json:"text_entry" validate:"required"`
	ProjectID string `json:"project_id" validate:"required"`
}

type Line struct {
	Text string          `json:"text" validate:"required"`
	Type model.AudioLine `json:"type" validate:"required,oneof=TTS SFX"`
}

type ScriptManualCreate struct {
	Lines     []Line `json:"lines" validate:"required,dive"`
	ProjectID string `json:"project_id" validate:"required,uuid"`
}

type ScriptUpdate struct {
	TextEntry *string `json:"text_entry"`
}

type ScriptReponse struct {
	ID                string `json:"id"`
	Prompt_Tokens     uint32 `json:"promt_tokens"`
	Completion_Tokens uint32 `json:"completion_tokens"`
	Total_Token       uint32 `json:"total_token"`
	Total_Cuentoken   uint   `json:"total_cuentoken"`
	State             string `json:"state"`
	Text_Entry        string `json:"text_entry"`
	Processed_Text    string `json:"processed_text"`
	Mixed_Audio       string `json:"mixed_audio"`
	Mixed_Media       string `json:"mixed_media"`

	Assets []asset.AssetResponse `json:"assets,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func ScriptToDTO(u *model.Script) ScriptReponse {
	var deletedAt *time.Time
	if u.DeletedAt.Valid {
		t := u.DeletedAt.Time
		deletedAt = &t
	}

	var assets []asset.AssetResponse
	if len(u.Assets) > 0 {
		assets = make([]asset.AssetResponse, 0, len(u.Assets))
		for i := range u.Assets {
			assets = append(assets, asset.AssetToDto(&u.Assets[i]))
		}
	}

	return ScriptReponse{
		ID:                u.ID.String(),
		Prompt_Tokens:     u.Prompt_Tokens,
		Completion_Tokens: u.Completion_Tokens,
		Total_Token:       u.Total_Tokens,
		Total_Cuentoken:   u.Total_Cuentoken,
		State:             string(u.State),
		Text_Entry:        u.Text_Entry,
		Processed_Text:    u.Processed_Text,
		Mixed_Audio:       u.Mixed_Audio,
		Mixed_Media:       u.Mixed_Media,
		Assets:            assets,

		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func ScriptToListDTO(list []model.Script) []ScriptReponse {
	out := make([]ScriptReponse, len(list))
	for i := range list {
		out[i] = ScriptToDTO(&list[i])
	}
	return out
}
