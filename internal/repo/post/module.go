package post

import (
	"context"
	"errors"
	"fmt"
	"twitter-api/pkg/db"
)

type Repo interface {
	GetAll(ctx context.Context, limit, offset int) ([]*GetAllPostsDTO, error)
	GetByID(ctx context.Context, id int) (*GetAllPostsDTO, error)
	GetByUserID(ctx context.Context, userID int) ([]*PostInfo, error)
	Create(ctx context.Context, userID int, content, fileURL string) (int, error)
	SoftDelete(ctx context.Context, id int) error
	HardDelete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, content, fileURL string) (*PostInfo, error)
	UpdateContent(ctx context.Context, id int, content string) error
	UpdateFileURL(ctx context.Context, id int, fileURL string) error
	IsAuthor(ctx context.Context, postID, userID int) (bool, error)
}
type repo struct {
	db db.Pool
}

func NewRepo(pool db.Pool) (Repo, error) {
	return &repo{
		db: pool,
	}, nil
}

func (r *repo) GetAll(ctx context.Context, limit, offset int) ([]*GetAllPostsDTO, error) {
	query := `
		select p.id, p.user_id, p.content, p.file_url, p.created_at, p.updated_at, u.email
		from post p
		join "user" u on p.user_id = u.id
		where p.deleted_at is null
 		order by created_at desc
		limit $1 offset $2;
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get all posts: %w", err)
	}
	defer rows.Close()

	var posts []*GetAllPostsDTO
	for rows.Next() {
		post := &GetAllPostsDTO{}
		err := rows.Scan(
			&post.ID,
			&post.UserId,
			&post.Content,
			&post.FileURL,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Email,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over posts: %w", err)
	}

	return posts, nil
}

func (r *repo) GetByID(ctx context.Context, id int) (*GetAllPostsDTO, error) {
	query := `
		select p.id, p.user_id, p.content, p.file_url, p.created_at, p.updated_at, u.email
		from post p
		join "user" u on p.user_id = u.id
		where p.id = $1 and p.deleted_at is null;
	`

	post := &GetAllPostsDTO{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&post.ID,
		&post.UserId,
		&post.Content,
		&post.FileURL,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Email,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get post by id: %w", err)
	}

	return post, nil
}

func (r *repo) GetByUserID(ctx context.Context, userID int) ([]*PostInfo, error) {
	query := `
		select id, user_id, content, file_url, created_at, updated_at
		from post
		where user_id = $1 and deleted_at is null
		order by created_at desc;
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts by user_id: %w", err)
	}
	defer rows.Close()

	var posts []*PostInfo
	for rows.Next() {
		post := &PostInfo{}
		err := rows.Scan(
			&post.ID,
			&post.UserId,
			&post.Content,
			&post.FileURL,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over posts: %w", err)
	}

	return posts, nil
}

func (r *repo) Create(ctx context.Context, userID int, content, fileURL string) (int, error) {
	var id int
	query := `
		insert into post (user_id, content, file_url)
		values ($1, $2, $3)
		returning id;
	`

	err := r.db.QueryRow(
		ctx,
		query,
		userID,
		content,
		fileURL,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create post: %w", err)
	}

	return id, nil
}

func (r *repo) SoftDelete(ctx context.Context, id int) error {
	query := `
		update post
		set deleted_at = now()
		where id = $1;
	`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to soft delete post by id: %w", err)
	}

	return nil
}

func (r *repo) HardDelete(ctx context.Context, id int) error {
	query := `delete from post
    			where id = $1;`

	err := r.db.QueryRow(ctx, query, id)
	if err != nil {
		return errors.New("failed to delete comment by id")
	}

	return nil
}

func (r *repo) Update(ctx context.Context, id int, content, fileURL string) (*PostInfo, error) {
	query := `update post
		set content=$1, file_url=$2, updated_at=now()
		where id = $3;`
	post := &PostInfo{}
	err := r.db.QueryRow(ctx, query, content, fileURL, id).Scan(
		&post.ID,
		&post.UserId,
		&post.Content,
		&post.FileURL,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update post by id: %w", err)
	}

	return post, nil
}

func (r *repo) UpdateContent(ctx context.Context, id int, content string) error {
	query := `
		update post
		set content = $1, updated_at = now()
		where id = $2
	`

	_, err := r.db.Exec(ctx, query, content, id)
	if err != nil {
		return fmt.Errorf("failed to update post content: %w", err)
	}

	return nil
}

func (r *repo) UpdateFileURL(ctx context.Context, id int, fileURL string) error {
	query := `
		update post
		set file_url = $1, updated_at = now()
		where id = $2
	`

	_, err := r.db.Exec(ctx, query, fileURL, id)
	if err != nil {
		return fmt.Errorf("failed to update post file URL: %w", err)
	}

	return nil
}

func (r *repo) IsAuthor(ctx context.Context, postID, userID int) (bool, error) {
	query := `
		select exists(
			select 1
			from post
			where id = $1 and user_id = $2
		);
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, postID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if user is author: %w", err)
	}

	return exists, nil
}
