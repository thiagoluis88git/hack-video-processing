package usecase

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/thiagoluis88git/hack-video-processing/internal/domain/entity"
	"github.com/thiagoluis88git/hack-video-processing/internal/domain/repository"
	"github.com/thiagoluis88git/hack-video-processing/pkg/queue"
	"github.com/thiagoluis88git/hack-video-processing/pkg/responses"
	"github.com/thiagoluis88git/hack-video-processing/pkg/utils"
	videoprocess "github.com/thiagoluis88git/hack-video-processing/pkg/video-process"
)

type ProcessVideoUseCase interface {
	Execute(ctx context.Context, chanMessage *types.Message) error
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

func (uc *ProcessVideoUseCaseImpl) Execute(ctx context.Context, chanMessage *types.Message) error {
	trackingID := *chanMessage.Body

	// Gerar um endpoint para presigned url para PUT
	// Criar um outro servico para enviar a URL com o arquivo lá dentro
	// A cada erro, enviar para a fila de erros com a mensagem do erro

	file, err := uc.repo.GetFile(ctx, trackingID)

	if err != nil {
		errorProccess := responses.Wrap("use case: error when getting file", err)
		uc.writeErrorMessage(errorProccess.Error())
		return errorProccess
	}

	err = uc.videProcess.ExtractFrames(file.Name, trackingID)

	if err != nil {
		errorProccess := responses.Wrap("use case: error when extracting frames", err)
		uc.writeErrorMessage(errorProccess.Error())
		return errorProccess
	}

	zipFileName := fmt.Sprintf("%v.zip", trackingID)
	zippedFile, err := uc.videProcess.ZipFiles(trackingID, zipFileName)

	if err != nil {
		errorProccess := responses.Wrap("use case: error when zipping file", err)
		uc.writeErrorMessage(errorProccess.Error())
		return errorProccess
	}

	zipURL, err := uc.repo.UploadFile(ctx, zipFileName, zippedFile, "arquivo ZIP")

	if err != nil {
		errorProccess := responses.Wrap("use case: error when uploading zip file", err)
		uc.writeErrorMessage(errorProccess.Error())
		return errorProccess
	}

	newMessage := entity.Message{
		ZippedURL:     zipURL,
		TrackingID:    trackingID,
		ReceiptHandle: *chanMessage.ReceiptHandle,
	}

	uc.queueManager.WriteMessage(newMessage)

	// remove files
	err = utils.RemoveContentsOfFile(trackingID)

	if err != nil {
		return responses.Wrap("use case: error when deleting folder with files", err)
	}

	err = utils.RemoveFile(trackingID)

	if err != nil {
		return responses.Wrap(fmt.Sprintf("use case: error when deleting %v file", trackingID), err)
	}

	err = utils.RemoveFile(file.Name)

	if err != nil {
		return responses.Wrap(fmt.Sprintf("use case: error when deleting %v file", file.Name), err)
	}

	err = utils.RemoveFile(zipFileName)

	if err != nil {
		return responses.Wrap(fmt.Sprintf("use case: error when deleting %v file", zipFileName), err)
	}

	return nil
}

func (uc *ProcessVideoUseCaseImpl) writeErrorMessage(message string) {
	newMessage := entity.ErrorMessage{
		Message: message,
	}

	uc.queueManager.WriteErrorMessage(newMessage)
}
