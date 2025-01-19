package usecase

import (
	"context"
	"fmt"

	"github.com/thiagoluis88git/hack-video-processing/internal/domain/repository"
	"github.com/thiagoluis88git/hack-video-processing/pkg/queue"
	"github.com/thiagoluis88git/hack-video-processing/pkg/responses"
	videoprocess "github.com/thiagoluis88git/hack-video-processing/pkg/video-process"
)

type ProcessVideoUseCase interface {
	Execute(ctx context.Context, videoID string) error
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

func (uc *ProcessVideoUseCaseImpl) Execute(ctx context.Context, videoID string) error {
	file, err := uc.repo.GetFile(ctx, videoID)

	if err != nil {
		return responses.Wrap("use case: error when getting file", err)
	}

	err = uc.videProcess.ExtractFrames(file.Name, videoID)

	if err != nil {
		return responses.Wrap("use case: error when extracting frames", err)
	}

	zippedFile, err := uc.videProcess.ZipFiles(fmt.Sprintf("output-%v", videoID), "files.zip")

	if err != nil {
		return responses.Wrap("use case: error when extracting frames", err)
	}

	//Depois, falta enviar o ZIP no S3 e mandar uma mensagem no SQS
	print(zippedFile)

	return nil
}
