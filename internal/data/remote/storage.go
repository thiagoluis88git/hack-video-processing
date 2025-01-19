package remote

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/thiagoluis88git/hack-video-processing/internal/data/model"
	"github.com/thiagoluis88git/hack-video-processing/pkg/responses"
)

type StorageRemoteDataSource interface {
	GetFiles(ctx context.Context, key string) (model.S3File, error)
}

type AWSS3StorageRemoteDataSourceImpl struct {
	session s3iface.S3API
	bucket  string
}

func NewStorageRemoteDataSource(session s3iface.S3API, bucket string) StorageRemoteDataSource {
	return &AWSS3StorageRemoteDataSourceImpl{
		session: session,
		bucket:  bucket,
	}
}

func (ds *AWSS3StorageRemoteDataSourceImpl) GetFiles(ctx context.Context, key string) (model.S3File, error) {
	downloader := s3manager.NewDownloaderWithClient(ds.session)

	file, err := os.Create(fmt.Sprintf("%v.mp4", key))

	if err != nil {
		return model.S3File{}, responses.Wrap("remote: could not create file", err)
	}

	defer file.Close()

	_, err = downloader.DownloadWithContext(ctx, file, &s3.GetObjectInput{
		Bucket: aws.String(ds.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return model.S3File{}, responses.Wrap("remote: AWS S3 upload error", err)
	}

	return model.S3File{
		File: file,
		Name: file.Name(),
	}, nil
}
