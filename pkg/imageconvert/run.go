package imageconvert

import (
	"fmt"
	"os"

	"github.com/kmulvey/goutils"
	"github.com/kmulvey/path"
	log "github.com/sirupsen/logrus"
)

func (ic *ImageConverter) Start(results chan ConversionResult) (int, int, int, int, map[string]int, error) {

	var resultChans = make([]chan ConversionResult, ic.Threads)
	var i uint8
	for i = 0; i < ic.Threads; i++ {
		resultChans[i] = make(chan ConversionResult)
	}

	// variables only for slice mode, these variables are returned so totals can be printed by the caller
	var compressedTotal, renamedTotal, resizedTotal int
	var conversionTypeTotals = make(map[string]int)
	var processAndWaitErrors = make(chan error)
	var processedLogHanlde *os.File
	var err error

	// handle how results are processed, give them to the caller or log them
	if results != nil {
		for result := range goutils.MergeChannels(resultChans...) {
			results <- result
		}
	} else {
		processedLogHanlde, err = os.OpenFile(ic.SkipMapEntry.String(), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
		if err != nil {
			return 0, 0, 0, 0, nil, fmt.Errorf("unable to open processed log file: %s, err: %w", ic.SkipMapEntry.String(), err)
		}

		go func() {
			compressedTotal, renamedTotal, resizedTotal, err = processAndWaitForResults(resultChans, conversionTypeTotals, len(ic.InputFiles), processedLogHanlde)
			if err != nil {
				processAndWaitErrors <- err
			}
			close(processAndWaitErrors)
		}()
	}

	// and away we go
	if ic.Watch {
		ic.startWatch(resultChans...)
	} else {
		go ic.startSlice(resultChans...)
	}

	if err := <-processAndWaitErrors; err != nil {
		return 0, 0, 0, 0, nil, err
	}

	return compressedTotal, renamedTotal, resizedTotal, len(ic.InputFiles), conversionTypeTotals, processedLogHanlde.Close()
}

// processAndWaitForResults reads the results chan which is a stream of files that have been converted by the worker.
// The files name is added to the processed log to prevent further processing in the future as well as tallying up some stats for logging.
func processAndWaitForResults(resultChans []chan ConversionResult, conversionTypeTotals map[string]int, totalFiles int, processedLog *os.File) (int, int, int, error) {

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
				return 0, 0, 0, fmt.Errorf("error writing to log file, error: %w", err)
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

	return compressedTotal, renamedTotal, resizedTotal, nil
}

// startWatch flow:
// ic.watchDir() watches for fs events and writes the to the watchEvents chan
// waitTilFileWritesComplete() debounces watchEvents and writes them to watchEventsDebounced
// we read and dedup files from watchEventsDebounced and writes them to conversionChan
// conversionWorker reads conversionWorker, converts the images and writes the results to the results chan
// if no resultChans are passed in then processAndWaitForResults is used to read the results and log them
func (ic *ImageConverter) startWatch(resultChans ...chan ConversionResult) {

	var conversionChan = make(chan path.Entry)
	var i uint8
	for i = 0; i < ic.Threads; i++ {
		go ic.conversionWorker(conversionChan, resultChans[i])
	}

	log.Infof("watrching dir: %s", ic.InputEntry.String())
	var watchEvents = make(chan path.WatchEvent)
	var watchEventsDebounced = make(chan path.WatchEvent)

	// we need to debounce file writes because fsnotify does not tell us when the file has finished being written to (closed)
	// so if we dont debounce the WRITE events we will try to open the file before its ready.
	go waitTilFileWritesComplete(watchEvents, watchEventsDebounced)

	// dedup files seen
	go func() {
		var seenFiles = make(map[string]struct{})
		for file := range watchEventsDebounced {
			if _, found := seenFiles[file.Entry.AbsolutePath]; !found {
				conversionChan <- file.Entry
				seenFiles[file.Entry.AbsolutePath] = struct{}{}
			}
		}
		close(conversionChan)
	}()

	// watch the dir for events
	go ic.watchDir(watchEvents)
}

// startSlice flow:
// we loop over the ic.InputFiles and write them to conversionChan
// conversionWorker reads conversionWorker, converts the images and writes the results to the results chan
// if no resultChans are passed in then processAndWaitForResults is used to read the results and log them
func (ic *ImageConverter) startSlice(resultChans ...chan ConversionResult) {

	var conversionChan = make(chan path.Entry)
	var i uint8
	for i = 0; i < ic.Threads; i++ {
		go ic.conversionWorker(conversionChan, resultChans[i])
	}

	log.Info("beginning ", len(ic.InputFiles), " conversions")

	for _, file := range ic.InputFiles {
		conversionChan <- file
	}

	close(conversionChan)
}
