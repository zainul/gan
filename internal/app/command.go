package app

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"

	"github.com/zainul/gan/internal/constant"
	"github.com/zainul/gan/internal/entity"
	"github.com/zainul/gan/internal/io"
	"github.com/zainul/gan/internal/log"
)

// Migration is type for creating thing that related with migration database
type MigrationCommand interface {
	CreateFile(name string, extention string, fileType string, additionalInfo ...interface{}) error
	Migrate(status string)
	Seed()
}


type storeMigration struct {
	// dir is directory for migrations file will be placed
	Dir string `json:"dir"`
	// conn is connection string to DB
	Conn string `json:"conn"`
	// seedDir  is directory for seed file will be placed
	SeedDir          string            `json:"seed_dir"`
	ProjectStructure *entity.ProjectStructure `json:"project_structure"`

	ProjectPackage string `json:"project_package"`
}

// NewMigration ..
func NewMigration(dir string, conn string, seedDir string, projectStructere ...interface{}) MigrationCommand {
	os.Setenv(constant.CONNDB, conn)
	os.Setenv(constant.DIR, dir)

	var pj *entity.ProjectStructure
	var pk string

	if projectStructere != nil && len(projectStructere) > 0 {
		item := projectStructere[0]
		pj = item.(*entity.ProjectStructure)
	}

	if len(projectStructere) > 1 {
		pk = projectStructere[1].(string)
		os.Setenv(constant.THORPACKAGE, pk)
	}
	return &storeMigration{
		Dir:              dir,
		Conn:             conn,
		SeedDir:          seedDir,
		ProjectStructure: pj,
		ProjectPackage:   pk,
	}
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

func (s *storeMigration) CreateFile(name string, extention string, fileType string, additionalInfo ...interface{}) error {
	var process string
	var customTemplateInput string
	var ReqTemplateInput string
	var tableName string
	structName := name
	name = strings.ToLower(name)

	AppPath := fmt.Sprintf("%v/%v", os.Getenv("GOPATH"), constant.PathAppName)

	sourceFilename := fmt.Sprintf("%v/internal/templates/%v.tpl", AppPath, fileType)

	destinationFilename := fmt.Sprintf("%v/%v.%v", s.Dir, name, extention)

	if fileType == constant.FileTypeCreationSeed {
		destinationFilename = fmt.Sprintf("%v/%v.%v", s.SeedDir, name, extention)
	} else if fileType == constant.FileTypeReverse {

		if additionalInfo != nil && len(additionalInfo) > 1 {
			process = additionalInfo[1].(string)
		} else {
			log.Error("Please provide valid config for reverse")
			os.Exit(2)
		}

		customTemplateInput = additionalInfo[0].(string)
		tableName = additionalInfo[2].(string)
		ReqTemplateInput = customTemplateInput

		if process == constant.CreateEntity {
			destinationFilename = fmt.Sprintf("%v/%v.%v", s.ProjectStructure.Entity.Dir, name, extention)
			sourceFilename = fmt.Sprintf("%v/internal/templates/%v.tpl", AppPath, "entity")

			if strings.Contains(customTemplateInput, "time.Time") {
				customTemplateInput = `
				import (
					"time"
				) 

				//` + structName + ` ....
				` + customTemplateInput
			} else {
				customTemplateInput = `
				//` + structName + ` ....
				` + customTemplateInput
			}

			customTemplateInput = strings.Replace(customTemplateInput, "Id", "ID", -1)

		} else if process == constant.CreateUseCase {
			destinationFilename = fmt.Sprintf("%v/%v.%v", s.ProjectStructure.UseCase.Dir, name, extention)
			sourceFilename = fmt.Sprintf("%v/internal/templates/%v.tpl", AppPath, "usecase")
		} else if process == constant.CreateStore {
			destinationFilename = fmt.Sprintf("%v/%v.%v", s.ProjectStructure.Store.Dir, name, extention)
			sourceFilename = fmt.Sprintf("%v/internal/templates/%v.tpl", AppPath, "store")
		} else if process == constant.CreateStoreImpl {
			os.Mkdir(fmt.Sprintf("%v", s.ProjectStructure.Store.Dir+"/store/"), 0700)
			destinationFilename = fmt.Sprintf("%v/%v.%v", s.ProjectStructure.Store.Dir+"/store", name, extention)
			sourceFilename = fmt.Sprintf("%v/internal/templates/%v.tpl", AppPath, "implementation_store_pg")
		}

	}

	// detect if file exists
	var _, err = os.Stat(destinationFilename)
	log.Info("will be create file in directory : ", destinationFilename)

	// create file if not exists
	if os.IsNotExist(err) {
		tmpl, err := template.ParseFiles(sourceFilename)

		if err != nil {
			log.Error("Failed create file ...", err)
			return err
		}

		ReqTemplateInput = strings.Replace(ReqTemplateInput, "type ", "type Request", -1)

		data := struct {
			Key                     string
			KeyLowerCase            string
			CustomTemplateFromInput string
			Name                    string
			Package                 string
			ReqTemplate             string
			TableName               string
		}{
			Key:                     strings.Title(name),
			KeyLowerCase:            strings.ToLower(name),
			CustomTemplateFromInput: customTemplateInput,
			Name:                    structName,
			Package:                 os.Getenv(constant.THORPACKAGE),
			ReqTemplate:             ReqTemplateInput,
			TableName:               ToSnakeCase(tableName),
		}

		var tpl bytes.Buffer
		err = tmpl.Execute(&tpl, data)

		if err != nil {
			log.Error("failed creating the file ", err)
		}

		tplStr := tpl.String()
		io.WriteFile(destinationFilename, tplStr)
	} else {
		log.Warning("File already exist")
		return nil
	}

	log.Info("done creating file ", destinationFilename)

	cmd := exec.Command("gofmt", "-l", "-w", destinationFilename)

	if out, err := cmd.CombinedOutput(); err != nil {
		log.Error("========================================================")
		log.Error("Go Fmt Error")
		log.Error("========================================================")
		log.Error(string(out), err)
		log.Error("========================================================")
	}

	log.Info("Go fmt file done ...")
	return nil
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
