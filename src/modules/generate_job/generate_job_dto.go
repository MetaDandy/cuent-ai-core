package generatejob

import (
	"time"

	"github.com/MetaDandy/cuent-ai-core/src/model"
)

type GeneratedJobResponse struct {
	ID            string  `json:"id"`
	Provider      string  `json:"provider"`
	Model         string  `json:"model"`
	Token_Spent   string  `json:"token_spent"`
	Chars_Used    uint    `json:"chars_used"`
	State         string  `json:"state"`
	Error_Message string  `json:"error_message"`
	Cost          float64 `json:"cost"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func GeneratedJobToDto(u *model.GeneratedJob) GeneratedJobResponse {
	var deletedAt *time.Time
	if u.DeletedAt.Valid {
		t := u.DeletedAt.Time
		deletedAt = &t
	}

	return GeneratedJobResponse{
		ID:            u.ID.String(),
		Provider:      string(u.Provider),
		Model:         u.Model,
		Token_Spent:   u.Token_Spent,
		Chars_Used:    u.Chars_Used,
		State:         string(u.State),
		Error_Message: u.Error_Message,
		Cost:          u.Cost,

		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func GeneratedJobToLisDTO(list []model.GeneratedJob) []GeneratedJobResponse {
	out := make([]GeneratedJobResponse, len(list))
	for i := range list {
		out[i] = GeneratedJobToDto(&list[i])
	}
	return out
}
