package script

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

func (r *Repository) Create(script *model.Script) error {
	return r.db.Create(script).Error
}

func (r *Repository) Update(script *model.Script) error {
	return r.db.Save(script).Error
}

func (r *Repository) FindAll(opts *helper.FindAllOptions) ([]model.Script, int64, error) {
	var finded []model.Script
	query := r.db.Model(model.Script{})
	var total int64
	query, total = helper.ApplyFindAllOptions(query, opts)

	err := query.Find(&finded).Error
	return finded, total, err
}

func (r *Repository) FindById(id string) (*model.Script, error) {
	var script model.Script
	err := r.db.First(&script, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &script, nil
}

func (r *Repository) FindByIdWithAssets(id string) (*model.Script, error) {
	var script model.Script
	err := r.db.Preload("Assets").First(&script, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &script, nil
}

func (r *Repository) FindByIdUnscoped(id string) (*model.Script, error) {
	var script model.Script
	err := r.db.Unscoped().First(&script, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &script, nil
}

func (r *Repository) SoftDelete(id string) error {
	return r.db.Delete(&model.Script{}, "id = ?", id).Error
}

func (r *Repository) Restore(id string) error {
	return r.db.Unscoped().
		Model(&model.Script{}).
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}
