package post

import (
	"time"
)

type Post struct {
	ID        int        `json:"id"`
	UserId    int        `json:"user_id"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
type PostDTO struct {
	Content string `json:"content"`
}

type PostInfo struct {
	ID        int       `json:"id"`
	UserId    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
