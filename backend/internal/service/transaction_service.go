package service

import (
	"shosha-finance/internal/models"
	"shosha-finance/internal/repository"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type TransactionService interface {
	Create(req *models.TransactionRequest, branchID uuid.UUID) (*models.Transaction, error)
	GetByID(id uuid.UUID) (*models.Transaction, error)
	GetAll(page, limit int, branchID uuid.UUID) ([]models.Transaction, int64, error)
	GetDashboardSummary(branchID uuid.UUID) (*repository.DashboardSummary, error)
	GetUnsynced(limit int) ([]models.Transaction, error)
	MarkAsSynced(ids []uuid.UUID) error
	GetUnsyncedCount() (int64, error)
}

type transactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{repo: repo}
}

func (s *transactionService) Create(req *models.TransactionRequest, branchID uuid.UUID) (*models.Transaction, error) {
	tx := &models.Transaction{
		ID:          uuid.New(),
		BranchID:    branchID,
		Type:        req.Type,
		Category:    req.Category,
		Amount:      req.Amount,
		Description: req.Description,
		IsSynced:    false,
	}

	err := s.repo.Create(tx)
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

func (s *transactionService) GetAll(page, limit int, branchID uuid.UUID) ([]models.Transaction, int64, error) {
	return s.repo.FindAll(page, limit, branchID)
}

func (s *transactionService) GetDashboardSummary(branchID uuid.UUID) (*repository.DashboardSummary, error) {
	return s.repo.GetDashboardSummary(branchID)
}

func (s *transactionService) GetUnsynced(limit int) ([]models.Transaction, error) {
	return s.repo.FindUnsynced(limit)
}

func (s *transactionService) MarkAsSynced(ids []uuid.UUID) error {
	return s.repo.MarkAsSynced(ids)
}

func (s *transactionService) GetUnsyncedCount() (int64, error) {
	return s.repo.GetUnsyncedCount()
}
