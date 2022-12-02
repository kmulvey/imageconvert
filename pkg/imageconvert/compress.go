package imageconvert

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// QualityCheck uses imagemagick to determine the quality of the image
// and returns true if the quality is above a given threshold
func QualityCheck(maxQuality int, imagePath string) (bool, error) {
	imagePath = EscapeFilePath(imagePath)
	cmd := fmt.Sprintf("identify -format %s %s", "'%Q'", imagePath)
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return false, fmt.Errorf("error running identify on image: %s, error: %s, output: %s", imagePath, err.Error(), out)
	}
	imageQuality, err := strconv.ParseInt(string(out), 10, 0)
	if err != nil {
		return false, fmt.Errorf("error parsing int for quality on image: %s, error: %s", imagePath, err.Error())
	}

	return int64(maxQuality) >= imageQuality, nil
}

// CompressJPEG uses jpegoptim to compress the image.
// Return values:
// 1. was it compressed? jpegoptim may not be able to compress it any further
// 2. jpegoptim output (if you want to log it)
// 3. error
func CompressJPEG(quality int, imagePath string) (bool, string, error) {
	// have to escape the file spaces for the exec call
	var escapedImagePath = EscapeFilePath(imagePath)
	var cmdStr = fmt.Sprintf("jpegoptim -p -o -m%d %s", quality, escapedImagePath)

	output, err := exec.Command("bash", "-c", cmdStr).Output()
	if err != nil {
		return false, string(output), fmt.Errorf("error running jpegoptim on image: %s, error: %s, output: %s", imagePath, err.Error(), output)
	}

	if strings.Contains(string(output), "optimized.") {
		return true, string(output), nil
	}

	return false, string(output), nil
}
