package database

import (
	"fmt"
	"log"

	"github.com/gigi434/sample-grpc-server/internal/config"
	"github.com/gigi434/sample-grpc-server/internal/modules/user/domain/entity"
	"gorm.io/gorm"
)

// Migrate runs database migrations
func Migrate() error {
	db, err := config.GetDB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	// List of models to migrate
	models := []interface{}{
		&entity.User{},
		// Add other models here as they are created
	}

	// Run migrations
	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// DropTables drops all tables (use with caution!)
func DropTables() error {
	db, err := config.GetDB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	// List of models to drop
	models := []interface{}{
		&entity.User{},
		// Add other models here as they are created
	}

	// Drop tables
	if err := db.Migrator().DropTable(models...); err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}

	log.Println("Tables dropped successfully")
	return nil
}

// GetConnection returns the database connection for direct access
// This should be used sparingly, prefer using repositories
func GetConnection() (*gorm.DB, error) {
	return config.GetDB()
}
