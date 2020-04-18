package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"os"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/urfave/cli"

	"github.com/zainul/gan/internal/app"
	app_creator "github.com/zainul/gan/internal/app-creator"
	"github.com/zainul/gan/internal/constant"
	"github.com/zainul/gan/internal/database"
	"github.com/zainul/gan/internal/entity"
	"github.com/zainul/gan/internal/io"
	"github.com/zainul/gan/internal/log"
)

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: constant.Migrate, Description: "Migrate migrations script"},
		{Text: constant.Seed, Description: "Seed the data from file"},
		{Text: constant.ReverseForSeed, Description: "Reverse table to struct and the added to seeder package"},
		{Text: constant.CreateSeed, Description: "Create seed template file"},
		{Text: constant.CreateFromFile, Description: "Create migration from SQL file"},
		{Text: constant.Create, Description: "Create migration file"},
		{Text: constant.CreateApp, Description: "Create starter apps"},
		{Text: constant.ReverseDB, Description: "Reverse DB to Entity and Repo"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func completerEmpty(d prompt.Document) []prompt.Suggest {
	s := make([]prompt.Suggest, 0)
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func completerConfig(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "migrations/config.json", Description: "By Author"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func completerConfigCreateApp(d prompt.Document) []prompt.Suggest {
	s := make([]prompt.Suggest, 0)
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
				if cfg, err := io.OpenConfigFile(config); err != nil {
					log.Error(fmt.Sprintf("Error opening file %+v", err))
				} else {
					mig := app.NewMigration(cfg.Dir, cfg.Conn, cfg.SeedDir)
					mig.Migrate(constant.StatusUp)
				}
				return nil
			},
		},
		{
			Name:  constant.Seed,
			Usage: "Seed the data from file",
			Action: func(c *cli.Context) error {
				if cfg, err := io.OpenConfigFile(config); err != nil {
					log.Error(fmt.Sprintf("Error opening file %+v", err))
				} else {
					mig := app.NewMigration(cfg.Dir, cfg.Conn, cfg.SeedDir)
					mig.Seed()
				}
				return nil
			},
		},
		{
			Name:  constant.ReverseForSeed,
			Usage: "Reverse table to struct and the added to seeder package",
			Action: func(c *cli.Context) error {
				if cfg, err := io.OpenConfigFile(config); err != nil {
					log.Error(fmt.Sprintf("Error opening file %+v", err))
				} else {
					reverseSeed(cfg)
				}
				return nil
			},
		},
		{
			Name:  constant.CreateSeed,
			Usage: "Create seed template file",
			Action: func(c *cli.Context) error {
				var cfg entity.Config
				var err error
				if cfg, err = io.OpenConfigFile(config); err != nil {
					log.Error(fmt.Sprintf("Error opening file %+v", err))
				}

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
				var cfg entity.Config
				var err error
				if cfg, err = io.OpenConfigFile(config); err != nil {
					log.Error(fmt.Sprintf("Error opening file %+v", err))
				}

				if c.Args().First() == "" {
					log.Error("Please provide 1st argument")
					return errors.New("argument not completed")
				}

				mig := app.NewMigration(cfg.Dir, cfg.Conn, cfg.SeedDir)
				_ = mig.CreateFile(
					fmt.Sprintf("%v_%v_%v",
						c.Args().First(),
						time.Now().Format("20060102_150405"),
						time.Now().UnixNano(),
					),
					constant.DotGo,
					constant.FileTypeMigrationFromFile,
				)
				log.Info("completed task: ", c.Args().First())
				return nil
			},
		},
		{
			Name:  constant.Create,
			Usage: "Create migration file",
			Action: func(c *cli.Context) error {
				var cfg entity.Config
				var err error
				if cfg, err = io.OpenConfigFile(config); err != nil {
					log.Error(fmt.Sprintf("Error opening file %+v", err))
				}

				if c.Args().First() == "" {
					log.Error("Please provide 1st argument")
					return errors.New("argument not completed")
				}

				migrationFile(cfg, c.Args().First())

				return nil
			},
		},
		{
			Name:  constant.CreateApp,
			Usage: "Create Starter Apps service",
			Action: func(c *cli.Context) error {
				fmt.Println("Hi , I will serve service for you with ♥ ")
				return nil
			},
		},
		{
			Name:  constant.ReverseDB,
			Usage: "Reverse DB to Entity Repo",
			Action: func(c *cli.Context) error {
				fmt.Println("Hi , I will serve you with ♥ ")
				return nil
			},
		},
	}

	err := appCli.Run(os.Args)

	fmt.Println("Please select action that you want")
	t := prompt.Input("> ", completer)
	t = strings.TrimSpace(t)
	config = prompt.Input("Where is the config file stored ?", completerConfig)
	var cfg entity.Config

	if cfg, err = io.OpenConfigFile(config); err != nil {
		log.Error(fmt.Sprintf("Error opening file %+v", err))
	}

	switch t {
	case constant.Migrate:
		mig := app.NewMigration(cfg.Dir, cfg.Conn, cfg.SeedDir)
		mig.Migrate(constant.StatusUp)
		break
	case constant.Seed:
		mig := app.NewMigration(cfg.Dir, cfg.Conn, cfg.SeedDir)
		mig.Seed()
		break
	case constant.ReverseForSeed:
		reverseSeed(cfg)
		break
	case constant.CreateSeed:
		var entity string
		for {
			entity = prompt.Input("What name for entity that you want to seed ?", completerEmpty)
			if strings.TrimSpace(entity) != "" {
				break
			}
		}
		if cfg, err := io.OpenConfigFile(config); err != nil {
			log.Error(fmt.Sprintf("Error opening file %+v", err))
		} else {
			seedTemplate(cfg, entity)
		}
		break
	case constant.CreateFromFile:
		var entity string
		for {
			entity = prompt.Input("What name of migration ?", completerEmpty)
			if strings.TrimSpace(entity) != "" {
				break
			}
		}
		if cfg, err := io.OpenConfigFile(config); err != nil {
			log.Error(fmt.Sprintf("Error opening file %+v", err))
		} else {
			mig := app.NewMigration(cfg.Dir, cfg.Conn, cfg.SeedDir)
			mig.Migrate(constant.StatusUp)
			mig.CreateFile(
				fmt.Sprintf("%v_%v_%v",
					entity,
					time.Now().Format("20060102_150405"),
					time.Now().UnixNano(),
				),
				constant.DotGo,
				constant.FileTypeMigrationFromFile,
			)
		}

		break
	case constant.Create:
		var entity string
		for {
			entity = prompt.Input("What name of migration ?", completerEmpty)
			if strings.TrimSpace(entity) != "" {
				break
			}
		}
		if cfg, err := io.OpenConfigFile(config); err != nil {
			log.Error(fmt.Sprintf("Error opening file %+v", err))
		} else {
			migrationFile(cfg, entity)
		}
		break
	case constant.CreateApp:
		name := prompt.Input("What is your service name ? (ex: danisa) ", completerConfigCreateApp)
		packageName := prompt.Input("What your root of your package you want ? (ex: github.com/zainul) ", completerConfigCreateApp)
		yourPath := prompt.Input(fmt.Sprintf("Your GOPATH ? ex : /home/zainul/go"), completerConfigCreateApp)
		fmt.Println("Hi , I will serve service " + name + " for you with ♥ ")

		c := app_creator.Creator{
			PackageURL:         cfg.CreatorApp.PackageURL,
			StarterProjectName: cfg.CreatorApp.StarterProjectName,
			StarterPackageName: cfg.CreatorApp.StarterPackageName,
			ProjectName:        name,
			Package:            packageName,
			Path:               yourPath,
		}

		err := c.CopyFindAndReplace()

		if err != nil {
			log.Error("Error when creating app :( ", err)
		}

		break
	case constant.ReverseDB:
		var (
			cfg entity.Config
			err error
		)

		if cfg, err = io.OpenConfigFile(config); err != nil {
			log.Error(fmt.Sprintf("Error opening file %+v", err))
			return
		}

		mig := app.NewMigration(cfg.Dir, cfg.Conn, cfg.SeedDir, cfg.ProjectStructure, cfg.ProjectPackage)

		if cfg.ProjectStructure != (entity.ProjectStructure{}) {
			conn, err := sql.Open("postgres", cfg.Conn)
			db := database.NewDB(conn)
			structs, err := db.GetEntityWithoutTableName()
			if err != nil {
				log.Error(err)
				os.Exit(2)
			}

			for _, strc := range structs {
				mig.CreateFile(strc.TableName, constant.DotGo, constant.FileTypeReverse, strc.Models, constant.CreateEntity, strc.TableName)
			}

			for _, strc := range structs {
				mig.CreateFile(strc.TableName, constant.DotGo, constant.FileTypeReverse, strc.Models, constant.CreateUseCase, strc.TableName)
			}

			for _, strc := range structs {
				mig.CreateFile(strc.TableName, constant.DotGo, constant.FileTypeReverse, strc.Models, constant.CreateStore, strc.TableName)
			}

			for _, strc := range structs {
				mig.CreateFile(strc.TableName, constant.DotGo, constant.FileTypeReverse, strc.Models, constant.CreateStoreImpl, strc.TableName)
			}

		}
		break
	}
}
