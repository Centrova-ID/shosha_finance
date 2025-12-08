package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Branch struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	Code        string     `gorm:"type:varchar(20);uniqueIndex;not null" json:"code"`
	Name        string     `gorm:"type:varchar(100);not null" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	IsSynced    bool       `gorm:"default:false" json:"is_synced"`
	SyncedAt    *time.Time `json:"synced_at"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (b *Branch) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

type BranchRequest struct {
	Code        string `json:"code" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}
