package io

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/zainul/gan/internal/app/log"
)

// OpenFile ...
func OpenFile(path string) ([]byte, error) {
	log.Info("will be open file at ", path)
	jsonFile, err := os.Open(path)

	if err != nil {
		log.Error(err)
		return nil, errors.New("failed to open file")
	}

	log.Info("Successfully Opened file ", path)

	byteJSON, err := ioutil.ReadAll(jsonFile)

	defer jsonFile.Close()

	if err != nil {
		log.Error(err)
	}

	return byteJSON, err
}
