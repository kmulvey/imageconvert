package imageconvert

import (
	"fmt"
	"os"

	"github.com/kmulvey/goutils"
	"github.com/kmulvey/path"
	log "github.com/sirupsen/logrus"
)

func (ic *ImageConverter) Start(fileChan chan string) int {

	var totalFiles int // only relevant for non-watch mode

	if ic.Watch {
		log.Infof("watrching dir: %s", ic.InputPath)
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
		go ic.watchDir(inputPath, watchEvents, tr, force, processedLog)

	} else {

		var files, err = ic.getFileList()
		if err != nil {
			return 0, err
		}
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
func processAndWaitForResults(resultChans []chan conversionResult, conversionTypeTotals map[string]int, totalFiles int, processedLog *os.File) (int, int, int) {

	var compressedTotal, renamedTotal, resizedTotal, fileCount int
	for result := range goutils.MergeChannels(resultChans...) {

		fileCount++

		if result.Error != nil {
			log.Error(result.Error)
		} else {

			conversionTypeTotals[result.ImageType]++

			if result.Compressed {
				compressedTotal++
			}

			if result.Renamed {
				renamedTotal++
			}

			if result.Resized {
				resizedTotal++
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
				"resized":            result.Resized,
			}
			if result.Compressed {
				fields["compressed output"] = result.CompressOutput
			}
			log.WithFields(fields).Info("Converted")
		}
	}

	return compressedTotal, renamedTotal, resizedTotal
}
