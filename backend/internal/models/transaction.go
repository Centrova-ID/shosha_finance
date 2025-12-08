package models

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionTypeIN  TransactionType = "IN"
	TransactionTypeOUT TransactionType = "OUT"
)

type Transaction struct {
	ID          uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	BranchID    uuid.UUID       `gorm:"type:uuid;index;not null" json:"branch_id"`
	Type        TransactionType `gorm:"type:varchar(10);not null" json:"type"`
	Category    string          `gorm:"type:varchar(50);not null" json:"category"`
	Amount      int64           `gorm:"not null" json:"amount"`
	Description string          `gorm:"type:text" json:"description"`
	CreatedAt   time.Time       `gorm:"autoCreateTime" json:"created_at"`
	IsSynced    bool            `gorm:"default:false" json:"is_synced"`
	SyncedAt    *time.Time      `json:"synced_at"`
	Branch      Branch          `gorm:"foreignKey:BranchID" json:"-"`
}

func (t *Transaction) BeforeCreate() error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type TransactionRequest struct {
	Type        TransactionType `json:"type" validate:"required,oneof=IN OUT"`
	Category    string          `json:"category" validate:"required"`
	Amount      int64           `json:"amount" validate:"required,gt=0"`
	Description string          `json:"description"`
}

type TransactionResponse struct {
	ID          uuid.UUID       `json:"id"`
	BranchID    uuid.UUID       `json:"branch_id"`
	Type        TransactionType `json:"type"`
	Category    string          `json:"category"`
	Amount      int64           `json:"amount"`
	Description string          `json:"description"`
	CreatedAt   time.Time       `json:"created_at"`
	IsSynced    bool            `json:"is_synced"`
	SyncedAt    *time.Time      `json:"synced_at"`
}
