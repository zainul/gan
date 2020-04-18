package main

import (
	"fmt"
	"os"

	"github.com/zainul/gan/internal/app"
)

func main() {
	db := app.GetDB()
	gopath := os.Getenv("GOPATH")
	mainDir := fmt.Sprintf("%v/src/github.com/zainul/gan/examples/json", gopath)

	house, houses := NewHouse(db)

	app.Seed(fmt.Sprintf("%v/house.json", mainDir), house, houses)
}
