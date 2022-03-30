package imageconvert

import (
	"errors"
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
// the second string returned is the type of image (png, webp)
func Convert(from string) (string, string) {
	var origFile, err = os.Open(from)
	HandleErr("img open", err)

	var ext = filepath.Ext(from)
	var newFile = strings.Replace(from, ext, ".jpg", 1)

	imgData, imageType, err := image.Decode(origFile)
	HandleErr("img decode", err)

	// dont bother converting jpegs
	if imageType == "jpeg" {
		return from, ""
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
			return from, ""
		}
	}

	err = origFile.Close()
	HandleErr("input img close", err)

	err = os.Remove(from)
	HandleErr("remove input file", err)

	out, err := os.Create(newFile)
	HandleErr("new jpg create", err)

	err = jpeg.Encode(out, imgData, &jpeg.Options{Quality: 85})
	HandleErr("jpg encode", err)

	err = out.Close()
	HandleErr("new jpg close", err)

	log.WithFields(log.Fields{
		"from": from,
		"to":   newFile,
	}).Info("Converted")

	return newFile, imageType
}

// WouldOverwrite looks to see if the file were to be converted to a jpeg,
// would it overwite an existing jpg file with the same name
func WouldOverwrite(path, imageType string) bool {
	var ext = filepath.Ext(path)
	var jpgPath = strings.Replace(path, ext, ".jpg", 1)

	// find an existing file
	if _, err := os.Stat(jpgPath); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
