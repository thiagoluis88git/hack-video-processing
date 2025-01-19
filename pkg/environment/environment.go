package environment

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	RedocFolderPath *string = flag.String("PATH_REDOC_FOLDER", "/docs/swagger.json", "Swagger docs folder")

	localDev = flag.String("localDev", "false", "local development")
)

const (
	VideoProcessingInputQueue = "VIDEO_PROCESSING_INPUT_QUEUE"
	VideoProcessedOutpuQueue  = "VIDEO_PROCESSED_OUTPUT_QUEUE"
	Bucket                    = "AWS_S3_BUCKET"
	BucketZip                 = "AWS_S3_BUCKET_ZIP"
	Region                    = "AWS_REGION"
)

type Environment struct {
	VideoProcessingInputQueue string
	VideoProcessedOutputQueue string
	Bucket                    string
	BucketZip                 string
	Region                    string
}

func LoadEnvironmentVariables() Environment {
	flag.Parse()

	if localFlag := *localDev; localFlag != "false" {
		err := godotenv.Load()

		if err != nil {
			log.Fatal("Error loading .env file", err.Error())
		}
	}

	region := getEnvironmentVariable(Region)
	videoProcessingInputQueue := getEnvironmentVariable(VideoProcessingInputQueue)
	videoProcessedOutputQueue := getEnvironmentVariable(VideoProcessedOutpuQueue)
	bucket := getEnvironmentVariable(Bucket)
	bucketZip := getEnvironmentVariable(BucketZip)

	return Environment{
		VideoProcessingInputQueue: videoProcessingInputQueue,
		VideoProcessedOutputQueue: videoProcessedOutputQueue,
		Region:                    region,
		Bucket:                    bucket,
		BucketZip:                 bucketZip,
	}
}

func getEnvironmentVariable(key string) string {
	value, hashKey := os.LookupEnv(key)

	if !hashKey {
		log.Fatalf("There is no %v environment variable", key)
	}

	return value
}
