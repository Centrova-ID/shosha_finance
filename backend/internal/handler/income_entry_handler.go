package handler

import (
	"shosha-finance/internal/models"
	"shosha-finance/internal/response"
	"shosha-finance/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type IncomeEntryHandler struct {
	service service.IncomeEntryService
}

func NewIncomeEntryHandler(svc service.IncomeEntryService) *IncomeEntryHandler {
	return &IncomeEntryHandler{service: svc}
}

func (h *IncomeEntryHandler) Create(c *fiber.Ctx) error {
	var req models.IncomeEntryRequest

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	entry, err := h.service.Create(&req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, "Income entry created successfully", entry)
}

func (h *IncomeEntryHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid ID")
	}

	var req models.IncomeEntryRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	entry, err := h.service.Update(id, &req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Income entry updated successfully", entry)
}

func (h *IncomeEntryHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid ID")
	}

	entry, err := h.service.GetByID(id)
	if err != nil {
		return response.NotFound(c, "Income entry not found")
	}

	return response.Success(c, "Success", entry)
}

func (h *IncomeEntryHandler) GetByDateRange(c *fiber.Ctx) error {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		return response.BadRequest(c, "start_date and end_date are required")
	}

	entries, err := h.service.GetByDateRange(startDate, endDate)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Success", entries)
}

func (h *IncomeEntryHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	entries, total, err := h.service.GetAll(page, limit)
	if err != nil {
		return response.InternalError(c, "Failed to get income entries")
	}

	return response.Paginated(c, "Success", entries, page, limit, total)
}

func (h *IncomeEntryHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid ID")
	}

	if err := h.service.Delete(id); err != nil {
		return response.InternalError(c, "Failed to delete income entry")
	}

	return response.Success(c, "Income entry deleted successfully", nil)
}

type ExpenseEntryHandler struct {
	service service.ExpenseEntryService
}

func NewExpenseEntryHandler(svc service.ExpenseEntryService) *ExpenseEntryHandler {
	return &ExpenseEntryHandler{service: svc}
}

func (h *ExpenseEntryHandler) Create(c *fiber.Ctx) error {
	var req models.ExpenseEntryRequest

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	entry, err := h.service.Create(&req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, "Expense entry created successfully", entry)
}

func (h *ExpenseEntryHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid ID")
	}

	var req models.ExpenseEntryRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	entry, err := h.service.Update(id, &req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Expense entry updated successfully", entry)
}

func (h *ExpenseEntryHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid ID")
	}

	entry, err := h.service.GetByID(id)
	if err != nil {
		return response.NotFound(c, "Expense entry not found")
	}

	return response.Success(c, "Success", entry)
}

func (h *ExpenseEntryHandler) GetByDateRange(c *fiber.Ctx) error {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		return response.BadRequest(c, "start_date and end_date are required")
	}

	entries, err := h.service.GetByDateRange(startDate, endDate)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Success", entries)
}

func (h *ExpenseEntryHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	entries, total, err := h.service.GetAll(page, limit)
	if err != nil {
		return response.InternalError(c, "Failed to get expense entries")
	}

	return response.Paginated(c, "Success", entries, page, limit, total)
}

func (h *ExpenseEntryHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.BadRequest(c, "Invalid ID")
	}

	if err := h.service.Delete(id); err != nil {
		return response.InternalError(c, "Failed to delete expense entry")
	}

	return response.Success(c, "Expense entry deleted successfully", nil)
}
