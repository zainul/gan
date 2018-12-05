package main

import (
	"github.com/zainul/gan/internal/app"
	"github.com/zainul/gan/internal/app/constant"
)

func main() {
	m := app.Migration{}
	m.Exec(constant.StatusDown)
}
