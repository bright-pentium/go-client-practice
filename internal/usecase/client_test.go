package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bright-pentium/go-client-practice/internal/domain"
	mockRepo "github.com/bright-pentium/go-client-practice/internal/repository/client"
	"github.com/bright-pentium/go-client-practice/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestListClientsByUser(t *testing.T) {
	mockRepo := new(mockRepo.MockClientRepository)
	uc := usecase.NewClientUseCase(mockRepo)

	userID := uuid.New()
	expectedClients := []domain.Client{{ID: uuid.New(), UserID: userID}}

	mockRepo.On("ListClientsByUser", mock.Anything, userID).Return(expectedClients, nil)

	ctx := context.Background()
	clients, err := uc.ListClientsByUser(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedClients, clients)
	mockRepo.AssertExpectations(t)
}

func TestCreateClient(t *testing.T) {
	mockRepo := new(mockRepo.MockClientRepository)
	uc := usecase.NewClientUseCase(mockRepo)

	clientID := uuid.New()
	userID := uuid.New()
	scope := []domain.Permission{"read", "write"}

	var capturedHash []byte
	expectedClient := &domain.Client{ID: clientID, UserID: userID, Scope: scope}

	mockRepo.
		On("CreateClient", mock.Anything, clientID, userID, scope, mock.AnythingOfType("[]uint8")).
		Run(func(args mock.Arguments) {
			capturedHash = args.Get(4).([]byte)
		}).
		Return(expectedClient, nil)

	ctx := context.Background()
	client, secret, err := uc.CreateClient(ctx, clientID, userID, scope)

	assert.NoError(t, err)
	assert.NotEmpty(t, secret)
	assert.Equal(t, expectedClient, client)
	assert.NotNil(t, capturedHash)
	mockRepo.AssertExpectations(t)
}

func TestClientLoginSuccess(t *testing.T) {
	mockRepo := new(mockRepo.MockClientRepository)
	uc := usecase.NewClientUseCase(mockRepo)

	clientID := uuid.New()
	secret := "supersecret"
	hash, _ := bcryptGenerateWithPepper(secret)

	expectedClient := &domain.Client{ID: clientID, SecretHash: hash}

	mockRepo.On("GetClientByID", mock.Anything, clientID).Return(expectedClient, nil)

	ctx := context.Background()
	client, err := uc.ClientLogin(ctx, clientID, secret)

	assert.NoError(t, err)
	assert.Equal(t, expectedClient, client)
	mockRepo.AssertExpectations(t)
}

func TestClientLoginFailure(t *testing.T) {
	mockRepo := new(mockRepo.MockClientRepository)
	uc := usecase.NewClientUseCase(mockRepo)

	clientID := uuid.New()
	incorrectSecret := "wrong"
	correctSecret := "correct"
	hash, _ := bcryptGenerateWithPepper(correctSecret)

	expectedClient := &domain.Client{ID: clientID, SecretHash: hash}

	mockRepo.On("GetClientByID", mock.Anything, clientID).Return(expectedClient, nil)

	ctx := context.Background()
	client, err := uc.ClientLogin(ctx, clientID, incorrectSecret)

	assert.ErrorIs(t, err, domain.ErrClientLoginFail)
	assert.Nil(t, client)
	mockRepo.AssertExpectations(t)
}

func TestClientLoginRepoError(t *testing.T) {
	mockRepo := new(mockRepo.MockClientRepository)
	uc := usecase.NewClientUseCase(mockRepo)

	clientID := uuid.New()

	mockRepo.On("GetClientByID", mock.Anything, clientID).Return(nil, errors.New("db error"))

	ctx := context.Background()
	client, err := uc.ClientLogin(ctx, clientID, "whatever")

	assert.Error(t, err)
	assert.Nil(t, client)
	mockRepo.AssertExpectations(t)
}

// Helper function to generate bcrypt hash with pepper
func bcryptGenerateWithPepper(secret string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(usecase.ClientPepper+secret), bcrypt.MinCost)
}
