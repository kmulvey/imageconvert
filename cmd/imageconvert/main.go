package main

import (
	"bufio"
	"flag"
	"os"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/kmulvey/imageconvert/pkg/imageconvert"
)

const staticBool = false

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// get the user options
	var rootDir string
	var processedLogFile string
	var compress bool
	var threads int

	flag.StringVar(&rootDir, "dir", "", "directory (abs path), could also be a single file")
	flag.StringVar(&processedLogFile, "log-file", "processed.log", "the file to write processes images to, so that we dont processes them again next time")
	flag.BoolVar(&compress, "compress", false, "compress")
	flag.IntVar(&threads, "threads", 1, "number of threads to use, >1 only useful when rebuilding the cache")
	flag.Parse()
	if strings.TrimSpace(rootDir) == "" {
		log.Fatal("directory not provided")
	}
	if threads <= 0 || threads > runtime.GOMAXPROCS(0) {
		threads = 1
	}

	// open the file
	log.Info("reading log file")
	var processedLog, err = os.OpenFile(processedLogFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatalf("processedLog open, error: %s", err.Error())
	}
	defer func() {
		err = processedLog.Close()
		if err != nil {
			log.Fatalf("processedLog close: error: %s", err.Error())
		}
	}()

	var skipMap = getSkipMap(processedLog)

	// Did they give us a dir or file?
	log.Info("building file list")
	fileInfo, err := os.Stat(rootDir)
	if err != nil {
		log.Fatal("could not stat file/dir ", err)
	}
	var files = make([]string, 1) // we will always have at least one
	if !fileInfo.IsDir() {
		files[0] = rootDir
	} else {
		// these are all the files all the way down the dir tree
		fileInfos, err := imageconvert.ListFiles(rootDir)
		if err != nil {
			log.Fatalf("error listing files: dir: %s, error: %s", rootDir, err.Error())
		}

		fileInfos = imageconvert.FilterFilesBySkipMap(fileInfos, skipMap)
		files = imageconvert.FileInfoToString(fileInfos)
	}

	// spin up goroutines to do the work
	log.Info("spinning up ", threads, " workers")
	var conversionTotals = make(map[string]int)
	var compressedTotal int
	var renamedTotal int
	var fileChan = make(chan string)
	var resultChans = make([]chan conversionResult, threads)
	for i := 0; i < threads; i++ {
		var results = make(chan conversionResult)
		resultChans[i] = results
		go conversionWorker(fileChan, results, compress)
	}

	log.Info("beginning ", len(files), " conversions")
	go func() {
		for _, file := range files {
			fileChan <- file
		}
		close(fileChan)
	}()

	// process results of our goroutines, every error is fatal
	log.Info("waiting for workers to complete")
	for result := range mergeResults(resultChans...) {
		if result.Error != nil {
			log.Fatal(result.Error)
		} else {
			conversionTotals[result.ImageType]++
			if result.Compressed {
				compressedTotal++
			}
			if result.Renamed {
				renamedTotal++
			}
			_, err = processedLog.WriteString(result.ConvertedFileName + "\n")
			if err != nil {
				log.Fatalf("error writing to log file, error: %s", err.Error())
			}

			log.WithFields(log.Fields{
				"original file name": result.OriginalFileName,
				"new file name":      result.ConvertedFileName,
				"type":               result.ImageType,
				"compressed":         result.Compressed,
				"renamed":            result.Renamed,
			}).Info("Converted")
		}
	}

	log.WithFields(log.Fields{
		"converted pngs":  conversionTotals["png"],
		"converted webps": conversionTotals["webp"],
		"compressed":      compressedTotal,
		"jpegs renamed":   renamedTotal,
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
