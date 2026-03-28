package imageconvert

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/kmulvey/resize"
)

func Resize(filename string, resizeWidthThreshold, resizeHeightThreshold int, resizeWidth, resizeHeight uint) (bool, error) {

	// open file
	file, err := os.OpenFile(filename, os.O_RDONLY, 0600) //nolint:gosec // filename is a caller-supplied image path
	if err != nil {
		return false, fmt.Errorf("error opening file for resizing: %w", err)
	}

	// get image config so we can look at height and width
	config, err := jpeg.DecodeConfig(file)
	if err != nil {
		return false, fmt.Errorf("error decoding config for resizing: %w", err)
	}

	if config.Width < resizeWidthThreshold && config.Height < resizeHeightThreshold {
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
	err = file.Close()
	if err != nil {
		return false, fmt.Errorf("error closing original image file: %w", err)
	}

	// preserve aspect ratio
	var resizedImage image.Image
	if config.Width > config.Height {
		resizedImage = resize.Resize(resizeWidth, 0, img, resize.Lanczos3)
	} else {
		resizedImage = resize.Resize(0, resizeHeight, img, resize.Lanczos3)
	}

	out, err := os.Create(filepath.Clean(filename))
	if err != nil {
		return false, fmt.Errorf("error opening file for resizing: %w", err)
	}

	// write new image to file
	err = jpeg.Encode(out, resizedImage, nil)
	if err != nil {
		_ = out.Close()
		return false, fmt.Errorf("error encoding resized image: %w", err)
	}

	if err = out.Close(); err != nil {
		return false, fmt.Errorf("error closing resized image file: %w", err)
	}

	return true, nil
}
