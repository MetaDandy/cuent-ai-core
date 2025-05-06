package project

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

func (r *Repository) Create(project *model.Project) error {
	return r.db.Create(project).Error
}

func (r *Repository) Update(project *model.Project) error {
	return r.db.Save(project).Error
}

func (r *Repository) FindAll(opts *helper.FindAllOptions) ([]model.Project, int64, error) {
	var finded []model.Project
	query := r.db.Model(model.User{})
	var total int64
	query, total = helper.ApplyFindAllOptions(query, opts)

	err := query.Find(&finded).Error
	return finded, total, err
}

func (r *Repository) FindById(id string) (*model.Project, error) {
	var project model.Project
	err := r.db.First(&project, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &project, err
}

func (r *Repository) FindByIdUnscoped(id string) (*model.Project, error) {
	var project model.Project
	err := r.db.Unscoped().First(&project, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &project, err
}

func (r *Repository) SoftDelete(id string) error {
	return r.db.Delete(&model.Project{}, "id = ?", id).Error
}

func (r *Repository) Restore(id string) error {
	return r.db.Unscoped().
		Model(&model.Project{}).
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}
