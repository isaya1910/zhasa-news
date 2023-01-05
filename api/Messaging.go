package api

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	db "github.com/isaya1910/zhasa-news/db/sqlc"
	"log"
)

func SendPostPush(post db.Post) error {
	//config := &firebase.Config{ProjectID: "zhasa-7a01b"}
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("error initializing app: %v", err)
	}
	fcmClient, err := app.Messaging(context.Background())

	if err != nil {
		log.Println(err)

		return fmt.Errorf("error initializing app: %v", err)
	}
	message := &messaging.Message{
		Data: map[string]string{
			"title": post.Title,
			"body":  post.Body,
		},
		Topic: "news",
	}

	response, err := fcmClient.Send(context.Background(), message)
	if err != nil {
		log.Println(err)

		fmt.Print(err)
		return err
	}
	fmt.Print(response)
	return err
}
