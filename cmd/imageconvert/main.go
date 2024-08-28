package main

import (
	"flag"
	"os"
	"runtime"
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
	var inputPath, processedLogFile, resizeThreshold, resizeSize string
	var compress, force, watch, v, h bool
	var threads, directoryDepth int
	var tr humantime.TimeRange

	flag.StringVar(&inputPath, "path", "", "path to files, globbing must be quoted")
	flag.StringVar(&processedLogFile, "processed-file", "processed.log", "the file to write processes images to, so that we dont processes them again next time")
	flag.StringVar(&resizeThreshold, "resize-threshold", "", "the min size to consider for resizing in the formate [width]x[height] e.g. 2560x1440")
	flag.StringVar(&resizeSize, "resize-size", "", "the size to resize the images to while preserving the aspect ratio [width]x[height] e.g. 5120x2880")
	flag.IntVar(&threads, "threads", 1, "number of threads to use")
	flag.IntVar(&directoryDepth, "depth", 1, "number levels to search directories for images")
	flag.Var(&tr, "time-range", "process files chnaged since this time")
	flag.BoolVar(&compress, "compress", true, "compress")
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

	if threads <= 0 || threads > runtime.GOMAXPROCS(0) {
		threads = 1
		log.Infof("invalid thread count: %d, setting threads to 1", threads)
	}

	log.Infof("Config: dir: %s, log file: %s, compress: %t, force: %t, watch: %t, threads: %d, modified-since: %s", inputPath, processedLogFile, compress, force, watch, threads, tr)

	var ic, err = imageconvert.NewWithDefaults(inputPath, processedLogFile, uint8(directoryDepth))
	if err != nil {
		log.Fatalf("error starting: %s", err)
	}

	if compress {
		ic.WithCompression(90)
	}

	if force {
		ic.WithForce()
	}

	resizeThreshold = strings.TrimSpace(resizeThreshold)
	if resizeThreshold != "" {

		var thresholdArr = strings.Split(resizeThreshold, "x")
		if len(thresholdArr) != 2 {
			log.Fatalf("resize threshold not in the format: [width]x[height] e.g. 230x400, input: %s, error: %s", resizeThreshold, err)
		}

		var sizeArr = strings.Split(resizeSize, "x")
		if len(sizeArr) != 2 {
			log.Fatalf("resize size not in the format: [width]x[height] e.g. 230x400, input: %s, error: %s", sizeArr, err)
		}

		ic.WithResize(getResizeValue(sizeArr[0]), getResizeValue(sizeArr[1]), getResizeValue(thresholdArr[0]), getResizeValue(thresholdArr[1]))
	}

	if watch {
		ic.WithWatch()
	}

	if threads > 1 {
		ic.WithThreads(threads)
	}

	if tr.From != imageconvert.NilTime || tr.To != imageconvert.NilTime {
		ic.WithTimeRange(tr)
	}

	compressedTotal, renamedTotal, resizedTotal, totalFiles, conversionTypeTotals, err := ic.Start(nil)
	if err != nil {
		log.Error(err)
	}

	log.WithFields(log.Fields{
		"converted pngs":   conversionTypeTotals["png"],
		"converted webps":  conversionTypeTotals["webp"],
		"compressed":       compressedTotal,
		"jpegs renamed":    renamedTotal,
		"resized":          resizedTotal,
		"total files seen": totalFiles,
	}).Info("Done")
}

func getResizeValue(str string) uint16 {

	var num, err = strconv.ParseUint(str, 10, 16)
	if err != nil {
		log.Fatalf("error resize value is not a number: '%s', err: %s", str, err.Error())
	}

	return uint16(num)
}
