package domain

import (
	"errors"

	"github.com/google/uuid"
)

type Client struct {
	ID         uuid.UUID    `json:"id" example:"11111111-2222-4444-3333-555555555555"`
	UserID     uuid.UUID    `json:"userId" example:"11111111-2222-4444-3333-555555555555"`
	SecretHash []byte       `json:"-"`
	Scope      []Permission `json:"scope" example:"*"`
}

var (
	// Returned when a client with the same ID or name already exists
	ErrClientAlreadyExists = errors.New("client already exists")

	// Returned when provided client data violates constraints (e.g., empty name, invalid format)
	ErrInvalidClientData = errors.New("invalid client data")

	// Returned when a client doesnot exists
	ErrClientNotFound = errors.New("client is not found")

	// Returned when hashing gots wrong, it is an internal error.
	ErrClientHashFail = errors.New("client generate hash fails")

	// Returned when hashing gots wrong, it is an internal error.
	ErrClientLoginFail = errors.New("client either not found or password is wrong")

	// Returned when multiple client found but expected only one
	ErrMultipleClientFound = errors.New("mutiple client is found")

	// other error occured in client domain, including pg system error
	ErrGeneralClient = errors.New("general client data")
)
