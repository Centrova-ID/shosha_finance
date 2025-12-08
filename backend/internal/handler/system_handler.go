package handler

import (
	"net/http"
	"time"

	"shosha-finance/internal/response"
	"shosha-finance/internal/service"

	"github.com/gofiber/fiber/v2"
)

type SystemHandler struct {
	txService service.TransactionService
	cloudURL  string
}

type SystemStatus struct {
	UnsyncedCount int64  `json:"unsynced_count"`
	Status        string `json:"status"`
	Timestamp     string `json:"timestamp"`
}

func NewSystemHandler(txService service.TransactionService, cloudURL string) *SystemHandler {
	return &SystemHandler{
		txService: txService,
		cloudURL:  cloudURL,
	}
}

func (h *SystemHandler) GetStatus(c *fiber.Ctx) error {
	count, err := h.txService.GetUnsyncedCount()
	if err != nil {
		return response.InternalError(c, "Failed to get sync status")
	}

	status := "offline"
	if h.checkOnline() {
		status = "online"
	}

	return response.Success(c, "Success", SystemStatus{
		UnsyncedCount: count,
		Status:        status,
		Timestamp:     time.Now().Format(time.RFC3339),
	})
}

func (h *SystemHandler) checkOnline() bool {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("https://www.google.com")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (h *SystemHandler) HealthCheck(c *fiber.Ctx) error {
	return response.Success(c, "OK", map[string]string{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
