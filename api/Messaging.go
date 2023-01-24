package api

import (
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	db "github.com/isaya1910/zhasa-news/db/sqlc"
	"google.golang.org/api/option"
	"log"
)

func SendPostPush(opt option.ClientOption, post db.Post) error {
	app, err := firebase.NewApp(context.Background(), nil, opt)
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
		Notification: &messaging.Notification{
			Title: post.Title,
			Body:  post.Body,
		},
		Topic: "news",
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					ContentAvailable: true,
					Alert: &messaging.ApsAlert{
						Title: post.Title,
						Body:  post.Body,
					},
					CustomData: map[string]interface{}{
						"deeplink": "news",
					},
					Sound: "default",
				},
			},
		},
	}
	s, _ := json.MarshalIndent(message, "", "\t")
	fmt.Printf("%+v\n", string(s))

	response, err := fcmClient.Send(context.Background(), message)
	if err != nil {
		log.Println(err)

		fmt.Print(err)
		return err
	}
	fmt.Print(response)
	return err
}
