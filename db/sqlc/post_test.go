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

	postsAndAuthors, err := testQueries.GetPostsAndPostAuthors(context.Background(), GetPostsAndPostAuthorsParams{10, 0})
	require.NoError(t, err)
	require.NotEmpty(t, postsAndAuthors)
	postAndAuthor := postsAndAuthors[0]
	require.Equal(t, postAndAuthor.PostID, testPost.ID)
	require.Equal(t, postAndAuthor.UserID, testUser.ID)
	require.Equal(t, postAndAuthor.Body, testPost.Body)

	require.Equal(t, postAndAuthor.Title, testPost.Title)
	require.Equal(t, postAndAuthor.CreatedAt, testPost.CreatedAt)
	require.Equal(t, postAndAuthor.FirstName, testUser.FirstName)

	err = testQueries.DeletePost(context.Background(), testPost.ID)
	require.NoError(t, err)

	err = testQueries.DeleteUser(context.Background(), testUser.ID)
	require.NoError(t, err)
}
