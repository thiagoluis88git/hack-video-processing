package utils

import (
	"os"
	"path/filepath"
)

func RemoveContentsOfFile(path string) (err error) {
	contents, err := filepath.Glob(path + "/*")

	if err != nil {
		return
	}

	for _, item := range contents {
		err = os.RemoveAll(item)
		if err != nil {
			return
		}
	}

	return
}

func RemoveFile(path string) error {
	err := os.RemoveAll(path)

	if err != nil {
		return err
	}

	return nil
}
