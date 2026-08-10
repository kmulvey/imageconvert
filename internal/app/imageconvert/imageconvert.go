// Package imageconvert provides the core application logic for converting, compressing, and resizing images.
package imageconvert

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/kmulvey/goutils"
	"github.com/kmulvey/humantime"
	"github.com/kmulvey/path"
)

// ImageConverter is the main config.
type ImageConverter struct {
	humantime.TimeRange

	CompressQuality       uint8
	Force                 bool
	Watch                 bool
	ResizeWidth           uint16
	ResizeWidthThreshold  uint16
	ResizeHeight          uint16
	ResizeHeightThreshold uint16
	Threads               int
	InputEntry            path.Entry
	InputFiles            []path.Entry
	SkipMapEntry          path.Entry
	SkipMap               map[string]struct{}
	ShutdownTrigger       chan struct{}
	ShutdownCompleted     []chan struct{}
	configErrors          []error
}

// ConfigFunc is used to configure ImageConverter, see examples below.
type ConfigFunc func(*ImageConverter)

// New returns a new ImageConverter with conservative defaults. Use ConfigFunc functions to further configure.
// inputPath may be an empty string if WithDirectory or WithFiles is used instead.
func New(inputPath, skipFile string, directoryDepth uint8, configs ...ConfigFunc) (*ImageConverter, error) {

	var ic = &ImageConverter{
		Threads:           1,
		ShutdownCompleted: make([]chan struct{}, 1),
	}
	var err error

	if strings.TrimSpace(skipFile) == "" {
		skipFile = "processed.log"
	}

	handle, err := os.OpenFile(skipFile, os.O_RDWR|os.O_CREATE, 0600) //nolint:gosec // skipFile is a user-configured log path
	if err != nil {
		return ic, fmt.Errorf("error opening skip file: %s, err: %w", skipFile, err)
	}
	if err := handle.Close(); err != nil {
		return ic, fmt.Errorf("error closing handle to skip file: %s, err: %w", skipFile, err)
	}

	ic.SkipMapEntry, err = path.NewEntry(skipFile, 0)
	if err != nil {
		return ic, fmt.Errorf("error opening skip file: %s, err: %w", skipFile, err)
	}

	if strings.TrimSpace(inputPath) != "" {
		ic.InputEntry, err = path.NewEntry(inputPath, directoryDepth)
		if err != nil {
			return ic, fmt.Errorf("unable to create new entry for path: %s, err: %w", inputPath, err)
		}
	}

	// apply configs before getFileList so that WithTimeRange, WithDirectory, and WithFiles take effect
	for _, config := range configs {
		config(ic)
	}

	if len(ic.configErrors) > 0 {
		return ic, errors.Join(ic.configErrors...)
	}

	// only scan the directory if WithFiles has not already populated InputFiles
	if len(ic.InputFiles) == 0 {
		ic.InputFiles, err = ic.getFileList()
		if err != nil {
			return ic, err
		}
	}

	return ic, nil
}

// WithCompression will compress the images.
func WithCompression(quality uint8) func(*ImageConverter) {
	return func(ic *ImageConverter) {
		ic.CompressQuality = quality
	}
}

// WithForce will process files even if there are present in the skip file.
func WithForce() func(*ImageConverter) {
	return func(ic *ImageConverter) {
		ic.Force = true
	}
}

// WithResize resizes images down to a size given by width X height greater than a threshold
// given by widthThreshold X heightThreshold.
func WithResize(width, height, widthThreshold, heightThreshold uint16) func(*ImageConverter) {
	return func(ic *ImageConverter) {
		ic.ResizeWidth = width
		ic.ResizeWidthThreshold = widthThreshold
		ic.ResizeHeight = height
		ic.ResizeHeightThreshold = heightThreshold
	}
}

// WithWatch enables watching a directory for new or modified files.
func WithWatch() func(*ImageConverter) {
	return func(ic *ImageConverter) {
		ic.Watch = true
	}
}

// WithThreads specifies the number of CPU threads to use. The default is one but increacing this
// will significaltny improve performance epsically when compressing images. Pass a positive number
// of threads you wish to use, if 0 is passed, num cores - 1 will be set.
func WithThreads(threads int) func(*ImageConverter) {
	return func(ic *ImageConverter) {
		if threads == 0 {
			ic.Threads = runtime.NumCPU() - 1
		} else {
			ic.Threads = threads
		}
		ic.ShutdownCompleted = make([]chan struct{}, ic.Threads)
	}
}

// WithTimeRange will set a time range within images must have been last modified in order to be considered for processing.
func WithTimeRange(tr humantime.TimeRange) func(*ImageConverter) {
	return func(ic *ImageConverter) {
		ic.TimeRange = tr
	}
}

// WithDirectory sets a directory as the input source, overriding the inputPath argument passed to New.
// Use this as an alternative to passing inputPath directly when constructing via ConfigFuncs only.
func WithDirectory(dirPath string, depth uint8) ConfigFunc {
	return func(ic *ImageConverter) {
		entry, err := path.NewEntry(dirPath, depth)
		if err != nil {
			ic.configErrors = append(ic.configErrors, fmt.Errorf("unable to create entry for directory: %s, err: %w", dirPath, err))
			return
		}
		ic.InputEntry = entry
	}
}

// WithFiles sets an explicit list of files to process, bypassing directory scanning entirely.
// When this option is used the inputPath argument passed to New and any WithDirectory option are ignored.
func WithFiles(files ...string) ConfigFunc {
	return func(ic *ImageConverter) {
		for _, f := range files {
			entry, err := path.NewEntry(f, 0)
			if err != nil {
				ic.configErrors = append(ic.configErrors, fmt.Errorf("unable to create entry for file: %s, err: %w", f, err))
				return
			}
			ic.InputFiles = append(ic.InputFiles, entry)
		}
	}
}

// Shutdown gracefully closes all chans and quits.
func (ic *ImageConverter) Shutdown() {
	close(ic.ShutdownTrigger)
	<-goutils.MergeChannels(ic.ShutdownCompleted...)
}
