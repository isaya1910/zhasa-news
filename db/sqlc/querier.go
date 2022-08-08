// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package db

import (
	"context"
)

type Querier interface {
	AddLike(ctx context.Context, arg AddLikeParams) (Like, error)
	CreateComment(ctx context.Context, arg CreateCommentParams) (Comment, error)
	CreateOrUpdateUser(ctx context.Context, arg CreateOrUpdateUserParams) (User, error)
	CreatePost(ctx context.Context, arg CreatePostParams) (Post, error)
	DeleteComment(ctx context.Context, id int32) error
	DeleteLike(ctx context.Context, arg DeleteLikeParams) error
	DeletePost(ctx context.Context, id int32) error
	DeleteUser(ctx context.Context, id int32) error
	// Example queries for sqlc
	GetCommentById(ctx context.Context, id int32) (Comment, error)
	GetCommentsAndAuthorsByPostId(ctx context.Context, postID int32) ([]GetCommentsAndAuthorsByPostIdRow, error)
	GetPostById(ctx context.Context, id int32) (Post, error)
	GetPostLikedUsers(ctx context.Context, arg GetPostLikedUsersParams) ([]GetPostLikedUsersRow, error)
	GetPostLikesCount(ctx context.Context, postID int32) (int64, error)
	GetPostsAndPostAuthors(ctx context.Context, arg GetPostsAndPostAuthorsParams) ([]GetPostsAndPostAuthorsRow, error)
	// Example queries for sqlc
	GetUser(ctx context.Context, id int32) (User, error)
	GetUserPostLike(ctx context.Context, arg GetUserPostLikeParams) (int32, error)
	ListPosts(ctx context.Context) ([]Post, error)
	ListUsers(ctx context.Context) ([]User, error)
}

var _ Querier = (*Queries)(nil)
