package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/kmulvey/imageconvert/pkg/imageconvert"
)

var incorrectSuffixRegex = regexp.MustCompile(".*.jpeg$|.*.png$|.*.webp$|.*.JPG$|.*.JPEG$|.*.PNG$|.*.WEBP$")

type conversionResult struct {
	OriginalFileName  string
	ConvertedFileName string
	ImageType         string
	CompressOutput    string
	Error             error
	Compressed        bool
	Renamed           bool
}

// conversionWorker reads from the file chan and does all the conversion work.
func conversionWorker(files chan string, results chan conversionResult, compress bool) {
	defer close(results)

	for file := range files {
		results <- convertImage(file, compress)
	}
}

// convertImage  converts, compresses and renames images.
// This is broken out from the conversionWorker for ease of testing.
// 1. does a file already exist with the output name? yes = skip
// 2. convert it to jpg
// 3. compress it (if enabled)
// 4. reset the mod time
func convertImage(file string, compress bool) conversionResult {

	var result = conversionResult{
		OriginalFileName: file,
	}

	var originalFileStat, err = os.Stat(file)
	if err != nil {
		result.Error = fmt.Errorf("error stat'ing file: %s, error: %w", file, err)
		return result
	}

	// CONVERT IT
	var imageType string
	result.ConvertedFileName, imageType, err = imageconvert.Convert(file)
	if err != nil {
		result.Error = fmt.Errorf("error converting image: %s, error: %w", file, err)
		return result
	}
	result.ImageType = imageType

	// COMPRESS IT
	if compress {
		compressed, stdout, err := imageconvert.CompressJPEG(85, result.ConvertedFileName)
		if err != nil {
			result.Error = fmt.Errorf("error compressing image: %s, error: %w", file, err)
			return result
		}
		result.Compressed = compressed

		if compressed {
			result.CompressOutput = stdout
		}
	}

	// make sure every file ends in ".jpg".
	if incorrectSuffixRegex.MatchString(filepath.Base(result.ConvertedFileName)) {
		var renamedFile = strings.Replace(result.ConvertedFileName, filepath.Ext(result.ConvertedFileName), ".jpg", 1)
		if err := os.Rename(result.ConvertedFileName, renamedFile); err != nil {
			result.Error = fmt.Errorf("could rename file: %s, err: %w", result.ConvertedFileName, err)
			return result
		}
		result.ConvertedFileName = renamedFile
		result.Renamed = true
	}

	// RESET MODTIME
	err = os.Chtimes(result.ConvertedFileName, originalFileStat.ModTime(), originalFileStat.ModTime())
	if err != nil {
		result.Error = fmt.Errorf("could not reset mod time of file: %s, err: %w", result.ConvertedFileName, err)
		return result
	}
	return result
}
