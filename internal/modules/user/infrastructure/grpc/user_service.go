package grpc

import (
	"context"

	"github.com/gigi434/sample-grpc-server/internal/modules/user/application/dto"
	"github.com/gigi434/sample-grpc-server/internal/modules/user/application/mapper"
	"github.com/gigi434/sample-grpc-server/internal/modules/user/application/usecase"
	commonpb "github.com/gigi434/sample-grpc-server/pkg/generated/common"
	pb "github.com/gigi434/sample-grpc-server/pkg/generated/v1/user"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserServiceServer implements the UserService gRPC server
type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	userUseCase *usecase.UserUseCase
}

// NewUserServiceServer creates a new UserServiceServer instance
func NewUserServiceServer(userUseCase *usecase.UserUseCase) *UserServiceServer {
	return &UserServiceServer{
		userUseCase: userUseCase,
	}
}

// CreateUser creates a new user
func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Validate request
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	// Convert request to DTO
	createDTO := mapper.CreateUserRequestToDTO(req)

	// Create user
	userDTO, err := s.userUseCase.CreateUser(ctx, createDTO)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert DTO to proto
	return &pb.CreateUserResponse{
		User: mapper.UserDTOToProto(userDTO),
	}, nil
}

// GetUser retrieves a user by ID
func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// Get user
	userDTO, err := s.userUseCase.GetUser(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	// Convert DTO to proto
	return &pb.GetUserResponse{
		User: mapper.UserDTOToProto(userDTO),
	}, nil
}

