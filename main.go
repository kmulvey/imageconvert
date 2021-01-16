package main

import (
	"flag"
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

func main() {
	var rootDir string
	flag.StringVar(&rootDir, "dir", "", "directory (abs path)")
	flag.Parse()
	if strings.TrimSpace(rootDir) == "" {
		log.Fatal("directory not provided")
	}

	var files, err = listFiles(rootDir)
	handleErr("list", err)

	for file := range files {
		if strings.HasSuffix(file, ".png") {
			if _, err := os.Stat(strings.Replace(file, ".png", ".jpg", 1)); err == nil {
				convertPng(file, strings.Replace(file, ".png", "-2.jpg", 1))
			} else {
				convertPng(file, strings.Replace(file, ".png", ".jpg", 1))
			}
		} else {
			if _, err := os.Stat(strings.Replace(file, ".webp", ".jpg", 1)); err == nil {
				convertWebp(file, strings.Replace(file, ".webp", "-2.jpg", 1))
			} else {
				convertWebp(file, strings.Replace(file, ".webp", ".jpg", 1))
			}
		}
		err = os.Remove(file)
		handleErr("remove", err)
		fmt.Println("converted", file)
	}
}

func convertPng(from, to string) {
	var pngFile, err = os.Open(from)
	handleErr("png open", err)

	pngData, err := png.Decode(pngFile)
	handleErr("png decode", err)

	out, err := os.Create(to)
	handleErr("png create", err)

	err = jpeg.Encode(out, pngData, &jpeg.Options{Quality: 100})
	handleErr("jpg encode", err)

	err = out.Close()
	handleErr("png-jpg close", err)
	err = pngFile.Close()
	handleErr("png close", err)
}

func convertWebp(from, to string) {
	var webpFile, err = os.Open(from)
	handleErr("webp open", err)

	webpData, err := webp.Decode(webpFile)
	handleErr("webp decode", err)

	out, err := os.Create(to)
	handleErr("webp create", err)

	err = jpeg.Encode(out, webpData, &jpeg.Options{Quality: 100})
	handleErr("jpg encode", err)

	err = out.Close()
	handleErr("webp-jpg close", err)
	err = webpFile.Close()
	handleErr("webp close", err)
}

func listFiles(root string) (map[string]bool, error) {
	var allFiles = make(map[string]bool)
	var staticBool bool
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return allFiles, err
	}
	for _, file := range files {
		if file.IsDir() {
			var subFiles, err = listFiles(path.Join(root, file.Name()))
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

func handleErr(prefix string, err error) {
	if err != nil {
		log.Fatal(fmt.Errorf("%s: %w", prefix, err))
	}
}
