package token

import (
	"context"
	"fmt"
	"twitter-api/pkg/types"
	"twitter-api/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) GetByToken(ctx context.Context, token string) (*Token, error) {
	query := `
		select id, token, user_id, role, expires_at
		from refresh_token
		where token = $1;
	`

	tokenData := &Token{}
	err := r.db.QueryRow(ctx, query, token).Scan(
		&tokenData.ID,
		&tokenData.Token,
		&tokenData.UserID,
		&tokenData.Role,
		&tokenData.ExpiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	return tokenData, nil
}

func (r *Repo) Create(
	ctx context.Context,
	userID int,
	role types.UserRole,
) (
	accessToken, refreshToken string,
	err error,
) {
	query := `
		insert into refresh_token (token, user_id, role, expires_at)
		values ($1, $2, $3, $4);
	`

	accessToken, refreshToken, err = utils.GenerateJwtTokens(userID, role)
	if err != nil {
		return "", "", fmt.Errorf("failed to created tokens: %w", err)
	}

	_, err = r.db.Exec(ctx,
		query,
		refreshToken,
		userID,
		role,
		utils.GetRefreshExpireTime(),
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to save refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}
