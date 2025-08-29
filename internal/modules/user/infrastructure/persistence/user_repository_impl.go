package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/gigi434/sample-grpc-server/internal/config"
	"github.com/gigi434/sample-grpc-server/internal/modules/user/domain/entity"
	"github.com/gigi434/sample-grpc-server/internal/modules/user/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// userRepository implements repository.UserRepository
type userRepository struct {
	// We don't store the DB connection here, we get it from config singleton
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository() repository.UserRepository {
	return &userRepository{}
}

// getDB gets the database connection from the singleton
func (r *userRepository) getDB() (*gorm.DB, error) {
	return config.GetDB()
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	db, err := r.getDB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	if err := db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	var user entity.User
	if err := db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	var user entity.User
	if err := db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	var user entity.User
	if err := db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return &user, nil
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	db, err := r.getDB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	if err := db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Delete soft deletes a user
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	result := db.WithContext(ctx).Delete(&entity.User{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return entity.ErrUserNotFound
	}
	return nil
}

// List retrieves users with pagination
func (r *userRepository) List(ctx context.Context, offset, limit int) ([]*entity.User, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	var users []*entity.User
	if err := db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}

// Count returns the total number of users
func (r *userRepository) Count(ctx context.Context) (int64, error) {
	db, err := r.getDB()
	if err != nil {
		return 0, fmt.Errorf("failed to get database connection: %w", err)
	}

	var count int64
	if err := db.WithContext(ctx).Model(&entity.User{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

// Exists checks if a user exists by ID
func (r *userRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	db, err := r.getDB()
	if err != nil {
		return false, fmt.Errorf("failed to get database connection: %w", err)
	}

	var count int64
	if err := db.WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", id).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return count > 0, nil
}

// ExistsByEmail checks if a user exists by email
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	db, err := r.getDB()
	if err != nil {
		return false, fmt.Errorf("failed to get database connection: %w", err)
	}

	var count int64
	if err := db.WithContext(ctx).
		Model(&entity.User{}).
		Where("email = ?", email).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	return count > 0, nil
}

// ExistsByUsername checks if a user exists by username
func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	db, err := r.getDB()
	if err != nil {
		return false, fmt.Errorf("failed to get database connection: %w", err)
	}

	var count int64
	if err := db.WithContext(ctx).
		Model(&entity.User{}).
		Where("username = ?", username).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}
	return count > 0, nil
}

// ListWithOptions retrieves users with advanced filtering and sorting
func (r *userRepository) ListWithOptions(ctx context.Context, opts *repository.ListOptions) ([]*entity.User, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	query := db.WithContext(ctx).Model(&entity.User{})

	// Apply filters
	if opts.Filter != nil {
		if opts.Filter.Email != nil {
			query = query.Where("email LIKE ?", "%"+*opts.Filter.Email+"%")
		}
		if opts.Filter.Username != nil {
			query = query.Where("username LIKE ?", "%"+*opts.Filter.Username+"%")
		}
		if opts.Filter.IsActive != nil {
			query = query.Where("is_active = ?", *opts.Filter.IsActive)
		}
		if opts.Filter.IsAdmin != nil {
			query = query.Where("is_admin = ?", *opts.Filter.IsAdmin)
		}
	}

	// Apply sorting
	if opts.Sort != nil {
		order := fmt.Sprintf("%s %s", opts.Sort.Field, opts.Sort.Order)
		query = query.Order(order)
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	var users []*entity.User
	if err := query.
		Offset(opts.Offset).
		Limit(opts.Limit).
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to list users with options: %w", err)
	}

	return users, nil
}

// CountWithFilter returns the count of users matching the filter
func (r *userRepository) CountWithFilter(ctx context.Context, filter *repository.UserFilter) (int64, error) {
	db, err := r.getDB()
	if err != nil {
		return 0, fmt.Errorf("failed to get database connection: %w", err)
	}

	query := db.WithContext(ctx).Model(&entity.User{})

	// Apply filters
	if filter != nil {
		if filter.Email != nil {
			query = query.Where("email LIKE ?", "%"+*filter.Email+"%")
		}
		if filter.Username != nil {
			query = query.Where("username LIKE ?", "%"+*filter.Username+"%")
		}
		if filter.IsActive != nil {
			query = query.Where("is_active = ?", *filter.IsActive)
		}
		if filter.IsAdmin != nil {
			query = query.Where("is_admin = ?", *filter.IsAdmin)
		}
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count users with filter: %w", err)
	}

	return count, nil
}
