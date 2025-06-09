package usecase

import (
	"context"
	"fmt"

	"github.com/bright-pentium/go-client-practice/internal/domain"
	"github.com/bright-pentium/go-client-practice/internal/repository/user"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type SysUserUseCase struct {
	repo user.IUserRepository
}

func NewSysUserUseCase(repo user.IUserRepository) *SysUserUseCase {
	return &SysUserUseCase{repo: repo}
}

func (u *SysUserUseCase) CreateUser(ctx context.Context, name string, account string, password string) (*domain.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(UserPepper+password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", domain.ErrUserHashFail, err)
	}
	return u.repo.CreateUser(ctx, uuid.New(), name, account, passwordHash)
}

func (u *SysUserUseCase) GetUserByID(ctx context.Context, ID uuid.UUID) (*domain.User, error) {
	return u.repo.GetUserByID(ctx, ID)
}

func (u *SysUserUseCase) UpdateUserByID(ctx context.Context, ID uuid.UUID, name string, password string) (*domain.User, error) {
	var passwordHash []byte
	var err error
	if password != "" {
		passwordHash, err = bcrypt.GenerateFromPassword([]byte(UserPepper+password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", domain.ErrUserHashFail, err)
		}
	}
	return u.repo.UpdateUserByID(ctx, ID, name, passwordHash)
}

func (u *SysUserUseCase) DeleteUserByID(ctx context.Context, ID uuid.UUID) error {
	return u.repo.DeleteUserByID(ctx, ID)
}
