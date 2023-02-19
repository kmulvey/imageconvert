package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/kmulvey/humantime"
	"github.com/kmulvey/imageconvert/pkg/imageconvert"
	"github.com/kmulvey/path"
	log "github.com/sirupsen/logrus"
	"go.szostok.io/version"
	"go.szostok.io/version/printer"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// get the user options
	var inputPath string
	var processedLogFile string
	var compress, force, watch, v, h bool
	var threads int
	var tr humantime.TimeRange

	flag.StringVar(&inputPath, "path", "", "path to files, globbing must be quoted")
	flag.StringVar(&processedLogFile, "log-file", "processed.log", "the file to write processes images to, so that we dont processes them again next time")
	flag.BoolVar(&compress, "compress", true, "compress")
	flag.BoolVar(&force, "force", false, "force")
	flag.BoolVar(&watch, "watch", false, "watch the dir")
	flag.IntVar(&threads, "threads", 1, "number of threads to use")
	flag.Var(&tr, "time-range", "process files chnaged since this time")
	flag.BoolVar(&v, "version", false, "print version")
	flag.BoolVar(&v, "v", false, "print version")
	flag.BoolVar(&h, "help", false, "print options")
	flag.Parse()

	if h {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if v {
		var verPrinter = printer.New()
		var info = version.Get()
		if err := verPrinter.PrintInfo(os.Stdout, info); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	var files, err = path.List(inputPath, 2, path.NewRegexEntitiesFilter(imageconvert.ImageExtensionRegex))
	if err != nil {
		log.Fatal(err)
	}

	if len(files) == 0 && !watch {
		log.Error(" input path does not have any files")
		return
	}

	if threads <= 0 || threads > runtime.GOMAXPROCS(0) {
		threads = 1
		log.Infof("invalid thread count: %d, setting threads to 1", threads)
	}
	log.Infof("Config: dir: %s, log file: %s, compress: %t, force: %t, watch: %t, threads: %d, modified-since: %s", inputPath, processedLogFile, compress, force, watch, threads, tr)

	processedLog, err := os.OpenFile(processedLogFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatalf("processedLog open, error: %s", err.Error())
	}

	// spin up goroutines to do the work
	var fileChan = make(chan string)
	var conversionTotals = make(map[string]int)
	var compressedTotal, renamedTotal, totalFiles int // totalFiles is only for non-watch
	var resultChans = make([]chan conversionResult, threads)
	for i := 0; i < threads; i++ {
		var results = make(chan conversionResult)
		resultChans[i] = results
		go conversionWorker(fileChan, results, compress)
	}

	// start er up
	totalFiles = start(watch, force, inputPath, files, tr, fileChan, processedLog)

	// wait for and process results from our worker goroutines
	var fileCount = processAndWaitForResults(resultChans, conversionTotals, compressedTotal, renamedTotal, totalFiles, processedLog)

	log.WithFields(log.Fields{
		"converted pngs":   conversionTotals["png"],
		"converted webps":  conversionTotals["webp"],
		"compressed":       compressedTotal,
		"jpegs renamed":    renamedTotal,
		"total files seen": fileCount,
	}).Info("Done")

	err = processedLog.Close()
	if err != nil {
		log.Fatalf("error closing log file: %s", err.Error())
	}
}

func start(watch, force bool, inputPath string, inputFiles []path.Entry, tr humantime.TimeRange, fileChan chan string, processedLog *os.File) int {

	var totalFiles int // only relevant for non-watch mode

	if watch {
		log.Infof("watrching dir: %s", inputPath)
		var watchEvents = make(chan path.WatchEvent)
		var watchEventsDebounced = make(chan path.WatchEvent)

		// we need to debounce file writes because fsnotify does not tell us when the file has finished being written to (closed)
		// so if we dont debounce the WRITE events we will try to open the file before its ready.
		go waitTilFileWritesComplete(watchEvents, watchEventsDebounced)

		go func() {
			var seenFiles = make(map[string]struct{})
			for file := range watchEventsDebounced {
				if _, found := seenFiles[file.Entry.AbsolutePath]; !found {
					fileChan <- file.Entry.AbsolutePath
					seenFiles[file.Entry.AbsolutePath] = struct{}{}
				}
			}
			close(fileChan)
		}()
		go watchDir(inputPath, watchEvents, tr, force, processedLog)

	} else {

		var files = getFileList(inputFiles, tr, force, processedLog)
		totalFiles = len(files)
		log.Info("beginning ", len(files), " conversions")

		go func() {
			for _, file := range files {
				fileChan <- file
			}
			close(fileChan)
		}()
	}

	return totalFiles
}

// processAndWaitForResults reads the results chan which is a stream of files that have been converted by the worker.
// The files name is added to the processed log to prevent further processing in the future as well as tallying up some stats for logging.
func processAndWaitForResults(resultChans []chan conversionResult, conversionTotals map[string]int, compressedTotal, renamedTotal, totalFiles int, processedLog *os.File) int {

	var fileCount int
	for result := range mergeResults(resultChans...) {

		fileCount++

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
			var _, err = processedLog.WriteString(result.ConvertedFileName + "\n")
			if err != nil {
				log.Fatalf("error writing to log file, error: %s", err.Error())
			}

			var fields = log.Fields{
				"original file name": result.OriginalFileName,
				"new file name":      result.ConvertedFileName,
				"type":               result.ImageType,
				"compressed":         result.Compressed,
				"progeress":          fmt.Sprintf("[%d/%d]", fileCount, totalFiles),
				"renamed":            result.Renamed,
			}
			if result.Compressed {
				fields["compressed output"] = result.CompressOutput
			}
			log.WithFields(fields).Info("Converted")
		}
	}
	return fileCount
}
