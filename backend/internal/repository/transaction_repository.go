package repository

import (
	"shosha-finance/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(tx *models.Transaction) error
	FindByID(id uuid.UUID) (*models.Transaction, error)
	FindAll(page, limit int, branchID uuid.UUID) ([]models.Transaction, int64, error)
	FindUnsynced(limit int) ([]models.Transaction, error)
	MarkAsSynced(ids []uuid.UUID) error
	GetDashboardSummary(branchID uuid.UUID) (*DashboardSummary, error)
	GetUnsyncedCount() (int64, error)
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

func (r *transactionRepository) FindAll(page, limit int, branchID uuid.UUID) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	offset := (page - 1) * limit

	query := r.db.Model(&models.Transaction{})
	if branchID != uuid.Nil {
		query = query.Where("branch_id = ?", branchID)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *transactionRepository) FindUnsynced(limit int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Where("is_synced = ?", false).Limit(limit).Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) MarkAsSynced(ids []uuid.UUID) error {
	return r.db.Model(&models.Transaction{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"is_synced": true,
			"synced_at": gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error
}

func (r *transactionRepository) GetDashboardSummary(branchID uuid.UUID) (*DashboardSummary, error) {
	var summary DashboardSummary

	query := r.db.Model(&models.Transaction{})
	if branchID != uuid.Nil {
		query = query.Where("branch_id = ?", branchID)
	}

	var totalIn, totalOut int64
	var countIn, countOut int64

	err := query.Where("type = ?", models.TransactionTypeIN).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalIn).Error
	if err != nil {
		return nil, err
	}

	err = query.Where("type = ?", models.TransactionTypeOUT).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalOut).Error
	if err != nil {
		return nil, err
	}

	err = query.Where("type = ?", models.TransactionTypeIN).Count(&countIn).Error
	if err != nil {
		return nil, err
	}

	err = query.Where("type = ?", models.TransactionTypeOUT).Count(&countOut).Error
	if err != nil {
		return nil, err
	}

	var unsyncCount int64
	err = r.db.Model(&models.Transaction{}).Where("is_synced = ?", false).Count(&unsyncCount).Error
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
