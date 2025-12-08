package response

import "github.com/gofiber/fiber/v2"

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PaginatedResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    Meta        `json:"meta"`
}

type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func Success(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Created(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func BadRequest(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
		Success: false,
		Message: message,
		Data:    nil,
	})
}

func Unauthorized(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(APIResponse{
		Success: false,
		Message: message,
		Data:    nil,
	})
}

func InternalError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
		Success: false,
		Message: message,
		Data:    nil,
	})
}

func NotFound(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(APIResponse{
		Success: false,
		Message: message,
		Data:    nil,
	})
}

func Paginated(c *fiber.Ctx, message string, data interface{}, page, limit int, total int64) error {
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return c.Status(fiber.StatusOK).JSON(PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta: Meta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}
