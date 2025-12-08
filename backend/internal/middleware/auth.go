package middleware

import (
	"strings"

	"shosha-finance/internal/response"
	"shosha-finance/internal/service"

	"github.com/gofiber/fiber/v2"
)

func BranchAuth(branchService service.BranchService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Unauthorized(c, "Missing authorization header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.Unauthorized(c, "Invalid authorization format")
		}

		apiKey := parts[1]
		branch, err := branchService.GetByAPIKey(apiKey)
		if err != nil {
			return response.Unauthorized(c, "Invalid API key")
		}

		c.Locals("branch", branch)
		c.Locals("branch_id", branch.ID.String())

		return c.Next()
	}
}
