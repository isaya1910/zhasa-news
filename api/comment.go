package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/isaya1910/zhasa-news/db/sqlc"
	"net/http"
	"strconv"
	"time"
)

type CreateCommentRequest struct {
	CommentBody string `json:"comment_body" binding:"required"`
	PostId      int32  `json:"post_id" binding:"required"`
}

type CommentResponse struct {
	Id          int32        `json:"id"`
	CommentBody string       `json:"comment_body"`
	PostId      int32        `json:"post_id"`
	CreatedAt   time.Time    `json:"created_at"`
	User        UserResponse `json:"user"`
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

	var commentsResponse []CommentResponse

	for _, value := range comments {
		comment := CommentResponse{
			Id:          value.CommentID,
			CommentBody: value.Body,
			PostId:      value.PostID,
			CreatedAt:   value.CreatedAt,
			User: UserResponse{
				FirstName: value.FirstName,
				LastName:  value.LastName,
				Role:      value.Bio,
				ID:        value.UserID,
				AvatarUrl: value.AvatarUrl,
			},
		}
		commentsResponse = append(commentsResponse, comment)
	}

	ctx.JSON(http.StatusOK, commentsResponse)
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

	userId := ctx.GetInt("user_id")

	createCommentParams := db.CreateCommentParams{
		Body:   createCommentRequest.CommentBody,
		UserID: int32(userId),
		PostID: createCommentRequest.PostId,
	}

	comment, err := server.store.CreateCommentTx(ctx, createCommentParams)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, comment)
}
