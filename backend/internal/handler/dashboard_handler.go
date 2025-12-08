package handler

import (
	"shosha-finance/internal/response"
	"shosha-finance/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DashboardHandler struct {
	txService service.TransactionService
	branchID  uuid.UUID
}

func NewDashboardHandler(txService service.TransactionService, branchID string) *DashboardHandler {
	bid, _ := uuid.Parse(branchID)
	return &DashboardHandler{
		txService: txService,
		branchID:  bid,
	}
}

func (h *DashboardHandler) GetSummary(c *fiber.Ctx) error {
	summary, err := h.txService.GetDashboardSummary(h.branchID)
	if err != nil {
		return response.InternalError(c, "Failed to get dashboard summary")
	}

	return response.Success(c, "Success", summary)
}
