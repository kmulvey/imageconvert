package imageconvert

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Kagami/go-avif"
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

	// ubuntu does not use the 'magick' prefix, everything else does
	var cmd = exec.Command("magick", "-version")
	var err = cmd.Run()
	if err != nil {
		cmd = exec.Command("identify", "-format", "'%Q'", imagePath)
	} else {
		cmd = exec.Command("magick", "identify", "-format", "'%Q'", imagePath)
	}
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return false, fmt.Errorf("error running identify on image: %s, error: %s, stderr: %s, output: %s", imagePath, err.Error(), stderr.String(), out.String())
	}

	var qualityStr = strings.ReplaceAll(out.String(), "'", "")
	imageQuality, err := strconv.ParseInt(qualityStr, 10, 0)
	if err != nil {
		return false, fmt.Errorf("error parsing int for quality on image: %s, quality: %s, error: %s", imagePath, qualityStr, err.Error())
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
	// nolint here because its worried about the Itoa ... but why?!?!
	// nolint:gosec
	var cmd = exec.Command("jpegoptim", "-p", "-o", "-v", "-m", strconv.Itoa(quality), escapedImagePath)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	var err = cmd.Run()
	outStr, errStr := stdout.String(), stderr.String()

	// check stdout first because we are getting some false 'exit status 1' in err
	if strings.Contains(outStr, "optimized.") {
		return true, outStr, nil
	}

	if err != nil {
		return false, outStr, fmt.Errorf("error running jpegoptim on image: %s, error: %s, stdErr: %s, output: %s", imagePath, err.Error(), errStr, outStr)
	}

	return false, outStr, nil
}

// quality 30
func CompressAVIF(quality, threads int, inputImagePath string) (string, error) {

	var src, err = os.Open(inputImagePath)
	if err != nil {
		return "", fmt.Errorf("error opening image: %s, error: %w", inputImagePath, err)
	}

	img, _, err := image.Decode(src)
	if err != nil {
		return "", fmt.Errorf("error decoding image: %s, error: %w", inputImagePath, err)
	}

	var outputImagePath = filepath.Base(strings.ReplaceAll(inputImagePath, filepath.Ext(inputImagePath), ".avif"))
	dst, err := os.Create(outputImagePath)
	if err != nil {
		return "", fmt.Errorf("error creating new image: %s, error: %w", outputImagePath, err)
	}

	var options = avif.Options{
		Threads:        threads,
		Speed:          avif.MaxSpeed,
		Quality:        quality,
		SubsampleRatio: nil,
	}
	err = avif.Encode(dst, img, &options)
	if err != nil {
		return "", fmt.Errorf("error encoding avif: %s, error: %w", outputImagePath, err)
	}

	return outputImagePath, nil
}

// EscapeFilePath escapes spaces in the filepath used for an exec() call.
func EscapeFilePath(file string) string {
	var r = strings.NewReplacer(" ", `\ `, "(", `\(`, ")", `\)`, "'", `\'`, "&", `\&`, "@", `\@`)
	return r.Replace(file)
}
