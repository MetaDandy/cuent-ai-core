package script

import (
	"time"

	"github.com/MetaDandy/cuent-ai-core/src/model"
)

type ScriptCreate struct {
	TextEntry string `json:"text_entry" validate:"required"`
	ProjectID string `json:"project_id" validate:"required"`
}

type ScriptUpdate struct {
	TextEntry *string `json:"text_entry"`
}

type ScriptReponse struct {
	ID                string  `json:"id"`
	Prompt_Tokens     uint32  `json:"promt_tokens"`
	Completion_Tokens uint32  `json:"completion_tokens"`
	Total_Tokens      uint32  `json:"total_tokens"`
	State             string  `json:"state"`
	Text_Entry        string  `json:"text_entry"`
	Processed_Text    string  `json:"processed_text"`
	Total_Cost        float64 `json:"total_cost"`
	Mixed_Audio       string  `json:"mixed_audio"`
	Mixed_Media       string  `json:"mixed_media"`

	//Poner proyecto
	//Poner assets

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

	return ScriptReponse{
		ID:                u.ID.String(),
		Prompt_Tokens:     u.Prompt_Tokens,
		Completion_Tokens: u.Completion_Tokens,
		Total_Tokens:      u.Total_Tokens,
		State:             string(u.State),
		Text_Entry:        u.Text_Entry,
		Processed_Text:    u.Processed_Text,
		Total_Cost:        u.Total_Cost,
		Mixed_Audio:       u.Mixed_Audio,
		Mixed_Media:       u.Mixed_Media,

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
