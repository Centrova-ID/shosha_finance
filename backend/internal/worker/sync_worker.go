package worker

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"shosha-finance/internal/config"
	"shosha-finance/internal/models"
	"shosha-finance/internal/service"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type SyncWorker struct {
	txService service.TransactionService
	cfg       *config.Config
	client    *http.Client
	stopChan  chan struct{}
}

type SyncPushRequest struct {
	Transactions []models.Transaction `json:"transactions"`
}

type SyncPushResponse struct {
	Success       bool        `json:"success"`
	ReceivedCount int         `json:"received_count"`
	SyncedIDs     []uuid.UUID `json:"synced_ids"`
}

func NewSyncWorker(txService service.TransactionService, cfg *config.Config) *SyncWorker {
	return &SyncWorker{
		txService: txService,
		cfg:       cfg,
		client:    &http.Client{Timeout: 30 * time.Second},
		stopChan:  make(chan struct{}),
	}
}

func (w *SyncWorker) Start() {
	log.Info().Int("interval", w.cfg.SyncInterval).Msg("Starting sync worker")

	go func() {
		ticker := time.NewTicker(time.Duration(w.cfg.SyncInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				w.sync()
			case <-w.stopChan:
				log.Info().Msg("Sync worker stopped")
				return
			}
		}
	}()
}

func (w *SyncWorker) Stop() {
	close(w.stopChan)
}

func (w *SyncWorker) sync() {
	if !w.isOnline() {
		log.Debug().Msg("Offline, skipping sync")
		return
	}

	transactions, err := w.txService.GetUnsynced(50)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get unsynced transactions")
		return
	}

	if len(transactions) == 0 {
		log.Debug().Msg("No unsynced transactions")
		return
	}

	log.Info().Int("count", len(transactions)).Msg("Syncing transactions")

	syncedIDs, err := w.pushToCloud(transactions)
	if err != nil {
		log.Error().Err(err).Msg("Failed to push to cloud")
		return
	}

	if len(syncedIDs) > 0 {
		err = w.txService.MarkAsSynced(syncedIDs)
		if err != nil {
			log.Error().Err(err).Msg("Failed to mark transactions as synced")
			return
		}
		log.Info().Int("count", len(syncedIDs)).Msg("Transactions synced successfully")
	}
}

func (w *SyncWorker) pushToCloud(transactions []models.Transaction) ([]uuid.UUID, error) {
	reqBody := SyncPushRequest{Transactions: transactions}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := w.cfg.CloudAPIURL + "/api/v1/sync/push"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+w.cfg.BranchAPIKey)

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn().Int("status", resp.StatusCode).Msg("Cloud API returned non-200 status")
		return nil, nil
	}

	var syncResp struct {
		Success bool `json:"success"`
		Data    struct {
			SyncedIDs []uuid.UUID `json:"synced_ids"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&syncResp); err != nil {
		return nil, err
	}

	return syncResp.Data.SyncedIDs, nil
}

func (w *SyncWorker) isOnline() bool {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("https://www.google.com")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
