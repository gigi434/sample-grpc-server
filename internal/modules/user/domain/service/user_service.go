package service

import (
	"context"
	"fmt"

	"github.com/gigi434/sample-grpc-server/internal/modules/user/domain/entity"
	"github.com/gigi434/sample-grpc-server/internal/modules/user/domain/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserService provides domain services for user operations
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user with password hashing
func (s *UserService) CreateUser(ctx context.Context, user *entity.User, plainPassword string) error {
	// Validate email
	email, err := entity.NewEmail(user.Email)
	if err != nil {
		return err
	}
	user.Email = email.Value()

	// Validate username
	username, err := entity.NewUsername(user.Username)
	if err != nil {
		return err
	}
	user.Username = username.Value()

	// Validate name
	name, err := entity.NewPersonName(user.FirstName, user.LastName)
	if err != nil {
		return err
	}
	user.FirstName = name.FirstName
	user.LastName = name.LastName

	// Hash password
	hashedPassword, err := s.HashPassword(plainPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = hashedPassword

	// Check if user already exists
	emailExists, err := s.userRepo.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("failed to check email existence: %w", err)
	}
	if emailExists {
		return fmt.Errorf("%w: email already in use", entity.ErrUserAlreadyExists)
	}

	usernameExists, err := s.userRepo.ExistsByUsername(ctx, user.Username)
	if err != nil {
		return fmt.Errorf("failed to check username existence: %w", err)
	}
	if usernameExists {
		return fmt.Errorf("%w: username already in use", entity.ErrUserAlreadyExists)
	}

	// Create user
	if err := s.userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, userID uuid.UUID, updates *entity.User) error {
	// Get existing user
	existingUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if existingUser == nil {
		return entity.ErrUserNotFound
	}

	// Update fields if provided
	if updates.Email != "" && updates.Email != existingUser.Email {
		email, err := entity.NewEmail(updates.Email)
		if err != nil {
			return err
		}

		// Check if new email is already in use
		emailExists, err := s.userRepo.ExistsByEmail(ctx, email.Value())
		if err != nil {
			return fmt.Errorf("failed to check email existence: %w", err)
		}
		if emailExists {
			return fmt.Errorf("%w: email already in use", entity.ErrUserAlreadyExists)
		}

		existingUser.Email = email.Value()
	}

	if updates.Username != "" && updates.Username != existingUser.Username {
		username, err := entity.NewUsername(updates.Username)
		if err != nil {
			return err
		}

		// Check if new username is already in use
		usernameExists, err := s.userRepo.ExistsByUsername(ctx, username.Value())
		if err != nil {
			return fmt.Errorf("failed to check username existence: %w", err)
		}
		if usernameExists {
			return fmt.Errorf("%w: username already in use", entity.ErrUserAlreadyExists)
		}

		existingUser.Username = username.Value()
	}

	if updates.FirstName != "" {
		existingUser.FirstName = updates.FirstName
	}

	if updates.LastName != "" {
		existingUser.LastName = updates.LastName
	}

	// Update user
	if err := s.userRepo.Update(ctx, existingUser); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return entity.ErrUserNotFound
	}

	// Verify old password
	if err := s.VerifyPassword(user.Password, oldPassword); err != nil {
		return entity.ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := s.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	user.Password = hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// Authenticate authenticates a user with email/username and password
func (s *UserService) Authenticate(ctx context.Context, identifier, password string) (*entity.User, error) {
	var user *entity.User
	var err error

	// Try to find user by email first
	if _, emailErr := entity.NewEmail(identifier); emailErr == nil {
		user, err = s.userRepo.GetByEmail(ctx, identifier)
	}

	// If not found by email, try username
	if user == nil {
		user, err = s.userRepo.GetByUsername(ctx, identifier)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if user == nil {
		return nil, entity.ErrInvalidCredentials
	}

	// Verify password
	if err := s.VerifyPassword(user.Password, password); err != nil {
		return nil, entity.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("user account is not active")
	}

	return user, nil
}

// HashPassword hashes a plain text password
func (s *UserService) HashPassword(password string) (string, error) {
	if len(password) < 8 {
		return "", entity.ErrPasswordTooShort
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// VerifyPassword verifies a password against a hash
func (s *UserService) VerifyPassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
