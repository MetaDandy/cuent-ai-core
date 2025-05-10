package asset

import (
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
	query := r.db.Model(model.Asset{})
	var total int64
	query, total = helper.ApplyFindAllOptions(query, opts)

	err := query.Find(&finded).Error
	return finded, total, err
}

func (r *Repository) FindById(id string) (*model.Asset, error) {
	var asset model.Asset
	err := r.db.First(&asset, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (r *Repository) FindByIdWithGeneratedJobs(id string) (*model.Asset, error) {
	var asset model.Asset
	err := r.db.Preload("GeneratedJobs").First(&asset, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (r *Repository) FindByIdUnscoped(id string) (*model.Asset, error) {
	var asset model.Asset
	err := r.db.Unscoped().First(&asset, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (r *Repository) FindByScriptID(scriptID string) ([]model.Asset, error) {
	var assets []model.Asset
	err := r.db.
		Where("script_id = ?", scriptID).
		Order("position").
		Find(&assets).Error
	return assets, err
}

func (r *Repository) FindByScriptIDWithGeneratedJobs(scriptID string) ([]model.Asset, error) {
	var assets []model.Asset
	err := r.db.
		Preload("GeneratedJobs").
		Where("script_id = ?", scriptID).
		Order("position").
		Find(&assets).Error
	return assets, err
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
