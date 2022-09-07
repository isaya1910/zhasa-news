package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	CreatePostTx(ctx context.Context, postArg CreatePostParams, imageUrl string, userArg CreateOrUpdateUserParams) (Post, User, error)
	CreateCommentTx(ctx context.Context, commentArg CreateCommentParams, userArg CreateOrUpdateUserParams) (Comment, User, error)
	AddLikeTx(ctx context.Context, postId int32, userArg CreateOrUpdateUserParams) (Like, error)
	DeleteLikeTx(ctx context.Context, postId int32, userArg CreateOrUpdateUserParams) error
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

func (store *SQLStore) CreatePostTx(ctx context.Context, postArg CreatePostParams, imageUrl string, userArg CreateOrUpdateUserParams) (Post, User, error) {
	var resultPost Post
	var resultUser User
	err := store.execTx(ctx, func(queries *Queries) error {
		var err error
		resultUser, err = queries.CreateOrUpdateUser(ctx, userArg)
		if err != nil {
			return err
		}
		postArg.UserID = resultUser.ID
		resultPost, err = queries.CreatePost(ctx, postArg)
		if len(imageUrl) == 0 {
			return nil
		}
		createImageParams := CreatePostImageParams{PostID: resultPost.ID, ImageUrl: imageUrl}
		_, err = queries.CreatePostImage(ctx, createImageParams)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return resultPost, resultUser, err
	}
	return resultPost, resultUser, err
}

func (store *SQLStore) CreateCommentTx(ctx context.Context, commentArg CreateCommentParams, userArg CreateOrUpdateUserParams) (Comment, User, error) {
	var resultComment Comment
	var resultUser User
	err := store.execTx(ctx, func(queries *Queries) error {
		var err error
		resultUser, err = queries.CreateOrUpdateUser(ctx, userArg)
		if err != nil {
			return err
		}
		commentArg.UserID = resultUser.ID
		resultComment, err = queries.CreateComment(ctx, commentArg)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return resultComment, resultUser, err
	}
	return resultComment, resultUser, err
}

func (store *SQLStore) AddLikeTx(ctx context.Context, postId int32, userArg CreateOrUpdateUserParams) (Like, error) {
	var resultLike Like
	err := store.execTx(ctx, func(queries *Queries) error {
		var err error
		var resultUser User
		resultUser, err = queries.CreateOrUpdateUser(ctx, userArg)
		if err != nil {
			return err
		}
		var addLikeParams AddLikeParams
		addLikeParams = AddLikeParams{
			UserID: resultUser.ID,
			PostID: postId,
		}
		resultLike, err = store.AddLike(ctx, addLikeParams)
		if err != nil {
			return err
		}
		return nil
	})
	return resultLike, err
}

func (store *SQLStore) DeleteLikeTx(ctx context.Context, postId int32, userArg CreateOrUpdateUserParams) error {
	err := store.execTx(ctx, func(queries *Queries) error {
		var err error
		var resultUser User
		resultUser, err = queries.CreateOrUpdateUser(ctx, userArg)
		if err != nil {
			return err
		}
		var addLikeParams DeleteLikeParams
		addLikeParams = DeleteLikeParams{
			UserID: resultUser.ID,
			PostID: postId,
		}
		err = store.DeleteLike(ctx, addLikeParams)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
