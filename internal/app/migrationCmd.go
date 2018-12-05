package app

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/zainul/gan/internal/app/constant"
)

// Migration is type for creating thing that related with migration database
type MigrationCommand interface {
	CreateFile(name string, extention string, fileType string) error
}

type storeMigration struct{}

// NewMigration ..
func NewMigration() MigrationCommand {
	return &storeMigration{}
}

func (s *storeMigration) CreateFile(name string, extention string, fileType string) error {

	AppPath := fmt.Sprintf("%v/%v", os.Getenv("GOPATH"), constant.PathAppName)

	sourceFilename := fmt.Sprintf("%v/internal/app/templates/%v.tpl", AppPath, fileType)
	destinationFilename := fmt.Sprintf("%v/%v.%v", os.Getenv("GANDIR"), name, extention)

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

		var tpl bytes.Buffer
		err = tmpl.Execute(&tpl, nil)

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
