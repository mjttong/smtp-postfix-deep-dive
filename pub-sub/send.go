package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"time"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/pubsub/v2"
)

type EmailData struct {
	To      string `json:"to"`
	From    string `json:"from"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func sendEmail(email EmailData) error {
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		email.From, email.To, email.Subject, email.Body)

	return smtp.SendMail("stdinout.com:587", nil, email.From, []string{email.To}, []byte(message))
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	projectID, err := metadata.ProjectIDWithContext(ctx)
	if err != nil {
		log.Fatalf("Failed to get project ID from metadata: %v", err)
	}
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	sub := client.Subscriber("smtp-topic-sub")

	log.Println("Listening for messages...")
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		log.Printf("Got message: %s", string(msg.Data))

		var email EmailData
		if json.Unmarshal(msg.Data, &email) != nil {
			log.Printf("json.Unmarshal failed: %v", err)
			msg.Nack()
			return
		}

		if sendEmail(email) != nil {
			log.Printf("sendEmail failed: %v", err)
			msg.Nack()
			return
		}

		log.Printf("Email sent to %s successfully.", email.To)
		msg.Ack()
	})

	if err != nil && err != context.DeadlineExceeded {
		log.Fatalf("sub.Receive finished with an unexpected error: %v", err)
	}

	log.Println("Finished processing messages or timed out. Shutting down.")
}
