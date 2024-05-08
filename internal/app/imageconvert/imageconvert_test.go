package imageconvert

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/kmulvey/humantime"
	"github.com/kmulvey/path"
	"github.com/stretchr/testify/assert"
)

func TestNewImageConverterGood(t *testing.T) {
	t.Parallel()

	var testdir = makeTestDir(t)
	var testImage = moveImage(t, testdir, testPair{Name: "./testimages/realjpg.jpg", Type: "jpeg"})
	var testImageEntry, err = path.NewEntry(testImage, 0)
	assert.NoError(t, err)

	var config = &ImageConverterConfig{
		OriginalImages: []string{testImage},
		Threads:        1,
		Quality:        30,
	}
	ic, err := NewImageConverter(config)
	assert.NoError(t, err)
	assert.Equal(t, []path.Entry{testImageEntry}, ic.OriginalImagesEntries)
	assert.Equal(t, "", ic.WatchDir)
	assert.Equal(t, 1, ic.Threads)
	assert.Equal(t, "processed.log", ic.SkipMapFileHandle.Name())
	assert.Equal(t, 0, len(ic.SkipMap))
	assert.Equal(t, 30, ic.Quality)
	assert.False(t, ic.Force)
	assert.False(t, ic.DeleteOriginal)
	assert.Equal(t, time.Time{}, ic.TimeRange.From)
	assert.Equal(t, time.Time{}, ic.TimeRange.To)
	assert.Equal(t, uint16(0), ic.ResizeWidth)
	assert.Equal(t, uint16(0), ic.ResizeHeight)
	assert.Equal(t, uint16(0), ic.ResizeWidthThreshold)
	assert.Equal(t, uint16(0), ic.ResizeHeightThreshold)

	assert.NoError(t, ic.Shutdown())
	assert.NoError(t, os.RemoveAll("processed.log"))

	// threads 0 & entire dir
	var testdirTwo = makeTestDir(t)
	ht, err := humantime.NewString2Time(time.UTC)
	assert.NoError(t, err)
	tr, err := ht.Parse("from May 8, 2009 5:57:51 PM to Sep 12, 2027 3:21:22 PM")
	assert.NoError(t, err)

	config = &ImageConverterConfig{
		Threads:        0,
		Depth:          2,
		Quality:        30,
		OriginalImages: []string{testdir, testdirTwo},
		TimeRange:      *tr,
	}
	ic, err = NewImageConverter(config)
	assert.NoError(t, err)
	assert.Equal(t, 10, len(ic.OriginalImagesEntries))
	assert.Equal(t, time.Date(2009, time.May, 8, 17, 57, 51, 0, time.UTC), ic.TimeRange.From)
	assert.Equal(t, time.Date(2027, time.September, 12, 15, 21, 22, 0, time.UTC), ic.TimeRange.To)
	assert.NoError(t, ic.Shutdown())

	assert.NoError(t, os.RemoveAll("processed.log"))
	assert.NoError(t, os.RemoveAll(testdir))
	assert.NoError(t, os.RemoveAll(testdirTwo))
}

func TestNewImageConverterBasicErrors(t *testing.T) {
	t.Parallel()

	var testdir = makeTestDir(t)
	var testImage = moveImage(t, testdir, testPair{Name: "./testimages/realjpg.jpg", Type: "jpeg"})

	// bad quality
	var config = &ImageConverterConfig{
		OriginalImages: []string{testImage},
		Threads:        1,
		Quality:        255,
	}
	_, err := NewImageConverter(config)
	assert.Error(t, err)
	assert.Equal(t, "quality: 255 is not in range 0-63", err.Error())

	// bad skipmap
	config = &ImageConverterConfig{
		OriginalImages: []string{testImage},
		SkipMapFile:    "/proc/kmsg",
		Threads:        2,
		Quality:        30,
	}
	_, err = NewImageConverter(config)
	assert.Error(t, err)
	assert.Equal(t, "error opening skip file: /proc/kmsg, err: open /proc/kmsg: permission denied", err.Error())

	// bad threads
	config = &ImageConverterConfig{
		OriginalImages: []string{testImage},
		Threads:        255,
		Quality:        30,
	}
	_, err = NewImageConverter(config)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "threads: 255 is not in range 0-"))

	// bad resize
	config = &ImageConverterConfig{
		OriginalImages:       []string{testImage},
		Threads:              2,
		Quality:              30,
		ResizeWidth:          300,
		ResizeWidthThreshold: 30,
	}
	_, err = NewImageConverter(config)
	assert.Error(t, err)
	assert.Equal(t, "resize height and width must be less than resize height and width thresholds", err.Error())

	assert.NoError(t, os.RemoveAll(testdir))
}
