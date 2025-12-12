package service

import (
	"errors"
	"shosha-finance/internal/models"
	"shosha-finance/internal/repository"
	"time"

	"github.com/google/uuid"
)

type IncomeEntryService interface {
	Create(req *models.IncomeEntryRequest) (*models.IncomeEntryResponse, error)
	Update(id uuid.UUID, req *models.IncomeEntryRequest) (*models.IncomeEntryResponse, error)
	GetByID(id uuid.UUID) (*models.IncomeEntryResponse, error)
	GetByDateRange(startDate, endDate string) ([]models.IncomeEntryResponse, error)
	GetAll(page, limit int) ([]models.IncomeEntryResponse, int64, error)
	Delete(id uuid.UUID) error
}

type incomeEntryService struct {
	repo       repository.IncomeEntryRepository
	branchRepo repository.BranchRepository
}

func NewIncomeEntryService(repo repository.IncomeEntryRepository, branchRepo repository.BranchRepository) IncomeEntryService {
	return &incomeEntryService{
		repo:       repo,
		branchRepo: branchRepo,
	}
}

func (s *incomeEntryService) Create(req *models.IncomeEntryRequest) (*models.IncomeEntryResponse, error) {
	branchID, err := uuid.Parse(req.BranchID)
	if err != nil {
		return nil, errors.New("invalid branch ID")
	}

	branch, err := s.branchRepo.GetByID(branchID)
	if err != nil {
		return nil, errors.New("branch not found")
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	// Check if entry already exists for this branch and date
	existing, err := s.repo.GetByBranchAndDate(branchID, date)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("entry already exists for this branch and date")
	}

	entry := &models.IncomeEntry{
		BranchID:      branchID,
		Date:          date,
		Omzet:         req.Omzet,
		PemasukanToru: req.PemasukanToru,
		QrisBCA:       req.QrisBCA,
		QrisBNI:       req.QrisBNI,
		QrisBRI:       req.QrisBRI,
		TransferBCA:   req.TransferBCA,
		TransferBNI:   req.TransferBNI,
		TransferBRI:   req.TransferBRI,
	}

	// Validate before saving
	if !entry.IsValid() {
		return nil, errors.New("total payments exceed pemasukan cash")
	}

	if err := s.repo.Create(entry); err != nil {
		return nil, err
	}

	return s.toResponse(entry, branch), nil
}

func (s *incomeEntryService) Update(id uuid.UUID, req *models.IncomeEntryRequest) (*models.IncomeEntryResponse, error) {
	entry, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	branchID, err := uuid.Parse(req.BranchID)
	if err != nil {
		return nil, errors.New("invalid branch ID")
	}

	branch, err := s.branchRepo.GetByID(branchID)
	if err != nil {
		return nil, errors.New("branch not found")
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	entry.BranchID = branchID
	entry.Date = date
	entry.Omzet = req.Omzet
	entry.PemasukanToru = req.PemasukanToru
	entry.QrisBCA = req.QrisBCA
	entry.QrisBNI = req.QrisBNI
	entry.QrisBRI = req.QrisBRI
	entry.TransferBCA = req.TransferBCA
	entry.TransferBNI = req.TransferBNI
	entry.TransferBRI = req.TransferBRI

	if !entry.IsValid() {
		return nil, errors.New("total payments exceed pemasukan cash")
	}

	if err := s.repo.Update(entry); err != nil {
		return nil, err
	}

	return s.toResponse(entry, branch), nil
}

func (s *incomeEntryService) GetByID(id uuid.UUID) (*models.IncomeEntryResponse, error) {
	entry, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(entry, &entry.Branch), nil
}

func (s *incomeEntryService) GetByDateRange(startDate, endDate string) ([]models.IncomeEntryResponse, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, errors.New("invalid start date format")
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, errors.New("invalid end date format")
	}

	entries, err := s.repo.GetByDateRange(start, end)
	if err != nil {
		return nil, err
	}

	responses := make([]models.IncomeEntryResponse, len(entries))
	for i, entry := range entries {
		responses[i] = *s.toResponse(&entry, &entry.Branch)
	}

	return responses, nil
}

func (s *incomeEntryService) GetAll(page, limit int) ([]models.IncomeEntryResponse, int64, error) {
	entries, total, err := s.repo.GetAll(page, limit)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]models.IncomeEntryResponse, len(entries))
	for i, entry := range entries {
		responses[i] = *s.toResponse(&entry, &entry.Branch)
	}

	return responses, total, nil
}

func (s *incomeEntryService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *incomeEntryService) toResponse(entry *models.IncomeEntry, branch *models.Branch) *models.IncomeEntryResponse {
	return &models.IncomeEntryResponse{
		ID:            entry.ID,
		BranchID:      entry.BranchID,
		BranchName:    branch.Name,
		BranchCode:    branch.Code,
		Date:          entry.Date.Format("2006-01-02"),
		Omzet:         entry.Omzet,
		PemasukanToru: entry.PemasukanToru,
		PemasukanCash: entry.PemasukanCash,
		QrisBCA:       entry.QrisBCA,
		QrisBNI:       entry.QrisBNI,
		QrisBRI:       entry.QrisBRI,
		TransferBCA:   entry.TransferBCA,
		TransferBNI:   entry.TransferBNI,
		TransferBRI:   entry.TransferBRI,
		TotalPayments: entry.GetTotalPayments(),
		CreatedAt:     entry.CreatedAt,
		UpdatedAt:     entry.UpdatedAt,
	}
}

