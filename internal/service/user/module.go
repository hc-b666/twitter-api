package user

import (
	"context"
	"fmt"
	"twitter-api/internal/repo/user"
	"twitter-api/pkg/errs"
	"twitter-api/pkg/utils"
)

type Service struct {
	userRepo *user.Repo
}

func NewService(userRepo *user.Repo) *Service {
	return &Service{
		userRepo: userRepo,
	}
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
