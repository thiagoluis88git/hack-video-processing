package model

import "os"

type S3File struct {
	File *os.File
	Name string
}
