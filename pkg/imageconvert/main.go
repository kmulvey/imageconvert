package imageconvert

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"os"

	log "github.com/sirupsen/logrus"

	"golang.org/x/image/webp"
)

func ConvertPng(from, to string) {
	var pngFile, err = os.Open(from)
	HandleErr("png open", err)

	pngData, err := png.Decode(pngFile)
	HandleErr("png decode", err)

	out, err := os.Create(to)
	HandleErr("png create", err)

	err = jpeg.Encode(out, pngData, &jpeg.Options{Quality: 85})
	HandleErr("jpg encode", err)

	err = out.Close()
	HandleErr("png-jpg close", err)

	err = pngFile.Close()
	HandleErr("png close", err)
}

func ConvertWebp(from, to string) {
	var webpFile, err = os.Open(from)
	HandleErr("webp open", err)

	webpData, err := webp.Decode(webpFile)
	HandleErr("webp decode", err)

	out, err := os.Create(to)
	HandleErr("webp create", err)

	err = jpeg.Encode(out, webpData, &jpeg.Options{Quality: 85})
	HandleErr("jpg encode", err)

	err = out.Close()
	HandleErr("webp-jpg close", err)

	err = webpFile.Close()
	HandleErr("webp close", err)
}

func HandleErr(prefix string, err error) {
	if err != nil {
		log.Fatal(fmt.Errorf("%s: %w", prefix, err))
	}
}
