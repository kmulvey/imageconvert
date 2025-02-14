package imageconvert

import (
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
	CompressQuality       uint8
	Force                 bool
	Watch                 bool
	ResizeWidth           uint16
	ResizeWidthThreshold  uint16
	ResizeHeight          uint16
	ResizeHeightThreshold uint16
	Threads               uint8
	InputEntry            path.Entry
	InputFiles            []path.Entry
	SkipMapEntry          path.Entry
	SkipMap               map[string]struct{}
	humantime.TimeRange
	ShutdownTrigger   chan struct{}
	ShutdownCompleted []chan struct{}
}

// ConfigFunc is used to configure ImageConverter, see examples below.
type ConfigFunc func(*ImageConverter)

// New returns a new ImageConverter with conservative defaults. Use ConfigFunc functions to further configure.
func New(inputPath, skipFile string, directoryDepth uint8, configs ...ConfigFunc) (*ImageConverter, error) {

	var ic = &ImageConverter{
		Threads:           1,
		ShutdownCompleted: make([]chan struct{}, 1),
	}
	var err error

	ic.InputEntry, err = path.NewEntry(inputPath, directoryDepth)
	if err != nil {
		return ic, fmt.Errorf("unable to create new entry for path: %s, err: %w", inputPath, err)
	}

	if strings.TrimSpace(skipFile) == "" {
		skipFile = "processed.log"
	}

	handle, err := os.OpenFile(skipFile, os.O_RDWR|os.O_CREATE, 0755)
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

	ic.InputFiles, err = ic.getFileList()
	if err != nil {
		return ic, err
	}

	for _, config := range configs {
		config(ic)
	}

	return ic, nil
}

// AddCompression will compress the images.
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
func WithThreads(threads uint8) func(*ImageConverter) {
	return func(ic *ImageConverter) {
		if threads == 0 {
			ic.Threads = uint8(runtime.NumCPU() - 1)
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

// Shutdown gracefully closes all chans and quits.
func (ic *ImageConverter) Shutdown() {
	close(ic.ShutdownTrigger)
	<-goutils.MergeChannels(ic.ShutdownCompleted...)
}
