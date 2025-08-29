package dto

import (
	"time"

	"github.com/gigi434/sample-grpc-server/internal/modules/user/domain/entity"
	"github.com/google/uuid"
)

// CreateUserDTO represents the data transfer object for creating a user
type CreateUserDTO struct {
	Email     string
	Username  string
	Password  string
	FirstName string
	LastName  string
	IsActive  bool
	IsAdmin   bool
}

// UpdateUserDTO represents the data transfer object for updating a user
type UpdateUserDTO struct {
	ID        uuid.UUID
	Email     *string
	Username  *string
	FirstName *string
	LastName  *string
	IsActive  *bool
	IsAdmin   *bool
}

// UserDTO represents the data transfer object for a user
type UserDTO struct {
	ID        uuid.UUID
	Email     string
	Username  string
	FirstName string
	LastName  string
	FullName  string
	IsActive  bool
	IsAdmin   bool
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// ToEntity converts CreateUserDTO to User entity
func (dto *CreateUserDTO) ToEntity() *entity.User {
	return &entity.User{
		Email:     dto.Email,
		Username:  dto.Username,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		IsActive:  dto.IsActive,
		IsAdmin:   dto.IsAdmin,
	}
}

// FromEntity creates a UserDTO from User entity
func FromEntity(user *entity.User) *UserDTO {
	dto := &UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		FullName:  user.GetFullName(),
		IsActive:  user.IsActive,
		IsAdmin:   user.IsAdmin,
		Status:    string(user.GetStatus()),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if user.DeletedAt.Valid {
		dto.DeletedAt = &user.DeletedAt.Time
	}

	return dto
}

// ListUsersDTO represents the data transfer object for listing users
type ListUsersDTO struct {
	Users      []*UserDTO
	Page       int
	PageSize   int
	TotalItems int
	TotalPages int
}

// SearchUsersDTO represents the data transfer object for searching users
type SearchUsersDTO struct {
	Query    string
	Page     int
	PageSize int
	Filter   *FilterDTO
}

// FilterDTO represents filter options for users
type FilterDTO struct {
	Email    *string
	Username *string
	IsActive *bool
	IsAdmin  *bool
}

// ChangePasswordDTO represents the data transfer object for changing password
type ChangePasswordDTO struct {
	UserID      uuid.UUID
	OldPassword string
	NewPassword string
}

// AuthenticateDTO represents the data transfer object for authentication
type AuthenticateDTO struct {
	Identifier string // Email or username
	Password   string
}
