package models

import (
	"time"
)

type User struct {
	ID       int       `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Password string    `json:"-"` // Don't include in JSON responses
	Created  time.Time `json:"created"`
}

type Session struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	Created   time.Time `json:"created"`
}