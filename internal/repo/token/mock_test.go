package token

import (
	"context"
	"github.com/stretchr/testify/mock"
	"twitter-api/pkg/types"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) GetByToken(ctx context.Context, token string) (*Token, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*Token), args.Error(1)
}
func (m *MockRepo) Create(ctx context.Context, userID int,
	role types.UserRole) (accessToken, refreshToken string, err error) {
	args := m.Called(ctx, userID, role)
	return args.String(0), args.String(1), args.Error(2)
}
