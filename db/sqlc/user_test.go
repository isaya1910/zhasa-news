package db

import (
	"context"
	"database/sql"
	"github.com/isaya1910/zhasa-news/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateUser(t *testing.T) User {
	arg := CreateUserParams{
		ID:        util.RandomInt(1, 1000),
		FirstName: util.RandomName(),
		LastName:  util.RandomName(),
		Bio:       util.RandomBio(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.LastName, user.LastName)

	require.NotZero(t, user.ID)
	return user
}

func TestDeleteUser(t *testing.T) {
	testUser := CreateUser(t)
	err := testQueries.DeleteUser(context.Background(), testUser.ID)
	require.NoError(t, err)

	testUser1, err := testQueries.GetUser(context.Background(), testUser.ID)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, testUser1)
}

func TestGetUser(t *testing.T) {
	testUser := CreateUser(t)

	actual, err := testQueries.GetUser(context.Background(), testUser.ID)

	require.NoError(t, err)

	require.Equal(t, actual.FirstName, testUser.FirstName)
	require.Equal(t, actual.LastName, testUser.LastName)
	require.Equal(t, actual.Bio, testUser.Bio)
	require.Equal(t, actual.ID, testUser.ID)

	err = testQueries.DeleteUser(context.Background(), testUser.ID)
	require.NoError(t, err)
}
