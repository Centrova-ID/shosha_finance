package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	Amount      float64         `gorm:"type:decimal(15,2);not null" json:"amount"`
	Description string          `gorm:"type:text" json:"description"`
	CreatedAt   time.Time       `gorm:"autoCreateTime" json:"created_at"`
	IsSynced    bool            `gorm:"default:false" json:"is_synced"`
	SyncedAt    *time.Time      `json:"synced_at"`
	Branch      Branch          `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type TransactionRequest struct {
	BranchID    string          `json:"branch_id" validate:"required"`
	Type        TransactionType `json:"type" validate:"required,oneof=IN OUT"`
	Category    string          `json:"category" validate:"required"`
	Amount      float64         `json:"amount" validate:"required,gt=0"`
	Description string          `json:"description"`
}

type TransactionResponse struct {
	ID          uuid.UUID       `json:"id"`
	BranchID    uuid.UUID       `json:"branch_id"`
	Type        TransactionType `json:"type"`
	Category    string          `json:"category"`
	Amount      float64         `json:"amount"`
	Description string          `json:"description"`
	CreatedAt   time.Time       `json:"created_at"`
}
