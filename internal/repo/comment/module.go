package comment

import (
	"context"
	"fmt"
	"twitter-api/pkg/db"
)

type Repo interface {
	Create(ctx context.Context, userID, postID int, comment *CommentDTO) (int, error)
	GetByID(ctx context.Context, id int) (*CommentInfo, error)
	GetByUserID(ctx context.Context, userID int) ([]*UserComment, error)
	GetAllCommentsToPost(ctx context.Context, postId, limit, offset int) ([]*GetAllCommentsDTO, error)
	GetAllComments(ctx context.Context, limit, offset int) ([]*Comment, error)
	Update(ctx context.Context, id int, content string) error
	HardDelete(ctx context.Context, id int) error
	SoftDelete(ctx context.Context, id int) error
	IsAuthor(ctx context.Context, userID, commentID int) (bool, error)
}

type repo struct {
	db db.Pool
}

func NewRepo(pool db.Pool) (Repo, error) {
	return &repo{
		db: pool,
	}, nil
}
func (r *repo) GetByID(ctx context.Context, id int) (*CommentInfo, error) {
	query := `
		select id, user_id, post_id,content, created_at, updated_at
		from comment
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

func (r *repo) GetByUserID(ctx context.Context, userID int) ([]*UserComment, error) {
	query := `
		select c.id, c.user_id, c.post_id, c.content, c.created_at, u.email 
		from comment c
		join "user" u on c.user_id = u.id
		where user_id = $1 and c.deleted_at is null;
	`

	comments := []*UserComment{}
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments by user id: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		comment := &UserComment{}
		err := rows.Scan(
			&comment.ID,
			&comment.UserId,
			&comment.PostId,
			&comment.Content,
			&comment.CreatedAt,
			&comment.Email,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over comments: %w", err)
	}

	return comments, nil
}

func (r *repo) Create(ctx context.Context, userID, postID int, comment *CommentDTO) (int, error) {
	var id int
	query := `
		insert into comment (user_id, post_id, content)
		values ($1, $2, $3)
		returning id;
	`
	err := r.db.QueryRow(
		ctx,
		query,
		userID,
		postID,
		comment.Content,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create comment: %w", err)
	}

	return id, nil
}

func (r *repo) SoftDelete(ctx context.Context, id int) error {
	query := `
		update comment
    set deleted_at = now()
		where id = $1;
	`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to soft delete comment by id: %w", err)
	}

	return nil
}

func (r *repo) HardDelete(ctx context.Context, id int) error {
	query := `delete from comment
			where id = $1;`

	err := r.db.QueryRow(ctx, query, id).Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to delete comment by id: %w", err)
	}

	return nil
}

func (r *repo) Update(ctx context.Context, id int, content string) error {
	query := `
		update comment
		set content = $1, updated_at = now()
		where id = $2;
	`

	_, err := r.db.Exec(ctx, query, content, id)
	if err != nil {
		return fmt.Errorf("failed to update comment by id: %w", err)
	}

	return nil
}

func (r *repo) GetAllCommentsToPost(ctx context.Context, postId, limit, offset int) ([]*GetAllCommentsDTO, error) {
	query := `
		select c.id, c.user_id, c.content, c.created_at, c.updated_at, u.email
		from comment c 
		join "user" u on c.user_id = u.id
		where c.post_id = $1 and c.deleted_at IS NULL
 		order by created_at desc
		limit $2 offset $3;
	`

	rows, err := r.db.Query(ctx, query, postId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get all comments to posts: %w", err)
	}
	defer rows.Close()

	var comments []*GetAllCommentsDTO
	for rows.Next() {
		comment := &GetAllCommentsDTO{}
		err := rows.Scan(
			&comment.ID,
			&comment.UserId,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.Email,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over posts: %w", err)
	}

	return comments, nil
}

func (r *repo) GetAllComments(ctx context.Context, limit, offset int) ([]*Comment, error) {
	query := `
		select id, user_id, post_id,content, created_at, updated_at, deleted_at	
		from comment
		order by created_at desc
		limit $1 offset $2;
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get all comments: %w", err)
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		comment := &Comment{}
		err := rows.Scan(
			&comment.ID,
			&comment.UserId,
			&comment.PostId,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over posts: %w", err)
	}

	return comments, nil
}

func (r *repo) IsAuthor(ctx context.Context, userID, commentID int) (bool, error) {
	query := `
		select exists(
			select 1 from comment
			where id = $1 and user_id = $2
		);
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, commentID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if user is author: %w", err)
	}

	return exists, nil
}
