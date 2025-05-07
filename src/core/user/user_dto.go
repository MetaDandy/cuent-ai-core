package user

import (
	"time"

	"github.com/MetaDandy/cuent-ai-core/src/model"
)

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func UserToDTO(u *model.User) UserResponse {
	var deletedAt *time.Time
	if u.DeletedAt.Valid {
		t := u.DeletedAt.Time
		deletedAt = &t
	}

	return UserResponse{
		ID:        u.ID.String(),
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func UsersToListDTO(list []model.User) []UserResponse {
	out := make([]UserResponse, len(list))
	for i := range list {
		out[i] = UserToDTO(&list[i])
	}
	return out
}
