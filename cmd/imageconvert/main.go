package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"time"

	"github.com/kmulvey/humantime"
	"github.com/kmulvey/path"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// get the user options
	var inputPath path.Path
	var processedLogFile string
	var compress bool
	var threads int
	var tr humantime.TimeRange

	flag.Var(&inputPath, "path", "path to files, globbing must be quoted")
	flag.StringVar(&processedLogFile, "log-file", "processed.log", "the file to write processes images to, so that we dont processes them again next time")
	flag.BoolVar(&compress, "compress", false, "compress")
	flag.IntVar(&threads, "threads", 1, "number of threads to use")
	flag.Var(&tr, "modified-since", "process files chnaged since this time")
	flag.Parse()
	if len(inputPath.Files) == 0 {
		log.Error("path not provided")
		flag.PrintDefaults()
		return
	}
	if threads <= 0 || threads > runtime.GOMAXPROCS(0) {
		threads = 1
		log.Infof("invalid thread count: %d, setting threads to 1", threads)
	}
	log.Infof("Config: dir: %s, log file: %s, compress: %t, threads: %d, modified-since: %s", inputPath.Input, processedLogFile, compress, threads, tr)

	log.Info("reading processed log file")
	var processedLog, err = os.OpenFile(processedLogFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatalf("processedLog open, error: %s", err.Error())
	}

	log.Info("building file list")
	files, err := getFileList(inputPath, tr, processedLog)
	if err != nil {
		log.Fatal(err)
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
			log.Error(result.Error)
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

	err = processedLog.Close()
	if err != nil {
		log.Fatalf("error closing log file: %s", err.Error())
	}
}

// getSkipMap read the log from the last time this was run and
// puts those filenames in a map so we dont have to process them again
// If you want to reprocess, just delete the file
func getSkipMap(processedImages *os.File) map[string]struct{} {

	var scanner = bufio.NewScanner(processedImages)
	scanner.Split(bufio.ScanLines)
	var compressedFiles = make(map[string]struct{})

	for scanner.Scan() {
		compressedFiles[scanner.Text()] = struct{}{}
	}

	return compressedFiles
}

// getFileList filters the file list
func getFileList(inputPath path.Path, modSince humantime.TimeRange, processedLog *os.File) ([]string, error) {

	var nilTime = time.Time{}
	var err error
	var trimmedFileList []path.File

	if modSince.From != nilTime {
		trimmedFileList, err = path.FilterFilesByDateRange(inputPath.Files, modSince.From, modSince.To)
		if err != nil {
			return nil, fmt.Errorf("unable to filter files by skip map")
		}
	} else {
		trimmedFileList = path.FilterFilesBySkipMap(inputPath.Files, getSkipMap(processedLog))
	}

	trimmedFileList = path.FilterFilesByRegex(trimmedFileList, regexp.MustCompile(".*.jpg$|.*.jpeg$|.*.png$|.*.webp$"))

	// these are all the files all the way down the dir tree
	return path.DirEntryToString(trimmedFileList), nil
}
