package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kmulvey/imageconvert/pkg/imageconvert"
	log "github.com/sirupsen/logrus"
)

type conversionResult struct {
	OriginalFileName  string
	ConvertedFileName string
	ImageType         string
	Error             error
	Compressed        bool
	Renamed           bool
}

func conversionWorker(files chan string, results chan conversionResult, compress bool) {

	for file := range files {
		var result = conversionResult{
			OriginalFileName: file,
		}

		stat, err := os.Stat(file)
		if err != nil {
			result.Error = fmt.Errorf("error stat'ing file: %s, error: %w", file, err)
			results <- result
			continue
		}

		convertedFileName, imageType, err := imageconvert.Convert(file)
		if err != nil {
			result.Error = fmt.Errorf("error converting image: %s, error: %w", file, err)
			results <- result
			continue
		}
		result.ImageType = imageType
		result.ConvertedFileName = convertedFileName

		if compress {
			converted, err := imageconvert.CompressJPEG(85, convertedFileName)
			if err != nil {
				result.Error = fmt.Errorf("error compressing image: %s, error: %w", file, err)
				results <- result
				continue
			}
			result.Compressed = converted
		}

		var fileExt = strings.ReplaceAll(filepath.Ext(result.ConvertedFileName), ".", "")
		if strings.HasSuffix(convertedFileName, ".jpeg") {
			renamed, err := rename(convertedFileName, ".jpeg", ".jpg")
			if err != nil {
				results <- result
				continue
			}
			if renamed != "" {
				result.Renamed = true
				result.ConvertedFileName = renamed
			}
		} else if result.ImageType != fileExt && fileExt != "jpg" {
			renamed, err := rename(convertedFileName, filepath.Ext(result.ConvertedFileName), ".jpg")
			if err != nil {
				results <- result
				continue
			}
			if renamed != "" {
				result.Renamed = true
				result.ConvertedFileName = renamed
			}
		}

		// reset modtime
		err = os.Chtimes(result.ConvertedFileName, stat.ModTime(), stat.ModTime())
		if err != nil {
			result.Error = fmt.Errorf("could reset mod time of file: %s, err: %w", result.ConvertedFileName, err)
			results <- result
			continue
		}

		results <- result
	}
	close(results)
}

func rename(convertedFileName, from, to string) (string, error) {
	var renamed = strings.ReplaceAll(convertedFileName, from, to)
	if imageconvert.WouldOverwrite(convertedFileName, renamed) {
		log.Warnf("renaming %s would overwrite an existing jpeg, skipping", convertedFileName)
	}

	var err = os.Rename(convertedFileName, renamed)
	if err != nil {
		return "", fmt.Errorf("could rename file: %s, err: %w", convertedFileName, err)
	}

	return renamed, nil
}
