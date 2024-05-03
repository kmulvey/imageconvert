package imageconvert

import (
	"context"

	"github.com/kmulvey/path"
)

// Start begins the conversion process and returns counts of each type of operation preformed.
func (ic *ImageConverter) Start(results chan ConversionResult) (int, int, error) {

	if ic.Watch {
		var originalImages = make(chan path.WatchEvent)
		var errors = make(chan error)
		var done = make(chan struct{})

		go func() {
			var errorsOpen, originalImagesOpen = true, true
			for errorsOpen || originalImagesOpen {
				select {

				case err, open := <-errors:
					if !open {
						errorsOpen = false
						continue
					}
					results <- ConversionResult{
						Error: err,
					}

				case originalImage, open := <-originalImages:
					if !open {
						originalImagesOpen = false
						continue
					}
					if !originalImage.Entry.FileInfo.IsDir() {
						results <- ic.convertImage(originalImage.Entry)
					}
				}
			}
			close(done)
		}()

		go path.WatchDir(context.Background(), ic.OriginalImagesEntry.AbsolutePath, ic.Depth, true, originalImages, errors)
		<-done

	} else {

		originalImages, err := ic.OriginalImagesEntry.Flatten(true)
		if err != nil {
			return 0, 0, err
		}

		var resizedTotal, totalProcessed int
		for i, image := range originalImages {
			if !image.FileInfo.IsDir() {
				var result = ic.convertImage(image)
				if err != nil {
					return 0, 0, err
				}

				if result.Resized {
					resizedTotal++
				}
			}
			totalProcessed = i
		}

		return totalProcessed, resizedTotal, nil
	}

	return 0, 0, nil
}
