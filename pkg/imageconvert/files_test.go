package imageconvert

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kmulvey/path"
	"github.com/stretchr/testify/assert"
)

func TestParseSkipMap(t *testing.T) {
	t.Parallel()

	// setup
	var testdir = makeTestDir(t)
	var handle, err = os.OpenFile(filepath.Join(testdir, "skipFile"), os.O_RDWR|os.O_CREATE, 0755)
	assert.NoError(t, err)
	_, err = handle.WriteString("realjpg.jpg")
	assert.NoError(t, err)
	err = handle.Close()
	assert.NoError(t, err)

	ic, err := NewWithDefaults(testdir, filepath.Join(testdir, "skipFile"), 0)
	assert.NoError(t, err)

	skipMap, err := ic.ParseSkipMap()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(skipMap))

	assert.NoError(t, os.RemoveAll(testdir))
}

func TestHasEOI(t *testing.T) {
	t.Parallel()

	// setup
	var testdir = makeTestDir(t)
	var testImage = moveImage(t, testdir, testPair{Name: "./testimages/realjpg.jpg", Type: "jpeg"})
	assert.True(t, hasEOI(testImage))
	assert.False(t, hasEOI("./compress.go"))
	assert.False(t, hasEOI("./doesnotexist"))
	assert.NoError(t, os.RemoveAll(testdir))
}

func TestWaitTilFileWritesComplete(t *testing.T) {
	t.Parallel()

	var fileAbs, err = filepath.Abs("./convert.go")
	assert.NoError(t, err)

	var DummyEntry = path.Entry{
		FileInfo:     nil,
		AbsolutePath: fileAbs,
		Children:     []path.Entry{},
	}

	var eventsIn = make(chan path.WatchEvent)
	var eventsOut = make(chan path.WatchEvent)

	go waitTilFileWritesComplete(eventsIn, eventsOut)

	var create = path.WatchEvent{Entry: DummyEntry, Op: 1}
	var write = path.WatchEvent{Entry: DummyEntry, Op: 2}

	for i := 0; i < 1000; i++ {
		if i == 0 {
			eventsIn <- create
		} else {
			eventsIn <- write
		}
	}

	go func() {
		for e := range eventsOut {
			assert.True(t, strings.HasSuffix(e.Entry.AbsolutePath, "watch.go"))
		}
	}()

	time.Sleep(time.Second * 2)
	close(eventsIn)
}

func TestEscapeFilePath(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "/some/file/\\&name", EscapeFilePath("/some/file/&name"))
	assert.Equal(t, "/some/file/\\(name\\)", EscapeFilePath("/some/file/(name)"))
}
