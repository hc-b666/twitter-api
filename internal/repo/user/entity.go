package user

import (
	"time"
	"twitter-api/pkg/types"
)

type User struct {
	ID        int            `json:"id"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	Role      types.UserRole `json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt *time.Time     `json:"deleted_at"`
}

type RegisterUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserProfile struct {
	ID        int            `json:"id"`
	Email     string         `json:"email"`
	Role      types.UserRole `json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type AdminUser struct {
	ID            int            `json:"id"`
	Email         string         `json:"email"`
	Role          types.UserRole `json:"role"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     *time.Time     `json:"deleted_at"`
	PostsCount    int            `json:"posts_count"`
	CommentsCount int            `json:"comments_count"`
}
