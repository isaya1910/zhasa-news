package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	db "github.com/isaya1910/zhasa-news/db/sqlc"
	"net/http"
	"strconv"
)

type createPostRequest struct {
	Title    string `json:"title" binding:"required"`
	Body     string `json:"body" binding:"required"`
	ImageUrl string `json:"image_url"`
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

const (
	defaultSize   = 20
	defaultOffset = 0
)

func (server *Server) getPosts(ctx *gin.Context) {
	size, err := strconv.Atoi(ctx.Query("size"))
	token := ctx.GetHeader("Authorization")

	user, err := server.repository.GetUser(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	userId := *user.ID

	if userId <= 0 {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("user not found")))
		return
	}

	if err != nil {
		size = defaultSize
	}

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = defaultOffset
	}

	var posts []db.GetPostsAndPostAuthorsRow

	arg := db.GetPostsAndPostAuthorsParams{
		Limit:  int32(size),
		Offset: int32(page),
		UserID: userId,
	}
	posts, err = server.store.GetPostsAndPostAuthors(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

func (server *Server) deletePost(ctx *gin.Context) {
	postId, err := strconv.Atoi(ctx.Query("post_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("post_id required")))
		return
	}

	_, err = server.store.GetPostById(ctx, int32(postId))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("no post found with this id")))
		return
	}
	err = server.store.DeletePost(ctx, int32(postId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (server *Server) createPost(ctx *gin.Context) {
	var req *createPostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
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

	post, _, err := server.store.CreatePostTx(ctx, argPost, req.ImageUrl, argUser)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, post)
}
