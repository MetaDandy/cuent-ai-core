package user

import (
	"time"

	"github.com/MetaDandy/cuent-ai-core/src/core/subscription"
	"github.com/MetaDandy/cuent-ai-core/src/model"
)

type Singup struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Signin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/* Change Passoword*/
type ChangePassoword struct {
	Old_Password     string `json:"old_password"`
	New_Password     string `json:"new_password"`
	Confirm_Password string `json:"confirm_password"`
}

type Payment struct {
	UserSuscribedID string `json:"user_suscribed_id"`
	PriceID         string `json:"price_id"`
}

type PaymentResponse struct {
	Session   string `json:"session"`
	PaymentID string `json:"payment_id"`
}

type PaymentDetail struct {
	ID        string    `json:"id"`
	Amount    int       `json:"amount"`
	Currency  string    `json:"currency"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`

	Subscriptions *[]UserSubscriptionResponse `json:"all_subscriptions,omitempty"`

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

	var subsPtr *[]UserSubscriptionResponse
	if len(u.UsersSubscriptions) > 0 {
		subs := UserSubscriptionToListDTO(u.UsersSubscriptions)
		subsPtr = &subs
	}

	return UserResponse{
		ID:            u.ID.String(),
		Name:          u.Name,
		Email:         u.Email,
		Subscriptions: subsPtr,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
		DeletedAt:     deletedAt,
	}
}

func UsersToListDTO(list []model.User) []UserResponse {
	out := make([]UserResponse, len(list))
	for i := range list {
		out[i] = UserToDTO(&list[i])
	}
	return out
}

/*  User Subscription*/
type UserSubscriptionResponse struct {
	ID               string `json:"id"`
	Total_Cuentokens uint   `json:"total_Cuentokens"`
	Start_Date       string `json:"start_date"`
	End_Date         string `json:"end_date"`

	Subscription subscription.SubscriptionResponse
	Payments     []PaymentDetail `json:"payments,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func UserSubscriptionToDto(u *model.UserSubscribed) UserSubscriptionResponse {
	var deletedAt *time.Time
	if u.DeletedAt.Valid {
		t := u.DeletedAt.Time
		deletedAt = &t
	}

	payments := make([]PaymentDetail, len(u.Payments))
	if len(payments) > 0 {
		payments = PaymentsToListDTO(u.Payments)
	}

	return UserSubscriptionResponse{
		ID:               u.ID.String(),
		Total_Cuentokens: u.TokensRemaining,
		Start_Date:       u.StartDate.Local().String(),
		End_Date:         u.EndDate.Local().String(),
		Subscription:     subscription.SubscriptionToDTO(&u.Subscription),
		Payments:         payments,
		CreatedAt:        u.CreatedAt,
		UpdatedAt:        u.UpdatedAt,
		DeletedAt:        deletedAt,
	}
}

func UserSubscriptionToListDTO(list []model.UserSubscribed) []UserSubscriptionResponse {
	out := make([]UserSubscriptionResponse, len(list))
	for i := range list {
		out[i] = UserSubscriptionToDto(&list[i])
	}
	return out
}

// PaymentToDTO convierte un model.Payment en PaymentDetail
func PaymentToDTO(p *model.Payment) PaymentDetail {
	return PaymentDetail{
		ID:        p.ID.String(),
		Amount:    p.Amount,
		Currency:  p.Currency,
		Status:    string(p.Status),
		CreatedAt: p.CreatedAt,
	}
}

func PaymentsToListDTO(list []model.Payment) []PaymentDetail {
	out := make([]PaymentDetail, len(list))
	for i := range list {
		out[i] = PaymentToDTO(&list[i])
	}
	return out
}
