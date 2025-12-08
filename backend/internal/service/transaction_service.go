package service

import (
	"time"

	"shosha-finance/internal/models"
	"shosha-finance/internal/repository"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type TransactionService interface {
	Create(req *models.TransactionRequest) (*models.Transaction, error)
	GetByID(id uuid.UUID) (*models.Transaction, error)
	GetAll(page, limit int) ([]models.Transaction, int64, error)
	GetDashboardSummary(filter *repository.DashboardFilter) (*repository.DashboardSummary, error)
	GetUnsyncedCount() (int64, error)
	Upsert(tx *models.Transaction) error
	GetUpdatedAfter(since *time.Time) ([]models.Transaction, error)
}

type transactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{repo: repo}
}

func (s *transactionService) Create(req *models.TransactionRequest) (*models.Transaction, error) {
	branchID, err := uuid.Parse(req.BranchID)
	if err != nil {
		return nil, err
	}

	tx := &models.Transaction{
		ID:          uuid.New(),
		BranchID:    branchID,
		Type:        req.Type,
		Category:    req.Category,
		Amount:      req.Amount,
		Description: req.Description,
	}

	err = s.repo.Create(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create transaction")
		return nil, err
	}

	log.Info().Str("id", tx.ID.String()).Msg("Transaction created")
	return tx, nil
}

func (s *transactionService) GetByID(id uuid.UUID) (*models.Transaction, error) {
	return s.repo.FindByID(id)
}

func (s *transactionService) GetAll(page, limit int) ([]models.Transaction, int64, error) {
	return s.repo.FindAll(page, limit)
}

func (s *transactionService) GetDashboardSummary(filter *repository.DashboardFilter) (*repository.DashboardSummary, error) {
	return s.repo.GetDashboardSummary(filter)
}

func (s *transactionService) GetUnsyncedCount() (int64, error) {
	return s.repo.GetUnsyncedCount()
}

func (s *transactionService) Upsert(tx *models.Transaction) error {
	return s.repo.Upsert(tx)
}

func (s *transactionService) GetUpdatedAfter(since *time.Time) ([]models.Transaction, error) {
	return s.repo.GetUpdatedAfter(since)
}
