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

func (s *Service) SoftDeleteComment(ctx context.Context, id int) (string, error) {
	err := s.commentRepo.SoftDelete(ctx, id)

	if err != nil {
		return "", fmt.Errorf("err: %w", err)
	}

	msg := " comment successfully deleted"
	return msg, nil
}
func (s *Service) HardDeleteComment(ctx context.Context, id int) error {
	err := s.commentRepo.HardDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("err: %w", err)
	}

	return nil
}

func (s *Service) UpdateComment(ctx context.Context, id int, content string) (*comment.CommentInfo, error) {
	commentInfo, err := s.commentRepo.Update(ctx, id, content)
	if err != nil {
		return nil, fmt.Errorf("err: %w", err)
	}
	return commentInfo, nil
}

func (s *Service) GetALlPostComments(
	ctx context.Context,
	postId int) ([]*comment.GetAllCommentsDTO, error) {
	comments, err := s.commentRepo.GetAllCommentsToPost(ctx, postId)
	if err != nil {
		return nil, fmt.Errorf("err: %w", err)
	}

	return comments, nil
}

func (s *Service) GetALlCommentsByAdmin(
	ctx context.Context) ([]*comment.Comment, error) {
	comments, err := s.commentRepo.GetAllComments(ctx)
	if err != nil {
		return nil, fmt.Errorf("err: %w", err)
	}

	return comments, nil
}
