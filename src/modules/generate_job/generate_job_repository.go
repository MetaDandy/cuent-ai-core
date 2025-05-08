package generatejob

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

func (r *Repository) Create(project *model.GeneratedJob) error {
	return r.db.Create(project).Error
}

func (r *Repository) Update(generatedJob *model.GeneratedJob) error {
	return r.db.Save(generatedJob).Error
}

func (r *Repository) FindAll(opts *helper.FindAllOptions) ([]model.GeneratedJob, int64, error) {
	var finded []model.GeneratedJob
	query := r.db.Model(model.GeneratedJob{})
	var total int64
	query, total = helper.ApplyFindAllOptions(query, opts)

	err := query.Find(&finded).Error
	return finded, total, err
}

func (r *Repository) FindById(id string) (*model.GeneratedJob, error) {
	var generatedJob model.GeneratedJob
	err := r.db.First(&generatedJob, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &generatedJob, nil
}

func (r *Repository) FindByIdUnscoped(id string) (*model.GeneratedJob, error) {
	var generatedJob model.GeneratedJob
	err := r.db.Unscoped().First(&generatedJob, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &generatedJob, nil
}

func (r *Repository) SoftDelete(id string) error {
	return r.db.Delete(&model.GeneratedJob{}, "id = ?", id).Error
}

func (r *Repository) Restore(id string) error {
	return r.db.Unscoped().
		Model(&model.GeneratedJob{}).
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}
