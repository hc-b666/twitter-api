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
	u, err := s.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("err: %w", err)
	}

	return u, nil
}

func (s *Service) GetUserPosts(ctx context.Context, postDTO *post.PostDTO) (*post.Post, error) {
	u, err := s.postRepo.GetByUserID(ctx, postDTO.UserId)
	if err != nil {
		return nil, fmt.Errorf("err: %w", err)
	}

	return u, nil
}

func (s *Service) CreatePost(
	ctx context.Context,
	postDTO *post.PostDTO,
) (int, error) {
	id, err := s.postRepo.Create(ctx, postDTO)
	if err != nil {
		return 0, fmt.Errorf("err: %w", err)
	}

	return id, nil
}
