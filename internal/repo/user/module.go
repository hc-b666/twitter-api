package user

import (
	"context"
	"fmt"
	"twitter-api/pkg/logger"
	"twitter-api/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
	l  *logger.Logger
}

func NewRepo(db *pgxpool.Pool, l *logger.Logger) *Repo {
	return &Repo{
		db: db,
		l:  l,
	}
}

func (r *Repo) GetAll(ctx context.Context) ([]*UserProfile, error) {
	query := `
		select id, email, role, created_at, updated_at
		from "user"
		where deleted_at is null
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	defer rows.Close()

	var users []*UserProfile
	for rows.Next() {
		user := &UserProfile{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
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

func (r *Repo) GetByID(ctx context.Context, id int) (*UserProfile, error) {
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

func (r *Repo) GetByEmail(ctx context.Context, email string) (*User, error) {
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

func (r *Repo) Create(ctx context.Context, userDTO *RegisterUserDTO) (int, error) {
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

func (r *Repo) CreateAdmin(ctx context.Context, userDTO *RegisterUserDTO) (int, error) {
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
