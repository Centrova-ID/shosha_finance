package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleManager UserRole = "manager"
	RoleStaff   UserRole = "staff"
)

type User struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	Username  string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email     *string    `gorm:"type:varchar(100);uniqueIndex" json:"email,omitempty"`
	Password  string     `gorm:"type:varchar(255);not null" json:"-"`
	Name      string     `gorm:"type:varchar(100);not null" json:"name"`
	Role      UserRole   `gorm:"type:varchar(20);not null;default:'staff'" json:"role"`
	BranchID  *uuid.UUID `gorm:"type:uuid" json:"branch_id,omitempty"`
	IsActive  bool       `gorm:"default:true" json:"is_active"`
	IsSynced  bool       `gorm:"default:false" json:"is_synced"`
	SyncedAt  *time.Time `json:"synced_at"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

type UserResponse struct {
	ID       uuid.UUID  `json:"id"`
	Username string     `json:"username"`
	Email    *string    `json:"email,omitempty"`
	Name     string     `json:"name"`
	Role     UserRole   `json:"role"`
	BranchID *uuid.UUID `json:"branch_id,omitempty"`
	IsActive bool       `json:"is_active"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Name:     u.Name,
		Role:     u.Role,
		BranchID: u.BranchID,
		IsActive: u.IsActive,
	}
}
