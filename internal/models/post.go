package models

import (
	"time"
)

type Post struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	UserID      int       `json:"user_id"`
	Username    string    `json:"username"`
	Categories  []string  `json:"categories"`
	Likes       int       `json:"likes"`
	Dislikes    int       `json:"dislikes"`
	UserLiked   *bool     `json:"user_liked,omitempty"` // nil if not voted, true if liked, false if disliked
	Created     time.Time `json:"created"`
	CommentsCount int     `json:"comments_count"`
}

type Comment struct {
	ID       int       `json:"id"`
	PostID   int       `json:"post_id"`
	UserID   int       `json:"user_id"`
	Username string    `json:"username"`
	Content  string    `json:"content"`
	Likes    int       `json:"likes"`
	Dislikes int       `json:"dislikes"`
	UserLiked *bool    `json:"user_liked,omitempty"`
	Created  time.Time `json:"created"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Vote struct {
	ID     int  `json:"id"`
	UserID int  `json:"user_id"`
	PostID *int `json:"post_id,omitempty"`
	CommentID *int `json:"comment_id,omitempty"`
	IsLike bool `json:"is_like"`
}