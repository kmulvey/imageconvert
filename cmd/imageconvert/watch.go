package main

import (
	"context"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/kmulvey/humantime"
	"github.com/kmulvey/imageconvert/pkg/imageconvert"
	"github.com/kmulvey/path"
)

// TimeOfEntry is a wrapper struct for waitTilFileWritesComplete
type TimeOfEntry struct {
	path.WatchEvent
	time.Time
}

// watchDir watches the given dir in a blocking manner with optional filters. Results are sent back on the files chan.
func watchDir(inputPath string, files chan path.WatchEvent, tr humantime.TimeRange, force bool, processedLog *os.File) {

	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	var nilTime = time.Time{}

	if tr.To == nilTime {
		tr.To = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	}
	var dateFilter = path.NewDateWatchFilter(tr.From, tr.To)
	var skipFilter = path.NewSkipMapWatchFilter(getSkipMap(processedLog))
	var extensionFilter = path.NewRegexWatchFilter(imageconvert.ImageExtensionRegex)
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
