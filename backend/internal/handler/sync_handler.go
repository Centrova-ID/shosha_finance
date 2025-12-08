package handler

import (
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

type SyncPushRequest struct {
	Transactions []models.Transaction `json:"transactions"`
}

type SyncPushResponse struct {
	ReceivedCount int         `json:"received_count"`
	SyncedIDs     []uuid.UUID `json:"synced_ids"`
}

func NewSyncHandler(txService service.TransactionService, branchService service.BranchService) *SyncHandler {
	return &SyncHandler{
		txService:     txService,
		branchService: branchService,
	}
}

func (h *SyncHandler) Push(c *fiber.Ctx) error {
	var req SyncPushRequest

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if len(req.Transactions) == 0 {
		return response.BadRequest(c, "No transactions to sync")
	}

	syncedIDs := make([]uuid.UUID, 0)

	for _, tx := range req.Transactions {
		existing, _ := h.txService.GetByID(tx.ID)
		if existing != nil {
			syncedIDs = append(syncedIDs, tx.ID)
			continue
		}

		newTx := &models.TransactionRequest{
			Type:        tx.Type,
			Category:    tx.Category,
			Amount:      tx.Amount,
			Description: tx.Description,
		}

		created, err := h.txService.Create(newTx, tx.BranchID)
		if err != nil {
			continue
		}

		syncedIDs = append(syncedIDs, created.ID)
	}

	return response.Success(c, "Sync completed", SyncPushResponse{
		ReceivedCount: len(req.Transactions),
		SyncedIDs:     syncedIDs,
	})
}
