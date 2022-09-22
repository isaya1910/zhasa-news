package api

import (
	"bytes"
	"encoding/json"
	"errors"
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

type UserStubRepository struct{}

func (UserStubRepository) GetUser(token string) (CreateUserJson, error) {
	if len(token) == 0 {
		return CreateUserJson{}, errors.New("token is not valid")
	}
	firstName := util.RandomName()
	lastName := util.RandomName()
	bio := util.RandomBio()
	id := util.RandomInt(1, 1000)
	avatarUrl := ""
	return CreateUserJson{
		FirstName: &firstName,
		LastName:  &lastName,
		Bio:       &bio,
		ID:        &id,
		AvatarUrl: &avatarUrl,
	}, nil
}

func TestDeletePostApi(t *testing.T) {

	testCases := []struct {
		name          string
		postId        string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "Happy path",
			postId: "3",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetPostById(gomock.Any(), gomock.Any()).Times(1).Return(db.Post{}, nil)
				store.EXPECT().
					DeletePost(gomock.Any(), int32(3)).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Post id is empty",
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
			postId: "",
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

			url := fmt.Sprint("/news/posts?post_id=", tc.postId)

			request, err := http.NewRequest(http.MethodDelete, url, nil)

			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestCreatePostApi(t *testing.T) {

	testUser := CreateRandomUser()

	createPostRequest := createPostRequest{
		Title: util.RandomTitle(),
		Body:  util.RandomPostBody(),
	}

	testCases := []struct {
		name             string
		user             db.User
		postArgs         db.CreatePostParams
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
				return json.Marshal(&createPostRequest)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUserTx(gomock.Any(), gomock.Any()).Times(1).Return(db.User{
					ID: 1,
				}, nil)
				store.EXPECT().
					CreatePostTx(gomock.Any(), gomock.Any()).
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
				return json.Marshal(&createPostRequest)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreatePostTx(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/news/posts")

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
			request.Header.Set("Authorization", tc.tokenHeader)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
