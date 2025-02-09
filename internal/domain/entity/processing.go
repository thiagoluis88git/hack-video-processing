package entity

import (
	"encoding/json"
	"fmt"
	"os"
)

type S3File struct {
	File *os.File
	Name string
}

type Message struct {
	ZippedURL     string `json:"zippedURL"`
	TrackingID    string `json:"trackingID"`
	ReceiptHandle string `json:"receiptHandle"`
}

type ErrorMessage struct {
	Message string `json:"errorMessage"`
}

func (m *Message) GetJSON() *string {
	b, err := json.Marshal(m)

	if err != nil {
		fmt.Println(err)
	}

	result := string(b)
	return &result
}

func (m *ErrorMessage) GetJSON() *string {
	b, err := json.Marshal(m)

	if err != nil {
		fmt.Println(err)
	}

	result := string(b)
	return &result
}
