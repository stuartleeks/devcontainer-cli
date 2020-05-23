package ioutil

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// TODO - recursive copy folder (also add recursive symlink)

func CopyFolder(source string, target string) error {
	sourceItem, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("Error reading source folder: %s\n", err)
	}
	if err = os.Mkdir(target, sourceItem.Mode()); err != nil {
		return fmt.Errorf("Error creating directory '%s': %s", target, err)
	}

	sourceSubItems, err := ioutil.ReadDir(source)
	if err != nil {
		return fmt.Errorf("Error reading source folder contents: %s\n", err)
	}

	for _, sourceSubItem := range sourceSubItems {
		if sourceSubItem.IsDir() {
			CopyFolder(source+"/"+sourceSubItem.Name(), target+"/"+sourceSubItem.Name())
		} else {
			CopyFile(source+"/"+sourceSubItem.Name(), target+"/"+sourceSubItem.Name(), sourceSubItem.Mode())
		}
	}
	return nil
}
func CopyFile(source string, target string, perm os.FileMode) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	targetFile, err := os.OpenFile(target, os.O_RDWR|os.O_CREATE|os.O_EXCL|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer targetFile.Close()
	_, err = io.Copy(targetFile, sourceFile)
	return err
}
