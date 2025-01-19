package repository

import (
	"context"

	"github.com/thiagoluis88git/hack-video-processing/internal/domain/entity"
)

type ProcessingRepository interface {
	GetFile(ctx context.Context, key string) (entity.S3File, error)
}
