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

type ImageConverter struct {
	Compress              bool
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

func NewWithDefaults(inputPath, skipFile string, directoryDepth uint8) (ImageConverter, error) {

	var ic = ImageConverter{
		Threads:           uint8(runtime.NumCPU() - 1),
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

func (ic *ImageConverter) Shutdown() {
	close(ic.ShutdownTrigger)
	<-goutils.MergeChannels(ic.ShutdownCompleted...)
}

func (ic *ImageConverter) WithCompression() {
	ic.Compress = true
}

func (ic *ImageConverter) WithForce() {
	ic.Force = true
}

func (ic *ImageConverter) WithResize(width, height, widthThreshold, heightThreshold uint16) {
	ic.ResizeWidth = width
	ic.ResizeWidthThreshold = widthThreshold
	ic.ResizeHeight = height
	ic.ResizeHeightThreshold = heightThreshold
}

func (ic *ImageConverter) WithWatch() {
	ic.Watch = true
}

func (ic *ImageConverter) WithThreads(threads uint8) {
	ic.Threads = threads
	ic.ShutdownCompleted = make([]chan struct{}, threads)
}

func (ic *ImageConverter) WithTimeRange(tr humantime.TimeRange) {
	ic.TimeRange = tr
}
