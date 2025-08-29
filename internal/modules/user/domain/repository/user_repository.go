package repository

import (
	"context"

	"github.com/gigi434/sample-grpc-server/internal/modules/user/domain/entity"
	"github.com/google/uuid"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *entity.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	// GetByUsername retrieves a user by username
	GetByUsername(ctx context.Context, username string) (*entity.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *entity.User) error

	// Delete soft deletes a user
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves users with pagination
	List(ctx context.Context, offset, limit int) ([]*entity.User, error)

	// Count returns the total number of users
	Count(ctx context.Context) (int64, error)

	// Exists checks if a user exists by ID
	Exists(ctx context.Context, id uuid.UUID) (bool, error)

	// ExistsByEmail checks if a user exists by email
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// ExistsByUsername checks if a user exists by username
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}

// UserFilter represents filter options for listing users
type UserFilter struct {
	Email    *string
	Username *string
	IsActive *bool
	IsAdmin  *bool
}

// UserSortOptions represents sort options for listing users
type UserSortOptions struct {
	Field string // "created_at", "updated_at", "email", "username"
	Order string // "asc" or "desc"
}

// ListOptions represents options for listing users
type ListOptions struct {
	Offset int
	Limit  int
	Filter *UserFilter
	Sort   *UserSortOptions
}

// UserRepositoryWithFilters extends UserRepository with advanced filtering
type UserRepositoryWithFilters interface {
	UserRepository

	// ListWithOptions retrieves users with advanced filtering and sorting
	ListWithOptions(ctx context.Context, opts *ListOptions) ([]*entity.User, error)

	// CountWithFilter returns the count of users matching the filter
	CountWithFilter(ctx context.Context, filter *UserFilter) (int64, error)
}
