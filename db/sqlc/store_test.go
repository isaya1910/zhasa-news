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
		Title:  util.RandomTitle(),
		Body:   util.RandomPostBody(),
		UserID: createUserArg.ID,
	}

	post, err := store.CreatePostTx(context.Background(), createPostArg, "")

	require.NoError(t, err)
	require.NotEmpty(t, post)

	err = store.DeletePost(context.Background(), post.ID)
	require.NoError(t, err)
	require.NoError(t, err)
}
