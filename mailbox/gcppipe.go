package main

import (
	"context"
	"io"
	"log"
	"os"

	"cloud.google.com/go/compute/metadata"
	pubsub "cloud.google.com/go/pubsub/v2"
)

func main() {
	ctx := context.Background()

	projectID, err := metadata.ProjectIDWithContext(ctx)
	if err != nil {
		log.Fatalf("Failed to get project ID from metadata: %v", err)
	}

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create pubsub client: %v", err)
	}
	defer client.Close()

	publisher := client.Publisher("pubsub-topic-id")
	defer publisher.Stop()

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("read stdin: %v", err)
	}
	result := publisher.Publish(ctx, &pubsub.Message{Data: data})

	msgID, err := result.Get(ctx)
	if err != nil {
		log.Fatalf("Failed to publish message to pubsub: %v", err)
	}

	log.Printf("Published message with ID: %s", msgID)
}
