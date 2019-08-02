package main

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/zainul/gan/internal/app"
	"github.com/zainul/gan/internal/app/constant"
	"github.com/zainul/gan/internal/app/database"
	"github.com/zainul/gan/internal/app/log"
)

var importedPackage string

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

func copyFindAndReplace(name, pkg, path string) {
	gitubName := "arkana-kit"
	cmd := exec.Command("wget", "https://github.com/zainul/arkana-kit/archive/master.zip")
	fmt.Println("Downloading ...")
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Error(err)
		os.Exit(2)
	} else {
		log.Info(string(out))
	}

	fmt.Println("Download completed ...")

	fmt.Println("Unzipping ...")
	cmd = exec.Command("unzip", "master.zip")

	if out, err := cmd.CombinedOutput(); err != nil {
		log.Error(err)
		os.Exit(2)
	} else {
		log.Info(string(out))
	}

	os.Remove("master.zip")

	fmt.Println("Unzip completed")

	cmd = exec.Command("mv", gitubName+"-master", name)

	if out, err := cmd.CombinedOutput(); err != nil {
		log.Error(err)
		os.Exit(2)
	} else {
		log.Info(string(out))
	}

	err := os.MkdirAll(path+"/src/"+pkg+"/", os.ModePerm)

	if err != nil {
		log.Error(err)
		os.Exit(2)
	}

	importedPackage = pkg + "/" + name

	err = filepath.Walk(name, visit)
	if err != nil {
		log.Error(err)
		os.Exit(2)
	}

	fmt.Println("will be send a gift for you ...")

	log.Info("move", name, "to", path+"/src/"+pkg+"/")
	errMove := CopyDir(name, path+"/"+pkg+"/")

	if errMove != nil {
		log.Error(err)
		os.Exit(2)
	}

}

func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		fmt.Println("copyFile open", err)
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		fmt.Println("copyFile create", err)
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		fmt.Println("copyFile copy", err)
		return
	}

	err = out.Sync()
	if err != nil {
		fmt.Println("copyFile sync", err)
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		fmt.Println("copyFile stat", err)
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		fmt.Println("CopyDir stat", err)
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("CopyDir isnot exist", err)
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		fmt.Println("CopyDir mkdir all", err)
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		fmt.Println("CopyDir reaDir", err)
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				fmt.Println("CopyDir copy dir", err)
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				fmt.Println("CopyDir copy file", err)
				return
			}
		}
	}

	return
}

func visit(path string, fi os.FileInfo, err error) error {

	if err != nil {
		return err
	}

	if !!fi.IsDir() {
		return nil //
	}

	matched, err := filepath.Match("*.go", fi.Name())

	if err != nil {
		return err
	}

	if matched {
		read, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		//fmt.Println(string(read))

		newContents := strings.Replace(string(read), "github.com/zainul/arkana-kit", importedPackage, -1)

		err = ioutil.WriteFile(path, []byte(newContents), 0)
		if err != nil {
			return err
		}

	}

	return nil
}
