package entity

import "errors"

var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")
	
	// ErrUserAlreadyExists is returned when attempting to create a user that already exists
	ErrUserAlreadyExists = errors.New("user already exists")
	
	// ErrInvalidEmail is returned when an email is invalid
	ErrInvalidEmail = errors.New("invalid email address")
	
	// ErrInvalidUsername is returned when a username is invalid
	ErrInvalidUsername = errors.New("invalid username")
	
	// ErrPasswordTooShort is returned when a password is too short
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	
	// ErrInvalidCredentials is returned when login credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")
	
	// ErrInvalidUserID is returned when a user ID is invalid
	ErrInvalidUserID = errors.New("invalid user ID")
)