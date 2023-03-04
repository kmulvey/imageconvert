package imageconvert

import (
	"fmt"
	"os"
	"strings"

	"github.com/kmulvey/humantime"
	"github.com/kmulvey/path"
)

type ImageConverter struct {
	Compress              bool
	Force                 bool
	ResizeWidth           uint16
	ResizeWidthThreshold  uint16
	ResizeHeight          uint16
	ResizeHeightThreshold uint16
	Watch                 bool
	Threads               uint8
	InputEntry            path.Entry
	InputFiles            []path.Entry
	ProcessedLogFile      string
	SkipMapEntry          path.Entry
	SkipMap               map[string]struct{}
	humantime.TimeRange
}

func NewWithDefaults(inputPath, processedLogFile, skipMapFile string, directoryDepth uint8) (ImageConverter, error) {

	var ic = ImageConverter{
		Threads: 1,
	}
	var err error

	ic.InputEntry, err = path.NewEntry(inputPath, directoryDepth)
	if err != nil {
		return ic, err
	}

	ic.InputFiles, err = ic.getFileList()
	if err != nil {
		return ic, err
	}

	ic.SkipMap, err = ic.ParseSkipMap()
	if err != nil {
		return ic, err
	}

	if strings.TrimSpace(processedLogFile) != "" {
		handle, err := os.OpenFile(processedLogFile, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return ic, fmt.Errorf("error opening default log file: %s, err: %w", processedLogFile, err)
		}
		if err := handle.Close(); err != nil {
			return ic, fmt.Errorf("error closing handle to default log file: %s, err: %w", processedLogFile, err)
		}

		ic.ProcessedLogFile = processedLogFile
	}

	return ic, nil
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
}

func (ic *ImageConverter) WithTimeRange(tr humantime.TimeRange) {
	ic.TimeRange = tr
}
