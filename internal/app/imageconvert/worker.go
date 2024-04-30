package imageconvert

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kmulvey/path"
)

type ConversionResult struct {
	OriginalFileName  string
	ConvertedFileName string
	ImageType         string
	CompressOutput    string
	Error             error
	Compressed        bool
	Renamed           bool
	Resized           bool
}

// conversionWorker reads from the file chan and does all the conversion work.
func (ic *ImageConverter) conversionWorker(files chan path.Entry, results chan ConversionResult, done chan struct{}) {
	defer close(results)
	defer close(done)

	for {
		select {
		case _, open := <-ic.ShutdownTrigger:
			if !open {
				return
			}
		case file, open := <-files:
			if !open {
				return
			}
			results <- ic.convertImage(file)
		}
	}
}

// convertImage  converts, compresses and renames images.
// This is broken out from the conversionWorker for ease of testing.
// 1. does a file already exist with the output name? yes = skip
// 2. convert it to jpg
// 3. compress it (if enabled)
// 4. reset the mod time
func (ic *ImageConverter) convertImage(originalFile path.Entry) ConversionResult {

	var result = ConversionResult{
		OriginalFileName: originalFile.String(),
	}

	// CONVERT IT
	var imageType string
	var err error
	result.ConvertedFileName, imageType, err = Convert(originalFile.String())
	if err != nil {
		result.Error = fmt.Errorf("error converting image: %s, error: %w", originalFile, err)
		return result
	}
	result.ImageType = imageType

	// RESIZE IT
	if ic.ResizeWidth > 0 && ic.ResizeHeight > 0 {
		resized, err := ic.Resize(result.ConvertedFileName)
		if err != nil {
			result.Error = fmt.Errorf("error resizing image: %s, error: %w", originalFile, err)
			return result
		}
		result.Resized = resized
	}

	// COMPRESS IT
	if ic.Compress {
		compressed, stdout, err := CompressJPEG(85, result.ConvertedFileName)
		if err != nil {
			result.Error = fmt.Errorf("error compressing image: %s, error: %w", originalFile, err)
			return result
		}
		result.Compressed = compressed

		if compressed {
			result.CompressOutput = stdout
		}
	}

	// make sure every file ends in ".jpg".
	if RenameSuffixRegex.MatchString(filepath.Base(result.ConvertedFileName)) {
		var renamedFile = strings.Replace(result.ConvertedFileName, filepath.Ext(result.ConvertedFileName), ".jpg", 1)
		if err := os.Rename(result.ConvertedFileName, renamedFile); err != nil {
			result.Error = fmt.Errorf("could rename file: %s, err: %w", result.ConvertedFileName, err)
			return result
		}
		result.ConvertedFileName = renamedFile
		result.Renamed = true
	}

	// RESET MODTIME
	err = os.Chtimes(result.ConvertedFileName, originalFile.FileInfo.ModTime(), originalFile.FileInfo.ModTime())
	if err != nil {
		result.Error = fmt.Errorf("could not reset mod time of file: %s, err: %w", result.ConvertedFileName, err)
		return result
	}

	return result
}
