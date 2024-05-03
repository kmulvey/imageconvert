package imageconvert

import (
	"fmt"
	"os"

	"github.com/kmulvey/imageconvert/v2/pkg/imageconvert"
	"github.com/kmulvey/path"
)

// ConversionResult is all the information about the image that was converted.
type ConversionResult struct {
	OriginalFileName  string
	ConvertedFileName string
	Error             error
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

// convertImage resizes and converts images.
// This is broken out from the conversionWorker for ease of testing.
func (ic *ImageConverter) convertImage(originalFile path.Entry) ConversionResult {

	var result = ConversionResult{
		OriginalFileName: originalFile.String(),
	}

	// RESIZE IT
	if ic.ResizeWidth > 0 && ic.ResizeHeight > 0 {
		resized, err := imageconvert.Resize(result.ConvertedFileName, int(ic.ResizeWidthThreshold), int(ic.ResizeHeightThreshold), uint(ic.ResizeWidth), uint(ic.ResizeHeight))
		if err != nil {
			result.Error = fmt.Errorf("error resizing image: %s, error: %w", originalFile, err)
			return result
		}
		result.Resized = resized
	}

	// CONVERT IT
	var err error
	result.ConvertedFileName, err = imageconvert.CompressAVIF(ic.Quality, int(ic.Threads), originalFile.AbsolutePath)
	if err != nil {
		result.Error = fmt.Errorf("error converting image: %s, error: %w", originalFile, err)
		return result
	}

	// DELETE IT
	if ic.DeleteOriginal {
		if err := os.Remove(originalFile.AbsolutePath); err != nil {
			result.Error = fmt.Errorf("error removing original image: %s, error: %w", originalFile, err)
			return result
		}
	}

	// RESET MODTIME
	err = os.Chtimes(result.ConvertedFileName, originalFile.FileInfo.ModTime(), originalFile.FileInfo.ModTime())
	if err != nil {
		result.Error = fmt.Errorf("could not reset mod time of file: %s, err: %w", result.ConvertedFileName, err)
		return result
	}

	return result
}
