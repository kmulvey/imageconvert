package imageconvert

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"

	// make sure png is always imported here.
	_ "image/png"
	// make sure webp is always imported here.
	_ "golang.org/x/image/webp"
)

// make sure the following are always imported above, some editors may remove them
// _ "golang.org/x/image/webp"
// _ "image/png"
// "image/jpeg"

const imageExt = "jpeg"

// Convert converts pngs and webps to jpeg, if successful the inputFile is deleted.
// This first string returned is the name of the new file.
// The second string returned is the type of the input image (png, webp), as detected from its encoding, not file name.
func Convert(inputFile string) (string, string, error) {

	var ext = filepath.Ext(inputFile)
	var convertedFile = strings.Replace(inputFile, ext, ".jpg", 1)

	var origFile, err = os.OpenFile(inputFile, os.O_RDONLY, 0755)
	if err != nil {
		return "", "", fmt.Errorf("error opening file for conversion, image: %s, error: %w", inputFile, err)
	}
	defer origFile.Close()

	imgData, imageType, err := image.Decode(origFile)
	if err != nil {
		return "", "", fmt.Errorf("error decoding image: %s, error: %w", inputFile, err)
	}

	// dont bother converting jpegs
	if imageType == imageExt {
		return inputFile, imageExt, nil
	}

	// Dont convert images that would result in an overwrite
	// e.g. a.png a.jpg exist, thus converting a.png would overwrite a.jpg, so we let the user handle it.
	// EXCEPT fake jpgs:
	// A "fake jpeg" is an image that has the extension .jpg or .jpeg but is really a different format e.g. png image named "x.jpg".
	// Basically dont trust file extensions.
	var fakeJPG bool
	if ext == ".jpg" && imageType != imageExt {
		fakeJPG = true
	}

	if WouldOverwrite(inputFile) && !fakeJPG {
		// nolint: err113
		return inputFile, imageType, errors.New("converting " + inputFile + " would overwrite an existing jpeg, skipping")
	}

	out, err := os.Create(convertedFile)
	if err != nil {
		return "", "", fmt.Errorf("error creating new image: %s, error: %w", inputFile, err)
	}

	if err := jpeg.Encode(out, imgData, &jpeg.Options{Quality: 100}); err != nil {
		return "", "", fmt.Errorf("error encoding new image: %s, error: %w", inputFile, err)
	}

	if err := out.Close(); err != nil {
		return "", "", fmt.Errorf("error closing new image: %s, error: %w", inputFile, err)
	}

	if err := origFile.Close(); err != nil {
		return "", "", fmt.Errorf("error closing new image: %s, error: %w", inputFile, err)
	}

	if !fakeJPG {
		if err := os.Remove(inputFile); err != nil {
			return "", "", fmt.Errorf("error removing image: %s, error: %w", inputFile, err)
		}
	}

	return convertedFile, imageType, nil
}

// WouldOverwrite looks to see if the file were to be converted to a jpeg,
// would it overwite an existing jpg file with the same name.
func WouldOverwrite(path string) bool {

	var ext = filepath.Ext(path)
	var jpgPath = strings.Replace(path, ext, ".jpg", 1)

	// find an existing file
	if _, err := os.Stat(jpgPath); errors.Is(err, os.ErrNotExist) {
		return false // there may be other errors but we live on the edge
	}

	return true
}
