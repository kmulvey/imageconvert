package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	log "github.com/sirupsen/logrus"
)

// 1. make sure every file actually exists
// 2. dedup entries
func main() {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner
	s.Start()
	defer s.Stop()

	var oldFile, newFile string
	var h bool
	flag.StringVar(&oldFile, "old-log-file", "", "path to old log file")
	flag.StringVar(&newFile, "new-log-file", "", "path to new log file")
	flag.BoolVar(&h, "help", false, "print options")
	flag.Parse()

	if h {
		flag.PrintDefaults()
		return
	}

	if err := cleanLogFile(oldFile, newFile); err != nil {
		log.Error(err)
	}
}

func cleanLogFile(oldLog, newLog string) error {

	var oldFile, err = os.OpenFile(oldLog, os.O_RDONLY, 0755)
	if err != nil {
		return fmt.Errorf("error opening the old log file: %w", err)
	}
	defer oldFile.Close()

	newFile, err := os.Create(newLog)
	if err != nil {
		return fmt.Errorf("error opening the new log file: %w", err)
	}
	defer oldFile.Close()

	var fileScanner = bufio.NewScanner(oldFile)
	fileScanner.Split(bufio.ScanLines)

	// only keep filenames of files that exist
	var uniqueFiles = make(map[string]struct{})
	for fileScanner.Scan() {
		var filename = fileScanner.Text()
		if _, err := os.Stat(filename); err == nil {
			uniqueFiles[filename] = struct{}{}
		}
	}

	for filename := range uniqueFiles {
		if _, err := newFile.WriteString(filename + "\n"); err != nil {
			return fmt.Errorf("error writing to the new log file: %w", err)
		}
	}
	return nil
}
