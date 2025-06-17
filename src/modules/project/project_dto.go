package project

import (
	"time"

	"github.com/MetaDandy/cuent-ai-core/src/model"
)

type ProjectCreate struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	UserId      string `json:"user_id" validate:"required"`
}

type ProjectUpdate struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type ProjectResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Cuentokens  string `json:"cuentokens"`
	State       string `json:"state"`

	// Poner user cuando se cree si amerita
	Script []ScriptReponse `json:"scripts,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func ProjectToDTO(u *model.Project) ProjectResponse {
	var deletedAt *time.Time
	if u.DeletedAt.Valid {
		t := u.DeletedAt.Time
		deletedAt = &t
	}

	var scripts []ScriptReponse
	if len(u.Scripts) > 0 {
		scripts = make([]ScriptReponse, 0, len(u.Scripts))
		for i := range u.Scripts {
			scripts = append(scripts, ScriptToDTO(&u.Scripts[i]))
		}
	}

	return ProjectResponse{
		ID:          u.ID.String(),
		Name:        u.Name,
		Description: u.Description,
		Cuentokens:  u.Cuentokens,
		State:       string(u.State),
		Script:      scripts,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		DeletedAt:   deletedAt,
	}
}

func ProjectsToListDTO(list []model.Project) []ProjectResponse {
	out := make([]ProjectResponse, len(list))
	for i := range list {
		out[i] = ProjectToDTO(&list[i])
	}
	return out
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
		Total_Token:       u.Total_Tokens,
		Total_Cuentoken:   u.Total_Cuentoken,
		State:             string(u.State),
		Text_Entry:        u.Text_Entry,
		Processed_Text:    u.Processed_Text,
		Mixed_Audio:       u.Mixed_Audio,
		Mixed_Media:       u.Mixed_Media,

		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: deletedAt,
	}
}
