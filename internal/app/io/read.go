package io

import (
	"fmt"
	"io/ioutil"
	"os"
)

func OpenFile(path string) ([]byte, error) {
	jsonFile, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened file ", path)

	byteJSON, err := ioutil.ReadAll(jsonFile)

	defer jsonFile.Close()

	if err != nil {
		fmt.Println(err)
	}

	return byteJSON, err
}
