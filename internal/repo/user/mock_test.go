package user

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) GetAll(ctx context.Context) ([]*AdminUser, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*AdminUser), args.Error(1)
}
func (m *MockRepo) GetByID(ctx context.Context, id int) (*UserProfile, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*UserProfile), args.Error(1)
}
func (m *MockRepo) GetByEmail(ctx context.Context, email string) (*User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*User), args.Error(1)
}
func (m *MockRepo) Create(ctx context.Context, userDTO *RegisterUserDTO) (int, error) {
	args := m.Called(ctx, userDTO)
	return args.Int(0), args.Error(1)
}
func (m *MockRepo) CreateAdmin(ctx context.Context, userDTO *RegisterUserDTO) (int, error) {
	args := m.Called(ctx, userDTO)
	return args.Int(0), args.Error(1)
}
