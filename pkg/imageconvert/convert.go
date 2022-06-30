package imageconvert

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	_ "golang.org/x/image/webp"
)

// make sure the following are always imported above, some editors may remove them
// _ "golang.org/x/image/webp"
// _ "image/png"
// "image/jpeg"

// Convert converts pngs and webps to jpeg
// this first string returned is the name of the new file
// the second string returned is the type of the input image (png, webp), as detected from its encoding, not file name
func Convert(from string) (string, string, error) {
	var origFile, err = os.Open(from)
	if err != nil {
		return "", "", fmt.Errorf("error opening file for conversion, image: %s, error: %s", from, err.Error())
	}
	defer origFile.Close()

	var ext = filepath.Ext(from)
	var newFile = strings.Replace(from, ext, ".jpg", 1)

	imgData, imageType, err := image.Decode(origFile)
	if err != nil {
		return "", "", fmt.Errorf("error decoding image: %s, error: %s", from, err.Error())
	}

	// dont bother converting jpegs
	if imageType == "jpeg" {
		return from, "", nil
	}

	// dont convert images that would result in an overwrite
	// e.g. a.png a.jpg exist, thus converting a.png would overwrite a.jpg,
	// so we let the user handle it
	// EXCEPT fake jpgs:
	// a "fake jpg" is an image that has the extension .jpg or .jpeg but is
	// really a different format e.g. png image named "x.jpg"
	// basically we cant just trust file extensions
	if WouldOverwrite(from, imageType) {
		// we only warn if the detected image format has the corresponding extension
		if "."+imageType == ext {
			log.Warnf("converting %s would overwrite an existing jpeg, skipping", from)
			return from, "", nil
		}
	}

	err = os.Remove(from)
	if err != nil {
		return "", "", fmt.Errorf("error removing image: %s, error: %s", from, err.Error())
	}

	out, err := os.Create(newFile)
	if err != nil {
		return "", "", fmt.Errorf("error creating new image: %s, error: %s", from, err.Error())
	}

	err = jpeg.Encode(out, imgData, &jpeg.Options{Quality: 85})
	if err != nil {
		return "", "", fmt.Errorf("error encoding new image: %s, error: %s", from, err.Error())
	}

	err = out.Close()
	if err != nil {
		return "", "", fmt.Errorf("error closing new image: %s, error: %s", from, err.Error())
	}

	log.WithFields(log.Fields{
		"from": from,
		"to":   newFile,
	}).Info("Converted")

	return newFile, imageType, nil
}

// WouldOverwrite looks to see if the file were to be converted to a jpeg,
// would it overwite an existing jpg file with the same name
func WouldOverwrite(path, imageType string) bool {
	var ext = filepath.Ext(path)
	var jpgPath = strings.Replace(path, ext, ".jpg", 1)

	// find an existing file
	if _, err := os.Stat(jpgPath); errors.Is(err, os.ErrNotExist) {
		return false // there may be other errors but we live on the edge
	}
	return true
}
