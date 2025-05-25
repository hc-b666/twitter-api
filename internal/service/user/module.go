package user

import (
	"context"
	"fmt"
	"twitter-api/internal/repo/user"
	"twitter-api/pkg/errs"
	"twitter-api/pkg/logger"
	"twitter-api/pkg/utils"
)

type Service struct {
	userRepo *user.Repo
	l        *logger.Logger
}

func NewService(userRepo *user.Repo, l *logger.Logger) *Service {
	return &Service{
		userRepo: userRepo,
		l:        l,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]*user.AdminUser, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("err: %w", err)
	}

	return users, nil
}

func (s *Service) GetByID(ctx context.Context, id int) (*user.UserProfile, error) {
	u, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("err: %w", err)
	}

	return u, nil
}

func (s *Service) ValidateUser(ctx context.Context, userDTO *user.LoginUserDTO) (*user.User, error) {
	u, err := s.userRepo.GetByEmail(ctx, userDTO.Email)
	if err != nil {
		return nil, fmt.Errorf("err: %w", err)
	}

	ok := utils.ComparePassword(u.Password, userDTO.Password)
	if !ok {
		return nil, errs.ErrInvalidCredentials
	}

	return u, nil
}

func (s *Service) CreateUser(
	ctx context.Context,
	userDTO *user.RegisterUserDTO,
) (int, error) {
	id, err := s.userRepo.Create(ctx, userDTO)
	if err != nil {
		return 0, fmt.Errorf("err: %w", err)
	}

	return id, nil
}

func (s *Service) CreateAdmin(
	ctx context.Context,
	userDTO *user.RegisterUserDTO,
) (int, error) {
	id, err := s.userRepo.CreateAdmin(ctx, userDTO)
	if err != nil {
		s.l.Error("failed to create admin", err)
		return 0, fmt.Errorf("err: %w", err)
	}

	return id, nil
}