type ExpenseEntryService interface {
	Create(req *models.ExpenseEntryRequest) (*models.ExpenseEntryResponse, error)
	Update(id uuid.UUID, req *models.ExpenseEntryRequest) (*models.ExpenseEntryResponse, error)
	GetByID(id uuid.UUID) (*models.ExpenseEntryResponse, error)
	GetByDateRange(startDate, endDate string) ([]models.ExpenseEntryResponse, error)
	GetAll(page, limit int) ([]models.ExpenseEntryResponse, int64, error)
	Delete(id uuid.UUID) error
}

type expenseEntryService struct {
	repo       repository.ExpenseEntryRepository
	branchRepo repository.BranchRepository
}

func NewExpenseEntryService(repo repository.ExpenseEntryRepository, branchRepo repository.BranchRepository) ExpenseEntryService {
	return &expenseEntryService{
		repo:       repo,
		branchRepo: branchRepo,
	}
}

func (s *expenseEntryService) Create(req *models.ExpenseEntryRequest) (*models.ExpenseEntryResponse, error) {
	branchID, err := uuid.Parse(req.BranchID)
	if err != nil {
		return nil, errors.New("invalid branch ID")
	}

	branch, err := s.branchRepo.GetByID(branchID)
	if err != nil {
		return nil, errors.New("branch not found")
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	existing, err := s.repo.GetByBranchAndDate(branchID, date)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("entry already exists for this branch and date")
	}

	entry := &models.ExpenseEntry{
		BranchID:        branchID,
		Date:            date,
		Omzet:           req.Omzet,
		PengeluaranToru: req.PengeluaranToru,
		QrisBCA:         req.QrisBCA,
		QrisBNI:         req.QrisBNI,
		QrisBRI:         req.QrisBRI,
		TransferBCA:     req.TransferBCA,
		TransferBNI:     req.TransferBNI,
		TransferBRI:     req.TransferBRI,
	}

	if !entry.IsValid() {
		return nil, errors.New("total payments exceed pengeluaran cash")
	}

	if err := s.repo.Create(entry); err != nil {
		return nil, err
	}

	return s.toResponse(entry, branch), nil
}

func (s *expenseEntryService) Update(id uuid.UUID, req *models.ExpenseEntryRequest) (*models.ExpenseEntryResponse, error) {
	entry, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	branchID, err := uuid.Parse(req.BranchID)
	if err != nil {
		return nil, errors.New("invalid branch ID")
	}

	branch, err := s.branchRepo.GetByID(branchID)
	if err != nil {
		return nil, errors.New("branch not found")
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	entry.BranchID = branchID
	entry.Date = date
	entry.Omzet = req.Omzet
	entry.PengeluaranToru = req.PengeluaranToru
	entry.QrisBCA = req.QrisBCA
	entry.QrisBNI = req.QrisBNI
	entry.QrisBRI = req.QrisBRI
	entry.TransferBCA = req.TransferBCA
	entry.TransferBNI = req.TransferBNI
	entry.TransferBRI = req.TransferBRI

	if !entry.IsValid() {
		return nil, errors.New("total payments exceed pengeluaran cash")
	}

	if err := s.repo.Update(entry); err != nil {
		return nil, err
	}

	return s.toResponse(entry, branch), nil
}

func (s *expenseEntryService) GetByID(id uuid.UUID) (*models.ExpenseEntryResponse, error) {
	entry, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(entry, &entry.Branch), nil
}

func (s *expenseEntryService) GetByDateRange(startDate, endDate string) ([]models.ExpenseEntryResponse, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, errors.New("invalid start date format")
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, errors.New("invalid end date format")
	}

	entries, err := s.repo.GetByDateRange(start, end)
	if err != nil {
		return nil, err
	}

	responses := make([]models.ExpenseEntryResponse, len(entries))
	for i, entry := range entries {
		responses[i] = *s.toResponse(&entry, &entry.Branch)
	}

	return responses, nil
}

func (s *expenseEntryService) GetAll(page, limit int) ([]models.ExpenseEntryResponse, int64, error) {
	entries, total, err := s.repo.GetAll(page, limit)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]models.ExpenseEntryResponse, len(entries))
	for i, entry := range entries {
		responses[i] = *s.toResponse(&entry, &entry.Branch)
	}

	return responses, total, nil
}

func (s *expenseEntryService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *expenseEntryService) toResponse(entry *models.ExpenseEntry, branch *models.Branch) *models.ExpenseEntryResponse {
	return &models.ExpenseEntryResponse{
		ID:              entry.ID,
		BranchID:        entry.BranchID,
		BranchName:      branch.Name,
		BranchCode:      branch.Code,
		Date:            entry.Date.Format("2006-01-02"),
		Omzet:           entry.Omzet,
		PengeluaranToru: entry.PengeluaranToru,
		PengeluaranCash: entry.PengeluaranCash,
		QrisBCA:         entry.QrisBCA,
		QrisBNI:         entry.QrisBNI,
		QrisBRI:         entry.QrisBRI,
		TransferBCA:     entry.TransferBCA,
		TransferBNI:     entry.TransferBNI,
		TransferBRI:     entry.TransferBRI,
		TotalPayments:   entry.GetTotalPayments(),
		CreatedAt:       entry.CreatedAt,
		UpdatedAt:       entry.UpdatedAt,
	}
}
