package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/urfave/cli"
	"github.com/zainul/gan/internal/app"
	"github.com/zainul/gan/internal/app/constant"
)

type Config struct {
	Dir  string `json:"dir"`
	Conn string `json:"conn"`
}

func main() {

	var config string

	appCli := cli.NewApp()
	appCli.Name = "gan"
	appCli.Usage = "gan use for migrate and seed the database"
	appCli.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "Load configuration from `FILE`",
			Destination: &config,
		},
	}

	appCli.Action = func(c *cli.Context) error {
		cfg := openFile(config)
		mig := app.NewMigration(cfg.Dir, cfg.Conn)
		switch c.Args().Get(0) {
		case constant.Migrate:
			mig.Migrate(constant.StatusUp)
		case constant.SetDB:
			mig.SetConnectionString(c.Args().Get(1))
		case constant.SetDir:
			mig.SetMigrationDirectory(c.Args().Get(1))
		case constant.Seed:
		case constant.Create:
			mig.CreateFile(
				fmt.Sprintf("%v_%v_%v",
					c.Args().Get(1),
					time.Now().Format("20060102_150405"),
					time.Now().UnixNano(),
				),
				constant.DotGo,
				constant.FileTypeMigration,
			)
		default:
		}

		return nil
	}

	err := appCli.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func openFile(config string) Config {
	jsonFile, err := os.Open(config)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened config file ")

	byteJSON, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	defer jsonFile.Close()
	cfg := Config{}

	err = json.Unmarshal(byteJSON, &cfg)

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	if cfg.Dir == "" || cfg.Conn == "" {
		fmt.Println("Must set config first")
		os.Exit(2)
	}

	fmt.Println("Config ", cfg.Conn)
	fmt.Println("Directory ", cfg.Dir)
	return cfg
}
