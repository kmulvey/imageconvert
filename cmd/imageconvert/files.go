package main

import (
	"bufio"
	"os"
	"time"

	"github.com/kmulvey/humantime"
	"github.com/kmulvey/imageconvert/pkg/imageconvert"
	"github.com/kmulvey/path"
)

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
func getFileList(inputPath path.Path, tr humantime.TimeRange, force bool, processedLog *os.File) []string {

	var nilTime = time.Time{}
	var trimmedFileList []path.Entry

	switch {
	case force:
		trimmedFileList = inputPath.Files
	case tr.From != nilTime:
		trimmedFileList = path.FilterEntities(inputPath.Files, path.NewDateEntitiesFilter(tr.From, tr.To))
	default:
		trimmedFileList = path.FilterEntities(inputPath.Files, path.NewSkipMapEntitiesFilter(getSkipMap(processedLog)))
	}

	trimmedFileList = path.FilterEntities(trimmedFileList, path.NewRegexEntitiesFilter(imageconvert.ImageExtensionRegex))

	// these are all the files all the way down the dir tree
	return path.OnlyNames(trimmedFileList)
}
