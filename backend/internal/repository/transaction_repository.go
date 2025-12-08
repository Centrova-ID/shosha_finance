package repository

import (
	"time"

	"shosha-finance/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionRepository interface {
	Create(tx *models.Transaction) error
	FindByID(id uuid.UUID) (*models.Transaction, error)
	FindAll(page, limit int) ([]models.Transaction, int64, error)
	GetDashboardSummary(filter *DashboardFilter) (*DashboardSummary, error)
	GetUnsyncedCount() (int64, error)
	Upsert(tx *models.Transaction) error
	GetUpdatedAfter(since *time.Time) ([]models.Transaction, error)
}

type DashboardSummary struct {
	TotalIn     int64 `json:"total_in"`
	TotalOut    int64 `json:"total_out"`
	Balance     int64 `json:"balance"`
	CountIn     int64 `json:"count_in"`
	CountOut    int64 `json:"count_out"`
	UnsyncCount int64 `json:"unsync_count"`
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(tx *models.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *transactionRepository) FindByID(id uuid.UUID) (*models.Transaction, error) {
	var tx models.Transaction
	err := r.db.Where("id = ?", id).First(&tx).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepository) FindAll(page, limit int) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	offset := (page - 1) * limit

	err := r.db.Model(&models.Transaction{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

type DashboardFilter struct {
	BranchID  *uuid.UUID
	StartDate *time.Time
	EndDate   *time.Time
}

func (r *transactionRepository) GetDashboardSummary(filter *DashboardFilter) (*DashboardSummary, error) {
	var summary DashboardSummary

	var totalIn, totalOut int64
	var countIn, countOut, unsyncCount int64

	// Helper to apply filters
	applyFilter := func(query *gorm.DB) *gorm.DB {
		if filter != nil {
			if filter.BranchID != nil {
				query = query.Where("branch_id = ?", *filter.BranchID)
			}
			if filter.StartDate != nil {
				query = query.Where("created_at >= ?", *filter.StartDate)
			}
			if filter.EndDate != nil {
				query = query.Where("created_at < ?", *filter.EndDate)
			}
		}
		return query
	}

	// Total IN
	queryIn := applyFilter(r.db.Model(&models.Transaction{}).Where("type = ?", models.TransactionTypeIN))
	err := queryIn.Select("COALESCE(SUM(amount), 0)").Scan(&totalIn).Error
	if err != nil {
		return nil, err
	}

	// Total OUT
	queryOut := applyFilter(r.db.Model(&models.Transaction{}).Where("type = ?", models.TransactionTypeOUT))
	err = queryOut.Select("COALESCE(SUM(amount), 0)").Scan(&totalOut).Error
	if err != nil {
		return nil, err
	}

	// Count IN
	queryCountIn := applyFilter(r.db.Model(&models.Transaction{}).Where("type = ?", models.TransactionTypeIN))
	err = queryCountIn.Count(&countIn).Error
	if err != nil {
		return nil, err
	}

	// Count OUT
	queryCountOut := applyFilter(r.db.Model(&models.Transaction{}).Where("type = ?", models.TransactionTypeOUT))
	err = queryCountOut.Count(&countOut).Error
	if err != nil {
		return nil, err
	}

	// Unsync count (no date filter for this)
	queryUnsync := r.db.Model(&models.Transaction{}).Where("is_synced = ?", false)
	if filter != nil && filter.BranchID != nil {
		queryUnsync = queryUnsync.Where("branch_id = ?", *filter.BranchID)
	}
	err = queryUnsync.Count(&unsyncCount).Error
	if err != nil {
		return nil, err
	}

	summary.TotalIn = totalIn
	summary.TotalOut = totalOut
	summary.Balance = totalIn - totalOut
	summary.CountIn = countIn
	summary.CountOut = countOut
	summary.UnsyncCount = unsyncCount

	return &summary, nil
}

func (r *transactionRepository) GetUnsyncedCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.Transaction{}).Where("is_synced = ?", false).Count(&count).Error
	return count, err
}

func (r *transactionRepository) Upsert(tx *models.Transaction) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(tx).Error
}

func (r *transactionRepository) GetUpdatedAfter(since *time.Time) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := r.db.Model(&models.Transaction{})
	if since != nil {
		query = query.Where("updated_at > ? OR created_at > ?", since, since)
	}
	err := query.Find(&transactions).Error
	return transactions, err
}
