package db

import (
	"context"
	"database/sql"
	"github.com/isaya1910/zhasa-news/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func CreatePost(t *testing.T, userId int32) Post {
	arg := CreatePostParams{
		Title:  util.RandomTitle(),
		Body:   util.RandomPostBody(),
		UserID: userId,
	}
	testPost, err := testQueries.CreatePost(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, testPost)
	require.Equal(t, arg.Title, testPost.Title)
	require.Equal(t, arg.Body, testPost.Body)
	require.Equal(t, arg.UserID, testPost.UserID)
	return testPost
}

func TestDeletePost(t *testing.T) {
	testUser := CreateOrUpdateUser(t)
	testPost := CreatePost(t, testUser.ID)
	err := testQueries.DeletePost(context.Background(), testPost.ID)
	require.NoError(t, err)

	post, err := testQueries.GetPostById(context.Background(), testPost.ID)
	require.Empty(t, post)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	err = testQueries.DeleteUser(context.Background(), testPost.UserID)
	require.NoError(t, err)
}

func TestGetPostsAndAuthors(t *testing.T) {
	testUser := CreateOrUpdateUser(t)
	testPost := CreatePost(t, testUser.ID)

	postsAndAuthors, err := testQueries.GetPostsAndPostAuthors(context.Background(), GetPostsAndPostAuthorsParams{testUser.ID, 10, 0})
	require.NoError(t, err)
	require.NotEmpty(t, postsAndAuthors)

	err = testQueries.DeletePost(context.Background(), testPost.ID)
	require.NoError(t, err)

	err = testQueries.DeleteUser(context.Background(), testUser.ID)
	require.NoError(t, err)
}
