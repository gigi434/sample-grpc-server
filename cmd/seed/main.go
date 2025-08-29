package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gigi434/sample-grpc-server/internal/config"
	"github.com/gigi434/sample-grpc-server/internal/modules/user/domain/entity"
	"github.com/gigi434/sample-grpc-server/internal/modules/user/domain/service"
	"github.com/gigi434/sample-grpc-server/internal/modules/user/infrastructure/persistence"
	"github.com/gigi434/sample-grpc-server/internal/shared/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedUser represents user data for seeding
type SeedUser struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	IsActive  bool   `json:"is_active"`
	IsAdmin   bool   `json:"is_admin"`
}

// SeedData represents the structure of seed data
type SeedData struct {
	Users []SeedUser `json:"users"`
}

func main() {
	log.Println("Starting database seeding...")

	// Check for --clean flag
	cleanFlag := false
	for _, arg := range os.Args[1:] {
		if arg == "--clean" {
			cleanFlag = true
			break
		}
	}

	// Run migrations first
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Clean data if flag is set
	if cleanFlag {
		log.Println("Cleaning existing data...")
		if err := clearUsers(); err != nil {
			log.Printf("Warning: Failed to clear users: %v", err)
		}
	}

	// Load seed data
	seedData, err := loadSeedData("test/fixtures/users.json")
	if err != nil {
		log.Fatalf("Failed to load seed data: %v", err)
	}

	// Seed users
	if err := seedUsers(seedData.Users); err != nil {
		log.Fatalf("Failed to seed users: %v", err)
	}

	log.Println("Database seeding completed successfully!")
}

// loadSeedData loads seed data from a JSON file
func loadSeedData(filepath string) (*SeedData, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open seed file: %w", err)
	}
	defer file.Close()

	var seedData SeedData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&seedData); err != nil {
		return nil, fmt.Errorf("failed to decode seed data: %w", err)
	}

	return &seedData, nil
}

// seedUsers seeds user data
func seedUsers(users []SeedUser) error {
	db, err := config.GetDB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	// Create user repository and service
	userRepo := persistence.NewUserRepository()
	userService := service.NewUserService(userRepo)

	log.Printf("Seeding %d users...", len(users))

	for _, seedUser := range users {
		// Check if user already exists
		var existingUser entity.User
		err := db.Where("email = ? OR username = ?", seedUser.Email, seedUser.Username).First(&existingUser).Error
		if err == nil {
			log.Printf("User with email %s or username %s already exists, skipping...", seedUser.Email, seedUser.Username)
			continue
		} else if err != gorm.ErrRecordNotFound {
			return fmt.Errorf("failed to check existing user: %w", err)
		}

		// Create new user
		user := &entity.User{
			ID:        uuid.New(),
			Email:     seedUser.Email,
			Username:  seedUser.Username,
			FirstName: seedUser.FirstName,
			LastName:  seedUser.LastName,
			IsActive:  seedUser.IsActive,
			IsAdmin:   seedUser.IsAdmin,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Hash password
		hashedPassword, err := userService.HashPassword(seedUser.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password for user %s: %w", seedUser.Email, err)
		}
		user.Password = hashedPassword

		// Save user directly (bypass service validation for seeding)
		if err := db.Create(user).Error; err != nil {
			return fmt.Errorf("failed to create user %s: %w", seedUser.Email, err)
		}

		log.Printf("Created user: %s (%s)", user.Email, user.GetFullName())
	}

	// Print summary
	var count int64
	db.Model(&entity.User{}).Count(&count)
	log.Printf("Total users in database: %d", count)

	return nil
}

// clearUsers clears all users from the database (optional utility function)
func clearUsers() error {
	db, err := config.GetDB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	// Delete all users (including soft deleted)
	if err := db.Unscoped().Where("1 = 1").Delete(&entity.User{}).Error; err != nil {
		return fmt.Errorf("failed to clear users: %w", err)
	}

	log.Println("All users have been cleared from the database")
	return nil
}
