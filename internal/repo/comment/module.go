package comment

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) GetByID(ctx context.Context, id int) (*CommentInfo, error) {
	query := `
		select id, user_id, post_id,content, created_at, updated_at
		from "comment"
		where id = $1;
	`

	comment := &CommentInfo{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&comment.ID,
		&comment.UserId,
		&comment.PostId,
		&comment.Content,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get comment by id: %w", err)
	}

	return comment, nil
}

func (r *Repo) GetByUserID(ctx context.Context, userID int) (*Comment, error) {
	query := `
		select id, user_id,post_id, content, created_at, updated_at, deleted_at
		from "comment"
		where user_id = $1;
	`

	comment := &Comment{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&comment.ID,
		&comment.UserId,
		&comment.PostId,
		&comment.Content,
		&comment.CreatedAt,
		&comment.UpdatedAt,
		&comment.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get comment by user_id: %w", err)
	}

	return comment, nil
}

func (r *Repo) Create(ctx context.Context, comment *CommentDTO) (int, error) {
	var id int
	query := `
		insert into "comment" (user_id,post_id, content)
		values ($1, $2,$3)
		returning id;
	`

	err := r.db.QueryRow(
		ctx,
		query,
		comment.UserId,
		comment.PostId,
		comment.Content,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create comment: %w", err)
	}

	return id, nil
}
