package main

import (
	"context"
	"flag"
	"os"
	"sync"
	"time"

	"github.com/kmulvey/humantime"
	app "github.com/kmulvey/imageconvert/v2/internal/app/imageconvert"
	"github.com/kmulvey/imageconvert/v2/pkg/imageconvert"
	log "github.com/sirupsen/logrus"
	"go.szostok.io/version"
	"go.szostok.io/version/printer"
)

type quality struct {
	aboveThreshold bool
	currentQuality int
	inputPath      string
	err            error
}

// nolint: funlen
func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.TimeOnly,
	})

	var inputPath, processedLogFile string
	var qualityThreshold, directoryDepth, threads int
	var timerange humantime.TimeRange
	var v, h bool

	flag.StringVar(&inputPath, "path", "", "path to files, globbing must be quoted")
	flag.StringVar(&processedLogFile, "processed-file", "processed.log", "the file previously processed images were written to")
	flag.IntVar(&qualityThreshold, "quality", 90, "quality threshold: files at or above this quality will be reported as compressible")
	flag.IntVar(&directoryDepth, "depth", 1, "number of levels to search directories for images")
	flag.IntVar(&threads, "threads", 1, "number of threads to use for quality checking")
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

	log.Infof("Config: dir: %s, log file: %s, quality threshold: %d, modified-since: %s", inputPath, processedLogFile, qualityThreshold, timerange)

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
	var f, _ = os.Create("compresscheck.log")
	var wg = sync.WaitGroup{}
	var ctx = context.Background()
	var images = make(chan string)
	var results = make(chan quality)

	for range threads {
		wg.Add(1)
		go qualityCheck(ctx, &wg, qualityThreshold, images, results)
	}

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		for _, entry := range ic.InputFiles {
			images <- entry.AbsolutePath
		}
		close(images)
		wg.Done()
	}(&wg)

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if result.err != nil {
			log.Warnf("skipping %s: %s", result.inputPath, result.err)
			continue
		}
		if result.aboveThreshold {
			compressible++
			log.Infof("compressible: %s (current quality: %d, would save ~%d%%)", result.inputPath, result.currentQuality, result.currentQuality-qualityThreshold)
			if _, err := f.WriteString(result.inputPath + "\n"); err != nil {
				log.Warnf("error writing to log file: %s", err)
			}
		}
	}

	if err := f.Close(); err != nil {
		log.Warnf("error closing log file: %s", err)
	}
	log.Infof("Done: %d of %d files are compressible at quality threshold %d", compressible, len(ic.InputFiles), qualityThreshold)
}

func qualityCheck(ctx context.Context, wg *sync.WaitGroup, qualityThreshold int, images chan string, results chan quality) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case image, open := <-images:
			if !open {
				return
			}

			aboveThreshold, currentQuality, err := imageconvert.QualityCheck(qualityThreshold, image)
			if err != nil {
				results <- quality{
					inputPath: image,
					err:       err,
				}
			} else {
				results <- quality{
					inputPath:      image,
					aboveThreshold: aboveThreshold,
					currentQuality: currentQuality,
				}
			}
		}
	}
}
