package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateCommentRequest struct {
	CommentBody string `json:"comment_body" binding:"required"`
}

func (server Server) CreateComment(ctx *gin.Context) {
	var createCommentRequest CreateCommentRequest
	if err := ctx.ShouldBindJSON(createCommentRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
}
