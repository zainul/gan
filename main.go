package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli"
	"github.com/zainul/gan/internal/app"
	"github.com/zainul/gan/internal/app/constant"
	"github.com/zainul/gan/internal/app/database"
	"github.com/zainul/gan/internal/app/io"
	"github.com/zainul/gan/internal/app/log"
)

// Config ...
type Config struct {
	Dir     string `json:"dir"`
	Conn    string `json:"conn"`
	SeedDir string `json:"seed_dir"`
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
				mig := app.NewMigration(cfg.Dir, cfg.Conn, cfg.SeedDir)
				mig.Migrate(constant.StatusUp)
				return nil
			},
		},
		{
			Name:  constant.Seed,
			Usage: "Seed the data from file",
			Action: func(c *cli.Context) error {
				cfg := openFile(config)
				mig := app.NewMigration(cfg.Dir, cfg.Conn, cfg.SeedDir)
				mig.Seed()
				return nil
			},
		},
		{
			Name:  constant.ReverseForSeed,
			Usage: "Reverse table to struct and the added to seeder package",
			Action: func(c *cli.Context) error {
				cfg := openFile(config)
				mig := app.NewMigration(cfg.Dir, cfg.Conn, cfg.SeedDir)

				conn, err := sql.Open("postgres", os.Getenv(constant.CONNDB))

				if err != nil {
					log.Error("failed make connection to DB please configure right connection")
					os.Exit(2)
				}
				db := database.NewDB(conn)

				resp, err := db.GetEntity("all")

				if err != nil {
					log.Error("failed to execute get schema ", err)
					os.Exit(2)
				}

				if len(resp) == 0 {
					log.Error("Cannot find table name")
					os.Exit(2)
				}

				log.Warning("+++++++++++++++++++++++++++++++++++++++++++++++++")
				log.Warning("If you doesn't have main.go in your seed directory please copy the script below :")

				log.Info(
					`
					package main

					import (
						"fmt"
						"os"

						"github.com/zainul/gan/pkg/seed"
					)

					func main() {
						db := seed.GetDB()
						gopath := os.Getenv("GOPATH")
						mainDir := fmt.Sprintf("%v/src/github.com/your/directory/to/json", gopath)
					}

					`,
				)

				for _, val := range resp {
					var strOut string
					mig.CreateFile(
						val.TableName,
						constant.DotGo,
						constant.FileTypeCreationSeed,
						val.Models,
						val.StructName,
					)

					strctLower := strings.ToLower(val.StructName)
					strOut = fmt.Sprintf(
						`
					%v, %vs := New%v(db)
					seed.Seed(fmt.Sprintf("%v/%v.json", mainDir), %v, %vs)
					`, strctLower, strctLower, strctLower, "%v", val.StructName, strctLower, strctLower,
					)

					log.Warning("+++++++++++++++++++++++++++++++++++++++++++++++++")
					log.Info(strOut)
					log.Warning("+++++++++++++++++++++++++++++++++++++++++++++++++")
				}

				return nil
			},
		},
		{
			Name:  constant.CreateSeed,
			Usage: "Create seed template file",
			Action: func(c *cli.Context) error {
				cfg := openFile(config)
				mig := app.NewMigration(cfg.Dir, cfg.Conn, cfg.SeedDir)

				conn, err := sql.Open("postgres", os.Getenv(constant.CONNDB))

				if err != nil {
					log.Error("failed make connection to DB please configure right connection")
					os.Exit(2)
				}
				db := database.NewDB(conn)

				resp, err := db.GetEntity(c.Args().First())

				if err != nil {
					log.Error("failed to execute get schema ", err)
					os.Exit(2)
				}

				if len(resp) == 0 {
					log.Error("Cannot find table name")
					os.Exit(2)
				}

				var strOut string
				for _, val := range resp {
					mig.CreateFile(
						val.TableName,
						constant.DotGo,
						constant.FileTypeCreationSeed,
						val.Models,
						val.StructName,
					)

					strctLower := strings.ToLower(val.StructName)
					strOut = fmt.Sprintf(
						`
					%v, %vs := New%v(db)
					seed.Seed(fmt.Sprintf("%v/%v.json", mainDir), %v, %vs)
					`, strctLower, strctLower, strctLower, "%v", val.StructName, strctLower, strctLower,
					)
				}

				log.Warning("+++++++++++++++++++++++++++++++++++++++++++++++++")
				log.Warning("If you doesn't have main.go in your seed directory please copy the script below :")

				log.Info(
					`
					package main

					import (
						"fmt"
						"os"

						"github.com/zainul/gan/pkg/seed"
					)

					func main() {
						db := seed.GetDB()
						gopath := os.Getenv("GOPATH")
						mainDir := fmt.Sprintf("%v/src/github.com/your/directory/to/json", gopath)
					}

					`,
				)

				log.Warning("If already have main.go please add the script below")
				log.Info(strOut)
				log.Warning("+++++++++++++++++++++++++++++++++++++++++++++++++")
				log.Info("completed task: ", c.Args().First())
				return nil
			},
		},
		{
			Name:  constant.CreateFromFile,
			Usage: "Create migration from SQL file",
			Action: func(c *cli.Context) error {
				cfg := openFile(config)
				mig := app.NewMigration(cfg.Dir, cfg.Conn, cfg.SeedDir)
				mig.CreateFile(
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
				cfg := openFile(config)
				mig := app.NewMigration(cfg.Dir, cfg.Conn, cfg.SeedDir)
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
