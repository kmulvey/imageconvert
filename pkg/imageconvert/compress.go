package imageconvert

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// QualityCheck uses imagemagick to determine the quality of the image
// and returns true if the quality is above a given threshold
func QualityCheck(maxQuality int, imagePath string) (bool, error) {

	// lint input, helps prevent arbitrary code execution
	if _, err := os.Stat(imagePath); err != nil {
		return false, err
	}

	// have to escape the file spaces for the exec call
	imagePath = EscapeFilePath(imagePath)

	var identifyCmd = fmt.Sprintf("identify -format %s %s", "'%Q'", imagePath)
	if runtime.GOOS == "windows" {
		identifyCmd = "magick " + identifyCmd
	}

	var cmd = exec.Command("bash", "-c", identifyCmd)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	var err = cmd.Run()
	if err != nil {
		return false, fmt.Errorf("error running identify on image: %s, error: %s, stderr: %s, output: %s", imagePath, err.Error(), stderr.String(), out.String())
	}

	imageQuality, err := strconv.ParseInt(out.String(), 10, 0)
	if err != nil {
		return false, fmt.Errorf("error parsing int for quality on image: %s, error: %s", imagePath, err.Error())
	}

	return imageQuality >= int64(maxQuality), nil
}

// CompressJPEG uses jpegoptim to compress the image.
// Return values:
// 1. was it compressed? jpegoptim may not be able to compress it any further
// 2. jpegoptim output (if you want to log it)
// 3. error
func CompressJPEG(quality int, imagePath string) (bool, string, error) {

	// lint input, helps prevent arbitrary code execution
	if _, err := os.Stat(imagePath); err != nil {
		return false, "", err
	}

	// have to escape the file spaces for the exec call
	var escapedImagePath = EscapeFilePath(imagePath)
	var cmdStr = fmt.Sprintf("jpegoptim -p -o -v -m%d %s", quality, escapedImagePath)

	output, err := exec.Command("bash", "-c", cmdStr).Output()
	if err != nil {
		return false, string(output), fmt.Errorf("error running jpegoptim on image: %s, error: %s, output: %s", imagePath, err.Error(), string(output))
	}

	fmt.Println("output ", string(output))
	if strings.Contains(string(output), "optimized.") {
		return true, string(output), nil
	}

	return false, string(output), nil
}
