package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateUser(t *testing.T) User {
	arg := CreateUserParams{
		ID:         util.RandomInt(1, 1000),
		BBB:        util.RandomName(),
		SecondName: util.RandomName(),
		Bio:        util.RanodmBio(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.BBB, user.FirstName)
	require.Equal(t, arg.SecondName, user.SecondName)

	require.NotZero(t, user.ID)
	return user
}

func TestDeleteUser(t *testing.T) {
	testUser = createUser(t)
	err := testQueries.DeleteUser(context.Background(), user.ID)
	require.NotError(t, err)

	testUser1, err := testQueries.GetAccount(context.Background(), testUser.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, testUser1)
}
