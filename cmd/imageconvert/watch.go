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

var Close uint32 = 6

type TimeOfEntry struct {
	path.Entry
	time.Time
}

// getFileList filters the file list
func watchDir(inputPath path.Path, files chan path.WatchEvent, tr humantime.TimeRange, force bool, processedLog *os.File) {

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
			// log.Errorf("error from WatchDir: %s", err)
		}
	}()

	if force {
		path.WatchDir(ctx, inputPath.ComputedPath.AbsolutePath, true, files, errors, extensionFilter, path.NewOpWatchFilter(fsnotify.Create))
	} else {
		path.WatchDir(ctx, inputPath.ComputedPath.AbsolutePath, true, files, errors, dateFilter, skipFilter, extensionFilter, path.NewOpWatchFilter(fsnotify.Create))
	}
}

func waitTilFileWritesComplete(eventsIn, eventsOut chan path.WatchEvent) {

	var cache = make(map[string]TimeOfEntry)
	var ticker = time.NewTicker(500 * time.Millisecond)

	for {
		select {
		case event, open := <-eventsIn:
			if !open {
				close(eventsOut)
				return
			}
			cache[event.Entry.AbsolutePath] = TimeOfEntry{Entry: event.Entry, Time: time.Now()}

		case <-ticker.C:
			for filename, entry := range cache {
				if time.Since(entry.Time) > 200*time.Millisecond {
					eventsOut <- path.WatchEvent{Entry: entry.Entry, Op: 6}
					delete(cache, filename)
				}
			}
		}
	}
}
