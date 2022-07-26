package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	CreatePostTx(ctx context.Context, postArg CreatePostParams, userArg CreateOrUpdateUserParams) (Post, User, error)
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

func (store *SQLStore) CreatePostTx(ctx context.Context, postArg CreatePostParams, userArg CreateOrUpdateUserParams) (Post, User, error) {
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
