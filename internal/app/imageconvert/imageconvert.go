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

// ImageConverter is the main config
type ImageConverter struct {
	Compress              bool
	Force                 bool
	Watch                 bool
	ResizeWidth           uint16
	ResizeWidthThreshold  uint16
	ResizeHeight          uint16
	ResizeHeightThreshold uint16
	InputEntry            path.Entry
	InputFiles            []path.Entry
	SkipMapEntry          path.Entry
	SkipMap               map[string]struct{}
	humantime.TimeRange
	ShutdownTrigger   chan struct{}
	ShutdownCompleted []chan struct{}
	///////////
	Quality        int
	Threads        int
	DeleteOriginal bool
}

// NewWithDefaults returns a new ImageConverter with conservative defaults. Use the WithX() functions to
// further configure.
func NewWithDefaults(inputPath, skipFile string, directoryDepth uint8) (ImageConverter, error) {

	var ic = ImageConverter{
		Threads:           1,
		ShutdownCompleted: make([]chan struct{}, 1),
	}
	var err error

	ic.InputEntry, err = path.NewEntry(inputPath, directoryDepth)
	if err != nil {
		return ic, err
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

	return ic, nil
}

// Shutdown gracefully closes all chans and quits.
func (ic *ImageConverter) Shutdown() {
	close(ic.ShutdownTrigger)
	<-goutils.MergeChannels(ic.ShutdownCompleted...)
}

// WithCompression will compress the images.
func (ic *ImageConverter) WithCompression() {
	ic.Compress = true
}

// WithForce will process files even if there are present in the skip file.
func (ic *ImageConverter) WithForce() {
	ic.Force = true
}

// WithResize resizes images down to a size given by width X height greater than a threshold
// given by widthThreshold X heightThreshold.
func (ic *ImageConverter) WithResize(width, height, widthThreshold, heightThreshold uint16) {
	ic.ResizeWidth = width
	ic.ResizeWidthThreshold = widthThreshold
	ic.ResizeHeight = height
	ic.ResizeHeightThreshold = heightThreshold
}

// WithWatch enables watching a directory for new or modified files.
func (ic *ImageConverter) WithWatch() {
	ic.Watch = true
}

// WithThreads specifies the number of CPU threads to use. The default is one but increacing this
// will significaltny improve performance epsically when compressing images. Pass a positive number
// of threads you wish to use, if 0 is passed, num cores - 1 will be set.
func (ic *ImageConverter) WithThreads(threads int) {
	if threads == 0 {
		ic.Threads = runtime.NumCPU() - 1
	} else {
		ic.Threads = threads
	}
	ic.ShutdownCompleted = make([]chan struct{}, ic.Threads)
}

// WithTimeRange will set a time range within images must have been last modified in order to be considered for processing.
func (ic *ImageConverter) WithTimeRange(tr humantime.TimeRange) {
	ic.TimeRange = tr
}
