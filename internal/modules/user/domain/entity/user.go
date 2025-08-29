package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Username  string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	FirstName string         `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName  string         `gorm:"type:varchar(100);not null" json:"last_name"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	IsAdmin   bool           `gorm:"default:false" json:"is_admin"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for User entity
func (User) TableName() string {
	return "users"
}

// BeforeCreate hook to set UUID before creating
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// UserStatus represents the status of a user
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
)

// GetStatus returns the user's current status
func (u *User) GetStatus() UserStatus {
	if !u.IsActive {
		return UserStatusInactive
	}
	if u.DeletedAt.Valid {
		return UserStatusSuspended
	}
	return UserStatusActive
}

// GetFullName returns the user's full name
func (u *User) GetFullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return u.Username
	}
	return u.FirstName + " " + u.LastName
}

// Validate validates the user entity
func (u *User) Validate() error {
	if u.Email == "" {
		return ErrInvalidEmail
	}
	if u.Username == "" {
		return ErrInvalidUsername
	}
	if len(u.Password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}