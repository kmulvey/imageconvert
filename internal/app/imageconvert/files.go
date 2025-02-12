package imageconvert

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/kmulvey/path"
	log "github.com/sirupsen/logrus"
)

// ImageExtensionRegex captures file extensions we can work with.
var ImageExtensionRegex = regexp.MustCompile(".*.jpg$|.*.jpeg$|.*.png$|.*.webp$|.*.JPG$|.*.JPEG$|.*.PNG$|.*.WEBP$")

// RenameExtensionRegex captures file extensions that we would like to rename to the extensions above.
var RenameSuffixRegex = regexp.MustCompile(".*.jpeg$|.*.png$|.*.webp$|.*.JPG$|.*.JPEG$|.*.PNG$|.*.WEBP$")

// NilTime is 0
var NilTime = time.Time{}

// ParseSkipMap reads the log from the last time this was run and
// puts those filenames in a map so we dont have to process them again
// If you want to reprocess, just delete the file.
func (ic *ImageConverter) ParseSkipMap() (map[string]struct{}, error) {

	var processedImages, err = os.OpenFile(ic.SkipMapEntry.String(), os.O_RDONLY, 0755)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		// if the file doesnt exist its not really an error so we just return an empty map
		return make(map[string]struct{}), nil
	} else if err != nil {
		return nil, fmt.Errorf("unable to open skipMap file: %w", err)
	}

	var scanner = bufio.NewScanner(processedImages)
	scanner.Split(bufio.ScanLines)
	var compressedFiles = make(map[string]struct{})

	for scanner.Scan() {
		compressedFiles[scanner.Text()] = struct{}{}
	}

	if err := processedImages.Close(); err != nil {
		return nil, fmt.Errorf("unable to close processed images file: %s, err: %w", ic.SkipMapEntry.String(), err)
	}

	return compressedFiles, nil
}

// getFileList filters the file list
func (ic *ImageConverter) getFileList() ([]path.Entry, error) {

	if ic.TimeRange.To == NilTime {
		ic.TimeRange.To = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	}
	var dateFilter = path.NewDateEntitiesFilter(ic.TimeRange.From, ic.TimeRange.To)

	skipMap, err := ic.ParseSkipMap()
	if err != nil {
		return nil, err
	}
	var skipFilter = path.NewSkipMapEntitiesFilter(skipMap)

	var extensionFilter = path.NewRegexEntitiesFilter(ImageExtensionRegex)

	flattenedFileList, err := ic.InputEntry.Flatten(true)
	if err != nil {
		return nil, fmt.Errorf("error flatteneing entries: %w", err)
	}

	if !ic.Force {
		flattenedFileList = path.FilterEntities(flattenedFileList, dateFilter, skipFilter, extensionFilter)
	}

	// these are all the files all the way down the dir tree
	return flattenedFileList, nil
}

// TimeOfEntry is a wrapper struct for waitTilFileWritesComplete
type TimeOfEntry struct {
	path.WatchEvent
	time.Time
}

// watchDir watches the given dir in a blocking manner with optional filters. Results are sent back on the files chan.
func (ic *ImageConverter) watchDir(files chan path.WatchEvent) {

	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	if ic.TimeRange.To == NilTime {
		ic.TimeRange.To = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	}
	var dateFilter = path.NewDateWatchFilter(ic.TimeRange.From, ic.TimeRange.To)
	var skipFilter = path.NewSkipMapWatchFilter(ic.SkipMap)
	var extensionFilter = path.NewRegexWatchFilter(ImageExtensionRegex)
	var errors = make(chan error)

	go func() {
		for range errors {
			// we dont really care for now but must drain this chan
		}
	}()

	if ic.Force {
		path.WatchDir(ctx, ic.InputEntry.String(), 1, false, files, errors, extensionFilter, path.NewOpWatchFilter(fsnotify.Create))
	} else {
		path.WatchDir(ctx, ic.InputEntry.String(), 1, false, files, errors, dateFilter, skipFilter, extensionFilter, path.NewOpWatchFilter(fsnotify.Create))
	}
}

// waitTilFileWritesComplete is a way to debounce fs events because fsnotify does not send us an event when the file is
// closed after writing so we just get a CREATE and a lot of WRITES. We gather the WRITE events and wait 200ms to see if
// they finish for a given file and if they do we send the event on the eventsOut chan.
func waitTilFileWritesComplete(eventsIn, eventsOut chan path.WatchEvent) {

	var cache = make(map[string]TimeOfEntry)
	var ticker = time.NewTicker(time.Second)

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
				if hasEOI(filename) {
					eventsOut <- entry.WatchEvent
					delete(cache, filename)
				}
			}
		}
	}
}

// hasEOI looks for the End of Image marker at the end of the file. Im not crazy about logging these errors
// but also kinda dont want to bubble them up either. Something to reconsider in the future.
func hasEOI(filepath string) bool {

	var file, err = os.OpenFile(filepath, os.O_RDONLY, 0755)
	if err != nil {
		log.Errorf("error opening file: %s", err)
		return false
	}
	defer file.Close()

	buf := make([]byte, 2)

	stat, err := os.Stat(filepath)
	if err != nil {
		log.Errorf("error opening file: %s", err)
		return false
	}

	start := stat.Size() - 2

	_, err = file.ReadAt(buf, start)
	if err != nil {
		log.Errorf("error opening file: %s", err)
		return false
	}

	if buf[0] == 0xFF && buf[1] == 0xD9 {
		return true
	}

	return false
}
