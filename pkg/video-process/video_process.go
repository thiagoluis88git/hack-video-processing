package videoprocess

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/thiagoluis88git/hack-video-processing/pkg/responses"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type FFProbeOutput struct {
	Stream []struct {
		NBFrames string `json:"nb_frames"`
	} `json:"streams"`
}

type VideoProcessService interface {
	ExtractFrames(videoFile string, prefixName string) error
	ZipFiles(sourceDir string, zipFile string) ([]byte, error)
}

type VideoProcessServiceImpl struct{}

func NewVideoProcess() VideoProcessService {
	return &VideoProcessServiceImpl{}
}

func (service *VideoProcessServiceImpl) ExtractFrames(videoFile string, prefixName string) error {
	// maxFrames, err := getFrameCount(videoFile)

	// if err != nil {
	// 	return responses.Wrap("video processing: error when getting max frames", err)
	// }

	// fmt.Printf("Max quantity: %v\n", maxFrames)

	if err := os.MkdirAll(prefixName, os.ModePerm); err != nil {
		return responses.Wrap("video processing: error when creating folder", err)
	}

	// Use ffmpeg-go to extract frames
	return ffmpeg.Input(videoFile).
		Output(filepath.Join(prefixName, "image%d.jpg"), ffmpeg.KwArgs{"vf": "fps=1"}).
		Run()
}

// Function to zip files in a directory
func (service *VideoProcessServiceImpl) ZipFiles(sourceDir string, zipFile string) ([]byte, error) {
	zipFileHandle, err := os.Create(zipFile)

	if err != nil {
		return []byte{}, err
	}

	zipWriter := zip.NewWriter(zipFileHandle)

	err = filepath.Walk(sourceDir, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		zipEntry, err := zipWriter.Create(info.Name())
		if err != nil {
			return err
		}

		fileHandle, err := os.Open(file)
		if err != nil {
			return err
		}
		defer fileHandle.Close()

		_, err = io.Copy(zipEntry, fileHandle)
		return err
	})

	if err != nil {
		return []byte{}, responses.Wrap("video processing: error when zipping files", err)
	}

	err = zipWriter.Close()

	if err != nil {
		return []byte{}, responses.Wrap("video processing: error when closing zipWriter", err)
	}

	openZipFile, err := os.Open(zipFile)

	if err != nil {
		return []byte{}, responses.Wrap("video processing: error when openning", err)
	}

	defer openZipFile.Close()

	data, err := io.ReadAll(openZipFile)

	if err != nil {
		return []byte{}, responses.Wrap("video processing: error when reading data from file", err)
	}

	return data, nil
}

// func getFrameCount(videoFile string) (string, error) {
// 	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=nb_frames", "-of", "json", videoFile)

// 	output, err := cmd.Output()
// 	if err != nil {
// 		return "", err
// 	}

// 	var ffProbeOutput FFProbeOutput
// 	if err := json.Unmarshal(output, &ffProbeOutput); err != nil {
// 		return "", err
// 	}

// 	if len(ffProbeOutput.Stream) > 0 {
// 		return ffProbeOutput.Stream[0].NBFrames, nil
// 	}

// 	return "0", nil // If no frames found
// }
