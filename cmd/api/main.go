package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/thiagoluis88git/hack-video-processing/pkg/di"
	"github.com/thiagoluis88git/hack-video-processing/pkg/environment"
	"github.com/thiagoluis88git/hack-video-processing/pkg/queue"
	videoprocess "github.com/thiagoluis88git/hack-video-processing/pkg/video-process"
)

func main() {
	chnMessages := make(chan *types.Message)

	env := environment.LoadEnvironmentVariables()

	queueManager := queue.ConfigQueueManager(env)
	videoProcess := videoprocess.NewVideoProcess()

	remoteDS := di.ProvidesStorageRemoteDataSource(env)
	repo := di.ProvidesProcessingRepository(remoteDS)
	processVideoUseCase := di.ProvidesProcessVideoUseCase(videoProcess, repo, queueManager)

	go queueManager.PollMessages(chnMessages)

	for chanMessage := range chnMessages {
		if chanMessage == nil {
			return
		}

		err := processVideoUseCase.Execute(context.Background(), chanMessage)

		if err != nil {
			log.Println(fmt.Printf("main: error when processing video: %v", err.Error()))
		}
	}
}
