package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/isaya1910/zhasa-news/db/sqlc"
)

type Server struct {
	store      db.Store
	router     *gin.Engine
	repository UserRepository
}

// NewServer creates new http server and setup routing
func NewServer(store db.Store, repository UserRepository) *Server {
	server := &Server{store: store, repository: repository}
	router := gin.Default()

	router.POST("/posts", server.createPost)
	router.DELETE("/posts", server.deletePost)
	router.GET("/posts", server.getPosts)

	router.POST("/comments", server.createComment)
	router.DELETE("/comments", server.deleteComment)
	router.GET("/comments", server.getCommentsAndAuthorsByPostId)

	router.POST("/posts/likes", server.toggleLike)

	server.router = router
	return server
}

// Start runs the HTTP serveron a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
