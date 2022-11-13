package main

import (
	"fmt"
	"os"
	"path/filepath"

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

// conversionWorker reads from the file chan and does all the conversion work.
// 1. does a file already exist with the output name? yes = skip
// 2. convert it to jpg
// 3. compress it (if enabled)
// 4. reset the mod time
func conversionWorker(files chan string, results chan conversionResult, compress bool) {
	defer close(results)

	for file := range files {
		var result = conversionResult{
			OriginalFileName: file,
		}

		var originalFileStat, err = os.Stat(file)
		if err != nil {
			result.Error = fmt.Errorf("error stat'ing file: %s, error: %w", file, err)
			results <- result
			continue
		}

		// if a file already exists with the output name, we skip it as not to overwrite it
		if imageconvert.WouldOverwrite(file) {
			result.ConvertedFileName = file
			if filepath.Ext(file) != ".jpg" && filepath.Ext(file) != ".jpeg" { // if its a jpg we can still compress it
				log.Warnf("renaming %s will overwrite an existing jpeg, skipping", file)
				continue
			}
		} else {

			// CONVERT IT
			var imageType string
			result.ConvertedFileName, imageType, err = imageconvert.Convert(file)
			if err != nil {
				result.Error = fmt.Errorf("error converting image: %s, error: %w", file, err)
				results <- result
				continue
			}
			result.ImageType = imageType
		}

		// COMPRESS IT
		if compress {
			converted, err := imageconvert.CompressJPEG(85, result.ConvertedFileName)
			if err != nil {
				result.Error = fmt.Errorf("error compressing image: %s, error: %w", file, err)
				results <- result
				continue
			}
			result.Compressed = converted
		}

		// RESET MODTIME
		err = os.Chtimes(result.ConvertedFileName, originalFileStat.ModTime(), originalFileStat.ModTime())
		if err != nil {
			result.Error = fmt.Errorf("could not reset mod time of file: %s, err: %w", result.ConvertedFileName, err)
			results <- result
			continue
		}

		results <- result
	}
}
