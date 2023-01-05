package api

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	db "github.com/isaya1910/zhasa-news/db/sqlc"
	"google.golang.org/api/option"
	"log"
)

func SendPostPush(post db.Post) error {

	opt := option.WithCredentialsFile("/home/zhasprod/zhasa-news/serviceAccount.json")

	config := &firebase.Config{ProjectID: "zhasa-7a01b"}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("error initializing app: %v", err)
	}
	fcmClient, err := app.Messaging(context.Background())

	if err != nil {
		log.Fatal(err)

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
		log.Fatal(err)

		fmt.Print(err)
		return err
	}
	fmt.Print(response)
	return err
}
