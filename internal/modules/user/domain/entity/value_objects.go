package entity

import (
	"fmt"
	"regexp"
	"strings"
)

// Email represents a validated email address
type Email struct {
	value string
}

// NewEmail creates a new Email value object
func NewEmail(email string) (*Email, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if !isValidEmail(email) {
		return nil, fmt.Errorf("%w: %s", ErrInvalidEmail, email)
	}
	return &Email{value: email}, nil
}

// String returns the string representation of the email
func (e Email) String() string {
	return e.value
}

// Value returns the email value
func (e Email) Value() string {
	return e.value
}

// Username represents a validated username
type Username struct {
	value string
}

// NewUsername creates a new Username value object
func NewUsername(username string) (*Username, error) {
	username = strings.TrimSpace(username)
	if err := validateUsername(username); err != nil {
		return nil, err
	}
	return &Username{value: username}, nil
}

// String returns the string representation of the username
func (u Username) String() string {
	return u.value
}

// Value returns the username value
func (u Username) Value() string {
	return u.value
}

// Password represents a hashed password
type Password struct {
	hash string
}

// NewPassword creates a new Password value object (expects already hashed password)
func NewPassword(hash string) *Password {
	return &Password{hash: hash}
}

// Hash returns the password hash
func (p Password) Hash() string {
	return p.hash
}

// PersonName represents a person's name
type PersonName struct {
	FirstName string
	LastName  string
}

// NewPersonName creates a new PersonName value object
func NewPersonName(firstName, lastName string) (*PersonName, error) {
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)
	
	if firstName == "" && lastName == "" {
		return nil, fmt.Errorf("at least one name must be provided")
	}
	
	if len(firstName) > 100 || len(lastName) > 100 {
		return nil, fmt.Errorf("name too long")
	}
	
	return &PersonName{
		FirstName: firstName,
		LastName:  lastName,
	}, nil
}

// FullName returns the full name
func (n PersonName) FullName() string {
	if n.FirstName == "" {
		return n.LastName
	}
	if n.LastName == "" {
		return n.FirstName
	}
	return fmt.Sprintf("%s %s", n.FirstName, n.LastName)
}

// Helper functions

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

func isValidEmail(email string) bool {
	if len(email) > 255 {
		return false
	}
	return emailRegex.MatchString(email)
}

func validateUsername(username string) error {
	if len(username) < 3 {
		return fmt.Errorf("%w: username must be at least 3 characters", ErrInvalidUsername)
	}
	if len(username) > 100 {
		return fmt.Errorf("%w: username must not exceed 100 characters", ErrInvalidUsername)
	}
	
	// Username must start with a letter and contain only letters, numbers, and underscores
	usernameRegex := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
	if !usernameRegex.MatchString(username) {
		return fmt.Errorf("%w: username must start with a letter and contain only letters, numbers, and underscores", ErrInvalidUsername)
	}
	
	return nil
}