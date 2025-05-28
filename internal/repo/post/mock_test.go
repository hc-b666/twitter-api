package post

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) GetAll(ctx context.Context, limit, offset int) ([]*GetAllPostsDTO, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*GetAllPostsDTO), args.Error(1)
}
func (m *MockRepo) GetByID(ctx context.Context, id int) (*GetAllPostsDTO, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*GetAllPostsDTO), args.Error(1)
}
func (m *MockRepo) GetByUserID(ctx context.Context, userID int) ([]*PostInfo, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*PostInfo), args.Error(1)
}
func (m *MockRepo) Create(ctx context.Context, userID int, content, fileURL string) (int, error) {
	args := m.Called(ctx, userID, content, fileURL)
	return args.Int(0), args.Error(1)
}
func (m *MockRepo) SoftDelete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockRepo) HardDelete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockRepo) Update(ctx context.Context, id int, content, fileURL string) (*PostInfo, error) {
	args := m.Called(ctx, id, content, fileURL)
	return args.Get(0).(*PostInfo), args.Error(1)
}
func (m *MockRepo) UpdateContent(ctx context.Context, id int, content string) error {
	args := m.Called(ctx, id, content)
	return args.Error(0)
}
func (m *MockRepo) UpdateFileURL(ctx context.Context, id int, fileURL string) error {
	args := m.Called(ctx, id, fileURL)
	return args.Error(0)
}
func (m *MockRepo) IsAuthor(ctx context.Context, postID, userID int) (bool, error) {
	args := m.Called(ctx, postID, userID)
	return args.Bool(0), args.Error(1)
}
