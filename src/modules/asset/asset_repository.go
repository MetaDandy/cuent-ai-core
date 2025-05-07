package asset

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

func (r *Repository) Create(asset *model.Asset) error {
	return r.db.Create(asset).Error
}

func (r *Repository) Update(asset *model.Asset) error {
	return r.db.Save(asset).Error
}

func (r *Repository) FindAll(opts *helper.FindAllOptions) ([]model.Asset, int64, error) {
	var finded []model.Asset
	query := r.db.Model(model.User{})
	var total int64
	query, total = helper.ApplyFindAllOptions(query, opts)

	err := query.Find(&finded).Error
	return finded, total, err
}

func (r *Repository) FindById(id string) (*model.Asset, error) {
	var asset model.Asset
	err := r.db.First(&asset, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &asset, err
}

func (r *Repository) FindByIdUnscoped(id string) (*model.Asset, error) {
	var asset model.Asset
	err := r.db.Unscoped().First(&asset, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &asset, err
}

func (r *Repository) SoftDelete(id string) error {
	return r.db.Delete(&model.Asset{}, "id = ?", id).Error
}

func (r *Repository) Restore(id string) error {
	return r.db.Unscoped().
		Model(&model.Asset{}).
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}
