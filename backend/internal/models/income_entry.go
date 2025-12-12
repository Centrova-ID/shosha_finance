package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// IncomeEntry represents a daily income entry for a branch
type IncomeEntry struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	BranchID      uuid.UUID  `gorm:"type:uuid;index;not null" json:"branch_id"`
	Date          time.Time  `gorm:"type:date;not null;index" json:"date"`
	Omzet         float64    `gorm:"type:decimal(15,2);not null" json:"omzet"`
	PemasukanToru float64    `gorm:"type:decimal(15,2);not null" json:"pemasukan_toru"`
	PemasukanCash float64    `gorm:"type:decimal(15,2);not null" json:"pemasukan_cash"` // Auto-calculated: Omzet - PemasukanToru
	QrisBCA       float64    `gorm:"type:decimal(15,2);not null;default:0" json:"qris_bca"`
	QrisBNI       float64    `gorm:"type:decimal(15,2);not null;default:0" json:"qris_bni"`
	QrisBRI       float64    `gorm:"type:decimal(15,2);not null;default:0" json:"qris_bri"`
	TransferBCA   float64    `gorm:"type:decimal(15,2);not null;default:0" json:"transfer_bca"`
	TransferBNI   float64    `gorm:"type:decimal(15,2);not null;default:0" json:"transfer_bni"`
	TransferBRI   float64    `gorm:"type:decimal(15,2);not null;default:0" json:"transfer_bri"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	IsSynced      bool       `gorm:"default:false" json:"is_synced"`
	SyncedAt      *time.Time `json:"synced_at"`
	Branch        Branch     `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
}

func (i *IncomeEntry) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	// Auto-calculate PemasukanCash
	i.PemasukanCash = i.Omzet - i.PemasukanToru
	return nil
}

func (i *IncomeEntry) BeforeUpdate(tx *gorm.DB) error {
	// Auto-calculate PemasukanCash on update
	i.PemasukanCash = i.Omzet - i.PemasukanToru
	return nil
}

// GetTotalPayments returns sum of all payment methods
func (i *IncomeEntry) GetTotalPayments() float64 {
	return i.QrisBCA + i.QrisBNI + i.QrisBRI + i.TransferBCA + i.TransferBNI + i.TransferBRI
}

// IsValid checks if total payments doesn't exceed PemasukanCash
func (i *IncomeEntry) IsValid() bool {
	return i.GetTotalPayments() <= i.PemasukanCash
}

type IncomeEntryRequest struct {
	BranchID      string  `json:"branch_id" validate:"required"`
	Date          string  `json:"date" validate:"required"` // Format: YYYY-MM-DD
	Omzet         float64 `json:"omzet" validate:"required,gt=0"`
	PemasukanToru float64 `json:"pemasukan_toru" validate:"required,gte=0"`
	QrisBCA       float64 `json:"qris_bca" validate:"gte=0"`
	QrisBNI       float64 `json:"qris_bni" validate:"gte=0"`
	QrisBRI       float64 `json:"qris_bri" validate:"gte=0"`
	TransferBCA   float64 `json:"transfer_bca" validate:"gte=0"`
	TransferBNI   float64 `json:"transfer_bni" validate:"gte=0"`
	TransferBRI   float64 `json:"transfer_bri" validate:"gte=0"`
}

