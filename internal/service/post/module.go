package post

import (
	"context"
	"fmt"
	"twitter-api/internal/repo/post"
)

type Service struct {
	postRepo *post.Repo
}

func NewService(postRepo *post.Repo) *Service {
	return &Service{
		postRepo: postRepo,
	}
}

func (s *Service) GetByID(ctx context.Context, id int) (*post.PostInfo, error) {
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
	postDTO *post.PostDTO,
) (int, error) {
	id, err := s.postRepo.Create(ctx, userID, postDTO)
	if err != nil {
		return 0, fmt.Errorf("err: %w", err)
	}

	return id, nil
}
