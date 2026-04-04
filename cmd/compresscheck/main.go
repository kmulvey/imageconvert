package main

import (
	"flag"
	"os"
	"time"

	"github.com/kmulvey/humantime"
	app "github.com/kmulvey/imageconvert/v2/internal/app/imageconvert"
	"github.com/kmulvey/imageconvert/v2/pkg/imageconvert"
	log "github.com/sirupsen/logrus"
	"go.szostok.io/version"
	"go.szostok.io/version/printer"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.TimeOnly,
	})

	var inputPath, processedLogFile string
	var quality, directoryDepth int
	var timerange humantime.TimeRange
	var v, h bool

	flag.StringVar(&inputPath, "path", "", "path to files, globbing must be quoted")
	flag.StringVar(&processedLogFile, "processed-file", "processed.log", "the file previously processed images were written to")
	flag.IntVar(&quality, "quality", 90, "quality threshold: files at or above this quality will be reported as compressible")
	flag.IntVar(&directoryDepth, "depth", 1, "number of levels to search directories for images")
	flag.Var(&timerange, "time-range", "only consider files changed within this time range")
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

	log.Infof("Config: dir: %s, log file: %s, quality threshold: %d, modified-since: %s", inputPath, processedLogFile, quality, timerange)

	var configs []app.ConfigFunc
	if timerange.From != app.NilTime || timerange.To != app.NilTime {
		configs = append(configs, app.WithTimeRange(timerange))
	}

	// nolint:gosec // directoryDepth is bounded by flag input
	ic, err := app.New(inputPath, processedLogFile, uint8(directoryDepth), configs...)
	if err != nil {
		log.Fatalf("error initializing: %s", err)
	}

	var compressible int
	for _, entry := range ic.InputFiles {
		aboveThreshold, currentQuality, err := imageconvert.QualityCheck(quality, entry.AbsolutePath)
		if err != nil {
			log.Warnf("skipping %s: %s", entry.AbsolutePath, err)
			continue
		}
		if aboveThreshold {
			compressible++
			log.Infof("compressible: %s (current quality: %d, would save ~%d%%)", entry.AbsolutePath, currentQuality, currentQuality-quality)
		}
	}

	log.Infof("Done: %d of %d files are compressible at quality threshold %d", compressible, len(ic.InputFiles), quality)
}
