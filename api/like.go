package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	db "github.com/isaya1910/zhasa-news/db/sqlc"
	"net/http"
	"strconv"
)

func (server *Server) toggleLike(ctx *gin.Context) {
	post, err := strconv.Atoi(ctx.Query("post_id"))
	postId := int32(post)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(errors.New("post_id required")))
		return
	}
	if postId <= 0 {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("post_id should be greater than 0")))
		return
	}

	userId := ctx.GetInt("user_id")
	params := db.GetUserPostLikeParams{
		UserID: int32(userId),
		PostID: postId,
	}
	like, err := server.store.GetUserPostLike(ctx, params)

	if err != nil {
		addLikeParams := db.AddLikeParams{
			UserID: int32(userId),
			PostID: postId,
		}
		_, err = server.store.AddLikeTx(ctx, addLikeParams)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.Status(http.StatusOK)
		return
	}

	deleteLikeParams := db.DeleteLikeParams{
		UserID: int32(userId),
		PostID: postId,
	}

	err = server.store.DeleteLikeTx(ctx, deleteLikeParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, like)
}
