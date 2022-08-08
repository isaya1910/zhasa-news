package api

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	mockdb "github.com/isaya1910/zhasa-news/db/mock"
	db "github.com/isaya1910/zhasa-news/db/sqlc"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddLike(t *testing.T) {
	testCases := []struct {
		name          string
		postId        string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "post id required",
			postId: "",
			buildStubs: func(store *mockdb.MockStore) {

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "bad post id",
			postId: "-1",
			buildStubs: func(store *mockdb.MockStore) {

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "create like",
			postId: "1",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserPostLike(gomock.Any(), gomock.Any()).
					Times(1).
					Return(int32(0), errors.New("not found"))

				store.EXPECT().AddLike(gomock.Any(), gomock.Any()).Times(1).Return(db.Like{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for i := 0; i < len(testCases); i++ {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)
			defer ctrl.Finish()

			tc.buildStubs(store)

			server := NewServer(store, UserStubRepository{})
			recorder := httptest.NewRecorder()

			url := fmt.Sprint("/posts/likes?post_id=", tc.postId)

			request, err := http.NewRequest(http.MethodPost, url, nil)
			request.Header.Set("Authorization", "testToken")
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
