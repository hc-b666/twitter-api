package comment

import (
	"context"
	"fmt"
	"twitter-api/internal/repo/comment"
)

type Service struct {
	commentRepo *comment.Repo
}

func NewService(commentRepo *comment.Repo) *Service {
	return &Service{
		commentRepo: commentRepo,
	}
}

func (s *Service) GetByID(ctx context.Context, id int) (*comment.CommentInfo, error) {
	u, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("err: %w", err)
	}

	return u, nil
}

func (s *Service) GetUserComments(ctx context.Context, commentDTO *comment.CommentDTO) (*comment.Comment, error) {
	u, err := s.commentRepo.GetByUserID(ctx, commentDTO.UserId)
	if err != nil {
		return nil, fmt.Errorf("err: %w", err)
	}

	return u, nil
}

func (s *Service) CreateComment(
	ctx context.Context,
	commentDTO *comment.CommentDTO,
) (int, error) {
	id, err := s.commentRepo.Create(ctx, commentDTO)
	if err != nil {
		return 0, fmt.Errorf("err: %w", err)
	}

	return id, nil
}
