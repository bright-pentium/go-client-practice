package client

import (
	"context"

	"github.com/bright-pentium/go-client-practice/internal/domain"
	"github.com/google/uuid"
)

type IClientRepository interface {
	CreateClient(ctx context.Context, ID uuid.UUID, userID uuid.UUID, scope []domain.Permission, secretHash []byte) (*domain.Client, error)
	UpdateClientByIDandUser(ctx context.Context, ID uuid.UUID, userID uuid.UUID, scope []domain.Permission, secretHash []byte) (*domain.Client, error)

	// cqs
	GetClientByID(ctx context.Context, ID uuid.UUID) (*domain.Client, error)
	ListClientsByUser(ctx context.Context, userID uuid.UUID) ([]domain.Client, error)

	DeleteClientByIDandUser(ctx context.Context, ID uuid.UUID, userID uuid.UUID) error
}
