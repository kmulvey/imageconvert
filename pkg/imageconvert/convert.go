package imageconvert

import (
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

func Convert(from string) string {
	var origFile, err = os.Open(from)
	HandleErr("img open", err)

	var ext = filepath.Ext(from)
	var newFile = strings.Replace(from, ext, ".jpg", 1)

	imgData, _, err := image.Decode(origFile)
	HandleErr("img decode", err)

	err = origFile.Close()
	HandleErr("input img close", err)

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

	return newFile
}

func HandleErr(prefix string, err error) {
	if err != nil {
		log.Fatal(fmt.Errorf("%s: %w", prefix, err))
	}
}
