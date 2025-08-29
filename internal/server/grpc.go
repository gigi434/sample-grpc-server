package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// GRPCServer represents the gRPC server
type GRPCServer struct {
	server   *grpc.Server
	listener net.Listener
	port     int
}

// NewGRPCServer creates a new gRPC server instance
func NewGRPCServer(port int, opts ...grpc.ServerOption) (*GRPCServer, error) {
	// Create listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port %d: %w", port, err)
	}

	// Default server options
	defaultOpts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     15 * time.Second,
			MaxConnectionAge:      30 * time.Second,
			MaxConnectionAgeGrace: 5 * time.Second,
			Time:                  5 * time.Second,
			Timeout:               1 * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second,
			PermitWithoutStream: true,
		}),
	}

	// Combine default options with provided options
	serverOpts := append(defaultOpts, opts...)

	// Create gRPC server
	grpcServer := grpc.NewServer(serverOpts...)

	// Enable reflection for debugging
	reflection.Register(grpcServer)

	return &GRPCServer{
		server:   grpcServer,
		listener: listener,
		port:     port,
	}, nil
}

// GetServer returns the underlying gRPC server
func (s *GRPCServer) GetServer() *grpc.Server {
	return s.server
}

// Start starts the gRPC server
func (s *GRPCServer) Start() error {
	log.Printf("Starting gRPC server on port %d", s.port)
	return s.server.Serve(s.listener)
}

// Stop gracefully stops the gRPC server
func (s *GRPCServer) Stop(ctx context.Context) error {
	log.Println("Stopping gRPC server...")
	
	// Create a channel to signal when graceful stop is complete
	stopped := make(chan struct{})
	
	go func() {
		s.server.GracefulStop()
		close(stopped)
	}()
	
	select {
	case <-ctx.Done():
		log.Println("Force stopping gRPC server...")
		s.server.Stop()
		return ctx.Err()
	case <-stopped:
		log.Println("gRPC server stopped gracefully")
		return nil
	}
}

// GetPort returns the port the server is listening on
func (s *GRPCServer) GetPort() int {
	return s.port
}