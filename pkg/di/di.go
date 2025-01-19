package di

import (
	"fmt"

	"github.com/thiagoluis88git/hack-video-processing/internal/data/remote"
	dataRepo "github.com/thiagoluis88git/hack-video-processing/internal/data/repository"
	"github.com/thiagoluis88git/hack-video-processing/internal/domain/repository"
	"github.com/thiagoluis88git/hack-video-processing/internal/domain/usecase"
	"github.com/thiagoluis88git/hack-video-processing/pkg/environment"
	"github.com/thiagoluis88git/hack-video-processing/pkg/queue"
	"github.com/thiagoluis88git/hack-video-processing/pkg/storage"
	videoprocess "github.com/thiagoluis88git/hack-video-processing/pkg/video-process"
)

func ProvidesStorageRemoteDataSource(env environment.Environment) remote.StorageRemoteDataSource {
	s3, err := storage.NewAWSS3Session(env.Region)

	if err != nil {
		panic(fmt.Sprintf("error when getting S3 session: %v", err.Error()))
	}

	return remote.NewStorageRemoteDataSource(s3, env.Bucket, env.BucketZip)
}

func ProvidesProcessingRepository(ds remote.StorageRemoteDataSource) repository.ProcessingRepository {
	return dataRepo.NewProcessingRepository(ds)
}

func ProvidesProcessVideoUseCase(
	videProcess videoprocess.VideoProcessService,
	repo repository.ProcessingRepository,
	queueManager queue.QueueManager,
) usecase.ProcessVideoUseCase {
	return usecase.NewProcessVideoUseCase(videProcess, repo, queueManager)
}
