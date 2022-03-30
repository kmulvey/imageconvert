package main

import (
	"bufio"
	"flag"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/kmulvey/imageconvert/pkg/imageconvert"
)

const staticBool = false

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	var rootDir string
	var processedLogFile string
	var compress bool

	flag.StringVar(&rootDir, "dir", "", "directory (abs path), could also be a single file")
	flag.StringVar(&processedLogFile, "log-file", "processed.log", "the file to write processes images to, so that we dont processes them again next time")
	flag.BoolVar(&compress, "compress", false, "compress")
	flag.Parse()
	if strings.TrimSpace(rootDir) == "" {
		log.Fatal("directory not provided")
	}

	// open the file
	log.Info("reading log file")
	var processedLog, err = os.OpenFile(processedLogFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	imageconvert.HandleErr("processedLog open", err)
	defer func() {
		imageconvert.HandleErr("close processedLog", processedLog.Close())
	}()

	var skipMap = getSkipMap(processedLog)

	// Did they give us a dir or file?
	fileInfo, err := os.Stat(rootDir)
	if err != nil {
		log.Fatal("could not stat file/dir ", err)
	}
	var files = make([]string, 1) // we will always have at least one
	if !fileInfo.IsDir() {
		files[0] = rootDir
	} else {
		// these are all the files all the way down the dir tree
		files = imageconvert.ListFiles(rootDir, skipMap)
	}

	// convert all the images to jpg
	log.Info("converting images to jpeg")
	var conversionTotals = make(map[string]int)
	var imageType string
	for i, filename := range files {
		files[i], imageType = imageconvert.Convert(filename)
		conversionTotals[imageType]++
	}

	log.Info("compressing jpegs")
	var compressed int
	if compress {
		for _, filename := range files {
			if _, found := skipMap[filename]; !found {
				imageconvert.CompressJPEG(85, filename)
				_, err = processedLog.WriteString(filename + "\n")
				imageconvert.HandleErr("write to log", err)
				compressed++
			}
		}
	}

	log.WithFields(log.Fields{
		"converted pngs":  conversionTotals["png"],
		"converted webps": conversionTotals["webp"],
		"compressed":      compressed,
	}).Info("Done")
}

// getSkipMap read the log from the last time this was run and
// puts those filenames in a map so we dont have to process them again
// If you want to reprocess, just delete the file
func getSkipMap(processedImages *os.File) map[string]bool {

	var scanner = bufio.NewScanner(processedImages)
	scanner.Split(bufio.ScanLines)
	var compressedFiles = make(map[string]bool)

	for scanner.Scan() {
		compressedFiles[scanner.Text()] = staticBool
	}

	return compressedFiles
}
