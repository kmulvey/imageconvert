package imageconvert

import (
	"bufio"
	"os"
	"regexp"
	"time"

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

	var scanner = bufio.NewScanner(ic.SkipMapFileHandle)
	scanner.Split(bufio.ScanLines)
	var compressedFiles = make(map[string]struct{})

	for scanner.Scan() {
		compressedFiles[scanner.Text()] = struct{}{}
	}

	return compressedFiles, nil
}

// TimeOfEntry is a wrapper struct for waitTilFileWritesComplete
type TimeOfEntry struct {
	path.WatchEvent
	time.Time
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
