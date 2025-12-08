package service

import (
	"crypto/rand"
	"encoding/hex"

	"shosha-finance/internal/models"
	"shosha-finance/internal/repository"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type BranchService interface {
	Create(code, name string) (*models.Branch, error)
	GetByID(id uuid.UUID) (*models.Branch, error)
	GetByAPIKey(apiKey string) (*models.Branch, error)
	GetAll() ([]models.Branch, error)
}

type branchService struct {
	repo repository.BranchRepository
}

func NewBranchService(repo repository.BranchRepository) BranchService {
	return &branchService{repo: repo}
}

func (s *branchService) Create(code, name string) (*models.Branch, error) {
	apiKey, err := generateAPIKey()
	if err != nil {
		return nil, err
	}

	branch := &models.Branch{
		ID:     uuid.New(),
		Code:   code,
		Name:   name,
		APIKey: apiKey,
	}

	err = s.repo.Create(branch)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create branch")
		return nil, err
	}

	log.Info().Str("code", code).Msg("Branch created")
	return branch, nil
}

func (s *branchService) GetByID(id uuid.UUID) (*models.Branch, error) {
	return s.repo.FindByID(id)
}

func (s *branchService) GetByAPIKey(apiKey string) (*models.Branch, error) {
	return s.repo.FindByAPIKey(apiKey)
}

func (s *branchService) GetAll() ([]models.Branch, error) {
	return s.repo.FindAll()
}

func generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
