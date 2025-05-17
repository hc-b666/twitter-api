package user

import (
	"context"
	"fmt"
	"twitter-api/pkg/utils"

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
