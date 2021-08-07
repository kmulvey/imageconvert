package main

import (
	"flag"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/kmulvey/imageconvert/pkg/imageconvert"
)

func main() {
	var rootDir string
	var compress bool

	flag.StringVar(&rootDir, "dir", "", "directory (abs path)")
	flag.BoolVar(&compress, "compress", false, "compress")
	flag.Parse()
	if strings.TrimSpace(rootDir) == "" {
		log.Fatal("directory not provided")
	}

	// these are all the files all the way down the dir tree
	var files = imageconvert.ListFiles(rootDir)

	// consistant extention for jpg
	var jpegRename int
	for _, filename := range imageconvert.FilerJPG(files) {
		if strings.HasSuffix(filename, ".JPG") {
			imageconvert.HandleErr("rename", os.Rename(filename, strings.ReplaceAll(filename, ".JPG", ".jpg")))
			jpegRename++
		} else if strings.HasSuffix(strings.ToLower(filename), ".jpeg") {
			var newFile = strings.ReplaceAll(filename, ".jpeg", ".jpg")
			newFile = strings.ReplaceAll(newFile, ".JPEG", ".jpg")
			imageconvert.HandleErr("rename", os.Rename(filename, newFile))
			jpegRename++
		}
	}

	// png -> jpg
	var pngs = imageconvert.FilerPNG(files)
	for _, filename := range pngs {
		// we dont want to overwite an existing jpg
		if _, err := os.Stat(strings.Replace(filename, ".png", ".jpg", 1)); err == nil {
			imageconvert.ConvertPng(filename, strings.Replace(filename, ".png", "-"+time.Now().String()+".jpg", 1))
		} else {
			imageconvert.ConvertPng(filename, strings.Replace(filename, ".png", ".jpg", 1))
		}
		imageconvert.HandleErr("remove", os.Remove(filename))
	}

	// webp -> jpg
	var webps = imageconvert.FilerWEBP(files)
	for _, filename := range webps {
		// we dont want to overwite an existing jpg
		if _, err := os.Stat(strings.Replace(filename, ".webp", ".jpg", 1)); err == nil {
			imageconvert.ConvertWebp(filename, strings.Replace(filename, ".webp", "-"+time.Now().String()+".jpg", 1))
		} else {
			imageconvert.ConvertWebp(filename, strings.Replace(filename, ".webp", ".jpg", 1))
		}
		imageconvert.HandleErr("remove", os.Remove(filename))
	}

	var compressed int
	if compress {
		// some files may have gotten renamed above so we call ListFiles again
		for _, filename := range imageconvert.FilerJPG(imageconvert.ListFiles(rootDir)) {
			imageconvert.CompressJPEG(85, filename)
			compressed++
		}
	}

	log.WithFields(log.Fields{
		"converted pngs":  len(pngs),
		"converted webps": len(webps),
		"compressed":      compressed,
		"jpgs renamed":    jpegRename,
	}).Info("Done")
}
