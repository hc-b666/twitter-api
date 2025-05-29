package token

import (
	"context"
	"fmt"
	"time"
	"twitter-api/internal/repo/token"
	"twitter-api/pkg/errs"
	"twitter-api/pkg/types"
	"twitter-api/pkg/utils"
)

type Service struct {
	tokenRepository token.Repo
}

func NewService(r token.Repo) *Service {
	return &Service{
		tokenRepository: r,
	}
}
func (s *Service) GetByToken(ctx context.Context, tokenStr string) (*token.Token, error) {
	tkn, err := s.tokenRepository.GetByToken(ctx, tokenStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if err = utils.VerifyRefreshToken(tokenStr, tkn.UserID, tkn.Role); err != nil {
		return nil, fmt.Errorf("failed to verify refresh token: %w", err)
	}

	if tkn.ExpiresAt.Before(time.Now()) {
		return nil, errs.ErrTokenExpired
	}

	return tkn, nil
}

func (s *Service) CreateTokens(
	ctx context.Context,
	userID int,
	role types.UserRole,
) (
	accessToken, refreshToken string,
	err error,
) {
	accessToken, refreshToken, err = s.tokenRepository.Create(ctx, userID, role)
	if err != nil {
		return "", "", fmt.Errorf("failed to create session: %w", err)
	}

	return accessToken, refreshToken, nil
}
