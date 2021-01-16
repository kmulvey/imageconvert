package imageconvert

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"golang.org/x/image/webp"
)

func ConvertPng(from, to string) {
	var pngFile, err = os.Open(from)
	HandleErr("png open", err)

	pngData, err := png.Decode(pngFile)
	HandleErr("png decode", err)

	out, err := os.Create(to)
	HandleErr("png create", err)

	err = jpeg.Encode(out, pngData, &jpeg.Options{Quality: 100})
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

	err = jpeg.Encode(out, webpData, &jpeg.Options{Quality: 100})
	HandleErr("jpg encode", err)

	err = out.Close()
	HandleErr("webp-jpg close", err)
	err = webpFile.Close()
	HandleErr("webp close", err)
}

func ListFiles(root string) (map[string]bool, error) {
	var allFiles = make(map[string]bool)
	var staticBool bool
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return allFiles, err
	}
	for _, file := range files {
		if file.IsDir() {
			var subFiles, err = ListFiles(path.Join(root, file.Name()))
			if err != nil {
				return allFiles, err
			}
			for subFile := range subFiles {
				allFiles[subFile] = staticBool
			}
		} else {
			if strings.HasSuffix(file.Name(), ".png") || strings.HasSuffix(file.Name(), ".webp") {
				allFiles[path.Join(root, file.Name())] = staticBool
			}
		}
	}
	return allFiles, nil
}

func HandleErr(prefix string, err error) {
	if err != nil {
		log.Fatal(fmt.Errorf("%s: %w", prefix, err))
	}
}
