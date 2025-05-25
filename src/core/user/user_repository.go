package user

import (
	"strings"
	"time"

	"github.com/MetaDandy/cuent-ai-core/helper"
	"github.com/MetaDandy/cuent-ai-core/src/model"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAll(opts *helper.FindAllOptions) ([]model.User, int64, error) {
	var finded []model.User
	query := r.db.Model(model.User{})
	var total int64
	query, total = helper.ApplyFindAllOptions(query, opts)

	err := query.Find(&finded).Error
	return finded, total, err
}

func (c *Repository) FindById(id string) (*model.User, error) {
	var user model.User
	err := c.db.
		Preload("UsersSubscriptions.Subscription").
		First(&user, "id = ?", id).Error

	return &user, err
}

func (c *Repository) FindSubscriptionById(id string) (*model.Subscription, error) {
	var sub model.Subscription
	err := c.db.
		First(&sub, "id = ?", id).Error

	return &sub, err
}

func (r *Repository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.
		Where("email = ?", email).
		Preload("UsersSubscriptions.Subscription").
		First(&user).Error

	return &user, err
}

func (r *Repository) Create(u *model.User) error {
	return r.db.Create(u).Error
}

func (r *Repository) Update(u *model.User) error {
	return r.db.Save(u).Error
}

func (r *Repository) CreateUserSubscription(sub *model.UserSubscribed) error {
	return r.db.Create(sub).Error
}

func (r *Repository) FindUserSuscribedByID(id string) (*model.UserSubscribed, error) {
	var sub *model.UserSubscribed
	err := r.db.
		Preload("Subscription").
		First(&sub, "id = ?", id).Error
	return sub, err
}

func (r *Repository) UpdateUserSuscribed(u *model.UserSubscribed) error {
	return r.db.Save(u).Error
}

func (r *Repository) FindPaymentID(id string) (*model.Payment, error) {
	var pay *model.Payment
	err := r.db.
		Preload("UserSuscribed").
		Where("id = ?", id).
		First(&pay).Error
	return pay, err
}

func (r *Repository) UpdatePayment(u *model.Payment) error {
	return r.db.Save(u).Error
}

func (r *Repository) ClosePreviousSubscriptions(userID string) error {
	return r.db.
		Model(&model.UserSubscribed{}).
		Where("user_id = ? AND end_date >= ?", userID, time.Now()).
		Update("end_date", time.Now()).
		Error
}

func (r *Repository) CreatePayment(u *model.Payment) error {
	return r.db.Create(u).Error
}

func (r *Repository) AddSubscription(sub *model.UserSubscribed) error {
	return r.db.Transaction(func(tx *gorm.DB) error { // asegura consistencia :contentReference[oaicite:3]{index=3}
		// 1. Anula cualquier suscripción activa anterior
		if err := tx.Model(&model.UserSubscribed{}).
			Where("user_id = ? AND end_date >= ?", sub.UserID, time.Now()).
			Update("end_date", time.Now()).Error; err != nil {
			return err
		}

		// 2. Crea la nueva
		if err := tx.Create(sub).Error; err != nil {
			return err
		}

		// 3. Preload para rellenar el struct
		return tx.Preload("Subscription").First(sub, "id = ?", sub.ID).Error
	})
}

func (r *Repository) GetActiveSubscription(userID string) (*model.UserSubscribed, error) {
	var us model.UserSubscribed
	err := r.db.
		Preload("Subscription").
		Preload("Payments").
		Where("user_id = ? AND end_date >= ?", userID, time.Now()).
		First(&us).Error
	return &us, err
}

func (r *Repository) FindByIdUnscoped(id string) (*model.User, error) {
	var user model.User
	err := r.db.Unscoped().First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindSubscriptionByName(name string) (*model.Subscription, error) {
	var sub model.Subscription
	err := r.db.First(&sub, "LOWER(name) = ?", strings.ToLower(name)).Error
	return &sub, err
}

func (r *Repository) SoftDelete(id string) error {
	return r.db.Delete(&model.User{}, "id = ?", id).Error
}

func (r *Repository) Restore(id string) error {
	return r.db.Unscoped().
		Model(&model.User{}).
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}
