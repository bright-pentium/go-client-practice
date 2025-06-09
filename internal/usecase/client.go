package usecase

import (
	"context"
	"fmt"

	"github.com/bright-pentium/go-client-practice/internal/domain"
	repo "github.com/bright-pentium/go-client-practice/internal/repository/client"
	"github.com/google/uuid"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
)

const (
	ClientPepper string = "3tw0d2o"
)

type ClientUseCase struct {
	repo repo.IClientRepository
}

func NewClientUseCase(repo repo.IClientRepository) *ClientUseCase {
	return &ClientUseCase{repo: repo}
}

func (u *ClientUseCase) ListClientsByUser(ctx context.Context, userID uuid.UUID) ([]domain.Client, error) {
	return u.repo.ListClientsByUser(ctx, userID)
}

func (u *ClientUseCase) CreateClient(ctx context.Context, ID uuid.UUID, userID uuid.UUID, scope []domain.Permission) (*domain.Client, string, error) {
	randomStrings, err := password.Generate(32, 10, 0, false, true)
	if err != nil {
		return nil, "", err
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(ClientPepper+randomStrings), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %s", domain.ErrClientHashFail, err)
	}

	client, err := u.repo.CreateClient(ctx, ID, userID, scope, passwordHash)
	if err != nil {
		return nil, "", err
	}
	return client, randomStrings, nil
}

func (u *ClientUseCase) ClientLogin(ctx context.Context, ID uuid.UUID, secret string) (*domain.Client, error) {
	client, err := u.repo.GetClientByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	if bcrypt.CompareHashAndPassword(client.SecretHash, []byte(ClientPepper+secret)) != nil {
		return nil, domain.ErrClientLoginFail
	}
	return client, nil
}

func (u *ClientUseCase) UpdateClientScope(ctx context.Context, ID uuid.UUID, userID uuid.UUID, scope []domain.Permission) (*domain.Client, error) {
	return u.repo.UpdateClientByIDandUser(ctx, ID, userID, scope, nil)
}

func (u *ClientUseCase) DeleteClientByIDandUser(ctx context.Context, ID uuid.UUID, userID uuid.UUID) error {
	return u.repo.DeleteClientByIDandUser(ctx, ID, userID)
}
