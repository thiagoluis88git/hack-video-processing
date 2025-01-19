package queue

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/thiagoluis88git/hack-video-processing/internal/domain/entity"
	"github.com/thiagoluis88git/hack-video-processing/pkg/environment"
)

const (
	waitTimeout         = 20
	maxNumberOfMessages = 1
)

type QueueManager struct {
	outputQueueURL    string
	outputQueueClient *sqs.Client
	inputQueueURL     string
	inputQueueClient  *sqs.Client
}

func ConfigQueueManager(environment environment.Environment) QueueManager {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("failed to load SDK configuration, %v", err)
	}

	// Create SQS clients
	inputQueueClient := sqs.NewFromConfig(cfg)
	outputQueueClient := sqs.NewFromConfig(cfg)

	return QueueManager{
		inputQueueURL:     environment.VideoProcessingInputQueue,
		outputQueueURL:    environment.VideoProcessedOutputQueue,
		outputQueueClient: outputQueueClient,
		inputQueueClient:  inputQueueClient,
	}
}

func (manager *QueueManager) PollMessages(
	chn chan<- *types.Message,
) {
	for {
		// Receive messages from source queue
		result, err := manager.inputQueueClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(manager.inputQueueURL),
			MaxNumberOfMessages: *aws.Int32(maxNumberOfMessages), // Receive only one message at a time
			WaitTimeSeconds:     *aws.Int32(waitTimeout),         // Wait for messages for 20 seconds
		})

		if err != nil {
			log.Fatalf("failed to receive messages, %v", err)
		}

		// Process received messages
		for _, message := range result.Messages {
			// Process the message content here (e.g., transform, enrich, etc.)
			chn <- &message
		}
	}
}

func (manager *QueueManager) WriteMessage(
	message entity.Message,
) {
	// Send processed message to destination queue
	_, err := manager.outputQueueClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(manager.outputQueueURL),
		MessageBody: message.Body,
	})

	if err != nil {
		log.Fatalf("failed to send message, %v", err)
	}

	// Delete message from source queue
	_, err = manager.inputQueueClient.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(manager.inputQueueURL),
		ReceiptHandle: message.ReceiptHandle,
	})

	if err != nil {
		log.Fatalf("failed to delete message, %v", err)
	}

	fmt.Println("Message processed and sent to destination queue")
}
