package repository

import (
	"time"

	"shosha-finance/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BranchRepository interface {
	Create(branch *models.Branch) error
	FindByID(id uuid.UUID) (*models.Branch, error)
	FindByCode(code string) (*models.Branch, error)
	FindAll() ([]models.Branch, error)
	FindActive() ([]models.Branch, error)
	Update(branch *models.Branch) error
	Delete(id uuid.UUID) error
	Count() (int64, error)
	Upsert(branch *models.Branch) error
	GetUpdatedAfter(since *time.Time) ([]models.Branch, error)
}

type branchRepository struct {
	db *gorm.DB
}

func NewBranchRepository(db *gorm.DB) BranchRepository {
	return &branchRepository{db: db}
}

func (r *branchRepository) Create(branch *models.Branch) error {
	return r.db.Create(branch).Error
}

func (r *branchRepository) FindByID(id uuid.UUID) (*models.Branch, error) {
	var branch models.Branch
	err := r.db.Where("id = ?", id).First(&branch).Error
	if err != nil {
		return nil, err
	}
	return &branch, nil
}

func (r *branchRepository) FindByCode(code string) (*models.Branch, error) {
	var branch models.Branch
	err := r.db.Where("code = ?", code).First(&branch).Error
	if err != nil {
		return nil, err
	}
	return &branch, nil
}

func (r *branchRepository) FindAll() ([]models.Branch, error) {
	var branches []models.Branch
	err := r.db.Order("name asc").Find(&branches).Error
	return branches, err
}

func (r *branchRepository) FindActive() ([]models.Branch, error) {
	var branches []models.Branch
	err := r.db.Where("is_active = ?", true).Order("name asc").Find(&branches).Error
	return branches, err
}

func (r *branchRepository) Update(branch *models.Branch) error {
	return r.db.Save(branch).Error
}

func (r *branchRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Branch{}, id).Error
}

func (r *branchRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Branch{}).Count(&count).Error
	return count, err
}

func (r *branchRepository) Upsert(branch *models.Branch) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(branch).Error
}

func (r *branchRepository) GetUpdatedAfter(since *time.Time) ([]models.Branch, error) {
	var branches []models.Branch
	query := r.db.Model(&models.Branch{})
	if since != nil {
		query = query.Where("updated_at > ? OR created_at > ?", since, since)
	}
	err := query.Find(&branches).Error
	return branches, err
}
