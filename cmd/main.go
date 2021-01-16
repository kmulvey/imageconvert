package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kmulvey/imageconvert/pkg/imageconvert"
)

func main() {
	var rootDir string
	flag.StringVar(&rootDir, "dir", "", "directory (abs path)")
	flag.Parse()
	if strings.TrimSpace(rootDir) == "" {
		log.Fatal("directory not provided")
	}

	var files, err = imageconvert.ListFiles(rootDir)
	imageconvert.HandleErr("list", err)

	for file := range files {
		if strings.HasSuffix(file, ".png") {
			if _, err := os.Stat(strings.Replace(file, ".png", ".jpg", 1)); err == nil {
				imageconvert.ConvertPng(file, strings.Replace(file, ".png", "-2.jpg", 1))
			} else {
				imageconvert.ConvertPng(file, strings.Replace(file, ".png", ".jpg", 1))
			}
		} else {
			if _, err := os.Stat(strings.Replace(file, ".webp", ".jpg", 1)); err == nil {
				imageconvert.ConvertWebp(file, strings.Replace(file, ".webp", "-2.jpg", 1))
			} else {
				imageconvert.ConvertWebp(file, strings.Replace(file, ".webp", ".jpg", 1))
			}
		}
		err = os.Remove(file)
		imageconvert.HandleErr("remove", err)
		fmt.Println("converted", file)
	}
}
