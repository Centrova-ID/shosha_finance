package models

import (
	"time"

	"github.com/google/uuid"
)

type Branch struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Code      string    `gorm:"type:varchar(10);uniqueIndex;not null" json:"code"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	APIKey    string    `gorm:"type:varchar(255);not null" json:"api_key"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (b *Branch) BeforeCreate() error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
