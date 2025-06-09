package domain

import "errors"

var (
	// Returned for unexpected/internal system errors (e.g., DB connection issues)
	ErrInternal = errors.New("internal db error")
)
