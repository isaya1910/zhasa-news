package api

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	db "github.com/isaya1910/zhasa-news/db/sqlc"
	"google.golang.org/api/option"
)

func SendPostPush(post db.Post) error {

	opt := option.WithCredentialsFile("serviceAccount.json")

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}
	fcmClient, err := app.Messaging(context.Background())

	if err != nil {
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
		fmt.Print(err)
		return err
	}
	fmt.Print(response)
	return err
}
