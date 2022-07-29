package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/isaya1910/zhasa-news/db/sqlc"
	"net/http"
	"strconv"
)

type CreateCommentRequest struct {
	CommentBody string `json:"comment_body" binding:"required"`
	PostId      int32  `json:"post_id" binding:"required"`
	UserId      int32  `json:"user_id" binding:"required"`
}

func (server *Server) getCommentsAndAuthorsByPostId(ctx *gin.Context) {
	postId, err := strconv.Atoi(ctx.Query("post_id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(buildArgumentRequiredError("post_id")))
		return
	}

	if postId <= 0 {
		ctx.JSON(http.StatusNotFound, errorResponse(buildArgumentRequiredError("post_id")))
		return
	}

	comments, err := server.store.GetCommentsAndAuthorsByPostId(ctx, int32(postId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, comments)
}

func (server *Server) deleteComment(ctx *gin.Context) {
	commentId, err := strconv.Atoi(ctx.Query("comment_id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(errors.New("comment_id required")))
		return
	}

	_, err = server.store.GetCommentById(ctx, int32(commentId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteComment(ctx, int32(commentId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.Status(http.StatusOK)
}

func (server *Server) createComment(ctx *gin.Context) {
	var createCommentRequest *CreateCommentRequest
	if err := ctx.ShouldBindJSON(&createCommentRequest); err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	token := ctx.GetHeader("Authorization")

	user, err := server.repository.GetUser(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("authorization required")))
		return
	}

	createCommentParams := db.CreateCommentParams{
		Body:   createCommentRequest.CommentBody,
		UserID: createCommentRequest.UserId,
		PostID: createCommentRequest.PostId,
	}

	argUser := db.CreateOrUpdateUserParams{
		FirstName: *user.FirstName,
		LastName:  *user.LastName,
		Bio:       *user.Bio,
		ID:        *user.ID,
	}
	comment, _, err := server.store.CreateCommentTx(ctx, createCommentParams, argUser)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, comment)
}
