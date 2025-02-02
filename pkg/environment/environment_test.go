package environment_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thiagoluis88git/hack-video-processing/pkg/environment"
)

func setup() {
	os.Setenv(environment.Region, "Region")
	os.Setenv(environment.VideoProcessedOutpuQueue, "Output")
	os.Setenv(environment.VideoProcessingInputQueue, "Input")
	os.Setenv(environment.Bucket, "S3Bucket")
	os.Setenv(environment.BucketZip, "S3Bucket-ZIP")
}

func TestEnvironment(t *testing.T) {
	t.Parallel()
	setup()

	t.Run("got success when loading variables", func(t *testing.T) {
		environment.LoadEnvironmentVariables()
	})

	t.Run("got success when initializing environment", func(t *testing.T) {
		env := environment.LoadEnvironmentVariables()

		assert.Equal(t, "Region", env.Region)
		assert.Equal(t, "Output", env.VideoProcessedOutputQueue)
		assert.Equal(t, "Input", env.VideoProcessingInputQueue)
	})
}
