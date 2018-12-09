package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli"
	"github.com/zainul/gan/internal/app"
	"github.com/zainul/gan/internal/app/constant"
	"github.com/zainul/gan/internal/app/io"
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
	appCli.Version = "0.0.1"
	appCli.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "Load configuration from `FILE`",
			Destination: &config,
		},
	}

	appCli.Commands = []cli.Command{
		{
			Name:  constant.Migrate,
			Usage: "Migrate migrations script",
			Action: func(c *cli.Context) error {
				cfg := openFile(config)
				mig := app.NewMigration(cfg.Dir, cfg.Conn)
				mig.Migrate(constant.StatusUp)
				return nil
			},
		},
		{
			Name:  constant.CreateFromFile,
			Usage: "Create migration from SQL file",
			Action: func(c *cli.Context) error {
				cfg := openFile(config)
				mig := app.NewMigration(cfg.Dir, cfg.Conn)
				mig.CreateFile(
					fmt.Sprintf("%v_%v_%v",
						c.Args().Get(1),
						time.Now().Format("20060102_150405"),
						time.Now().UnixNano(),
					),
					constant.DotGo,
					constant.FileTypeMigrationFromFile,
				)
				fmt.Println("completed task: ", c.Args().First())
				return nil
			},
		},
		{
			Name:  constant.Create,
			Usage: "Create migration file",
			Action: func(c *cli.Context) error {
				cfg := openFile(config)
				mig := app.NewMigration(cfg.Dir, cfg.Conn)
				mig.CreateFile(
					fmt.Sprintf("%v_%v_%v",
						c.Args().First(),
						time.Now().Format("20060102_150405"),
						time.Now().UnixNano(),
					),
					constant.DotGo,
					constant.FileTypeMigration,
				)
				return nil
			},
		},
	}

	err := appCli.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func openFile(config string) Config {
	byteJSON, err := io.OpenFile(config)

	if err != nil {
		fmt.Println("Failed to open file ", err)
		os.Exit(2)
	}

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
