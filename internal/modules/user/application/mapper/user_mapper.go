package mapper

import (
	"github.com/gigi434/sample-grpc-server/internal/modules/user/application/dto"
	"github.com/gigi434/sample-grpc-server/internal/modules/user/domain/entity"
	pb "github.com/gigi434/sample-grpc-server/pkg/generated/v1/user"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserToProto converts a User entity to proto message
func UserToProto(user *entity.User) *pb.User {
	if user == nil {
		return nil
	}

	protoUser := &pb.User{
		Id:        user.ID.String(),
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		FullName:  user.GetFullName(),
		IsActive:  user.IsActive,
		IsAdmin:   user.IsAdmin,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}

	// Set status
	switch user.GetStatus() {
	case entity.UserStatusActive:
		protoUser.Status = pb.UserStatus_USER_STATUS_ACTIVE
	case entity.UserStatusInactive:
		protoUser.Status = pb.UserStatus_USER_STATUS_INACTIVE
	case entity.UserStatusSuspended:
		protoUser.Status = pb.UserStatus_USER_STATUS_SUSPENDED
	default:
		protoUser.Status = pb.UserStatus_USER_STATUS_UNSPECIFIED
	}

	// Set deleted_at if soft deleted
	if user.DeletedAt.Valid {
		protoUser.DeletedAt = timestamppb.New(user.DeletedAt.Time)
	}

	return protoUser
}

// UserDTOToProto converts a UserDTO to proto message
func UserDTOToProto(dto *dto.UserDTO) *pb.User {
	if dto == nil {
		return nil
	}

	protoUser := &pb.User{
		Id:        dto.ID.String(),
		Email:     dto.Email,
		Username:  dto.Username,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		FullName:  dto.FullName,
		IsActive:  dto.IsActive,
		IsAdmin:   dto.IsAdmin,
		CreatedAt: timestamppb.New(dto.CreatedAt),
		UpdatedAt: timestamppb.New(dto.UpdatedAt),
	}

	// Set status
	switch dto.Status {
	case "active":
		protoUser.Status = pb.UserStatus_USER_STATUS_ACTIVE
	case "inactive":
		protoUser.Status = pb.UserStatus_USER_STATUS_INACTIVE
	case "suspended":
		protoUser.Status = pb.UserStatus_USER_STATUS_SUSPENDED
	default:
		protoUser.Status = pb.UserStatus_USER_STATUS_UNSPECIFIED
	}

	// Set deleted_at if present
	if dto.DeletedAt != nil {
		protoUser.DeletedAt = timestamppb.New(*dto.DeletedAt)
	}

	return protoUser
}

// CreateUserRequestToDTO converts CreateUserRequest to CreateUserDTO
func CreateUserRequestToDTO(req *pb.CreateUserRequest) *dto.CreateUserDTO {
	if req == nil {
		return nil
	}

	dto := &dto.CreateUserDTO{
		Email:     req.Email,
		Username:  req.Username,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,  // Default to true
		IsAdmin:   false, // Default to false
	}

	// Override defaults if provided
	if req.IsActive != nil {
		dto.IsActive = *req.IsActive
	}
	if req.IsAdmin != nil {
		dto.IsAdmin = *req.IsAdmin
	}

	return dto
}

// UpdateUserRequestToDTO converts UpdateUserRequest to UpdateUserDTO
func UpdateUserRequestToDTO(req *pb.UpdateUserRequest) (*dto.UpdateUserDTO, error) {
	if req == nil {
		return nil, nil
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	dto := &dto.UpdateUserDTO{
		ID: id,
	}

	// Check field mask to determine which fields to update
	if req.UpdateMask != nil {
		for _, path := range req.UpdateMask.Paths {
			switch path {
			case "email":
				if req.Email != nil {
					dto.Email = req.Email
				}
			case "username":
				if req.Username != nil {
					dto.Username = req.Username
				}
			case "first_name":
				if req.FirstName != nil {
					dto.FirstName = req.FirstName
				}
			case "last_name":
				if req.LastName != nil {
					dto.LastName = req.LastName
				}
			case "is_active":
				if req.IsActive != nil {
					dto.IsActive = req.IsActive
				}
			case "is_admin":
				if req.IsAdmin != nil {
					dto.IsAdmin = req.IsAdmin
				}
			}
		}
	} else {
		// If no field mask, update all provided fields
		dto.Email = req.Email
		dto.Username = req.Username
		dto.FirstName = req.FirstName
		dto.LastName = req.LastName
		dto.IsActive = req.IsActive
		dto.IsAdmin = req.IsAdmin
	}

	return dto, nil
}

// ListUsersFilterToDTO converts ListUsersFilter to FilterDTO
func ListUsersFilterToDTO(filter *pb.ListUsersFilter) *dto.FilterDTO {
	if filter == nil {
		return nil
	}

	return &dto.FilterDTO{
		Email:    filter.Email,
		Username: filter.Username,
		IsActive: filter.IsActive,
		IsAdmin:  filter.IsAdmin,
	}
}
