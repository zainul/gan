#!/bin/bash

go get github.com/jinzhu/gorm
go build
sudo cp ${GOPATH}/src/github.com/zainul/gan/gan /usr/bin/