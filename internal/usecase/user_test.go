package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bright-pentium/go-client-practice/internal/domain"

	userRepo "github.com/bright-pentium/go-client-practice/internal/repository/user"
	"github.com/bright-pentium/go-client-practice/internal/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// Helper function to hash passwords for testing
func hashPassword(password string) []byte {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(usecase.UserPepper+password), bcrypt.DefaultCost)
	return hashedPassword
}

func TestLoginUser(t *testing.T) {
	mockRepo := new(userRepo.MockUserRepository)
	userUseCase := usecase.NewUserUseCase(mockRepo)
	ctx := context.Background()

	// --- Test Case 1: Successful Login ---
	t.Run("successful login", func(t *testing.T) {
		account := "test@example.com"
		password := "correctpassword"
		hashedPass := hashPassword(password)

		expectedUser := &domain.User{
			ID:           uuid.New(),
			Account:      account,
			PasswordHash: hashedPass,
			Name:         "Test User",
		}

		// Set up the mock expectation: when GetUserByAccount is called, return expectedUser
		mockRepo.On("GetUserByAccount", ctx, account).Return(expectedUser, nil).Once()

		user, err := userUseCase.LoginUser(ctx, account, password)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.Account, user.Account)
		assert.Equal(t, expectedUser.Name, user.Name)
		mockRepo.AssertExpectations(t) // Verify that mock methods were called as expected
	})

	// --- Test Case 2: User Not Found ---
	t.Run("user not found", func(t *testing.T) {
		account := "nonexistent@example.com"
		password := "anypassword"

		// Set up the mock expectation: GetUserByAccount returns ErrUserNotFound
		mockRepo.On("GetUserByAccount", ctx, account).Return((*domain.User)(nil), domain.ErrUserNotFound).Once()

		user, err := userUseCase.LoginUser(ctx, account, password)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrUserLoginFail)) // Expect ErrUserLoginFail due to not found
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})

	// --- Test Case 3: Incorrect Password ---
	t.Run("incorrect password", func(t *testing.T) {
		account := "test@example.com"
		correctPassword := "correctpassword"
		incorrectPassword := "wrongpassword"
		hashedPass := hashPassword(correctPassword)

		existingUser := &domain.User{
			ID:           uuid.New(),
			Account:      account,
			PasswordHash: hashedPass,
			Name:         "Another User",
		}

		// Set up the mock expectation: GetUserByAccount returns the user
		mockRepo.On("GetUserByAccount", ctx, account).Return(existingUser, nil).Once()

		user, err := userUseCase.LoginUser(ctx, account, incorrectPassword)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrUserLoginFail)) // Expect ErrUserLoginFail due to wrong password
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})

	// --- Test Case 4: Repository Internal Error ---
	t.Run("repository internal error", func(t *testing.T) {
		account := "error@example.com"
		password := "anypassword"
		repoError := errors.New("database connection failed") // Simulate a generic DB error

		// Set up the mock expectation: GetUserByAccount returns a generic error
		mockRepo.On("GetUserByAccount", ctx, account).Return((*domain.User)(nil), repoError).Once()

		user, err := userUseCase.LoginUser(ctx, account, password)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, repoError))                // Expect the original repository error to be propagated
		assert.False(t, errors.Is(err, domain.ErrUserLoginFail)) // Should NOT be a login fail error
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})

	// --- Test Case 5: Hashed password is corrupted (very unlikely but good for robustness) ---
	t.Run("corrupted hashed password", func(t *testing.T) {
		account := "corrupt@example.com"
		password := "anypassword"
		// Simulate a corrupted hash that bcrypt.CompareHashAndPassword cannot process
		corruptedHash := []byte("not_a_valid_bcrypt_hash")

		existingUser := &domain.User{
			ID:           uuid.New(),
			Account:      account,
			PasswordHash: corruptedHash, // Corrupted hash
			Name:         "Corrupt User",
		}

		mockRepo.On("GetUserByAccount", ctx, account).Return(existingUser, nil).Once()

		user, err := userUseCase.LoginUser(ctx, account, password)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrUserLoginFail)) // bcrypt.CompareHashAndPassword will return error
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})
}
