package imageconvert

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/Kagami/go-avif"
	"github.com/kmulvey/humantime"
	"github.com/kmulvey/path"
)

type ImageConverterConfig struct {
	OriginalImages []string
	WatchDir       string
	SkipMapFile    string
	Force          bool
	DeleteOriginal bool
	humantime.TimeRange
	ResizeWidth           uint16
	ResizeWidthThreshold  uint16
	ResizeHeight          uint16
	ResizeHeightThreshold uint16
	Quality               uint8
	Threads               uint8
	Depth                 uint8
}

// ImageConverter is the main config
type ImageConverter struct {
	OriginalImagesEntries []path.Entry
	WatchDir              string // currently we only support one
	SkipMap               map[string]struct{}
	SkipMapFileHandle     *os.File
	Force                 bool
	DeleteOriginal        bool
	humantime.TimeRange
	Depth                 uint8
	ResizeWidth           uint16
	ResizeWidthThreshold  uint16
	ResizeHeight          uint16
	ResizeHeightThreshold uint16
	Quality               int
	Threads               int
}

// NewImageConverter returns a new ImageConverter.
func NewImageConverter(config *ImageConverterConfig) (*ImageConverter, error) {

	// copy basic configs that do not need to be checked
	var ic = &ImageConverter{
		Force:          config.Force,
		DeleteOriginal: config.DeleteOriginal,
		TimeRange:      config.TimeRange,
	}
	var err error

	// file / dir list
	for _, original := range config.OriginalImages {
		entry, err := path.NewEntry(original, config.Depth, path.NewRegexEntitiesFilter(ImageExtensionRegex), path.NewDateEntitiesFilter(config.TimeRange.From, config.TimeRange.To))
		if err != nil {
			return nil, err
		}

		entries, err := entry.Flatten(true)
		if err != nil {
			return nil, err
		}

		ic.OriginalImagesEntries = append(ic.OriginalImagesEntries, entries...)
	}

	// watch
	if _, err := os.Stat(config.WatchDir); err != nil {
		return ic, fmt.Errorf("error opening watch dir: %s, err: %w", config.WatchDir, err)
	} else {
		ic.WatchDir = config.WatchDir
	}

	// skip map & processed log
	if strings.TrimSpace(config.SkipMapFile) == "" {
		config.SkipMapFile = "processed.log"
	}

	ic.SkipMapFileHandle, err = os.OpenFile(config.SkipMapFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return ic, fmt.Errorf("error opening skip file: %s, err: %w", config.SkipMapFile, err)
	}

	ic.SkipMap, err = ic.ParseSkipMap()
	if err != nil {
		return ic, fmt.Errorf("error creating skipmap from file: %s, err: %w", config.SkipMapFile, err)
	}

	// quality
	if config.Quality < avif.MinQuality || config.Quality > avif.MaxQuality {
		return nil, fmt.Errorf("quality: %d is not in range %d-%d", config.Quality, avif.MinQuality, avif.MaxQuality)
	} else {
		ic.Quality = int(config.Quality)
	}

	// threads
	if config.Threads < 0 || config.Threads > uint8(runtime.NumCPU()) {
		return nil, fmt.Errorf("threads: %d is not in range %d-%d", config.Threads, 0, runtime.NumCPU())
	} else if config.Threads == 0 {
		ic.Threads = runtime.NumCPU() - 1
	} else {
		ic.Threads = int(config.Threads)
	}

	// resize
	if config.ResizeWidth > config.ResizeWidthThreshold || config.ResizeHeight > config.ResizeHeightThreshold {
		return nil, errors.New("resize height and width must be less than resize height and width thresholds")
	}

	return ic, nil
}

// Shutdown gracefully closes all chans and quits.
func (ic *ImageConverter) Shutdown() error {
	return ic.SkipMapFileHandle.Close()
}
