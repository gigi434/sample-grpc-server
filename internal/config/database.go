package config

import (
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
	dbErr  error
)

// GetDB returns the singleton database connection
func GetDB() (*gorm.DB, error) {
	dbOnce.Do(func() {
		db, dbErr = initDB()
	})
	return db, dbErr
}

// MustGetDB returns the singleton database connection or panics if there's an error
func MustGetDB() *gorm.DB {
	database, err := GetDB()
	if err != nil {
		panic(fmt.Sprintf("Failed to get database connection: %v", err))
	}
	return database
}

// initDB initializes the database connection
func initDB() (*gorm.DB, error) {
	cfg := GetConfig()

	// Configure GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		// Additional GORM configurations can be added here
	}

	// Open database connection
	database, err := gorm.Open(postgres.Open(cfg.Database.GetDSN()), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL database to configure connection pool
	sqlDB, err := database.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")
	return database, nil
}

// CloseDB closes the database connection
func CloseDB() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return fmt.Errorf("failed to get database instance: %w", err)
		}
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("failed to close database connection: %w", err)
		}
		log.Println("Database connection closed")
	}
	return nil
}

// ResetDB resets the database connection (useful for testing or reconnection)
func ResetDB() {
	dbOnce = sync.Once{}
	db = nil
	dbErr = nil
}