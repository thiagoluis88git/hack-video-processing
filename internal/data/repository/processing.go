package repository

import (
	"context"

	"github.com/thiagoluis88git/hack-video-processing/internal/data/remote"
	"github.com/thiagoluis88git/hack-video-processing/internal/domain/entity"
	"github.com/thiagoluis88git/hack-video-processing/internal/domain/repository"
	"github.com/thiagoluis88git/hack-video-processing/pkg/responses"
)

type ProcessingRepositoryImpl struct {
	ds remote.StorageRemoteDataSource
}

func NewProcessingRepository(ds remote.StorageRemoteDataSource) repository.ProcessingRepository {
	return &ProcessingRepositoryImpl{
		ds: ds,
	}
}

func (repo *ProcessingRepositoryImpl) GetFile(ctx context.Context, key string) (entity.S3File, error) {
	file, err := repo.ds.GetFiles(ctx, key)

	if err != nil {
		return entity.S3File{}, responses.Wrap("repository: error when getting s3 file", err)
	}

	return entity.S3File{
		File: file.File,
		Name: file.Name,
	}, nil
}

func (repo *ProcessingRepositoryImpl) UploadFile(ctx context.Context, key string, data []byte, description string) (string, error) {
	zipURL, err := repo.ds.UploadFile(ctx, key, data, description)

	if err != nil {
		return "", responses.Wrap("repository: error when uploading s3 file", err)
	}

	return zipURL, nil
}
