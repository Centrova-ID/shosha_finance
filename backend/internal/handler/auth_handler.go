package handler

import (
	"shosha-finance/internal/models"
	"shosha-finance/internal/response"
	"shosha-finance/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type LoginResponse struct {
	User  models.UserResponse `json:"user"`
	Token string              `json:"token"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if req.Identifier == "" || req.Password == "" {
		return response.BadRequest(c, "Identifier and password are required")
	}

	user, token, err := h.authService.Login(req.Identifier, req.Password)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			return response.Unauthorized(c, "Invalid credentials")
		case service.ErrUserNotActive:
			return response.Unauthorized(c, "User is not active")
		default:
			return response.InternalError(c, "Failed to authenticate")
		}
	}

	return response.Success(c, "Login successful", LoginResponse{
		User:  user.ToResponse(),
		Token: token,
	})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	return response.Success(c, "User retrieved successfully", user.ToResponse())
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	return response.Success(c, "Logout successful", nil)
}
