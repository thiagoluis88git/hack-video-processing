package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/thiagoluis88git/hack-video-processing/pkg/di"
	"github.com/thiagoluis88git/hack-video-processing/pkg/environment"
	"github.com/thiagoluis88git/hack-video-processing/pkg/httpserver"
	"github.com/thiagoluis88git/hack-video-processing/pkg/queue"
	videoprocess "github.com/thiagoluis88git/hack-video-processing/pkg/video-process"
)

func main() {
	chnMessages := make(chan *types.Message)

	env := environment.LoadEnvironmentVariables()

	// Config API. Must be async
	router := chi.NewRouter()
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.RealIP)
	router.Use(chiMiddleware.Recoverer)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(`{"message":"ok"}`)
	})

	server := httpserver.New(router)
	go server.Start()

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