// ListUsers retrieves a list of users with pagination
func (s *UserServiceServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	// Default pagination
	page := int(1)
	pageSize := int(10)

	if req.Pagination != nil {
		if req.Pagination.Page > 0 {
			page = int(req.Pagination.Page)
		}
		if req.Pagination.PageSize > 0 {
			pageSize = int(req.Pagination.PageSize)
		}
		// Limit page size to prevent abuse
		if pageSize > 100 {
			pageSize = 100
		}
	}

	// Convert filter
	var filter *dto.FilterDTO
	if req.Filter != nil {
		filter = mapper.ListUsersFilterToDTO(req.Filter)
	}

	// List users
	listDTO, err := s.userUseCase.ListUsers(ctx, page, pageSize, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert DTOs to proto
	users := make([]*pb.User, len(listDTO.Users))
	for i, userDTO := range listDTO.Users {
		users[i] = mapper.UserDTOToProto(userDTO)
	}

	// Create pagination response
	hasNext := listDTO.Page < listDTO.TotalPages
	hasPrevious := listDTO.Page > 1

	return &pb.ListUsersResponse{
		Users: users,
		Pagination: &commonpb.PaginationResponse{
			Page:        int32(listDTO.Page),
			PageSize:    int32(listDTO.PageSize),
			TotalItems:  int32(listDTO.TotalItems),
			TotalPages:  int32(listDTO.TotalPages),
			HasNext:     hasNext,
			HasPrevious: hasPrevious,
		},
	}, nil
}

// UpdateUser updates an existing user
func (s *UserServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// Convert request to DTO
	updateDTO, err := mapper.UpdateUserRequestToDTO(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Update user
	userDTO, err := s.userUseCase.UpdateUser(ctx, updateDTO)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert DTO to proto
	return &pb.UpdateUserResponse{
		User: mapper.UserDTOToProto(userDTO),
	}, nil
}

// DeleteUser deletes a user
func (s *UserServiceServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	// Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// Delete user
	err := s.userUseCase.DeleteUser(ctx, req.Id, req.HardDelete)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteUserResponse{
		Success: true,
		Message: "User deleted successfully",
	}, nil
}

// BatchGetUsers retrieves multiple users by IDs
func (s *UserServiceServer) BatchGetUsers(ctx context.Context, req *pb.BatchGetUsersRequest) (*pb.BatchGetUsersResponse, error) {
	// Validate request
	if len(req.Ids) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one id is required")
	}

	// Limit batch size to prevent abuse
	if len(req.Ids) > 100 {
		return nil, status.Error(codes.InvalidArgument, "batch size cannot exceed 100")
	}

	// Get users
	users, notFound, err := s.userUseCase.BatchGetUsers(ctx, req.Ids)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert DTOs to proto
	protoUsers := make(map[string]*pb.User)
	for id, userDTO := range users {
		protoUsers[id] = mapper.UserDTOToProto(userDTO)
	}

	return &pb.BatchGetUsersResponse{
		Users:    protoUsers,
		NotFound: notFound,
	}, nil
}

// SearchUsers searches users by various criteria
func (s *UserServiceServer) SearchUsers(ctx context.Context, req *pb.SearchUsersRequest) (*pb.SearchUsersResponse, error) {
	// Default pagination
	page := int(1)
	pageSize := int(10)

	if req.Pagination != nil {
		if req.Pagination.Page > 0 {
			page = int(req.Pagination.Page)
		}
		if req.Pagination.PageSize > 0 {
			pageSize = int(req.Pagination.PageSize)
		}
		// Limit page size to prevent abuse
		if pageSize > 100 {
			pageSize = 100
		}
	}

	// Create search DTO
	searchDTO := &dto.SearchUsersDTO{
		Query:    req.Query,
		Page:     page,
		PageSize: pageSize,
	}

	if req.Filter != nil {
		searchDTO.Filter = mapper.ListUsersFilterToDTO(req.Filter)
	}

	// Search users
	listDTO, err := s.userUseCase.SearchUsers(ctx, searchDTO)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert DTOs to proto
	users := make([]*pb.User, len(listDTO.Users))
	for i, userDTO := range listDTO.Users {
		users[i] = mapper.UserDTOToProto(userDTO)
	}

	// Create pagination response
	hasNext := listDTO.Page < listDTO.TotalPages
	hasPrevious := listDTO.Page > 1

	return &pb.SearchUsersResponse{
		Users: users,
		Pagination: &commonpb.PaginationResponse{
			Page:        int32(listDTO.Page),
			PageSize:    int32(listDTO.PageSize),
			TotalItems:  int32(listDTO.TotalItems),
			TotalPages:  int32(listDTO.TotalPages),
			HasNext:     hasNext,
			HasPrevious: hasPrevious,
		},
		TotalMatches: int32(listDTO.TotalItems),
	}, nil
}

// ChangePassword changes a user's password
func (s *UserServiceServer) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.OldPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "old_password is required")
	}
	if req.NewPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "new_password is required")
	}

	// Parse user ID
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	// Create DTO
	changeDTO := &dto.ChangePasswordDTO{
		UserID:      userID,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}

	// Change password
	if err := s.userUseCase.ChangePassword(ctx, changeDTO); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ChangePasswordResponse{
		Success: true,
		Message: "Password changed successfully",
	}, nil
}

// AuthenticateUser authenticates a user with email/username and password
func (s *UserServiceServer) AuthenticateUser(ctx context.Context, req *pb.AuthenticateUserRequest) (*pb.AuthenticateUserResponse, error) {
	// Validate request
	if req.Identifier == "" {
		return nil, status.Error(codes.InvalidArgument, "identifier is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	// Create DTO
	authDTO := &dto.AuthenticateDTO{
		Identifier: req.Identifier,
		Password:   req.Password,
	}

	// Authenticate user
	userDTO, err := s.userUseCase.AuthenticateUser(ctx, authDTO)
	if err != nil {
		return &pb.AuthenticateUserResponse{
			Success: false,
			Message: "Invalid credentials",
		}, nil
	}

	// Convert DTO to proto
	return &pb.AuthenticateUserResponse{
		User:    mapper.UserDTOToProto(userDTO),
		Success: true,
		Message: "Authentication successful",
		// TODO: Generate JWT token here
	}, nil
}
