package api

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	db "github.com/isaya1910/zhasa-news/db/sqlc"
	"net/http"
)

type createPostRequest struct {
	Title  string `json:"title" binding:"required"`
	Body   string `json:"body" binding:"required"`
	UserID int32  `json:"user_id"`
}

type createUserJson struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Bio       *string `json:"bio"`
	ID        *int32  `json:"id"`
}

func (u createUserJson) validateUserJson() error {
	if u.FirstName == nil {
		return errors.New("first_name required")
	}
	if u.LastName == nil {
		return errors.New("last_name required")
	}
	if u.ID == nil {
		return errors.New("id required")
	}
	if u.Bio == nil {
		return errors.New("bio required")
	}
	return nil
}

func (server *Server) createPost(ctx *gin.Context) {
	var req *createPostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userJsonString := ctx.GetHeader("user")
	var userJson createUserJson

	err := json.Unmarshal([]byte(userJsonString), &userJson)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("user not found")))
		return
	}

	err = userJson.validateUserJson()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	argPost := db.CreatePostParams{
		Title:  req.Title,
		Body:   req.Body,
		UserID: req.UserID,
	}

	argUser := db.CreateOrUpdateUserParams{
		FirstName: *userJson.FirstName,
		LastName:  *userJson.LastName,
		Bio:       *userJson.Bio,
		ID:        *userJson.ID,
	}

	post, _, err := server.store.CreatePostTx(ctx, argPost, argUser)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, post)
}
