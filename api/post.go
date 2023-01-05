package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	db "github.com/isaya1910/zhasa-news/db/sqlc"
	"net/http"
	"strconv"
	"time"
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
	AvatarUrl *string `json:"avatar_url"`
}

type PostResponse struct {
	Id            int32        `json:"id"`
	Title         string       `json:"title"`
	Body          string       `json:"body"`
	LikesCount    int64        `json:"likes_count"`
	CommentsCount int64        `json:"comments_count"`
	IsLiked       bool         `json:"is_liked"`
	ImageUrls     []string     `json:"image_urls"`
	User          UserResponse `json:"user"`
	CreatedAt     time.Time    `json:"created_at"`
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

	userId := ctx.GetInt("user_id")

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
		Offset: int32(page * size),
		UserID: int32(userId),
	}

	posts, err = server.store.GetPostsAndPostAuthors(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	postsResponse := make([]PostResponse, 0)
	for _, value := range posts {
		postResponse := PostResponse{
			Id:            value.ID,
			Title:         value.Title,
			Body:          value.Body,
			ImageUrls:     value.ImageUrls,
			LikesCount:    value.LikesCount,
			CommentsCount: value.CommentsCount,
			IsLiked:       value.IsLiked,
			CreatedAt:     value.CreatedAt,
			User: UserResponse{
				FirstName: value.FirstName,
				LastName:  value.LastName,
				Role:      value.Bio,
				ID:        value.UserID,
				AvatarUrl: value.AvatarUrl,
			},
		}
		postsResponse = append(postsResponse, postResponse)
	}

	ctx.JSON(http.StatusOK, postsResponse)
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

	userId := ctx.GetInt("user_id")

	argPost := db.CreatePostParams{
		Title:  req.Title,
		Body:   req.Body,
		UserID: int32(userId),
	}

	if userId <= 0 {
		ctx.JSON(http.StatusUnauthorized, errorResponse(buildArgumentRequiredError("Authorization required")))
		return
	}

	post, err := server.store.CreatePostTx(ctx, argPost, req.ImageUrl)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	SendPostPush(server.opt, post)
	ctx.JSON(http.StatusOK, post)
}
