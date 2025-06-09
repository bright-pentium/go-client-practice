package usecase_test

import (
	"context"
	"testing"

	"github.com/bright-pentium/go-client-practice/internal/domain"
	mockRepo "github.com/bright-pentium/go-client-practice/internal/repository/user"
	"github.com/bright-pentium/go-client-practice/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const UserPepper = "3tw0d2o" // Ensure this matches the actual constant

func TestCreateUserSuccess(t *testing.T) {
	mockRepo := new(mockRepo.MockUserRepository)
	useCase := usecase.NewSysUserUseCase(mockRepo)

	name := "Test User"
	account := "testuser"
	password := "password123"
	var capturedHash []byte

	mockRepo.On("CreateUser", mock.Anything, mock.Anything, name, account, mock.AnythingOfType("[]uint8")).
		Run(func(args mock.Arguments) {
			capturedHash = args.Get(4).([]byte)
		}).
		Return(&domain.User{ID: uuid.New(), Name: name, Account: account}, nil)

	ctx := context.Background()
	user, err := useCase.CreateUser(ctx, name, account, password)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotNil(t, capturedHash)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID(t *testing.T) {
	mockRepo := new(mockRepo.MockUserRepository)
	useCase := usecase.NewSysUserUseCase(mockRepo)

	id := uuid.New()
	expectedUser := &domain.User{ID: id, Name: "Alice"}

	mockRepo.On("GetUserByID", mock.Anything, id).Return(expectedUser, nil)

	ctx := context.Background()
	user, err := useCase.GetUserByID(ctx, id)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUserByIDSuccess(t *testing.T) {
	mockRepo := new(mockRepo.MockUserRepository)
	useCase := usecase.NewSysUserUseCase(mockRepo)

	id := uuid.New()
	name := "Updated Name"
	password := "newpassword"
	var capturedHash []byte

	mockRepo.On("UpdateUserByID", mock.Anything, id, name, mock.AnythingOfType("[]uint8")).
		Run(func(args mock.Arguments) {
			capturedHash = args.Get(3).([]byte)
		}).
		Return(&domain.User{ID: id, Name: name}, nil)

	ctx := context.Background()
	user, err := useCase.UpdateUserByID(ctx, id, name, password)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotNil(t, capturedHash)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUserByID(t *testing.T) {
	mockRepo := new(mockRepo.MockUserRepository)
	useCase := usecase.NewSysUserUseCase(mockRepo)

	id := uuid.New()
	mockRepo.On("DeleteUserByID", mock.Anything, id).Return(nil)

	ctx := context.Background()
	err := useCase.DeleteUserByID(ctx, id)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateUserHashFail(t *testing.T) {
	useCase := usecase.NewSysUserUseCase(nil)

	// Override bcrypt to fail intentionally via an invalid cost (not directly mockable)
	longPassword := string(make([]byte, 1<<20)) // huge password to likely trigger bcrypt error
	ctx := context.Background()

	user, err := useCase.CreateUser(ctx, "Foo", "bar", longPassword)

	assert.Nil(t, user)
	assert.Error(t, err)
}
