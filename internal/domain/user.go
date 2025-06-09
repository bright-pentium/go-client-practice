package domain

import (
	"errors"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" example:"11111111-2222-4444-3333-555555555555"`
	Name         string    `json:"name" example:"John Doe"`
	Account      string    `json:"account" example:"johndoe123"`
	PasswordHash []byte    `json:"-"` // typically not included in JSON responses
}

var (
	// Returned when a user with the same ID or name already exists
	ErrUserAlreadyExists = errors.New("user already exists")

	// Returned when provided user data violates constraints (e.g., empty name, invalid format)
	ErrInvalidUserData = errors.New("invalid user data")

	// Returned when a user doesnot exists
	ErrUserNotFound = errors.New("user is not found")

	// Returned when hashing gots wrong, it is an internal error.
	ErrUserHashFail = errors.New("user generate hash fails")

	// Returned when hashing gots wrong, it is an internal error.
	ErrUserLoginFail = errors.New("user either not found or password is wrong")

	// Returned when multiple user found but expected only one
	ErrMultipleUserFound = errors.New("mutiple user is found")

	// other error occured in user domain, including pg system error
	ErrGeneralUser = errors.New("general user data")
)
