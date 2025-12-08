package handler

import (
	"time"

	"shosha-finance/internal/repository"
	"shosha-finance/internal/response"
	"shosha-finance/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DashboardHandler struct {
	txService service.TransactionService
}

func NewDashboardHandler(txService service.TransactionService) *DashboardHandler {
	return &DashboardHandler{txService: txService}
}

func (h *DashboardHandler) GetSummary(c *fiber.Ctx) error {
	filter := &repository.DashboardFilter{}

	// Optional branch_id filter
	branchIDParam := c.Query("branch_id")
	if branchIDParam != "" {
		id, err := uuid.Parse(branchIDParam)
		if err != nil {
			return response.BadRequest(c, "Invalid branch_id")
		}
		filter.BranchID = &id
	}

	// Optional date filter (format: YYYY-MM-DD)
	dateParam := c.Query("date")
	if dateParam != "" {
		date, err := time.Parse("2006-01-02", dateParam)
		if err != nil {
			return response.BadRequest(c, "Invalid date format. Use YYYY-MM-DD")
		}
		// Start of day
		startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		// End of day (start of next day)
		endOfDay := startOfDay.AddDate(0, 0, 1)
		filter.StartDate = &startOfDay
		filter.EndDate = &endOfDay
	}

	summary, err := h.txService.GetDashboardSummary(filter)
	if err != nil {
		return response.InternalError(c, "Failed to get dashboard summary")
	}

	return response.Success(c, "Success", summary)
}
