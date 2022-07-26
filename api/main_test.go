package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/isaya1910/zhasa-news/db/sqlc"
	"github.com/isaya1910/zhasa-news/util"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func CreateRandomUser() db.User {
	return db.User{
		FirstName: util.RandomName(),
		LastName:  util.RandomName(),
		ID:        util.RandomInt(1, 1000),
		Bio:       util.RandomBio(),
	}
}

func CreateRandomPost(userId int32) db.Post {
	return db.Post{
		Title:     util.RandomTitle(),
		Body:      util.RandomPostBody(),
		ID:        util.RandomInt(1, 1000),
		CreatedAt: time.Now(),
		UserID:    userId,
	}
}
