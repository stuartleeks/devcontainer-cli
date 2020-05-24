package ioutil

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// TODO - recursive copy folder (also add recursive symlink)

func CopyFolder(source string, target string) error {
	copy := func(sourceFolder string, targetFolder string, item os.FileInfo) error {
		return CopyFile(sourceFolder+"/"+item.Name(), targetFolder+"/"+item.Name(), item.Mode())
	}
	return processFolder(source, target, copy)
}
func LinkFolder(source string, target string) error {
	symlink := func(sourceFolder string, targetFolder string, item os.FileInfo) error {
		return os.Symlink(sourceFolder+"/"+item.Name(), targetFolder+"/"+item.Name())
	}
	return processFolder(source, target, symlink)
}

func processFolder(source string, target string, fileHandler func(sourceFolder string, targetFolder string, item os.FileInfo) error) error {
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
			if err = processFolder(source+"/"+sourceSubItem.Name(), target+"/"+sourceSubItem.Name(), fileHandler); err != nil {
				return err
			}
		} else {
			if err = fileHandler(source, target, sourceSubItem); err != nil {
				return err
			}
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
