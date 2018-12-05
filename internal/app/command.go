package app

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/zainul/gan/internal/app/constant"
)

// Migration is type for creating thing that related with migration database
type MigrationCommand interface {
	CreateFile(name string, extention string, fileType string) error
	Migrate(status string)
	SetMigrationDirectory(dir string)
	SetConnectionString(conn string)
}

type storeMigration struct {
	dir  string `json:"dir"`
	conn string `json:"conn"`
}

// NewMigration ..
func NewMigration(dir string, conn string) MigrationCommand {
	os.Setenv(constant.CONNDB, conn)
	return &storeMigration{dir, conn}
}

// SetMigrationDirectory ...
func (s *storeMigration) SetMigrationDirectory(dir string) {
	// if dir != "" {
	// 	err := os.Setenv(constant.DIR, dir)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// }
	// fmt.Println("migration directory has been set to ", dir)
}

// SetConnectionString ...
func (s *storeMigration) SetConnectionString(conn string) {
	// if conn != "" {

	// 	err := os.Setenv(constant.CONNDB, conn)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// }
	// fmt.Println("connection setting already set ", s.conn)
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
	changeDirectory(s.dir)
	cmd := exec.Command("go", "build", "-o", "ganrun")

	if _, err := cmd.CombinedOutput(); err != nil {
		// TODO: make remove temp binary
		deleteTempFile(s.dir, "main.go")
		deleteTempFile(s.dir, "ganrun")
		os.Exit(2)
	}

	// 3. run the binary
	changeDirectory(s.dir)
	cmd = exec.Command("./ganrun")

	if out, err := cmd.CombinedOutput(); err != nil {
		// TODO: make remove temp binary
		deleteTempFile(s.dir, "main.go")
		deleteTempFile(s.dir, "ganrun")
		os.Exit(2)
	} else {
		fmt.Println("========================================================")
		fmt.Println("MIGRATION START")
		fmt.Println("========================================================")
		fmt.Println(string(out))
		fmt.Println("========================================================")
	}

	// 4. delete the binary and main.go
	deleteTempFile(s.dir, "main.go")
	deleteTempFile(s.dir, "ganrun")
	fmt.Println("Gan Migration success ...  !!!")
	os.Exit(2)

}

func changeDirectory(dir string) {
	if err := os.Chdir(dir); err != nil {
		fmt.Println("Could not find migration directory: %s", err)
	}
}

func deleteTempFile(dir, file string) {
	changeDirectory(dir)
	if err := os.Remove(file); err != nil {
		fmt.Println("Could not remove temporary file: %s", err)
	}
}

func (s *storeMigration) CreateFile(name string, extention string, fileType string) error {

	AppPath := fmt.Sprintf("%v/%v", os.Getenv("GOPATH"), constant.PathAppName)

	sourceFilename := fmt.Sprintf("%v/internal/app/templates/%v.tpl", AppPath, fileType)
	destinationFilename := fmt.Sprintf("%v/%v.%v", s.dir, name, extention)

	// detect if file exists
	var _, err = os.Stat(destinationFilename)
	fmt.Println("will be create file in directory ...", destinationFilename)

	// create file if not exists
	if os.IsNotExist(err) {
		tmpl, err := template.ParseFiles(sourceFilename)

		if err != nil {
			fmt.Println("Failed create file ...", err)
			return err
		}

		data := struct {
			Key string
		}{
			Key: name,
		}

		var tpl bytes.Buffer
		err = tmpl.Execute(&tpl, data)

		if err != nil {
			fmt.Println("failed creating the file ", err)
		}

		file, err := os.Create(destinationFilename)
		tplStr := tpl.String()

		_, err = file.Write([]byte(tplStr))

		if err != nil {
			fmt.Println("failed write content to the file ", err)
		}
		file.Close()
	}

	fmt.Println("done creating file ")
	return nil
}
