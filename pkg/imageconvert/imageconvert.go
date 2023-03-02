package imageconvert

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/kmulvey/humantime"
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
	InputPath             string
	ProcessedLogFile      string
	SkipMapFile           string
	humantime.TimeRange
}

func NewWithDefaults(inputPath string) (ImageConverter, error) {

	if _, err := os.Stat(inputPath); errors.Is(err, fs.ErrNotExist) {
		return ImageConverter{}, fmt.Errorf("%s does not exist", inputPath)
	}

	return ImageConverter{
		Threads:          1,
		InputPath:        inputPath,
		ProcessedLogFile: "./processed.log",
	}, nil
}

func (ic ImageConverter) WithCompression() ImageConverter {
	ic.Compress = true
	return ic
}

func (ic ImageConverter) WithForce() ImageConverter {
	ic.Force = true
	return ic
}

func (ic ImageConverter) WithResize(width, height, widthThreshold, heightThreshold uint16) ImageConverter {
	ic.ResizeWidth = width
	ic.ResizeWidthThreshold = widthThreshold
	ic.ResizeHeight = height
	ic.ResizeHeightThreshold = heightThreshold
	return ic
}

func (ic ImageConverter) WithWatch() ImageConverter {
	ic.Watch = true
	return ic
}

func (ic ImageConverter) WithThreads(threads uint8) ImageConverter {
	ic.Threads = threads
	return ic
}

func (ic ImageConverter) WithProcessedLogFile(logFile string) ImageConverter {
	ic.ProcessedLogFile = logFile
	return ic
}

func (ic ImageConverter) WithSkipMap(skipFile string) ImageConverter {
	ic.SkipMapFile = skipFile
	return ic
}

func (ic ImageConverter) WithTimeRange(from, to time.Time) ImageConverter {
	ic.TimeRange = humantime.TimeRange{From: from, To: to}
	return ic
}
