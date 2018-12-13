package io

import (
	"fmt"
	"os"
)

// WriteFile ...
func WriteFile(path string, content string) {
	file, err := os.Create(path)

	if err != nil {
		fmt.Println("failed to create file ", err)
	}

	_, err = file.Write([]byte(content))

	if err != nil {
		fmt.Println("failed write content to the file ", err)
	}
	file.Close()
}
