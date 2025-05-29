package post

import (
	"context"
	"fmt"
	"twitter-api/internal/repo/post"
)

type Service struct {
	postRepo post.Repo
}

func NewService(r post.Repo) *Service {
	return &Service{
		postRepo: r,
	}
}

func (s *Service) GetAll(ctx context.Context, limit, offset int) ([]*post.GetAllPostsDTO, error) {
	posts, err := s.postRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("err: %w", err)
	}

	return posts, nil
}

func (s *Service) GetByID(ctx context.Context, id int) (*post.GetAllPostsDTO, error) {
	p, err := s.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("err: %w", err)
	}

	return p, nil
}

func (s *Service) GetUserPosts(ctx context.Context, userID int) ([]*post.PostInfo, error) {
	posts, err := s.postRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("err: %w", err)
	}

	return posts, nil
}

func (s *Service) CreatePost(
	ctx context.Context,
	userID int,
	content string,
	fileURL string,
) (int, error) {
	id, err := s.postRepo.Create(ctx, userID, content, fileURL)
	if err != nil {
		return 0, fmt.Errorf("err: %w", err)
	}

	return id, nil
}
func (s *Service) UpdatePost(
	ctx context.Context,
	postID int,
	content string,
	fileURL string,
) (*post.PostInfo, error) {
	p, err := s.postRepo.Update(ctx, postID, content, fileURL)
	if err != nil {
		return nil, fmt.Errorf("err: %w", err)
	}

	return p, nil
}
func (s *Service) SoftDeletePost(ctx context.Context, id int) error {
	err := s.postRepo.SoftDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("err: %w", err)
	}

	return nil
}
func (s *Service) HardDeletePost(ctx context.Context, id int) (string, error) {
	err := s.postRepo.HardDelete(ctx, id)
	if err != nil {
		return "", fmt.Errorf("err: %w", err)
	}

	return "post is deleted successfully", nil
}

func (s *Service) UpdatePostContent(ctx context.Context, postID int, content string) error {
	err := s.postRepo.UpdateContent(ctx, postID, content)
	if err != nil {
		return fmt.Errorf("err: %w", err)
	}

	return nil
}

func (s *Service) UpdatePostFileURL(ctx context.Context, postID int, fileURL string) error {
	err := s.postRepo.UpdateFileURL(ctx, postID, fileURL)
	if err != nil {
		return fmt.Errorf("err: %w", err)
	}

	return nil
}

func (s *Service) IsAuthor(ctx context.Context, postID, userID int) (bool, error) {
	isAuthor, err := s.postRepo.IsAuthor(ctx, postID, userID)
	if err != nil {
		return false, fmt.Errorf("err: %w", err)
	}

	return isAuthor, nil
}
