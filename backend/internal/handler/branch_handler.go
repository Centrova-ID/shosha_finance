package handler

import (
	"shosha-finance/internal/models"
	"shosha-finance/internal/response"
	"shosha-finance/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type BranchHandler struct {
	branchService service.BranchService
}

func NewBranchHandler(branchService service.BranchService) *BranchHandler {
	return &BranchHandler{branchService: branchService}
}

func (h *BranchHandler) Create(c *fiber.Ctx) error {
	log.Info().Msg("Branch Create handler called")
	
	var req models.BranchRequest
	if err := c.BodyParser(&req); err != nil {
		log.Error().Err(err).Msg("Failed to parse request body")
		return response.BadRequest(c, "Invalid request body")
	}

	log.Info().Str("code", req.Code).Str("name", req.Name).Msg("Parsed request")

	if req.Code == "" || req.Name == "" {
		return response.BadRequest(c, "Code and name are required")
	}

	branch, err := h.branchService.Create(&req)
	if err != nil {
		log.Error().Err(err).Msg("Service Create failed")
		return response.InternalError(c, "Failed to create branch: "+err.Error())
	}

	log.Info().Str("id", branch.ID.String()).Msg("Branch created in handler")
	return response.Created(c, "Branch created successfully", branch)
}

func (h *BranchHandler) GetAll(c *fiber.Ctx) error {
	branches, err := h.branchService.GetAll()
	if err != nil {
		return response.InternalError(c, "Failed to get branches")
	}

	return response.Success(c, "Branches retrieved successfully", branches)
}

func (h *BranchHandler) GetActive(c *fiber.Ctx) error {
	branches, err := h.branchService.GetActive()
	if err != nil {
		return response.InternalError(c, "Failed to get branches")
	}

	return response.Success(c, "Branches retrieved successfully", branches)
}

func (h *BranchHandler) GetByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return response.BadRequest(c, "Invalid branch ID")
	}

	branch, err := h.branchService.GetByID(id)
	if err != nil {
		return response.NotFound(c, "Branch not found")
	}

	return response.Success(c, "Branch retrieved successfully", branch)
}

func (h *BranchHandler) Update(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return response.BadRequest(c, "Invalid branch ID")
	}

	var req models.BranchRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if req.Code == "" || req.Name == "" {
		return response.BadRequest(c, "Code and name are required")
	}

	branch, err := h.branchService.Update(id, &req)
	if err != nil {
		return response.InternalError(c, "Failed to update branch")
	}

	return response.Success(c, "Branch updated successfully", branch)
}

func (h *BranchHandler) Delete(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return response.BadRequest(c, "Invalid branch ID")
	}

	err = h.branchService.Delete(id)
	if err != nil {
		return response.InternalError(c, "Failed to delete branch")
	}

	return response.Success(c, "Branch deleted successfully", nil)
}
