package entity

import "os"

type S3File struct {
	File *os.File
	Name string
}

type Message struct {
	Body          *string
	ReceiptHandle *string
}
