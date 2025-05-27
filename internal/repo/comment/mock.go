package comment

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Create(ctx context.Context, userID, postID int, comment *CommentDTO) (int, error) {
	args := m.Called(userID, postID, comment)
	return args.Int(0), args.Error(1)
}
func (m *MockRepo) GetByID(ctx context.Context, id int) (*CommentInfo, error) {
	args := m.Called(id)
	return args.Get(0).(*CommentInfo), args.Error(1)
}
func (m *MockRepo) GetByUserID(ctx context.Context, userID int) ([]*UserComment, error) {
	args := m.Called(userID)
	return args.Get(0).([]*UserComment), args.Error(1)
}
func (m *MockRepo) GetAllCommentsToPost(ctx context.Context, postId, limit, offset int) ([]*GetAllCommentsDTO, error) {
	args := m.Called(postId, limit, offset)
	return args.Get(0).([]*GetAllCommentsDTO), args.Error(1)
}
func (m *MockRepo) GetAllComments(ctx context.Context, limit, offset int) ([]*Comment, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]*Comment), args.Error(1)
}
func (m *MockRepo) Update(ctx context.Context, id int, content string) error {
	args := m.Called(id, content)
	return args.Error(0)
}
func (m *MockRepo) HardDelete(ctx context.Context, id int) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockRepo) SoftDelete(ctx context.Context, id int) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockRepo) IsAuthor(ctx context.Context, userID, commentID int) (bool, error) {
	args := m.Called(userID, commentID)
	return args.Bool(0), args.Error(1)
}