type IncomeEntryResponse struct {
	ID            uuid.UUID `json:"id"`
	BranchID      uuid.UUID `json:"branch_id"`
	BranchName    string    `json:"branch_name"`
	BranchCode    string    `json:"branch_code"`
	Date          string    `json:"date"`
	Omzet         float64   `json:"omzet"`
	PemasukanToru float64   `json:"pemasukan_toru"`
	PemasukanCash float64   `json:"pemasukan_cash"`
	QrisBCA       float64   `json:"qris_bca"`
	QrisBNI       float64   `json:"qris_bni"`
	QrisBRI       float64   `json:"qris_bri"`
	TransferBCA   float64   `json:"transfer_bca"`
	TransferBNI   float64   `json:"transfer_bni"`
	TransferBRI   float64   `json:"transfer_bri"`
	TotalPayments float64   `json:"total_payments"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ExpenseEntry represents a daily expense entry for a branch
type ExpenseEntry struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	BranchID        uuid.UUID  `gorm:"type:uuid;index;not null" json:"branch_id"`
	Date            time.Time  `gorm:"type:date;not null;index" json:"date"`
	Omzet           float64    `gorm:"type:decimal(15,2);not null" json:"omzet"`
	PengeluaranToru float64    `gorm:"type:decimal(15,2);not null" json:"pengeluaran_toru"`
	PengeluaranCash float64    `gorm:"type:decimal(15,2);not null" json:"pengeluaran_cash"` // Auto-calculated
	QrisBCA         float64    `gorm:"type:decimal(15,2);not null;default:0" json:"qris_bca"`
	QrisBNI         float64    `gorm:"type:decimal(15,2);not null;default:0" json:"qris_bni"`
	QrisBRI         float64    `gorm:"type:decimal(15,2);not null;default:0" json:"qris_bri"`
	TransferBCA     float64    `gorm:"type:decimal(15,2);not null;default:0" json:"transfer_bca"`
	TransferBNI     float64    `gorm:"type:decimal(15,2);not null;default:0" json:"transfer_bni"`
	TransferBRI     float64    `gorm:"type:decimal(15,2);not null;default:0" json:"transfer_bri"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	IsSynced        bool       `gorm:"default:false" json:"is_synced"`
	SyncedAt        *time.Time `json:"synced_at"`
	Branch          Branch     `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
}

func (e *ExpenseEntry) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	e.PengeluaranCash = e.Omzet - e.PengeluaranToru
	return nil
}

func (e *ExpenseEntry) BeforeUpdate(tx *gorm.DB) error {
	e.PengeluaranCash = e.Omzet - e.PengeluaranToru
	return nil
}

func (e *ExpenseEntry) GetTotalPayments() float64 {
	return e.QrisBCA + e.QrisBNI + e.QrisBRI + e.TransferBCA + e.TransferBNI + e.TransferBRI
}

func (e *ExpenseEntry) IsValid() bool {
	return e.GetTotalPayments() <= e.PengeluaranCash
}

type ExpenseEntryRequest struct {
	BranchID        string  `json:"branch_id" validate:"required"`
	Date            string  `json:"date" validate:"required"`
	Omzet           float64 `json:"omzet" validate:"required,gt=0"`
	PengeluaranToru float64 `json:"pengeluaran_toru" validate:"required,gte=0"`
	QrisBCA         float64 `json:"qris_bca" validate:"gte=0"`
	QrisBNI         float64 `json:"qris_bni" validate:"gte=0"`
	QrisBRI         float64 `json:"qris_bri" validate:"gte=0"`
	TransferBCA     float64 `json:"transfer_bca" validate:"gte=0"`
	TransferBNI     float64 `json:"transfer_bni" validate:"gte=0"`
	TransferBRI     float64 `json:"transfer_bri" validate:"gte=0"`
}

type ExpenseEntryResponse struct {
	ID              uuid.UUID `json:"id"`
	BranchID        uuid.UUID `json:"branch_id"`
	BranchName      string    `json:"branch_name"`
	BranchCode      string    `json:"branch_code"`
	Date            string    `json:"date"`
	Omzet           float64   `json:"omzet"`
	PengeluaranToru float64   `json:"pengeluaran_toru"`
	PengeluaranCash float64   `json:"pengeluaran_cash"`
	QrisBCA         float64   `json:"qris_bca"`
	QrisBNI         float64   `json:"qris_bni"`
	QrisBRI         float64   `json:"qris_bri"`
	TransferBCA     float64   `json:"transfer_bca"`
	TransferBNI     float64   `json:"transfer_bni"`
	TransferBRI     float64   `json:"transfer_bri"`
	TotalPayments   float64   `json:"total_payments"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
