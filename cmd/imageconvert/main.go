package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kmulvey/humantime"
	"github.com/kmulvey/imageconvert/v2/internal/app/imageconvert"
	log "github.com/sirupsen/logrus"
	"go.szostok.io/version"
	"go.szostok.io/version/printer"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.TimeOnly, // "15:04:05",
	})

	// get the user options
	var images, skipMapFile, resizeThreshold, resizeSize string
	var deleteOriginal, force, watch, v, h bool
	var threads, quality, directoryDepth int
	var tr humantime.TimeRange

	flag.StringVar(&images, "images", "", "path to images, globbing and multiple paths must be quoted")
	flag.StringVar(&skipMapFile, "skip-map-file", "processed.log", "the file to write processes images to, so that we dont processes them again next time")
	flag.StringVar(&resizeThreshold, "resize-threshold", "", "the min size to consider for resizing in the formate [width]x[height] e.g. 2560x1440")
	flag.StringVar(&resizeSize, "resize-size", "", "the size to resize the images to while preserving the aspect ratio [width]x[height] e.g. 5120x2880")
	flag.IntVar(&threads, "threads", 1, "number of threads to use")
	flag.IntVar(&quality, "quality", 30, "avif quality to use")
	flag.IntVar(&directoryDepth, "depth", 1, "number levels to search directories for images")
	flag.Var(&tr, "time-range", "process files chnaged since this time")
	flag.BoolVar(&deleteOriginal, "delete-original", false, "delete original files")
	flag.BoolVar(&force, "force", false, "force")
	flag.BoolVar(&watch, "watch", false, "watch the dir")
	flag.BoolVar(&v, "version", false, "print version")
	flag.BoolVar(&v, "v", false, "print version")
	flag.BoolVar(&h, "help", false, "print options")
	flag.Parse()

	if h {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if v {
		var verPrinter = printer.New()
		var info = version.Get()
		if err := verPrinter.PrintInfo(os.Stdout, info); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	var width, height, widthThreshold, heightThreshold, err = parseResize(resizeThreshold, resizeSize)
	if err != nil {
		log.Fatal(err)
	}

	var imagesArr = strings.Split(images, " ")

	var config = &imageconvert.ImageConverterConfig{
		OriginalImages:        imagesArr,
		Threads:               uint8(threads),
		Quality:               uint8(quality),
		SkipMapFile:           skipMapFile,
		ResizeWidth:           width,
		ResizeWidthThreshold:  widthThreshold,
		ResizeHeight:          height,
		ResizeHeightThreshold: heightThreshold,
		TimeRange:             tr,
		Force:                 force,
		DeleteOriginal:        deleteOriginal,
		// TODO: WatchDir:
	}

	log.Infof(`Config: 
OriginalImages: 		%+v,
SkipMapFile:			%s,
Threads:			%d,
Quality:			%d,
Resize:				%s,
ResizeThreshold:		%s,
TimeRange From:			%s,
TimeRange To:			%s,
Force:				%t,
DeleteOriginal:			%t,
	`, imagesArr, skipMapFile, threads, quality, resizeSize, resizeThreshold, tr.From, tr.To, force, deleteOriginal)

	ic, err := imageconvert.NewImageConverter(config)
	if err != nil {
		log.Fatal(err)
	}

	processedTotal, resizedTotal, err := ic.Start()
	if err != nil {
		log.Error(err)
	}

	log.WithFields(log.Fields{
		"resized":               resizedTotal,
		"total files processed": processedTotal,
	}).Info("Done")
}

func parseResize(resizeThreshold, resizeSize string) (uint16, uint16, uint16, uint16, error) {

	resizeThreshold = strings.TrimSpace(resizeThreshold)
	resizeSize = strings.TrimSpace(resizeSize)

	if resizeThreshold != "" {
		var thresholdArr = strings.Split(resizeThreshold, "x")
		if len(thresholdArr) != 2 {
			return 0, 0, 0, 0, fmt.Errorf("resize threshold not in the format: [width]x[height] e.g. 230x400, input: %s", resizeThreshold)
		}

		var sizeArr = strings.Split(resizeSize, "x")
		if len(sizeArr) != 2 {
			return 0, 0, 0, 0, fmt.Errorf("resize size not in the format: [width]x[height] e.g. 230x400, input: %s", sizeArr)
		}
		return getResizeValue(sizeArr[0]), getResizeValue(sizeArr[1]), getResizeValue(thresholdArr[0]), getResizeValue(thresholdArr[1]), nil
	}
	return 0, 0, 0, 0, nil
}

func getResizeValue(str string) uint16 {

	var num, err = strconv.ParseUint(str, 10, 16)
	if err != nil {
		log.Fatalf("error resize value is not a number: '%s', err: %s", str, err.Error())
	}

	return uint16(num)
}
