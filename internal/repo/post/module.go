package post

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

func (r *Repo) GetByID(ctx context.Context, id int) (*PostInfo, error) {
	query := `
		select id, user_id, content, created_at, updated_at
		from post
		where id = $1;
	`

	post := &PostInfo{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&post.ID,
		&post.UserId,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get post by id: %w", err)
	}

	return post, nil
}

func (r *Repo) GetByUserID(ctx context.Context, userID int) ([]*PostInfo, error) {
	query := `
		select id, user_id, content, created_at, updated_at
		from post
		where user_id = $1
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

func (r *Repo) Create(ctx context.Context, userID int, post *PostDTO) (int, error) {
	var id int
	query := `
		insert into post (user_id, content)
		values ($1, $2)
		returning id;
	`

	err := r.db.QueryRow(
		ctx,
		query,
		userID,
		post.Content,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create post: %w", err)
	}

	return id, nil
}
