// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: user.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  id, first_name, second_name, bio 
) VALUES (
  $1, $2, $3, $4 
)
RETURNING id, first_name, second_name, bio
`

type CreateUserParams struct {
	ID         int32  `json:"id"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Bio        string `json:"bio"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.FirstName,
		arg.SecondName,
		arg.Bio,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.SecondName,
		&i.Bio,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUser = `-- name: GetUser :one

SELECT id, first_name, second_name, bio FROM users
WHERE id = $1 LIMIT 1
`

// Example queries for sqlc
func (q *Queries) GetUser(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.SecondName,
		&i.Bio,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, first_name, second_name, bio FROM users
ORDER BY name
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.SecondName,
			&i.Bio,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
