package imageconvert

import (
	"os"
	"testing"
	"time"

	"github.com/kmulvey/humantime"
	"github.com/stretchr/testify/assert"
)

func TestNewWithDefaults(t *testing.T) {
	t.Parallel()

	// setup
	var testdir = makeTestDir(t)
	var testImage = moveImage(t, testdir, testPair{Name: "./testimages/realjpg.jpg", Type: "jpeg"})

	var ic, err = NewWithDefaults(testImage, "", 0)
	assert.NoError(t, err)
	assert.Equal(t, uint8(1), ic.Threads)

	ic.WithCompression()
	assert.True(t, ic.Compress)

	ic.WithWatch()
	assert.True(t, ic.Watch)

	ic.WithForce()
	assert.True(t, ic.Force)

	ic.WithResize(200, 100, 300, 200)
	assert.Equal(t, uint16(200), ic.ResizeWidth)
	assert.Equal(t, uint16(300), ic.ResizeWidthThreshold)
	assert.Equal(t, uint16(100), ic.ResizeHeight)
	assert.Equal(t, uint16(200), ic.ResizeHeightThreshold)

	ic.WithThreads(3)
	assert.Equal(t, uint8(3), ic.Threads)

	ic.WithTimeRange(humantime.TimeRange{From: time.Time{}, To: time.Now()})
	assert.Equal(t, time.Time{}, ic.TimeRange.From)
	assert.Equal(t, time.Now().Day(), ic.TimeRange.To.Day())

	assert.NoError(t, os.RemoveAll(testdir))
}
