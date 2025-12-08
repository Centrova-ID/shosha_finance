package worker

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"shosha-finance/internal/config"
	"shosha-finance/internal/models"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type SyncWorker struct {
	db       *gorm.DB
	cfg      *config.Config
	client   *http.Client
	stopChan chan struct{}
	isOnline bool
}

type SyncPushRequest struct {
	Branches     []models.Branch     `json:"branches"`
	Transactions []models.Transaction `json:"transactions"`
}

type SyncPullResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Branches     []models.Branch     `json:"branches"`
		Transactions []models.Transaction `json:"transactions"`
		LastSyncAt   string               `json:"last_sync_at"`
	} `json:"data"`
}

type SyncPushResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Branches     []uuid.UUID `json:"branches"`
		Transactions []uuid.UUID `json:"transactions"`
	} `json:"data"`
}

func NewSyncWorker(db *gorm.DB, cfg *config.Config) *SyncWorker {
	return &SyncWorker{
		db:       db,
		cfg:      cfg,
		client:   &http.Client{Timeout: 30 * time.Second},
		stopChan: make(chan struct{}),
		isOnline: false,
	}
}

func (w *SyncWorker) Start() {
	log.Info().Int("interval", w.cfg.SyncInterval).Msg("Starting sync worker")

	go func() {
		// Initial sync
		w.sync()

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

func (w *SyncWorker) IsOnline() bool {
	return w.isOnline
}

func (w *SyncWorker) sync() {
	online := w.checkOnline()
	w.isOnline = online

	if !online {
		log.Debug().Msg("Offline, skipping sync")
		return
	}

	// Pull first (get latest data from cloud)
	if err := w.pull(); err != nil {
		log.Error().Err(err).Msg("Failed to pull from cloud")
	}

	// Then push (send local unsynced data to cloud)
	if err := w.push(); err != nil {
		log.Error().Err(err).Msg("Failed to push to cloud")
	}
}

func (w *SyncWorker) pull() error {
	url := w.cfg.CloudAPIURL + "/api/v1/sync/pull"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+w.cfg.BranchAPIKey)

	resp, err := w.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn().Int("status", resp.StatusCode).Msg("Cloud API pull returned non-200 status")
		return nil
	}

	var pullResp SyncPullResponse
	if err := json.NewDecoder(resp.Body).Decode(&pullResp); err != nil {
		return err
	}

	if !pullResp.Success {
		return nil
	}

	// Upsert branches
	for _, branch := range pullResp.Data.Branches {
		branch.IsSynced = true
		now := time.Now()
		branch.SyncedAt = &now
		w.db.Save(&branch)
	}

	// Upsert transactions
	for _, tx := range pullResp.Data.Transactions {
		tx.IsSynced = true
		now := time.Now()
		tx.SyncedAt = &now
		w.db.Save(&tx)
	}

	log.Info().
		Int("branches", len(pullResp.Data.Branches)).
		Int("transactions", len(pullResp.Data.Transactions)).
		Msg("Pulled data from cloud")

	return nil
}

func (w *SyncWorker) push() error {
	// Get unsynced branches
	var branches []models.Branch
	w.db.Where("is_synced = ?", false).Find(&branches)

	// Get unsynced transactions
	var transactions []models.Transaction
	w.db.Where("is_synced = ?", false).Limit(100).Find(&transactions)

	log.Info().
		Int("unsynced_branches", len(branches)).
		Int("unsynced_transactions", len(transactions)).
		Msg("Checking unsynced data")

	if len(branches) == 0 && len(transactions) == 0 {
		log.Debug().Msg("No unsynced data to push")
		return nil
	}

	reqBody := SyncPushRequest{
		Branches:     branches,
		Transactions: transactions,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	url := w.cfg.CloudAPIURL + "/api/v1/sync/push"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+w.cfg.BranchAPIKey)

	resp, err := w.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn().Int("status", resp.StatusCode).Msg("Cloud API push returned non-200 status")
		return nil
	}

	var pushResp SyncPushResponse
	if err := json.NewDecoder(resp.Body).Decode(&pushResp); err != nil {
		log.Error().Err(err).Msg("Failed to decode push response")
		return err
	}

	log.Info().
		Bool("success", pushResp.Success).
		Int("synced_branches", len(pushResp.Data.Branches)).
		Int("synced_transactions", len(pushResp.Data.Transactions)).
		Msg("Push response received")

	if !pushResp.Success {
		log.Warn().Msg("Push response success=false")
		return nil
	}

	now := time.Now()

	// Mark branches as synced
	if len(pushResp.Data.Branches) > 0 {
		w.db.Model(&models.Branch{}).
			Where("id IN ?", pushResp.Data.Branches).
			Updates(map[string]interface{}{
				"is_synced": true,
				"synced_at": now,
			})
	}

	// Mark transactions as synced
	if len(pushResp.Data.Transactions) > 0 {
		w.db.Model(&models.Transaction{}).
			Where("id IN ?", pushResp.Data.Transactions).
			Updates(map[string]interface{}{
				"is_synced": true,
				"synced_at": now,
			})
	}

	log.Info().
		Int("branches", len(pushResp.Data.Branches)).
		Int("transactions", len(pushResp.Data.Transactions)).
		Msg("Pushed data to cloud")

	return nil
}

func (w *SyncWorker) checkOnline() bool {
	if w.cfg.CloudAPIURL == "" {
		return false
	}

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(w.cfg.CloudAPIURL + "/api/v1/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (w *SyncWorker) GetUnsyncedCount() (int64, error) {
	var count int64
	err := w.db.Model(&models.Transaction{}).Where("is_synced = ?", false).Count(&count).Error
	return count, err
}
