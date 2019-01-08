package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/urfave/cli"
	"github.com/zainul/gan/internal/app/constant"
	"github.com/zainul/gan/internal/app/io"
	"github.com/zainul/gan/internal/app/log"
)

// Config ...
type Config struct {
	Dir     string `json:"dir"`
	Conn    string `json:"conn"`
	SeedDir string `json:"seed_dir"`
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: constant.Migrate, Description: "Migrate migrations script"},
		{Text: constant.Seed, Description: "Seed the data from file"},
		{Text: constant.ReverseForSeed, Description: "Reverse table to struct and the added to seeder package"},
		{Text: constant.CreateSeed, Description: "Create seed template file"},
		{Text: constant.CreateFromFile, Description: "Create migration from SQL file"},
		{Text: constant.Create, Description: "Create migration file"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func completerEmpty(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func completerConfig(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "migrations/config.json", Description: "By Author"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
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

	fmt.Println(
		`
--------------------------
______ _______ __   _
|  ____ |_____| | \  |
|_____| |     | |  \_|
					  
--------------------------
`,
	)

	appCli.Commands = []cli.Command{
		{
			Name:  constant.Migrate,
			Usage: "Migrate migrations script",
			Action: func(c *cli.Context) error {
				cfg := openFile(config)
				migrate(cfg)
				return nil
			},
		},
		{
			Name:  constant.Seed,
			Usage: "Seed the data from file",
			Action: func(c *cli.Context) error {
				cfg := openFile(config)
				seedDataFromFile(cfg)
				return nil
			},
		},
		{
			Name:  constant.ReverseForSeed,
			Usage: "Reverse table to struct and the added to seeder package",
			Action: func(c *cli.Context) error {
				cfg := openFile(config)
				reverseSeed(cfg)
				return nil
			},
		},
		{
			Name:  constant.CreateSeed,
			Usage: "Create seed template file",
			Action: func(c *cli.Context) error {
				cfg := openFile(config)

				if c.Args().First() == "" {
					log.Error("Please provide 1st argument")
					return errors.New("argument not completed")
				}
				seedTemplate(cfg, c.Args().First())
				return nil
			},
		},
		{
			Name:  constant.CreateFromFile,
			Usage: "Create migration from SQL file",
			Action: func(c *cli.Context) error {
				cfg := openFile(config)
				if c.Args().First() == "" {
					log.Error("Please provide 1st argument")
					return errors.New("argument not completed")
				}
				migrationSQLFromFile(cfg, c.Args().First())
				return nil
			},
		},
		{
			Name:  constant.Create,
			Usage: "Create migration file",
			Action: func(c *cli.Context) error {
				cfg := openFile(config)
				if c.Args().First() == "" {
					log.Error("Please provide 1st argument")
					return errors.New("argument not completed")
				}
				migrationFile(cfg, c.Args().First())
				return nil
			},
		},
	}

	err := appCli.Run(os.Args)

	fmt.Println("Please select action that you want")
	t := prompt.Input("> ", completer)
	t = strings.TrimSpace(t)
	switch t {
	case constant.Migrate:
		config := prompt.Input("Where is the config file stored ?", completerConfig)
		cfg := openFile(config)
		migrate(cfg)
		break
	case constant.Seed:
		config := prompt.Input("Where is the config file stored ?", completerConfig)
		cfg := openFile(config)
		seedDataFromFile(cfg)
		break
	case constant.ReverseForSeed:
		config := prompt.Input("Where is the config file stored ?", completerConfig)
		cfg := openFile(config)
		reverseSeed(cfg)
		break
	case constant.CreateSeed:
		config := prompt.Input("Where is the config file stored ?", completerConfig)
		var entity string
		for {
			entity = prompt.Input("What name for entity that you want to seed ?", completerEmpty)
			if strings.TrimSpace(entity) != "" {
				break
			}
		}
		cfg := openFile(config)
		seedTemplate(cfg, entity)
		break
	case constant.CreateFromFile:
		config := prompt.Input("Where is the config file stored ?", completerConfig)
		var entity string
		for {
			entity = prompt.Input("What name of migration ?", completerEmpty)
			if strings.TrimSpace(entity) != "" {
				break
			}
		}
		cfg := openFile(config)
		migrationSQLFromFile(cfg, entity)
		break
	case constant.Create:
		config := prompt.Input("Where is the config file stored ?", completerConfig)
		var entity string
		for {
			entity = prompt.Input("What name of migration ?", completerEmpty)
			if strings.TrimSpace(entity) != "" {
				break
			}
		}
		cfg := openFile(config)
		migrationFile(cfg, entity)
		break
	}

	if err != nil {
		log.Error(err)
	}
}

func openFile(config string) Config {
	byteJSON, err := io.OpenFile(config)

	if err != nil {
		log.Error("Failed to open file ", err)
		os.Exit(2)
	}

	cfg := Config{}

	err = json.Unmarshal(byteJSON, &cfg)

	if err != nil {
		log.Error(err)
		os.Exit(2)
	}

	if cfg.Dir == "" || cfg.Conn == "" {
		log.Error("Must set config first")
		os.Exit(2)
	}

	log.Warning("Config ", cfg.Conn)
	log.Warning("Directory ", cfg.Dir)
	return cfg
}
