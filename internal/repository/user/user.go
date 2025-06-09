package user

import (
	"context"

	"github.com/bright-pentium/go-client-practice/internal/domain"
	"github.com/google/uuid"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, ID uuid.UUID, name string, account string, passwordHash []byte) (*domain.User, error)
	UpdateUserByID(ctx context.Context, ID uuid.UUID, name string, passwordHash []byte) (*domain.User, error)
	GetUserByID(ctx context.Context, ID uuid.UUID) (*domain.User, error)
	GetUserByAccount(ctx context.Context, account string) (*domain.User, error)
	// ListUser(ctx context.Context) ([]domain.User, error)
	DeleteUserByID(ctx context.Context, ID uuid.UUID) error
}
