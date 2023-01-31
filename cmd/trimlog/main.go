package main

import (
	"bufio"
	"flag"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/kmulvey/path"
	log "github.com/sirupsen/logrus"
)

// 1. make sure every file actually exists
// 2. dedup
func main() {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner
	s.Start()
	defer s.Stop()

	var inputPath path.Path
	var h bool
	flag.Var(&inputPath, "path", "path to log")
	flag.BoolVar(&h, "help", false, "print options")
	flag.Parse()

	if h {
		flag.PrintDefaults()
		os.Exit(0)
	}

	originalFile, err := os.Open(inputPath.ComputedPath.AbsolutePath)
	if err != nil {
		log.Fatal(err)
	}
	defer originalFile.Close()

	newFile, err := os.Create("./new.log")
	if err != nil {
		log.Fatal(err)
	}
	fileScanner := bufio.NewScanner(originalFile)
	fileScanner.Split(bufio.ScanLines)

	// only keep filenames of files that exist
	var uniqueFiles = make(map[string]struct{})
	for fileScanner.Scan() {
		var entry = fileScanner.Text()
		var _, err = path.NewEntry(entry)
		if err == nil {
			uniqueFiles[entry] = struct{}{}
		}
	}

	for filename := range uniqueFiles {
		if _, err := newFile.WriteString(filename + "\n"); err != nil {
			log.Fatal(err)
		}
	}
}
