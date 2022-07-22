package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/kmulvey/imageconvert/pkg/imageconvert"
	log "github.com/sirupsen/logrus"
)

type conversionResult struct {
	OriginalFileName  string
	ConvertedFileName string
	ImageType         string
	Compressed        bool
	Renamed           bool
	Error             error
}

func conversionWorker(files chan string, results chan conversionResult, compress bool) {

	for file := range files {
		var result = conversionResult{
			OriginalFileName: file,
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

		if strings.HasSuffix(convertedFileName, ".jpeg") {
			var renamed = strings.ReplaceAll(convertedFileName, ".jpeg", ".jpg")

			if imageconvert.WouldOverwrite(convertedFileName, renamed) {
				log.Warnf("renaming %s would overwrite an existing jpeg, skipping", convertedFileName)
				continue
			}

			err = os.Rename(convertedFileName, renamed)
			if err != nil {
				result.Error = fmt.Errorf("could rename file: %s, err: %w", convertedFileName, err)
				results <- result
				continue
			}
			result.Renamed = true
		}
		results <- result
	}
	close(results)
}
