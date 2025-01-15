package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type FFProbeOutput struct {
	Stream []struct {
		NBFrames string `json:"nb_frames"`
	} `json:"streams"`
}

func mainFFMPEG() {
	videoFile := "video.mp4" // Path to your video file
	outputDir := "output"    // Directory to store extracted images
	zipFile := "images.zip"  // Name of the zip file

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		fmt.Println("Error creating output directory:", err)
		return
	}

	// Extract frames from the video
	err := extractFrames(videoFile, outputDir)
	if err != nil {
		fmt.Println("Error extracting frames:", err)
		return
	}

	maxFrames, err := getFrameCount(videoFile)

	if err != nil {
		fmt.Println("Error getting max frame count:", err)
		return
	}

	fmt.Printf("Max quantity: %v\n", maxFrames)

	// Create a ZIP file with the extracted images
	err = zipFiles(outputDir, zipFile)
	if err != nil {
		fmt.Println("Error creating zip file:", err)
		return
	}

	fmt.Println("Images extracted and zipped successfully.")
}

// Function to extract frames using ffmpeg-go
func extractFrames(videoFile string, outputDir string) error {
	// Use ffmpeg-go to extract frames
	return ffmpeg.Input(videoFile).
		Output(filepath.Join(outputDir, "image%d.jpg"), ffmpeg.KwArgs{"vf": "fps=1"}).
		Run()
}

// Function to zip files in a directory
func zipFiles(sourceDir string, zipFile string) error {
	zipFileHandle, err := os.Create(zipFile)
	if err != nil {
		return err
	}
	defer zipFileHandle.Close()

	zipWriter := zip.NewWriter(zipFileHandle)
	defer zipWriter.Close()

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

	return err
}

func getFrameCount(videoFile string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=nb_frames", "-of", "json", videoFile)

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	var ffProbeOutput FFProbeOutput
	if err := json.Unmarshal(output, &ffProbeOutput); err != nil {
		return "", err
	}

	if len(ffProbeOutput.Stream) > 0 {
		return ffProbeOutput.Stream[0].NBFrames, nil
	}

	return "0", nil // If no frames found
}
