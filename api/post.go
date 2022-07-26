package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	db "github.com/isaya1910/zhasa-news/db/sqlc"
	"net/http"
)

type createPostRequest struct {
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
}

type CreateUserJson struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Bio       *string `json:"bio"`
	ID        *int32  `json:"id"`
}

func (u CreateUserJson) validateUserJson() error {
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
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("бляяя")))
		return
	}

	token := ctx.GetHeader("Authorization")

	user, err := server.repository.GetUser(token)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = user.validateUserJson()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	argPost := db.CreatePostParams{
		Title: req.Title,
		Body:  req.Body,
	}

	argUser := db.CreateOrUpdateUserParams{
		FirstName: *user.FirstName,
		LastName:  *user.LastName,
		Bio:       *user.Bio,
		ID:        *user.ID,
	}

	post, _, err := server.store.CreatePostTx(ctx, argPost, argUser)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, post)
}
