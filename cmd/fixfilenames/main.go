package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/kmulvey/imageconvert/v2/internal/app/imageconvert"
	"github.com/kmulvey/path"
	log "github.com/sirupsen/logrus"
)

func main() {

	var inputPath string
	var h bool
	flag.StringVar(&inputPath, "path", "", "path to directory")
	flag.BoolVar(&h, "help", false, "print options")
	flag.Parse()

	if h {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if err := renameFiles(inputPath); err != nil {
		log.Fatal(err)
	}
}

func renameFiles(inputPath string) error {
	var files, err = path.List(inputPath, 2, false, path.NewRegexEntitiesFilter(imageconvert.ImageExtensionRegex))
	if err != nil {
		return err
	}

	for _, file := range files {
		var newFileName, changed = changeFileName(file)
		if changed {
			var newPath = filepath.Join(filepath.Dir(file.AbsolutePath), newFileName+filepath.Ext(file.AbsolutePath))

			if _, err := os.Stat(newPath); errors.Is(err, os.ErrNotExist) {
				if err := moveFile(file.AbsolutePath, newPath); err != nil {
					return err
				}
				log.Infof("old name: %s, new name: %s", filepath.Base(file.AbsolutePath), newFileName)
			} else {
				log.Infof("already exists: %s", newPath)
			}
		}
	}
	return nil
}

func changeFileName(file path.Entry) (string, bool) {
	var filename = filepath.Base(file.AbsolutePath)
	var justName = strings.TrimSuffix(filename, filepath.Ext(filename))

	var newName = strings.Builder{}
	var changed bool
	for _, char := range justName {

		if !validCharacter(char) {
			newName.WriteString("a")
			changed = true
		} else {
			newName.WriteRune(char)
		}
	}
	return newName.String(), changed
}

func validCharacter(r rune) bool {
	if (r >= 48 && r <= 57) || // digits
		(r >= 65 && r <= 90) || // capital letters
		(r >= 97 && r <= 122) || // lower case letters
		(r == 45 || r == 95) { // - and _
		return true
	}

	return false
}

func moveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %v", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("Couldn't open dest file: %v", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return fmt.Errorf("Couldn't copy to dest from source: %v", err)
	}

	inputFile.Close() // for Windows, close before trying to remove: https://stackoverflow.com/a/64943554/246801

	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't remove source file: %v", err)
	}
	return nil
}
