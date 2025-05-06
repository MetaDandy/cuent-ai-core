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
	// Tambien poner script si amerita

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

	return ProjectResponse{
		ID:          u.ID.String(),
		Name:        u.Name,
		Description: u.Description,
		Cuentokens:  u.Cuentokens,
		State:       string(u.State),
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
