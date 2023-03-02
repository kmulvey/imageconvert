package imageconvert

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"

	"github.com/kmulvey/resize"
)

func (ic *ImageConverter) Resize(filename string) (bool, error) {

	// open file
	file, err := os.Open(filename)
	if err != nil {
		return false, fmt.Errorf("error opening file for resizing: %w", err)
	}

	// get image config so we can look at height and width
	config, err := jpeg.DecodeConfig(file)
	if err != nil {
		return false, fmt.Errorf("error decoding config for resizing: %w", err)
	}

	if config.Width < int(ic.ResizeWidthThreshold) && config.Height < int(ic.ResizeHeightThreshold) {
		return false, nil
	}

	// rewind reader
	_, err = file.Seek(0, 0)
	if err != nil {
		return false, fmt.Errorf("error rewinding file reader: %w", err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		return false, fmt.Errorf("error decoding image for resizing: %w", err)
	}
	file.Close()

	// preserve aspect ratio
	var resizedImage image.Image
	if config.Width > config.Height {
		resizedImage = resize.Resize(uint(ic.ResizeWidth), 0, img, resize.Lanczos3)
	} else {
		resizedImage = resize.Resize(0, uint(ic.ResizeHeight), img, resize.Lanczos3)
	}

	out, err := os.Create(filename)
	if err != nil {
		return false, fmt.Errorf("error opening file for resizing: %w", err)
	}
	defer out.Close()

	// write new image to file
	err = jpeg.Encode(out, resizedImage, nil)
	if err != nil {
		return false, fmt.Errorf("error encoding resized image: %w", err)
	}

	return true, nil
}
