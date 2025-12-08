package handler

import (
	"strconv"

	"shosha-finance/internal/models"
	"shosha-finance/internal/response"
	"shosha-finance/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(svc service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: svc}
}

func (h *TransactionHandler) Create(c *fiber.Ctx) error {
	var req models.TransactionRequest

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if req.BranchID == "" {
		return response.BadRequest(c, "Branch ID is required")
	}

	if req.Type != models.TransactionTypeIN && req.Type != models.TransactionTypeOUT {
		return response.BadRequest(c, "Type must be IN or OUT")
	}

	if req.Amount <= 0 {
		return response.BadRequest(c, "Amount must be greater than 0")
	}

	if req.Category == "" {
		return response.BadRequest(c, "Category is required")
	}

	tx, err := h.service.Create(&req)
	if err != nil {
		return response.InternalError(c, "Failed to create transaction")
	}

	return response.Created(c, "Transaction created successfully", tx)
}

func (h *TransactionHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	transactions, total, err := h.service.GetAll(page, limit)
	if err != nil {
		return response.InternalError(c, "Failed to get transactions")
	}

	return response.Paginated(c, "Success", transactions, page, limit, total)
}

func (h *TransactionHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid transaction ID")
	}

	tx, err := h.service.GetByID(id)
	if err != nil {
		return response.NotFound(c, "Transaction not found")
	}

	return response.Success(c, "Success", tx)
}
