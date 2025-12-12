package repository

import (
	"errors"
	"shosha-finance/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IncomeEntryRepository interface {
	Create(entry *models.IncomeEntry) error
	Update(entry *models.IncomeEntry) error
	GetByID(id uuid.UUID) (*models.IncomeEntry, error)
	GetByBranchAndDate(branchID uuid.UUID, date time.Time) (*models.IncomeEntry, error)
	GetByDateRange(startDate, endDate time.Time) ([]models.IncomeEntry, error)
	GetAll(page, limit int) ([]models.IncomeEntry, int64, error)
	Delete(id uuid.UUID) error
}

type incomeEntryRepository struct {
	db *gorm.DB
}

func NewIncomeEntryRepository(db *gorm.DB) IncomeEntryRepository {
	return &incomeEntryRepository{db: db}
}

func (r *incomeEntryRepository) Create(entry *models.IncomeEntry) error {
	return r.db.Create(entry).Error
}

func (r *incomeEntryRepository) Update(entry *models.IncomeEntry) error {
	return r.db.Save(entry).Error
}

func (r *incomeEntryRepository) GetByID(id uuid.UUID) (*models.IncomeEntry, error) {
	var entry models.IncomeEntry
	err := r.db.Preload("Branch").First(&entry, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("income entry not found")
		}
		return nil, err
	}
	return &entry, nil
}

func (r *incomeEntryRepository) GetByBranchAndDate(branchID uuid.UUID, date time.Time) (*models.IncomeEntry, error) {
	var entry models.IncomeEntry
	err := r.db.Preload("Branch").Where("branch_id = ? AND date = ?", branchID, date).First(&entry).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entry, nil
}

func (r *incomeEntryRepository) GetByDateRange(startDate, endDate time.Time) ([]models.IncomeEntry, error) {
	var entries []models.IncomeEntry
	err := r.db.Preload("Branch").
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Order("date DESC, branch_id ASC").
		Find(&entries).Error
	return entries, err
}

func (r *incomeEntryRepository) GetAll(page, limit int) ([]models.IncomeEntry, int64, error) {
	var entries []models.IncomeEntry
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&models.IncomeEntry{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("Branch").
		Order("date DESC, branch_id ASC").
		Offset(offset).
		Limit(limit).
		Find(&entries).Error

	return entries, total, err
}

func (r *incomeEntryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.IncomeEntry{}, "id = ?", id).Error
}

type ExpenseEntryRepository interface {
	Create(entry *models.ExpenseEntry) error
	Update(entry *models.ExpenseEntry) error
	GetByID(id uuid.UUID) (*models.ExpenseEntry, error)
	GetByBranchAndDate(branchID uuid.UUID, date time.Time) (*models.ExpenseEntry, error)
	GetByDateRange(startDate, endDate time.Time) ([]models.ExpenseEntry, error)
	GetAll(page, limit int) ([]models.ExpenseEntry, int64, error)
	Delete(id uuid.UUID) error
}

type expenseEntryRepository struct {
	db *gorm.DB
}

func NewExpenseEntryRepository(db *gorm.DB) ExpenseEntryRepository {
	return &expenseEntryRepository{db: db}
}

func (r *expenseEntryRepository) Create(entry *models.ExpenseEntry) error {
	return r.db.Create(entry).Error
}

func (r *expenseEntryRepository) Update(entry *models.ExpenseEntry) error {
	return r.db.Save(entry).Error
}

func (r *expenseEntryRepository) GetByID(id uuid.UUID) (*models.ExpenseEntry, error) {
	var entry models.ExpenseEntry
	err := r.db.Preload("Branch").First(&entry, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("expense entry not found")
		}
		return nil, err
	}
	return &entry, nil
}

func (r *expenseEntryRepository) GetByBranchAndDate(branchID uuid.UUID, date time.Time) (*models.ExpenseEntry, error) {
	var entry models.ExpenseEntry
	err := r.db.Preload("Branch").Where("branch_id = ? AND date = ?", branchID, date).First(&entry).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entry, nil
}

func (r *expenseEntryRepository) GetByDateRange(startDate, endDate time.Time) ([]models.ExpenseEntry, error) {
	var entries []models.ExpenseEntry
	err := r.db.Preload("Branch").
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Order("date DESC, branch_id ASC").
		Find(&entries).Error
	return entries, err
}

func (r *expenseEntryRepository) GetAll(page, limit int) ([]models.ExpenseEntry, int64, error) {
	var entries []models.ExpenseEntry
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&models.ExpenseEntry{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("Branch").
		Order("date DESC, branch_id ASC").
		Offset(offset).
		Limit(limit).
		Find(&entries).Error

	return entries, total, err
}

func (r *expenseEntryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.ExpenseEntry{}, "id = ?", id).Error
}
