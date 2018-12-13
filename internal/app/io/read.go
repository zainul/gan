package io

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

// OpenFile ...
func OpenFile(path string) ([]byte, error) {
	fmt.Println("will be open file at ", path)
	jsonFile, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
		return nil, errors.New("failed to open file")
	}

	fmt.Println("Successfully Opened file ", path)

	byteJSON, err := ioutil.ReadAll(jsonFile)

	defer jsonFile.Close()

	if err != nil {
		fmt.Println(err)
	}

	return byteJSON, err
}
