package user

import (
	"context"
	"fmt"
	"twitter-api/pkg/db"
	"twitter-api/pkg/logger"
	"twitter-api/pkg/utils"
)

type Repo interface {
	GetAll(ctx context.Context) ([]*AdminUser, error)
	GetByID(ctx context.Context, id int) (*UserProfile, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, userDTO *RegisterUserDTO) (int, error)
	CreateAdmin(ctx context.Context, userDTO *RegisterUserDTO) (int, error)
}
type repo struct {
	db db.Pool
	l  *logger.Logger
}

func NewRepo(pool db.Pool, l *logger.Logger) (Repo, error) {
	return &repo{
		db: pool,
		l:  l,
	}, nil
}

func (r *repo) GetAll(ctx context.Context) ([]*AdminUser, error) {
	query := `
		select 
			u.id, 
			u.email, 
			u.role, 
			u.created_at, 
			u.updated_at, 
			u.deleted_at, 
			count(distinct p.id) as posts_count, 
			count(distinct c.id) as comments_count
		from "user" u
		left join post p on p.user_id = u.id
		left join comment c on c.user_id = u.id
		group by u.id, u.email, u.role, u.created_at, u.updated_at, u.deleted_at
		order by u.created_at desc;
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	defer rows.Close()

	var users []*AdminUser
	for rows.Next() {
		user := &AdminUser{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.PostsCount,
			&user.CommentsCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over users: %w", err)
	}

	return users, nil
}

func (r *repo) GetByID(ctx context.Context, id int) (*UserProfile, error) {
	query := `
		select id, email, role, created_at, updated_at
		from "user"
		where id = $1;
	`

	user := &UserProfile{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

func (r *repo) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		select id, email, password, role, created_at, updated_at, deleted_at
		from "user"
		where email = $1;
	`

	user := &User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

func (r *repo) Create(ctx context.Context, userDTO *RegisterUserDTO) (int, error) {
	var id int
	query := `
		insert into "user" (email, password)
		values ($1, $2)
		returning id;
	`

	hashedPassword, err := utils.HashPassword(userDTO.Password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash: %w", err)
	}

	err = r.db.QueryRow(
		ctx,
		query,
		userDTO.Email,
		hashedPassword,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func (r *repo) CreateAdmin(ctx context.Context, userDTO *RegisterUserDTO) (int, error) {
	var id int
	query := `
		insert into "user" (email, password, role)
		values ($1, $2, 'admin')
		returning id;
	`

	hashedPassword, err := utils.HashPassword(userDTO.Password)
	if err != nil {
		r.l.Error("failed to hash password", err)
		return 0, fmt.Errorf("failed to hash: %w", err)
	}

	err = r.db.QueryRow(
		ctx,
		query,
		userDTO.Email,
		hashedPassword,
	).Scan(&id)
	if err != nil {
		r.l.Error("failed to create user", err)
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}
