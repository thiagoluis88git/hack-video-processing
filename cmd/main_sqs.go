package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func main() {
	chnMessages := make(chan *types.Message)

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load SDK configuration, %v", err)
	}

	// Create SQS clients
	sourceQueueClient := sqs.NewFromConfig(cfg)
	destinationQueueClient := sqs.NewFromConfig(cfg)

	// Get source queue URL from environment variable
	// sourceQueueURL := os.Getenv("SOURCE_QUEUE_URL")
	sourceQueueURL := "https://sqs.us-east-1.amazonaws.com/199710252834/teste-in"
	if sourceQueueURL == "" {
		log.Fatal("SOURCE_QUEUE_URL environment variable is not set")
	}

	// Get destination queue URL from environment variable
	// destinationQueueURL := os.Getenv("DESTINATION_QUEUE_URL")
	destinationQueueURL := "https://sqs.us-east-1.amazonaws.com/199710252834/teste-out"
	if destinationQueueURL == "" {
		log.Fatal("DESTINATION_QUEUE_URL environment variable is not set")
	}

	go pollMessages(sourceQueueURL, sourceQueueClient, chnMessages)

	for message := range chnMessages {
		if message == nil {
			return
		}

		// Send processed message to destination queue
		_, err = destinationQueueClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
			QueueUrl:    aws.String(destinationQueueURL),
			MessageBody: message.Body,
		})

		if err != nil {
			log.Fatalf("failed to send message, %v", err)
		}

		// Delete message from source queue
		_, err = sourceQueueClient.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
			QueueUrl:      aws.String(sourceQueueURL),
			ReceiptHandle: message.ReceiptHandle,
		})

		if err != nil {
			log.Fatalf("failed to delete message, %v", err)
		}

		fmt.Println("Message processed and sent to destination queue")
	}
}

func pollMessages(
	sourceQueueURL string,
	sourceQueueClient *sqs.Client,
	chn chan<- *types.Message,
) {
	for {
		// Receive messages from source queue
		result, err := sourceQueueClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(sourceQueueURL),
			MaxNumberOfMessages: *aws.Int32(1),  // Receive only one message at a time
			WaitTimeSeconds:     *aws.Int32(20), // Wait for messages for 20 seconds
		})

		if err != nil {
			log.Fatalf("failed to receive messages, %v", err)
		}

		if result.Messages == nil {
			fmt.Println("TIMEOUT")
		}

		// Process received messages
		for _, message := range result.Messages {
			// Process the message content here (e.g., transform, enrich, etc.)
			chn <- &message

		}
	}
}
