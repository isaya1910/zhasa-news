package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	CreateUserTx(ctx context.Context, userArg CreateOrUpdateUserParams) (User, error)
	CreatePostTx(ctx context.Context, postArg CreatePostParams) (Post, error)
	CreateCommentTx(ctx context.Context, commentArg CreateCommentParams) (Comment, error)
	AddLikeTx(ctx context.Context, params AddLikeParams) (Like, error)
	DeleteLikeTx(ctx context.Context, params DeleteLikeParams) error
}

// SQLStore SQLSQLStore Store provides all functions to executed queries transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *SQLStore {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

func (store *SQLStore) CreatePostTx(ctx context.Context, postArg CreatePostParams) (Post, error) {
	var resultPost Post
	err := store.execTx(ctx, func(queries *Queries) error {
		var err error
		resultPost, err = queries.CreatePost(ctx, postArg)
		if err != nil {
			return err
		}
		createImageParams := CreatePostImageParams{PostID: resultPost.ID, ImageUrl: postArg.ImageUrl}
		if len(createImageParams.ImageUrl) == 0 {
			return nil
		}
		_, err = queries.CreatePostImage(ctx, createImageParams)
		if err != nil {
			return err
		}
		return nil
	})
	return resultPost, err
}

func (store *SQLStore) CreateCommentTx(ctx context.Context, commentArg CreateCommentParams) (Comment, error) {
	var resultComment Comment
	err := store.execTx(ctx, func(queries *Queries) error {
		var err error
		if err != nil {
			return err
		}
		resultComment, err = queries.CreateComment(ctx, commentArg)
		if err != nil {
			return err
		}
		return nil
	})
	return resultComment, err
}

func (store *SQLStore) AddLikeTx(ctx context.Context, params AddLikeParams) (Like, error) {
	var resultLike Like
	err := store.execTx(ctx, func(queries *Queries) error {
		var err error
		resultLike, err = store.AddLike(ctx, params)
		if err != nil {
			return err
		}
		return nil
	})
	return resultLike, err
}

func (store *SQLStore) DeleteLikeTx(ctx context.Context, params DeleteLikeParams) error {
	err := store.execTx(ctx, func(queries *Queries) error {
		var err error
		err = store.DeleteLike(ctx, params)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (store *SQLStore) CreateUserTx(ctx context.Context, userArg CreateOrUpdateUserParams) (User, error) {
	var resultUser User
	err := store.execTx(ctx, func(queries *Queries) error {
		var err error
		resultUser, err = queries.CreateOrUpdateUser(ctx, userArg)
		if err != nil {
			return err
		}
		return nil
	})
	return resultUser, err
}
