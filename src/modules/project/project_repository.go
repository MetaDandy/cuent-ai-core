package project

import (
	"github.com/MetaDandy/cuent-ai-core/helper"
	"github.com/MetaDandy/cuent-ai-core/src/model"
	"gorm.io/gorm"
)

type Repository interface {
	Create(project *model.Project) error
	Update(project *model.Project) error
	FindAll(opts *helper.FindAllOptions) ([]model.Project, int64, error)
	FindById(id string) (*model.Project, error)
	FindByIdUnscoped(id string) (*model.Project, error)
	SoftDelete(id string) error
	Restore(id string) error
}

type PostgresRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(project *model.Project) error {
	return r.db.Create(project).Error
}

func (r *PostgresRepository) Update(project *model.Project) error {
	return r.db.Save(project).Error
}

func (r *PostgresRepository) FindAll(opts *helper.FindAllOptions) ([]model.Project, int64, error) {
	var finded []model.Project
	query := r.db.Model(model.Project{})
	var total int64
	query, total = helper.ApplyFindAllOptions(query, opts)

	err := query.Find(&finded).Error
	return finded, total, err
}

func (r *PostgresRepository) FindById(id string) (*model.Project, error) {
	var project model.Project
	err := r.db.Preload("Scripts").First(&project, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *PostgresRepository) FindByIdUnscoped(id string) (*model.Project, error) {
	var project model.Project
	err := r.db.Unscoped().First(&project, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *PostgresRepository) SoftDelete(id string) error {
	return r.db.Delete(&model.Project{}, "id = ?", id).Error
}

func (r *PostgresRepository) Restore(id string) error {
	return r.db.Unscoped().
		Model(&model.Project{}).
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}
