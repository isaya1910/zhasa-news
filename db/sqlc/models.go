// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package db

import (
	"time"
)

type Comment struct {
	ID        int32     `json:"id"`
	Body      string    `json:"body"`
	UserID    int32     `json:"user_id"`
	PostID    int32     `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Post struct {
	ID        int32     `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	UserID    int32     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID        int32  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Bio       string `json:"bio"`
}
