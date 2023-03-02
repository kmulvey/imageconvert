package imageconvert

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/kmulvey/humantime"
	"github.com/kmulvey/path"
)

// ImageExtensionRegex captures file extensions we can work with.
var ImageExtensionRegex = regexp.MustCompile(".*.jpg$|.*.jpeg$|.*.png$|.*.webp$|.*.JPG$|.*.JPEG$|.*.PNG$|.*.WEBP$")

// RenameExtensionRegex captures file extensions that we would like to rename to the extensions above.
var RenameSuffixRegex = regexp.MustCompile(".*.jpeg$|.*.png$|.*.webp$|.*.JPG$|.*.JPEG$|.*.PNG$|.*.WEBP$")

// EscapeFilePath escapes spaces in the filepath used for an exec() call.
func EscapeFilePath(file string) string {
	var r = strings.NewReplacer(" ", `\ `, "(", `\(`, ")", `\)`, "'", `\'`, "&", `\&`, "@", `\@`)
	return r.Replace(file)
}

var nilTime = time.Time{}

// ParseSkipMap read the log from the last time this was run and
// puts those filenames in a map so we dont have to process them again
// If you want to reprocess, just delete the file
func (ic *ImageConverter) ParseSkipMap() (map[string]struct{}, error) {

	var processedImages, err = os.Open(ic.SkipMapFile)
	if err != nil {
		return nil, fmt.Errorf("unable to open skipMap file: %w", err)
	}

	var scanner = bufio.NewScanner(processedImages)
	scanner.Split(bufio.ScanLines)
	var compressedFiles = make(map[string]struct{})

	for scanner.Scan() {
		compressedFiles[scanner.Text()] = struct{}{}
	}

	return compressedFiles, nil
}

// getFileList filters the file list
func (ic *ImageConverter) getFileList() ([]string, error) {

	var trimmedFileList []path.Entry

	if ic.TimeRange.To == nilTime {
		ic.TimeRange.To = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	}
	var dateFilter = path.NewDateEntitiesFilter(ic.TimeRange.From, ic.TimeRange.To)

	skipMap, err := ic.ParseSkipMap()
	if err != nil {
		return nil, err
	}
	var skipFilter = path.NewSkipMapEntitiesFilter(skipMap)

	var extensionFilter = path.NewRegexEntitiesFilter(ImageExtensionRegex)

	if ic.Force {
		trimmedFileList, err = path.List(ic.InputPath, 2, false)
		if err != nil {
			return nil, fmt.Errorf("error getting file list for: %s, err: %w", ic.InputPath, err)
		}

	} else {
		trimmedFileList, err = path.List(ic.InputPath, 2, false, dateFilter, skipFilter, extensionFilter)
		if err != nil {
			return nil, fmt.Errorf("error getting file list for: %s, err: %w", ic.InputPath, err)
		}
	}

	// these are all the files all the way down the dir tree
	return path.OnlyNames(trimmedFileList), nil
}

// TimeOfEntry is a wrapper struct for waitTilFileWritesComplete
type TimeOfEntry struct {
	path.WatchEvent
	time.Time
}

// watchDir watches the given dir in a blocking manner with optional filters. Results are sent back on the files chan.
func (ic *ImageConverter) watchDir(inputPath string, files chan path.WatchEvent, tr humantime.TimeRange, force bool, processedLog *os.File) error {

	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	if tr.To == nilTime {
		tr.To = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	}
	var dateFilter = path.NewDateWatchFilter(ic.TimeRange.From, ic.TimeRange.To)
	skipMap, err := ic.ParseSkipMap()
	if err != nil {
		close(files)
		return err
	}
	var skipFilter = path.NewSkipMapWatchFilter(skipMap)
	var extensionFilter = path.NewRegexWatchFilter(ImageExtensionRegex)
	var errors = make(chan error)

	go func() {
		for range errors {
			// we dont really care for now but must drain this chan
		}
	}()

	if force {
		path.WatchDir(ctx, inputPath, 1, false, files, errors, extensionFilter, path.NewOpWatchFilter(fsnotify.Create))
	} else {
		path.WatchDir(ctx, inputPath, 1, false, files, errors, dateFilter, skipFilter, extensionFilter, path.NewOpWatchFilter(fsnotify.Create))
	}

	return nil
}

// waitTilFileWritesComplete is a way to debounce fs events because fsnotify does not send us an event when the file is
// closed after writing so we just get a CREATE and a lot of WRITES. We gather the WRITE events and wait 200ms to see if
// they finish for a given file and if they do we send the event on the eventsOut chan.
func waitTilFileWritesComplete(eventsIn, eventsOut chan path.WatchEvent) {

	var cache = make(map[string]TimeOfEntry)
	var ticker = time.NewTicker(200 * time.Millisecond)

	for {
		select {
		case event, open := <-eventsIn:
			if !open {
				close(eventsOut)
				return
			}

			cache[event.Entry.AbsolutePath] = TimeOfEntry{WatchEvent: event, Time: time.Now()}

		case <-ticker.C:

			for filename, entry := range cache {
				if time.Since(entry.Time) > 6*time.Second { // this is a long time because large files take a while to get written to spinning rust
					eventsOut <- entry.WatchEvent
					delete(cache, filename)
				}
			}
		}
	}
}
