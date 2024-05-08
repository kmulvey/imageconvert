package imageconvert

import (
	"context"

	"github.com/kmulvey/path"
)

func (ic *ImageConverter) Watch(ctx context.Context, results chan ConversionResult) {

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
					if result.Error == nil {
						if _, err := ic.SkipMapFileHandle.WriteString(originalImage.AbsolutePath + "\n"); err != nil {
							result.Error = err
						}
						ic.SkipMap[originalImage.AbsolutePath] = struct{}{}
					}
					results <- ic.convertImage(originalImage.Entry)
				}
			}
		}
		close(done)
	}()

	go path.WatchDir(ctx, ic.WatchDir, ic.Depth, true, originalImages, errors, path.NewDateWatchFilter(ic.TimeRange.From, ic.TimeRange.To))
	<-done
}

// Start begins the conversion process and returns counts of each type of operation preformed.
func (ic *ImageConverter) Start() (int, int, error) {

	var resizedTotal, processedTotal int
	for i, originalImage := range ic.OriginalImagesEntries {

		var _, exists = ic.SkipMap[originalImage.AbsolutePath]

		if (!originalImage.FileInfo.IsDir() && !exists) || ic.Force {
			var result = ic.convertImage(originalImage)
			if result.Error != nil {
				return 0, 0, result.Error
			}

			if _, err := ic.SkipMapFileHandle.WriteString(originalImage.AbsolutePath + "\n"); err != nil {
				return 0, 0, err
			}
			ic.SkipMap[originalImage.AbsolutePath] = struct{}{}

			if result.Resized {
				resizedTotal++
			}
			processedTotal = i + 1
		}
	}

	return processedTotal, resizedTotal, nil
}
