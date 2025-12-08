package handler

import (
	"time"

	"shosha-finance/internal/response"
	"shosha-finance/internal/service"
	"shosha-finance/internal/worker"

	"github.com/gofiber/fiber/v2"
)

type SystemHandler struct {
	txService  service.TransactionService
	syncWorker *worker.SyncWorker
}

func NewSystemHandler(txService service.TransactionService, syncWorker *worker.SyncWorker) *SystemHandler {
	return &SystemHandler{
		txService:  txService,
		syncWorker: syncWorker,
	}
}

type SystemStatus struct {
	Status        string `json:"status"`
	UnsyncedCount int64  `json:"unsynced_count"`
	Timestamp     string `json:"timestamp"`
}

func (h *SystemHandler) GetStatus(c *fiber.Ctx) error {
	status := "offline"
	if h.syncWorker != nil && h.syncWorker.IsOnline() {
		status = "online"
	}

	unsyncedCount, _ := h.txService.GetUnsyncedCount()

	return response.Success(c, "Success", SystemStatus{
		Status:        status,
		UnsyncedCount: unsyncedCount,
		Timestamp:     time.Now().Format(time.RFC3339),
	})
}

func (h *SystemHandler) HealthCheck(c *fiber.Ctx) error {
	return response.Success(c, "OK", map[string]string{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
