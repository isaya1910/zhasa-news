package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	mockdb "github.com/isaya1910/zhasa-news/db/mock"
	db "github.com/isaya1910/zhasa-news/db/sqlc"
	"github.com/isaya1910/zhasa-news/util"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetComments(t *testing.T) {
	var commentsAuthorsList []db.GetCommentsAndAuthorsByPostIdRow

	testCases := []struct {
		name          string
		postId        int32
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "Get not empty comments list",
			postId: int32(1),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCommentsAndAuthorsByPostId(gomock.Any(), int32(1)).Return(
					commentsAuthorsList, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			tc.buildStubs(store)

			server := NewServer(store, UserStubRepository{})
			recorder := httptest.NewRecorder()

			url := fmt.Sprint("/comments?post_id=", tc.postId)

			request, err := http.NewRequest(http.MethodGet, url, nil)

			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteCommentApi(t *testing.T) {

	testCases := []struct {
		name          string
		commendId     string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "Happy path",
			commendId: "3",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCommentById(gomock.Any(), gomock.Any()).Times(1).Return(db.Comment{}, nil)
				store.EXPECT().
					DeleteComment(gomock.Any(), int32(3)).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "comment_id is empty",
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
			commendId: "",
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			tc.buildStubs(store)

			server := NewServer(store, UserStubRepository{})
			recorder := httptest.NewRecorder()

			url := fmt.Sprint("/comments?comment_id=", tc.commendId)

			request, err := http.NewRequest(http.MethodDelete, url, nil)

			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestCreateCommentApi(t *testing.T) {

	testUser := CreateRandomUser()
	testPost := CreateRandomPost(testUser.ID)

	createCommentRequest := CreateCommentRequest{
		CommentBody: util.RandomPostBody(),
		PostId:      testPost.ID,
		UserId:      testUser.ID,
	}

	testCases := []struct {
		name             string
		user             db.User
		commentArgs      db.CreateCommentParams
		tokenHeader      string
		buildRequestBody func() ([]byte, error)
		buildStubs       func(store *mockdb.MockStore)
		checkResponse    func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "Happy path",
			user:        testUser,
			tokenHeader: "testToken",
			buildRequestBody: func() ([]byte, error) {
				return json.Marshal(&createCommentRequest)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateCommentTx(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:        "User not found",
			user:        testUser,
			tokenHeader: "",
			buildRequestBody: func() ([]byte, error) {
				return json.Marshal(&createCommentRequest)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreatePostTx(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			tc.buildStubs(store)

			server := NewServer(store, UserStubRepository{})
			recorder := httptest.NewRecorder()

			requestBody, err := tc.buildRequestBody()

			require.NoError(t, err)

			url := fmt.Sprintf("/comments")

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
			request.Header.Set("Authorization", tc.tokenHeader)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
