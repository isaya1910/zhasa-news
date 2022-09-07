package db

import (
	"context"
	"github.com/isaya1910/zhasa-news/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreatePostTx(t *testing.T) {
	store := NewStore(testDB)

	createUserArg := CreateOrUpdateUserParams{
		ID:        util.RandomInt(1, 1000),
		FirstName: util.RandomName(),
		LastName:  util.RandomName(),
		Bio:       util.RandomBio(),
	}

	createPostArg := CreatePostParams{
		Title: util.RandomTitle(),
		Body:  util.RandomPostBody(),
	}

	post, user, err := store.CreatePostTx(context.Background(), createPostArg, "", createUserArg)

	require.NoError(t, err)
	require.NotEmpty(t, post)
	require.NotEmpty(t, user)
	require.Equal(t, post.UserID, user.ID)

	err = store.DeletePost(context.Background(), post.ID)
	require.NoError(t, err)
	err = store.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)
}
