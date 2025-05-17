package subscription

import (
	"time"

	"github.com/MetaDandy/cuent-ai-core/src/model"
)

type SubscriptionResponse struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Cuentokens string    `json:"cuent_tokens"`
	Duration   time.Time `json:"duration"`

	// ponser user subscription si se necesita

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func SubscriptionToDTO(u *model.Subscription) SubscriptionResponse {
	var deletedAt *time.Time
	if u.DeletedAt.Valid {
		t := u.DeletedAt.Time
		deletedAt = &t
	}

	return SubscriptionResponse{
		ID:         u.ID.String(),
		Name:       u.Name,
		Cuentokens: u.Cuentokens,
		Duration:   u.Duration,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
		DeletedAt:  deletedAt,
	}
}

func SubscriptionToListDTO(list []model.Subscription) []SubscriptionResponse {
	out := make([]SubscriptionResponse, len(list))
	for i := range list {
		out[i] = SubscriptionToDTO(&list[i])
	}
	return out
}
