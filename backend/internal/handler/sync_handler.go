package handler

import (
	"time"

	"shosha-finance/internal/models"
	"shosha-finance/internal/response"
	"shosha-finance/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SyncHandler struct {
	txService     service.TransactionService
	branchService service.BranchService
}

func NewSyncHandler(txService service.TransactionService, branchService service.BranchService) *SyncHandler {
	return &SyncHandler{
		txService:     txService,
		branchService: branchService,
	}
}

type SyncPushRequest struct {
	Branches     []models.Branch      `json:"branches"`
	Transactions []models.Transaction `json:"transactions"`
}

type SyncPushResponse struct {
	Branches     []uuid.UUID `json:"branches"`
	Transactions []uuid.UUID `json:"transactions"`
}

type SyncPullResponse struct {
	Branches     []models.Branch      `json:"branches"`
	Transactions []models.Transaction `json:"transactions"`
	LastSyncAt   string               `json:"last_sync_at"`
}

// Push - receive data from local app and save to cloud
func (h *SyncHandler) Push(c *fiber.Ctx) error {
	var req SyncPushRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	syncedBranches := []uuid.UUID{}
	syncedTransactions := []uuid.UUID{}

	// Upsert branches
	for _, branch := range req.Branches {
		err := h.branchService.Upsert(&branch)
		if err != nil {
			continue
		}
		syncedBranches = append(syncedBranches, branch.ID)
	}

	// Upsert transactions
	for _, tx := range req.Transactions {
		err := h.txService.Upsert(&tx)
		if err != nil {
			continue
		}
		syncedTransactions = append(syncedTransactions, tx.ID)
	}

	return response.Success(c, "Data synced successfully", SyncPushResponse{
		Branches:     syncedBranches,
		Transactions: syncedTransactions,
	})
}

// Pull - send latest data to local app
func (h *SyncHandler) Pull(c *fiber.Ctx) error {
	// Get last sync time from query param (optional)
	lastSyncParam := c.Query("last_sync", "")
	var lastSync *time.Time
	if lastSyncParam != "" {
		t, err := time.Parse(time.RFC3339, lastSyncParam)
		if err == nil {
			lastSync = &t
		}
	}

	// Get branches updated after lastSync
	branches, err := h.branchService.GetUpdatedAfter(lastSync)
	if err != nil {
		return response.InternalError(c, "Failed to get branches")
	}

	// Get transactions updated after lastSync
	transactions, err := h.txService.GetUpdatedAfter(lastSync)
	if err != nil {
		return response.InternalError(c, "Failed to get transactions")
	}

	return response.Success(c, "Data retrieved successfully", SyncPullResponse{
		Branches:     branches,
		Transactions: transactions,
		LastSyncAt:   time.Now().Format(time.RFC3339),
	})
}
