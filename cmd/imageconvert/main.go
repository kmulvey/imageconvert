package main

import (
	"errors"
	"flag"
	"fmt"
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

// nolint: funlen
func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.TimeOnly, // "15:04:05",
	})

	// get the user options
	var inputPath, processedLogFile, resizeThreshold, resizeSize string
	var compress, force, watch, v, h bool
	var threads, directoryDepth int
	var timerange humantime.TimeRange

	flag.StringVar(&inputPath, "path", "", "path to files, globbing must be quoted")
	flag.StringVar(&processedLogFile, "processed-file", "processed.log", "the file to write processes images to, so that we dont processes them again next time")
	flag.StringVar(&resizeThreshold, "resize-threshold", "", "the min size to consider for resizing in the formate [width]x[height] e.g. 2560x1440")
	flag.StringVar(&resizeSize, "resize-size", "", "the size to resize the images to while preserving the aspect ratio [width]x[height] e.g. 5120x2880")
	flag.IntVar(&threads, "threads", 1, "number of threads to use")
	flag.IntVar(&directoryDepth, "depth", 1, "number levels to search directories for images")
	flag.Var(&timerange, "time-range", "process files chnaged since this time")
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

	log.Infof("Config: dir: %s, log file: %s, compress: %t, force: %t, watch: %t, threads: %d, modified-since: %s", inputPath, processedLogFile, compress, force, watch, threads, timerange)

	var configs, err = parseParams(compress, force, watch, uint8(threads), timerange, strings.TrimSpace(resizeThreshold), strings.TrimSpace(resizeSize))
	if err != nil {
		log.Fatalf("error parsing configs: %s", err)
	}

	// nolint here because of the uint8 conversion of directory depth
	// nolint:gosec
	ic, err := imageconvert.New(inputPath, processedLogFile, uint8(directoryDepth), configs...)
	if err != nil {
		log.Fatalf("error starting: %s", err)
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

func parseParams(compress, force, watch bool, threads uint8, timerange humantime.TimeRange, resizeThreshold, resizeSize string) ([]imageconvert.ConfigFunc, error) {

	var configs []imageconvert.ConfigFunc

	if compress {
		configs = append(configs, imageconvert.WithCompression(90))
	}

	if force {
		configs = append(configs, imageconvert.WithForce())
	}

	if watch {
		configs = append(configs, imageconvert.WithWatch())
	}

	if threads > 1 {
		if threads <= 0 || threads > uint8(runtime.GOMAXPROCS(0)) {
			return nil, errors.New(fmt.Sprintf("invalid number of threads: %d, min: 0, max: %d", threads, runtime.GOMAXPROCS(0)))
		}
		configs = append(configs, imageconvert.WithThreads(threads))
	}

	if timerange.From != imageconvert.NilTime || timerange.To != imageconvert.NilTime {
		configs = append(configs, imageconvert.WithTimeRange(timerange))
	}

	if resizeThreshold != "" {
		var thresholdArr = strings.Split(resizeThreshold, "x")
		if len(thresholdArr) != 2 {
			return nil, errors.New("resize threshold not in the format: [width]x[height] e.g. 230x400, input: " + resizeThreshold)
		}

		var sizeArr = strings.Split(resizeSize, "x")
		if len(sizeArr) != 2 {
			return nil, errors.New("resize size not in the format: [width]x[height] e.g. 230x400, input: " + resizeSize)
		}

		configs = append(configs, imageconvert.WithResize(getResizeValue(sizeArr[0]), getResizeValue(sizeArr[1]), getResizeValue(thresholdArr[0]), getResizeValue(thresholdArr[1])))
	}

	return configs, nil
}

func getResizeValue(str string) uint16 {

	var num, err = strconv.ParseUint(str, 10, 16)
	if err != nil {
		log.Fatalf("error resize value is not a number: '%s', err: %s", str, err.Error())
	}

	// nolint here because of the uint16 conversion
	// nolint:gosec
	return uint16(num)
}
