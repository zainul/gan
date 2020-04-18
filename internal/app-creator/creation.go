package app_creator

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zainul/gan/internal/log"
)

type Creator struct {
	PackageURL         string
	StarterProjectName string
	ProjectName        string
	Package            string
	Path               string
	StarterPackageName string
}

func (c *Creator) CopyFindAndReplace() error {
	cmd := exec.Command("wget", c.PackageURL)

	fmt.Println("Downloading ...")

	if out, err := cmd.CombinedOutput(); err != nil {
		log.Error(err)
		return err
	} else {
		log.Info(string(out))
	}

	fmt.Println("Download completed ...")

	fmt.Println("Unzipping ...")

	branchArr := strings.Split(c.PackageURL, "/")
	branchName := branchArr[len(branchArr)-1]

	cmd = exec.Command("unzip", branchName)

	if out, err := cmd.CombinedOutput(); err != nil {
		log.Error(err)
		return err
	} else {
		log.Info(string(out))
	}

	_ = os.Remove(branchName)

	fmt.Println("Unzip completed")

	cmd = exec.Command("mv", c.StarterProjectName+fmt.Sprintf("-%s", strings.Replace(branchName, ".zip", "", -1)), c.ProjectName)

	if out, err := cmd.CombinedOutput(); err != nil {
		log.Error(err)
		return err
	} else {
		log.Info(string(out))
	}

	err := os.MkdirAll(c.Path+"/src/"+c.Package+"/", os.ModePerm)

	if err != nil {
		log.Error(err)
		return err
	}

	err = filepath.Walk(c.ProjectName, c.visit)
	if err != nil {
		log.Error(err)
		return err
	}

	fmt.Println("will be send a gift for you ...")

	log.Info("move", c.ProjectName, "to", c.Path+"/src/"+c.Package+"/")
	errMove := CopyDir(c.ProjectName, c.Path+"/"+c.Package+"/")

	if errMove != nil {
		log.Error(err)
		return err
	}

	return nil
}

func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		fmt.Println("copyFile open", err)
		return
	}
	defer func() {
		_ = in.Close()
	}()

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
		fmt.Println("CopyDir is not exist", err)
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

func (c *Creator) visit(path string, fi os.FileInfo, err error) error {

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

		newContents := strings.Replace(string(read), c.StarterPackageName, c.Package+"/"+c.ProjectName, -1)

		err = ioutil.WriteFile(path, []byte(newContents), 0)
		if err != nil {
			return err
		}

	}

	return nil
}
