package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kmulvey/imageconvert/pkg/imageconvert"
	"github.com/kmulvey/path"
	log "github.com/sirupsen/logrus"
)

var randSource *rand.Rand

func init() {
	// rand here is used to generate random strings, it does not need to be crypto secure so we suppress the linter warning
	//nolint:gosec
	randSource = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func main() {

	var inputPath string
	var h bool
	flag.StringVar(&inputPath, "path", "", "path to log")
	flag.BoolVar(&h, "help", false, "print options")
	flag.Parse()

	if h {
		flag.PrintDefaults()
		os.Exit(0)
	}

	var files, err = path.List(inputPath, 2, path.NewRegexEntitiesFilter(imageconvert.ImageExtensionRegex))
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		var filename = filepath.Base(file.AbsolutePath)
		var justName = strings.TrimSuffix(filename, filepath.Ext(filename))

		var newName = strings.Builder{}
		var changed bool
		for _, char := range justName {

			if !validCharacter(char) {
				newName.WriteString(randomString())
				changed = true
			} else {
				newName.WriteRune(char)
			}
		}

		if changed {
			var newPath = filepath.Join(filepath.Dir(file.AbsolutePath), newName.String()+filepath.Ext(filename))

			if _, err := os.Stat(newPath); errors.Is(err, os.ErrNotExist) {
				fmt.Printf("old name: %s \nnew name: %s\n", justName, newName.String())
				fmt.Printf("%s \n\n", newPath)
			} else {
				log.Infof("already exists: %s", newPath)
			}
		}
	}
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

func randomString() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	return string(letters[randSource.Intn(len(letters))])
}
