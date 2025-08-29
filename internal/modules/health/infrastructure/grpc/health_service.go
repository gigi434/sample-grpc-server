package grpc

import (
	"context"
	"time"

	"github.com/gigi434/sample-grpc-server/internal/config"
	pb "github.com/gigi434/sample-grpc-server/pkg/generated/v1/health"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	// Server start time for uptime calculation
	serverStartTime = time.Now()
)

// HealthServiceServer implements the HealthService gRPC server
type HealthServiceServer struct {
	pb.UnimplementedHealthServiceServer
	version string
}

// NewHealthServiceServer creates a new HealthServiceServer instance
func NewHealthServiceServer(version string) *HealthServiceServer {
	return &HealthServiceServer{
		version: version,
	}
}

// Check performs a health check
func (s *HealthServiceServer) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	// Check database health
	dbHealth := s.checkDatabase()

	// Determine overall status
	overallStatus := pb.ServingStatus_SERVING_STATUS_SERVING
	if dbHealth.Status != pb.ServingStatus_SERVING_STATUS_SERVING {
		overallStatus = pb.ServingStatus_SERVING_STATUS_NOT_SERVING
	}

	// Calculate uptime
	uptime := time.Since(serverStartTime).Seconds()

	// Build response
	response := &pb.HealthCheckResponse{
		Status:        overallStatus,
		CheckedAt:     timestamppb.Now(),
		Version:       s.version,
		UptimeSeconds: int64(uptime),
		Services: map[string]*pb.ServiceHealth{
			"user_service": {
				Name:        "user_service",
				Status:      pb.ServingStatus_SERVING_STATUS_SERVING,
				Message:     "User service is healthy",
				LastChecked: timestamppb.Now(),
				Dependencies: map[string]*pb.DependencyHealth{
					"database": dbHealth,
				},
			},
			"health_service": {
				Name:        "health_service",
				Status:      pb.ServingStatus_SERVING_STATUS_SERVING,
				Message:     "Health service is healthy",
				LastChecked: timestamppb.Now(),
			},
		},
	}

	// If a specific service was requested, filter the response
	if req.Service != "" {
		if service, ok := response.Services[req.Service]; ok {
			response.Services = map[string]*pb.ServiceHealth{
				req.Service: service,
			}
		} else {
			return nil, status.Errorf(codes.NotFound, "service %s not found", req.Service)
		}
	}

	return response, nil
}

// Watch performs a health check and returns a stream of health status
func (s *HealthServiceServer) Watch(req *pb.HealthCheckRequest, stream pb.HealthService_WatchServer) error {
	// Send health status every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Send initial health status
	resp, err := s.Check(stream.Context(), req)
	if err != nil {
		return err
	}
	if err := stream.Send(resp); err != nil {
		return err
	}

	// Send periodic updates
	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case <-ticker.C:
			resp, err := s.Check(stream.Context(), req)
			if err != nil {
				return err
			}
			if err := stream.Send(resp); err != nil {
				return err
			}
		}
	}
}

// checkDatabase checks the database health
func (s *HealthServiceServer) checkDatabase() *pb.DependencyHealth {
	start := time.Now()

	// Try to get database connection
	db, err := config.GetDB()
	if err != nil {
		return &pb.DependencyHealth{
			Name:           "postgresql",
			Type:           "database",
			Status:         pb.ServingStatus_SERVING_STATUS_NOT_SERVING,
			ResponseTimeMs: time.Since(start).Milliseconds(),
			Message:        err.Error(),
		}
	}

	// Try to ping database
	sqlDB, err := db.DB()
	if err != nil {
		return &pb.DependencyHealth{
			Name:           "postgresql",
			Type:           "database",
			Status:         pb.ServingStatus_SERVING_STATUS_NOT_SERVING,
			ResponseTimeMs: time.Since(start).Milliseconds(),
			Message:        err.Error(),
		}
	}

	if err := sqlDB.Ping(); err != nil {
		return &pb.DependencyHealth{
			Name:           "postgresql",
			Type:           "database",
			Status:         pb.ServingStatus_SERVING_STATUS_NOT_SERVING,
			ResponseTimeMs: time.Since(start).Milliseconds(),
			Message:        err.Error(),
		}
	}

	return &pb.DependencyHealth{
		Name:           "postgresql",
		Type:           "database",
		Status:         pb.ServingStatus_SERVING_STATUS_SERVING,
		ResponseTimeMs: time.Since(start).Milliseconds(),
		Message:        "Database is healthy",
	}
}
