package domain

import (
	"github.com/google/uuid"
)

type Resource struct {
	ID uuid.UUID `json:"id" example:"11111111-2222-4444-3333-555555555555"`
}
