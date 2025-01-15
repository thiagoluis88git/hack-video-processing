package main

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/thiagoluis88git/hack-video-processing/pkg/environment"
	"github.com/thiagoluis88git/hack-video-processing/pkg/queue"
)

func main() {
	chnMessages := make(chan *types.Message)

	env := environment.LoadEnvironmentVariables()

	queueManager := queue.ConfigQueueManager(env)

	go queueManager.PollMessages(chnMessages)

	for message := range chnMessages {
		if message == nil {
			return
		}

		queueManager.WriteMessage(message)
	}
}
