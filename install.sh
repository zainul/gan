#!/bin/bash

go get github.com/jinzhu/gorm
go get github.com/c-bata/go-prompt
go build
sudo cp ${GOPATH}/src/github.com/zainul/gan/gan /usr/bin/