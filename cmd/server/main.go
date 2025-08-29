package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gigi434/sample-grpc-server/internal/config"
	healthgrpc "github.com/gigi434/sample-grpc-server/internal/modules/health/infrastructure/grpc"
	"github.com/gigi434/sample-grpc-server/internal/modules/user/application/usecase"
	"github.com/gigi434/sample-grpc-server/internal/modules/user/domain/service"
	usergrpc "github.com/gigi434/sample-grpc-server/internal/modules/user/infrastructure/grpc"
	"github.com/gigi434/sample-grpc-server/internal/modules/user/infrastructure/persistence"
	"github.com/gigi434/sample-grpc-server/internal/server"
	"github.com/gigi434/sample-grpc-server/internal/shared/database"
	healthpb "github.com/gigi434/sample-grpc-server/pkg/generated/v1/health"
	userpb "github.com/gigi434/sample-grpc-server/pkg/generated/v1/user"
	"google.golang.org/grpc"
)

const (
	version = "1.0.0"
)

func main() {
	// Initialize logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Starting gRPC server version %s", version)

	// Load configuration
	_ = config.GetConfig()

	// Get server port from environment or use default
	port := 50051
	if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
		fmt.Sscanf(portStr, "%d", &port)
	}

	// Initialize database
	log.Println("Initializing database connection...")
	_, err := config.GetDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	log.Println("Running database migrations...")
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := persistence.NewUserRepository()

	// Initialize domain services
	userService := service.NewUserService(userRepo)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo, userService)

	// Create gRPC service implementations
	userServiceServer := usergrpc.NewUserServiceServer(userUseCase)
	healthServiceServer := healthgrpc.NewHealthServiceServer(version)

	// Create gRPC server with interceptors
	grpcServer, err := server.NewGRPCServer(
		port,
		server.ChainUnaryInterceptors(
			server.RecoveryInterceptor(),
			server.LoggingInterceptor(),
			server.ValidationInterceptor(),
			// Uncomment to enable authentication
			// server.AuthInterceptor(),
		),
	)
	if err != nil {
		log.Fatalf("Failed to create gRPC server: %v", err)
	}

	// Register services
	userpb.RegisterUserServiceServer(grpcServer.GetServer(), userServiceServer)
	healthpb.RegisterHealthServiceServer(grpcServer.GetServer(), healthServiceServer)

	// Start server in a goroutine
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("gRPC server listening on port %d", port)
		log.Printf("Health check available at: grpc://localhost:%d/health.v1.HealthService/Check", port)
		log.Printf("User service available at: grpc://localhost:%d/user.v1.UserService/*", port)
		serverErrors <- grpcServer.Start()
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("Server error: %v", err)
	case sig := <-quit:
		log.Printf("Received signal: %v", sig)

		// Create a context with timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Stop gRPC server
		if err := grpcServer.Stop(ctx); err != nil {
			log.Printf("Failed to stop gRPC server: %v", err)
		}

		// Close database connection
		if err := config.CloseDB(); err != nil {
			log.Printf("Failed to close database connection: %v", err)
		}

		log.Println("Server shutdown complete")
	}
}

// Helper function to setup dependencies (for testing)
func setupDependencies() (*usecase.UserUseCase, *healthgrpc.HealthServiceServer, error) {
	// Initialize repositories
	userRepo := persistence.NewUserRepository()

	// Initialize domain services
	userService := service.NewUserService(userRepo)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo, userService)

	// Create health service
	healthServiceServer := healthgrpc.NewHealthServiceServer(version)

	return userUseCase, healthServiceServer, nil
}

// RegisterServices registers all gRPC services (for testing)
func RegisterServices(grpcServer *grpc.Server, userUseCase *usecase.UserUseCase, healthService *healthgrpc.HealthServiceServer) {
	userServiceServer := usergrpc.NewUserServiceServer(userUseCase)
	userpb.RegisterUserServiceServer(grpcServer, userServiceServer)
	healthpb.RegisterHealthServiceServer(grpcServer, healthService)
}
