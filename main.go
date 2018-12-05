package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli"
	"github.com/zainul/gan/internal/app"
	"github.com/zainul/gan/internal/app/constant"
)

func main() {
	mig := app.NewMigration()

	appCli := cli.NewApp()
	appCli.Name = "gan"
	appCli.Usage = "gan use for migrate and seed the database"

	appCli.Action = func(c *cli.Context) error {
		switch c.Args().Get(0) {
		case constant.Migrate:
		case constant.Seed:
		case constant.Create:
			err := mig.CreateFile(
				fmt.Sprintf("%v_%v", time.Now().Unix(), c.Args().Get(1)),
				constant.DotGo,
				constant.FileTypeMigration,
			)
			fmt.Println("err ", err)
		default:
		}

		return nil
	}

	err := appCli.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
