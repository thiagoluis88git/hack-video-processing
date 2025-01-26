package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/thiagoluis88git/hack-video-processing/pkg/responses"
)

func NewAWSS3Session(region string) (*s3.S3, error) {
	session, err := session.NewSession(
		&aws.Config{
			Region:                        aws.String(region),
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
	)

	if err != nil {
		return nil, responses.Wrap("AWS Session error", err)
	}

	return s3.New(session), nil
}
