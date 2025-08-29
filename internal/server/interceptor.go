package server

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor logs all incoming requests
func LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		
		// Get metadata
		md, _ := metadata.FromIncomingContext(ctx)
		
		// Log request
		log.Printf("[REQUEST] Method: %s, Metadata: %v", info.FullMethod, md)
		
		// Handle request
		resp, err := handler(ctx, req)
		
		// Log response
		duration := time.Since(start)
		if err != nil {
			log.Printf("[ERROR] Method: %s, Duration: %v, Error: %v", info.FullMethod, duration, err)
		} else {
			log.Printf("[RESPONSE] Method: %s, Duration: %v", info.FullMethod, duration)
		}
		
		return resp, err
	}
}

// RecoveryInterceptor recovers from panics and returns an error
func RecoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC] Method: %s, Panic: %v", info.FullMethod, r)
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()
		
		return handler(ctx, req)
	}
}

// ValidationInterceptor validates incoming requests
func ValidationInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// You can add custom validation logic here
		// For example, check if required fields are present
		
		if validator, ok := req.(interface{ Validate() error }); ok {
			if err := validator.Validate(); err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
			}
		}
		
		return handler(ctx, req)
	}
}

// AuthInterceptor handles authentication
func AuthInterceptor() grpc.UnaryServerInterceptor {
	// List of methods that don't require authentication
	publicMethods := map[string]bool{
		"/user.v1.UserService/AuthenticateUser": true,
		"/user.v1.UserService/CreateUser":       true,
		"/health.v1.HealthService/Check":        true,
	}
	
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Skip authentication for public methods
		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}
		
		// Get metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata not found")
		}
		
		// Check for authorization header
		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization header not found")
		}
		
		// TODO: Validate token
		// For now, we'll just check if a token is present
		// In production, you should validate the JWT token here
		
		return handler(ctx, req)
	}
}

// ChainUnaryInterceptors chains multiple unary interceptors
func ChainUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(interceptors...)
}