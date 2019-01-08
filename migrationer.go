package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/zainul/gan/internal/app"
	"github.com/zainul/gan/internal/app/constant"
	"github.com/zainul/gan/internal/app/database"
	"github.com/zainul/gan/internal/app/log"
)

func migrate(cfg Config) {
	mig := app.NewMigration(cfg.Dir, cfg.Conn, cfg.SeedDir)
	mig.Migrate(constant.StatusUp)
}

func seedDataFromFile(cfg Config) {
	mig := app.NewMigration(cfg.Dir, cfg.Conn, cfg.SeedDir)
	mig.Seed()
}

func reverseSeed(cfg Config) {
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
}

func seedTemplate(cfg Config, entityName string) {
	mig := app.NewMigration(cfg.Dir, cfg.Conn, cfg.SeedDir)

	conn, err := sql.Open("postgres", os.Getenv(constant.CONNDB))

	if err != nil {
		log.Error("failed make connection to DB please configure right connection")
		os.Exit(2)
	}
	db := database.NewDB(conn)

	resp, err := db.GetEntity(entityName)

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
	log.Info("completed task: ", entityName)
}

func migrationSQLFromFile(cfg Config, migrationName string) {

	mig := app.NewMigration(cfg.Dir, cfg.Conn, cfg.SeedDir)
	mig.CreateFile(
		fmt.Sprintf("%v_%v_%v",
			migrationName,
			time.Now().Format("20060102_150405"),
			time.Now().UnixNano(),
		),
		constant.DotGo,
		constant.FileTypeMigrationFromFile,
	)
	log.Info("completed task: ", migrationName)
}

func migrationFile(cfg Config, migrationName string) {
	mig := app.NewMigration(cfg.Dir, cfg.Conn, cfg.SeedDir)
	mig.CreateFile(
		fmt.Sprintf("%v_%v_%v",
			migrationName,
			time.Now().Format("20060102_150405"),
			time.Now().UnixNano(),
		),
		constant.DotGo,
		constant.FileTypeMigration,
	)
}
