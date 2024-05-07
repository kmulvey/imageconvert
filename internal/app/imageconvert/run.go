package imageconvert

import (
	"context"

	"github.com/kmulvey/path"
)

// Start begins the conversion process and returns counts of each type of operation preformed.
func (ic *ImageConverter) Start(ctx context.Context, results chan ConversionResult) (int, int, error) {

	if ic.WatchDir != "" {
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
					var _, exists = ic.SkipMap[originalImage.AbsolutePath]

					if (!originalImage.Entry.FileInfo.IsDir() && !exists) || ic.Force {

						var result = ic.convertImage(originalImage.Entry)
						results <- ic.convertImage(originalImage.Entry)
						if result.Error == nil {
							ic.SkipMapFileHandle.WriteString(originalImage.AbsolutePath + "\n")
							ic.SkipMap[originalImage.AbsolutePath] = struct{}{}
						}
					}
				}
			}
			close(done)
		}()

		go path.WatchDir(ctx, ic.WatchDir, ic.Depth, true, originalImages, errors, path.NewDateWatchFilter(ic.TimeRange.From, ic.TimeRange.To))
		<-done

	} else {

		var resizedTotal, totalProcessed int
		for i, originalImage := range ic.OriginalImagesEntries {

			var _, exists = ic.SkipMap[originalImage.AbsolutePath]

			if (!originalImage.FileInfo.IsDir() && !exists) || ic.Force {
				var result = ic.convertImage(originalImage)
				if result.Error != nil {
					return 0, 0, result.Error
				}
				ic.SkipMapFileHandle.WriteString(originalImage.AbsolutePath + "\n")
				ic.SkipMap[originalImage.AbsolutePath] = struct{}{}

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
