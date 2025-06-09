package usecase

import (
	"context"
	"errors"

	"github.com/bright-pentium/go-client-practice/internal/domain"
	repo "github.com/bright-pentium/go-client-practice/internal/repository/user"
	"golang.org/x/crypto/bcrypt"
)

const (
	UserPepper string = "3tw0d2o"
)

type UserUseCase struct {
	repo repo.IUserRepository
}

func NewUserUseCase(repo repo.IUserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (u *UserUseCase) LoginUser(ctx context.Context, account string, password string) (*domain.User, error) {
	user, err := u.repo.GetUserByAccount(ctx, account)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrUserLoginFail
		} else {
			return nil, err
		}
	}

	if bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(UserPepper+password)) != nil {
		return nil, domain.ErrUserLoginFail
	}
	return user, nil
}
