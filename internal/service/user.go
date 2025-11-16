package service

import (
	"context"

	"github.com/ynsssss/pr-manager/internal/domain"
)

// TODO: move
type UserRepository interface {
	GetByID(ctx context.Context, userID string) (domain.User, error)

	SetIsActive(ctx context.Context, userID string, isActive bool) (domain.User, error)
	// TODO: rename method
	UpsertUsers(ctx context.Context, users []domain.User) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) SetIsActive(
	ctx context.Context,
	userID string,
	active bool,
) (domain.User, error) {
	return s.repo.SetIsActive(ctx, userID, active)
}
