package io

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/zainul/gan/internal/entity"
	"github.com/zainul/gan/internal/log"
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

func OpenConfigFile(config string) (entity.Config, error) {
	byteJSON, err := OpenFile(config)

	if err != nil {
		return entity.Config{}, err
	}

	cfg := entity.Config{}

	err = json.Unmarshal(byteJSON, &cfg)

	if err != nil {
		return entity.Config{}, err
	}

	if cfg.Dir == "" || cfg.Conn == "" {
		return entity.Config{}, errors.New("must set config first")
	}
	return cfg, nil
}
