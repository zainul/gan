package io

import (
	"os"

	"github.com/zainul/gan/internal/app/log"
)

// WriteFile ...
func WriteFile(path string, content string) {
	file, err := os.Create(path)

	if err != nil {
		log.Error("failed to create file ", err)
	}

	_, err = file.Write([]byte(content))

	if err != nil {
		log.Error("failed write content to the file ", err)
	}
	file.Close()
}
