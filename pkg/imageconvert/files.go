package imageconvert

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

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
func (ic *ImageConverter) getFileList(inputFiles []path.Entry) ([]string, error) {

	var nilTime = time.Time{}
	var trimmedFileList, err = path.List(ic.InputPath, 2, false)
	if err != nil {
		return nil, fmt.Errorf("error getting file list for: %s, err: %w", ic.InputPath, err)
	}

	if ic.Force {
		goto BypassFilters
	}

	if ic.TimeRange.From != nilTime {
		trimmedFileList, err = path.List(ic.InputPath, 2, false, path.NewDateEntitiesFilter(ic.TimeRange.From, ic.TimeRange.To))
		if err != nil {
			return nil, fmt.Errorf("error getting file list for: %s, err: %w", ic.InputPath, err)
		}
	}

	if ic.SkipMapFile != "" {
		skipMap, err := ic.ParseSkipMap()
		if err != nil {
			return nil, err
		}

		trimmedFileList, err = path.List(ic.InputPath, 2, false, path.NewSkipMapEntitiesFilter(skipMap))
		if err != nil {
			return nil, fmt.Errorf("error getting file list for: %s, err: %w", ic.InputPath, err)
		}
	}

BypassFilters:

	trimmedFileList = path.FilterEntities(trimmedFileList, path.NewRegexEntitiesFilter(ImageExtensionRegex))

	// these are all the files all the way down the dir tree
	return path.OnlyNames(trimmedFileList), nil
}
