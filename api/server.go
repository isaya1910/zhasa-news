package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/isaya1910/zhasa-news/db/sqlc"
	"google.golang.org/api/option"
	"log"
	"net/http"
)

type Server struct {
	store      db.Store
	router     *gin.Engine
	repository UserRepository
	pushSender PushMessageSender
}

// NewServer creates new http server and setup routing
func NewServer(opt option.ClientOption, store db.Store, repository UserRepository) *Server {
	pushSender := FirebasePushMessageSender{opt: opt}
	server := &Server{store: store, repository: repository, pushSender: pushSender}
	router := gin.Default()
	router.Use(getAndSetUser(repository, store))

	router.POST("/news/posts", server.createPost)
	router.DELETE("/news/posts", server.deletePost)
	router.GET("/news/posts", server.getPosts)

	router.POST("/news/comments", server.createComment)
	router.DELETE("/news/comments", server.deleteComment)
	router.GET("/news/comments", server.getCommentsAndAuthorsByPostId)

	router.POST("/news/posts/likes", server.toggleLike)

	server.router = router
	return server
}

func getAndSetUser(repository UserRepository, store db.Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		userJson, err := repository.GetUser(token)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		argUser := db.CreateOrUpdateUserParams{
			FirstName: *userJson.FirstName,
			LastName:  *userJson.LastName,
			Bio:       *userJson.Bio,
			ID:        *userJson.ID,
			AvatarUrl: *userJson.AvatarUrl,
		}

		user, err := store.CreateUserTx(ctx, argUser)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		log.Print(user.ID)
		ctx.Set("user_id", int(user.ID))
		ctx.Next()
	}
}

// Start runs the HTTP server a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
