package user

import (
	"errors"

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
	err := c.db.First(&user, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}
