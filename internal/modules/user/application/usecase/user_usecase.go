package usecase

import (
	"context"
	"fmt"
	"math"

	"github.com/gigi434/sample-grpc-server/internal/modules/user/application/dto"
	"github.com/gigi434/sample-grpc-server/internal/modules/user/domain/entity"
	"github.com/gigi434/sample-grpc-server/internal/modules/user/domain/repository"
	"github.com/gigi434/sample-grpc-server/internal/modules/user/domain/service"
	"github.com/google/uuid"
)

// UserUseCase handles user-related business logic
type UserUseCase struct {
	userRepo    repository.UserRepository
	userService *service.UserService
}

// NewUserUseCase creates a new instance of UserUseCase
func NewUserUseCase(userRepo repository.UserRepository, userService *service.UserService) *UserUseCase {
	return &UserUseCase{
		userRepo:    userRepo,
		userService: userService,
	}
}

// CreateUser creates a new user
func (uc *UserUseCase) CreateUser(ctx context.Context, createDTO *dto.CreateUserDTO) (*dto.UserDTO, error) {
	// Convert DTO to entity
	user := createDTO.ToEntity()

	// Create user using domain service (handles validation and password hashing)
	if err := uc.userService.CreateUser(ctx, user, createDTO.Password); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Convert entity to DTO
	return dto.FromEntity(user), nil
}

// GetUser retrieves a user by ID
func (uc *UserUseCase) GetUser(ctx context.Context, id string) (*dto.UserDTO, error) {
	// Parse UUID
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Get user from repository
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, entity.ErrUserNotFound
	}

	// Convert entity to DTO
	return dto.FromEntity(user), nil
}

// ListUsers retrieves a list of users with pagination
func (uc *UserUseCase) ListUsers(ctx context.Context, page, pageSize int, filter *dto.FilterDTO) (*dto.ListUsersDTO, error) {
	// Calculate offset
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// Create repository filter
	var repoFilter *repository.UserFilter
	if filter != nil {
		repoFilter = &repository.UserFilter{
			Email:    filter.Email,
			Username: filter.Username,
			IsActive: filter.IsActive,
			IsAdmin:  filter.IsAdmin,
		}
	}

	// Check if repository supports advanced filtering
	if advRepo, ok := uc.userRepo.(repository.UserRepositoryWithFilters); ok {
		// Use advanced filtering
		opts := &repository.ListOptions{
			Offset: offset,
			Limit:  pageSize,
			Filter: repoFilter,
			Sort: &repository.UserSortOptions{
				Field: "created_at",
				Order: "desc",
			},
		}

		users, err := advRepo.ListWithOptions(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list users: %w", err)
		}

		// Get total count
		totalCount, err := advRepo.CountWithFilter(ctx, repoFilter)
		if err != nil {
			return nil, fmt.Errorf("failed to count users: %w", err)
		}

		// Convert entities to DTOs
		userDTOs := make([]*dto.UserDTO, len(users))
		for i, user := range users {
			userDTOs[i] = dto.FromEntity(user)
		}

		// Calculate total pages
		totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

		return &dto.ListUsersDTO{
			Users:      userDTOs,
			Page:       page,
			PageSize:   pageSize,
			TotalItems: int(totalCount),
			TotalPages: totalPages,
		}, nil
	}

	// Fallback to basic listing
	users, err := uc.userRepo.List(ctx, offset, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Get total count
	totalCount, err := uc.userRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Convert entities to DTOs
	userDTOs := make([]*dto.UserDTO, len(users))
	for i, user := range users {
		userDTOs[i] = dto.FromEntity(user)
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &dto.ListUsersDTO{
		Users:      userDTOs,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: int(totalCount),
		TotalPages: totalPages,
	}, nil
}

// UpdateUser updates an existing user
func (uc *UserUseCase) UpdateUser(ctx context.Context, updateDTO *dto.UpdateUserDTO) (*dto.UserDTO, error) {
	// Create an entity with update fields
	updates := &entity.User{}

	if updateDTO.Email != nil {
		updates.Email = *updateDTO.Email
	}
	if updateDTO.Username != nil {
		updates.Username = *updateDTO.Username
	}
	if updateDTO.FirstName != nil {
		updates.FirstName = *updateDTO.FirstName
	}
	if updateDTO.LastName != nil {
		updates.LastName = *updateDTO.LastName
	}
	if updateDTO.IsActive != nil {
		updates.IsActive = *updateDTO.IsActive
	}
	if updateDTO.IsAdmin != nil {
		updates.IsAdmin = *updateDTO.IsAdmin
	}

	// Update user using domain service (handles validation)
	if err := uc.userService.UpdateUser(ctx, updateDTO.ID, updates); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Get updated user
	user, err := uc.userRepo.GetByID(ctx, updateDTO.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated user: %w", err)
	}

	// Convert entity to DTO
	return dto.FromEntity(user), nil
}

// DeleteUser deletes a user (soft delete by default)
func (uc *UserUseCase) DeleteUser(ctx context.Context, id string, hardDelete bool) error {
	// Parse UUID
	userID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Check if user exists
	exists, err := uc.userRepo.Exists(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	if !exists {
		return entity.ErrUserNotFound
	}

	// Delete user
	if err := uc.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// BatchGetUsers retrieves multiple users by IDs
func (uc *UserUseCase) BatchGetUsers(ctx context.Context, ids []string) (map[string]*dto.UserDTO, []string, error) {
	users := make(map[string]*dto.UserDTO)
	notFound := make([]string, 0)

	for _, id := range ids {
		// Parse UUID
		userID, err := uuid.Parse(id)
		if err != nil {
			notFound = append(notFound, id)
			continue
		}

		// Get user from repository
		user, err := uc.userRepo.GetByID(ctx, userID)
		if err != nil || user == nil {
			notFound = append(notFound, id)
			continue
		}

		// Convert entity to DTO
		users[id] = dto.FromEntity(user)
	}

	return users, notFound, nil
}

// SearchUsers searches users by query
func (uc *UserUseCase) SearchUsers(ctx context.Context, searchDTO *dto.SearchUsersDTO) (*dto.ListUsersDTO, error) {
	// For now, we'll use the list functionality with filters
	// In a real implementation, you might want to use a search engine like Elasticsearch

	// Create filter based on search query
	filter := searchDTO.Filter
	if filter == nil {
		filter = &dto.FilterDTO{}
	}

	// Add search query to filter (search in email and username)
	if searchDTO.Query != "" {
		filter.Email = &searchDTO.Query
		// Note: In a real implementation, you'd want to search across multiple fields
	}

	return uc.ListUsers(ctx, searchDTO.Page, searchDTO.PageSize, filter)
}

// ChangePassword changes a user's password
func (uc *UserUseCase) ChangePassword(ctx context.Context, changeDTO *dto.ChangePasswordDTO) error {
	// Use domain service to change password
	if err := uc.userService.ChangePassword(ctx, changeDTO.UserID, changeDTO.OldPassword, changeDTO.NewPassword); err != nil {
		return fmt.Errorf("failed to change password: %w", err)
	}

	return nil
}

// AuthenticateUser authenticates a user with email/username and password
func (uc *UserUseCase) AuthenticateUser(ctx context.Context, authDTO *dto.AuthenticateDTO) (*dto.UserDTO, error) {
	// Use domain service to authenticate
	user, err := uc.userService.Authenticate(ctx, authDTO.Identifier, authDTO.Password)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Convert entity to DTO
	return dto.FromEntity(user), nil
}
