package main

import (
	"bufio"
	"flag"
	"os"
	"runtime"
	"time"

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
	var inputPath path.Path
	var processedLogFile string
	var compress bool
	var force bool
	var threads int
	var v bool
	var h bool
	var tr humantime.TimeRange

	flag.Var(&inputPath, "path", "path to files, globbing must be quoted")
	flag.StringVar(&processedLogFile, "log-file", "processed.log", "the file to write processes images to, so that we dont processes them again next time")
	flag.BoolVar(&compress, "compress", false, "compress")
	flag.BoolVar(&force, "force", false, "force")
	flag.IntVar(&threads, "threads", 1, "number of threads to use")
	flag.Var(&tr, "modified-since", "process files chnaged since this time")
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

	if len(inputPath.Files) == 0 {
		log.Error("input path does not have any files")
		flag.PrintDefaults()
		return
	}
	if threads <= 0 || threads > runtime.GOMAXPROCS(0) {
		threads = 1
		log.Infof("invalid thread count: %d, setting threads to 1", threads)
	}
	log.Infof("Config: dir: %s, log file: %s, compress: %t, threads: %d, modified-since: %s", inputPath.ComputedPath.AbsolutePath, processedLogFile, compress, threads, tr)

	log.Info("reading processed log file")
	var processedLog, err = os.OpenFile(processedLogFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatalf("processedLog open, error: %s", err.Error())
	}

	log.Info("building file list")
	var files = getFileList(inputPath, tr, force, processedLog)

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

	// process results of our goroutines
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
func getFileList(inputPath path.Path, modSince humantime.TimeRange, force bool, processedLog *os.File) []string {

	var nilTime = time.Time{}
	var trimmedFileList []path.Entry

	switch {
	case force:
		trimmedFileList = inputPath.Files
	case modSince.From != nilTime:
		trimmedFileList = path.FilterEntities(inputPath.Files, path.NewDateEntitiesFilter(modSince.From, modSince.To))
	default:
		trimmedFileList = path.FilterEntities(inputPath.Files, path.NewSkipMapEntitiesFilter(getSkipMap(processedLog)))
	}

	trimmedFileList = path.FilterEntities(trimmedFileList, path.NewRegexEntitiesFilter(imageconvert.ImageExtensionRegex))

	// these are all the files all the way down the dir tree
	return path.OnlyNames(trimmedFileList)
}
