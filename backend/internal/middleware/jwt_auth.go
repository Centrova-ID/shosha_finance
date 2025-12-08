package middleware

import (
	"strings"

	"shosha-finance/internal/response"
	"shosha-finance/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func JWTAuth(authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Unauthorized(c, "Missing authorization header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.Unauthorized(c, "Invalid authorization format")
		}

		tokenString := parts[1]
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			log.Error().Err(err).Msg("Token validation failed")
			return response.Unauthorized(c, "Invalid or expired token")
		}

		log.Debug().Str("user_id", claims.UserID).Str("username", claims.Username).Msg("Token validated")

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			log.Error().Err(err).Str("user_id", claims.UserID).Msg("Failed to parse user ID")
			return response.Unauthorized(c, "Invalid user ID in token")
		}

		user, err := authService.GetUserByID(userID)
		if err != nil {
			return response.Unauthorized(c, "User not found")
		}

		if !user.IsActive {
			return response.Unauthorized(c, "User is not active")
		}

		c.Locals("user", user)
		c.Locals("user_id", user.ID.String())
		c.Locals("user_role", string(user.Role))

		return c.Next()
	}
}

func RequireRoles(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("user_role").(string)

		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}

		return response.Unauthorized(c, "Insufficient permissions")
	}
}
