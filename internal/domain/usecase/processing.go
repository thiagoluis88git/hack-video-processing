package usecase

import (
	"context"
	"fmt"

	"github.com/thiagoluis88git/hack-video-processing/internal/domain/entity"
	"github.com/thiagoluis88git/hack-video-processing/internal/domain/repository"
	"github.com/thiagoluis88git/hack-video-processing/pkg/queue"
	"github.com/thiagoluis88git/hack-video-processing/pkg/responses"
	videoprocess "github.com/thiagoluis88git/hack-video-processing/pkg/video-process"
)

type ProcessVideoUseCase interface {
	Execute(ctx context.Context, message entity.Message) error
}

type ProcessVideoUseCaseImpl struct {
	videProcess  videoprocess.VideoProcessService
	repo         repository.ProcessingRepository
	queueManager queue.QueueManager
}

func NewProcessVideoUseCase(
	videProcess videoprocess.VideoProcessService,
	repo repository.ProcessingRepository,
	queueManager queue.QueueManager,
) ProcessVideoUseCase {
	return &ProcessVideoUseCaseImpl{
		videProcess:  videProcess,
		repo:         repo,
		queueManager: queueManager,
	}
}

func (uc *ProcessVideoUseCaseImpl) Execute(ctx context.Context, message entity.Message) error {
	file, err := uc.repo.GetFile(ctx, *message.Body)

	if err != nil {
		return responses.Wrap("use case: error when getting file", err)
	}

	err = uc.videProcess.ExtractFrames(file.Name, *message.Body)

	if err != nil {
		return responses.Wrap("use case: error when extracting frames", err)
	}

	zippedFile, err := uc.videProcess.ZipFiles(fmt.Sprintf("output-%v", *message.Body), "files.zip")

	if err != nil {
		return responses.Wrap("use case: error when zipping file", err)
	}

	zipURL, err := uc.repo.UploadFile(ctx, fmt.Sprintf("%v.zip", *message.Body), zippedFile, "arquivo ZIP")

	if err != nil {
		return responses.Wrap("use case: error when uploading zip file", err)
	}

	newMessage := entity.Message{
		Body:          &zipURL,
		ReceiptHandle: message.ReceiptHandle,
	}

	uc.queueManager.WriteMessage(newMessage)

	return nil
}
