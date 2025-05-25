package comment

import "time"

type Comment struct {
	ID        int        `json:"id"`
	UserId    int        `json:"user_id"`
	PostId    int        `json:"post_id"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
type CommentDTO struct {
	ID      int    `json:"id"`
	UserId  int    `json:"user_id"`
	PostId  int    `json:"post_id"`
	Content string `json:"content"`
}

type CommentInfo struct {
	ID        int       `json:"id"`
	UserId    int       `json:"user_id"`
	PostId    int       `json:"post_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type GetAllCommentsDTO struct {
	ID        int       `json:"id"`
	UserId    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}
