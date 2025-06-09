package usecase

import (
	"context"

	"github.com/bright-pentium/go-client-practice/internal/domain"
	"github.com/google/uuid"
)

type ResourceUseCase struct{}

func NewResourceUseCase() *ResourceUseCase {
	return &ResourceUseCase{}
}

func (u *ResourceUseCase) CreateResource(ctx context.Context) (*domain.Resource, error) {
	return &domain.Resource{ID: uuid.New()}, nil
}
