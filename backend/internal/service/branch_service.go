package service

import (
	"time"

	"shosha-finance/internal/models"
	"shosha-finance/internal/repository"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type BranchService interface {
	Create(req *models.BranchRequest) (*models.Branch, error)
	GetByID(id uuid.UUID) (*models.Branch, error)
	GetAll() ([]models.Branch, error)
	GetActive() ([]models.Branch, error)
	Update(id uuid.UUID, req *models.BranchRequest) (*models.Branch, error)
	Delete(id uuid.UUID) error
	Count() (int64, error)
	CreateDefaultBranches() error
	Upsert(branch *models.Branch) error
	GetUpdatedAfter(since *time.Time) ([]models.Branch, error)
}

type branchService struct {
	repo repository.BranchRepository
}

func NewBranchService(repo repository.BranchRepository) BranchService {
	return &branchService{repo: repo}
}

func (s *branchService) Create(req *models.BranchRequest) (*models.Branch, error) {
	log.Info().Str("code", req.Code).Str("name", req.Name).Msg("Creating branch...")
	
	branch := &models.Branch{
		ID:          uuid.New(),
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
	}

	log.Info().Str("id", branch.ID.String()).Msg("Branch object created, saving to DB...")

	err := s.repo.Create(branch)
	if err != nil {
		log.Error().Err(err).Str("code", req.Code).Msg("Failed to create branch in repository")
		return nil, err
	}

	log.Info().Str("code", req.Code).Msg("Branch created successfully")
	return branch, nil
}

func (s *branchService) GetByID(id uuid.UUID) (*models.Branch, error) {
	return s.repo.FindByID(id)
}

func (s *branchService) GetAll() ([]models.Branch, error) {
	return s.repo.FindAll()
}

func (s *branchService) GetActive() ([]models.Branch, error) {
	return s.repo.FindActive()
}

func (s *branchService) Update(id uuid.UUID, req *models.BranchRequest) (*models.Branch, error) {
	branch, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	branch.Code = req.Code
	branch.Name = req.Name
	branch.Description = req.Description

	err = s.repo.Update(branch)
	if err != nil {
		return nil, err
	}

	return branch, nil
}

func (s *branchService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *branchService) Count() (int64, error) {
	return s.repo.Count()
}

func (s *branchService) CreateDefaultBranches() error {
	count, err := s.repo.Count()
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	defaultBranches := []models.BranchRequest{
		{Code: "DAPUR", Name: "Dapur Pusat", Description: "Unit dapur dan produksi"},
		{Code: "OPS", Name: "Operasional", Description: "Unit operasional umum"},
		{Code: "OUTLET", Name: "Outlet", Description: "Unit penjualan"},
	}

	for _, b := range defaultBranches {
		if _, err := s.Create(&b); err != nil {
			return err
		}
	}

	return nil
}

func (s *branchService) Upsert(branch *models.Branch) error {
	return s.repo.Upsert(branch)
}

func (s *branchService) GetUpdatedAfter(since *time.Time) ([]models.Branch, error) {
	return s.repo.GetUpdatedAfter(since)
}
