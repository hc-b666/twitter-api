package token

import (
	"time"
	"twitter-api/pkg/types"
)

type Token struct {
	ID        int            `json:"id"`
	UserID    int            `json:"user_id"`
	Token     string         `json:"token"`
	Role      types.UserRole `json:"role"`
	ExpiresAt time.Time      `json:"expires_at"`
}
