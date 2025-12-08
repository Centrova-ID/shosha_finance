package repository

import (
	"shosha-finance/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BranchRepository interface {
	Create(branch *models.Branch) error
	FindByID(id uuid.UUID) (*models.Branch, error)
	FindByCode(code string) (*models.Branch, error)
	FindByAPIKey(apiKey string) (*models.Branch, error)
	FindAll() ([]models.Branch, error)
	Update(branch *models.Branch) error
	Delete(id uuid.UUID) error
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

func (r *branchRepository) FindByAPIKey(apiKey string) (*models.Branch, error) {
	var branch models.Branch
	err := r.db.Where("api_key = ?", apiKey).First(&branch).Error
	if err != nil {
		return nil, err
	}
	return &branch, nil
}

func (r *branchRepository) FindAll() ([]models.Branch, error) {
	var branches []models.Branch
	err := r.db.Find(&branches).Error
	return branches, err
}

func (r *branchRepository) Update(branch *models.Branch) error {
	return r.db.Save(branch).Error
}

func (r *branchRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Branch{}, id).Error
}
