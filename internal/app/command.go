package app

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/zainul/gan/internal/app/io"
	"github.com/zainul/gan/internal/app/log"

	"github.com/zainul/gan/internal/app/constant"
)

// Migration is type for creating thing that related with migration database
type MigrationCommand interface {
	CreateFile(name string, extention string, fileType string, additionalInfo ...string) error
	Migrate(status string)
	Seed()
}

type storeMigration struct {
	// dir is directory for migrations file will be placed
	Dir string `json:"dir"`
	// conn is connection string to DB
	Conn string `json:"conn"`
	// seedDir  is directory for seed file will be placed
	SeedDir string `json:"seed_dir"`
}

// NewMigration ..
func NewMigration(dir string, conn string, seedDir string) MigrationCommand {
	os.Setenv(constant.CONNDB, conn)
	os.Setenv(constant.DIR, dir)
	return &storeMigration{dir, conn, seedDir}
}

func (s *storeMigration) Seed() {
	// 1. make build by cmd
	// 2. run the binary
	// 3. delete the binary
	ganseed := "ganseed"

	// 1. make build by cmd
	changeDirectory(s.SeedDir)
	cmd := exec.Command("go", "build", "-o", ganseed)

	if _, err := cmd.CombinedOutput(); err != nil {
		log.Error("read binary error while seed ", err)
		deleteTempFile(s.SeedDir, ganseed)
		os.Exit(2)
	}

	log.Info("step 1. make build by cmd done ...")

	// 2. run the binary
	cmd = exec.Command(fmt.Sprintf("./%v", ganseed))
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Error("Error while run binary ", err, string(out))
		deleteTempFile(s.SeedDir, ganseed)
		os.Exit(2)
	} else {
		log.Info("========================================================")
		log.Info("SEEDER START")
		log.Info("========================================================")
		log.Info(string(out))
		log.Info("========================================================")
	}
	log.Info("step 2. run the binary done...")

	// 3. delete the binary
	deleteTempFile(s.SeedDir, ganseed)
	log.Info("step 3. delete the binary done...")
	os.Exit(2)
}

func (s *storeMigration) Migrate(status string) {
	// 1. create migrate by up or down in selected location
	// 2. make build by cmd
	// 3. run the binary
	// 4. delete the binary and main.go

	// implementation
	// 1. create migrate by up or down in selected location
	if status == constant.StatusUp {
		s.CreateFile("main", constant.DotGo, constant.FileTypeMigrationUp)
	} else if status == constant.StatusDown {
		s.CreateFile("main", constant.DotGo, constant.FileTypeMigrationUp)
	}

	// 2. make build by cmd
	changeDirectory(s.Dir)
	cmd := exec.Command("go", "build", "-o", "ganrun")

	if out, err := cmd.CombinedOutput(); err != nil {
		log.Error("========================================================")
		log.Error("Migration Error")
		log.Error("========================================================")
		log.Error(string(out), err)
		log.Error("========================================================")
		// TODO: make remove temp binary
		deleteTempFile(s.Dir, "main.go")
		deleteTempFile(s.Dir, "ganrun")
		os.Exit(2)
	}

	// 3. run the binary
	changeDirectory(s.Dir)
	cmd = exec.Command("./ganrun")

	if out, err := cmd.CombinedOutput(); err != nil {
		// TODO: make remove temp binary
		log.Error("========================================================")
		log.Error("MIGRATION ERROR")
		log.Error("========================================================")
		log.Error(err, string(out))
		deleteTempFile(s.Dir, "main.go")
		deleteTempFile(s.Dir, "ganrun")
		os.Exit(2)
	} else {
		log.Info("========================================================")
		log.Info("MIGRATION START")
		log.Info("========================================================")
		log.Info(string(out))
		log.Info("========================================================")
	}

	// 4. delete the binary and main.go
	deleteTempFile(s.Dir, "main.go")
	deleteTempFile(s.Dir, "ganrun")
	log.Info("Gan Migration success ...  !!!")
	os.Exit(2)

}

func changeDirectory(dir string) {
	if err := os.Chdir(dir); err != nil {
		log.Error("Could not find migration directory: %s", err)
	}
}

func deleteTempFile(dir, file string) {
	changeDirectory(dir)
	if err := os.Remove(file); err != nil {
		log.Error("Could not remove temporary file: %s", err)
	}
}

func (s *storeMigration) CreateFile(name string, extention string, fileType string, additionalInfo ...string) error {
	// additionalInfo :
	// 0. template
	tableName := name
	name = strings.ToLower(name)
	var structName string
	var structTemplate string
	var dep string

	AppPath := fmt.Sprintf("%v/%v", os.Getenv("GOPATH"), constant.PathAppName)

	sourceFilename := fmt.Sprintf("%v/internal/app/templates/%v.tpl", AppPath, fileType)

	destinationFilename := fmt.Sprintf("%v/%v.%v", s.Dir, name, extention)

	if fileType == constant.FileTypeCreationSeed {
		if len(additionalInfo) > 1 {
			structTemplate = additionalInfo[0]
			structName = additionalInfo[1]
		} else {
			log.Error("Please provide table name")
			os.Exit(2)
		}
		destinationFilename = fmt.Sprintf("%v/%v.%v", s.SeedDir, name, extention)

		if strings.Contains(structTemplate, "time.Time") {
			dep = `
			"time"
			`
		}

		structTemplate = strings.Replace(structTemplate, "Id", "ID", -1)

		structName = strings.Replace(structName, "Id", "ID", -1)
	}

	// detect if file exists
	var _, err = os.Stat(destinationFilename)
	log.Info("will be create file in directory ...", destinationFilename)

	// create file if not exists
	if os.IsNotExist(err) {
		tmpl, err := template.ParseFiles(sourceFilename)

		if err != nil {
			log.Error("Failed create file ...", err)
			return err
		}

		data := struct {
			Key             string
			KeyLowerCase    string
			StructName      string
			StructNameLower string
			StructTemplate  string
			TableName       string
			Dependencies    string
		}{
			Key:             strings.Title(name),
			KeyLowerCase:    strings.ToLower(name),
			StructName:      structName,
			StructNameLower: strings.ToLower(structName),
			StructTemplate:  structTemplate,
			TableName:       tableName,
			Dependencies:    dep,
		}

		var tpl bytes.Buffer
		err = tmpl.Execute(&tpl, data)

		if err != nil {
			log.Error("failed creating the file ", err)
		}

		tplStr := tpl.String()
		io.WriteFile(destinationFilename, tplStr)
	}

	log.Info("done creating file ")

	cmd := exec.Command("gofmt", "-l", "-w", destinationFilename)

	if out, err := cmd.CombinedOutput(); err != nil {
		log.Error("========================================================")
		log.Error("Go Fmt Error")
		log.Error("========================================================")
		log.Error(string(out), err)
		log.Error("========================================================")
	}

	log.Info("Formatting file done ...")

	return nil
}
