package imageconvert

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/hectane/go-acl"
	"github.com/kmulvey/imageconvert/v2/testimages"
	"github.com/kmulvey/path"
	"github.com/stretchr/testify/assert"
)

func TestParseSkipMap(t *testing.T) {
	t.Parallel()

	// setup
	var testdir = testimages.MakeTestDir(t)
	var handle, err = os.OpenFile(filepath.Join(testdir, "skipFile"), os.O_RDWR|os.O_CREATE, 0755)
	assert.NoError(t, err)
	_, err = handle.WriteString("realjpg.jpg")
	assert.NoError(t, err)
	err = handle.Close()
	assert.NoError(t, err)

	ic, err := New(testdir, filepath.Join(testdir, "skipFile"), 0)
	assert.NoError(t, err)

	skipMap, err := ic.ParseSkipMap()
	assert.NoError(t, err)
	assert.Len(t, skipMap, 1)

	ic.SkipMapEntry.AbsolutePath = ""
	skipMap, err = ic.ParseSkipMap()
	assert.NoError(t, err)
	assert.Empty(t, skipMap)

	// this is what will create the error in ParseSkipMap
	ic.SkipMapEntry.AbsolutePath = filepath.Join(testdir, "skipFile")
	assert.NoError(t, os.Chmod(ic.SkipMapEntry.AbsolutePath, 0000))
	assert.NoError(t, acl.Chmod(ic.SkipMapEntry.AbsolutePath, 0000)) // for windows
	skipMap, err = ic.ParseSkipMap()
	assert.Error(t, err)
	assert.Nil(t, skipMap)

	assert.NoError(t, acl.Chmod(ic.SkipMapEntry.AbsolutePath, fs.ModePerm)) // for windows
	assert.NoError(t, os.RemoveAll(testdir))
}

func TestGetFileList(t *testing.T) {
	t.Parallel()

	// setup
	var testdir = testimages.MakeTestDir(t)
	var handle, err = os.OpenFile(filepath.Join(testdir, "skipFile"), os.O_RDWR|os.O_CREATE, 0755)
	assert.NoError(t, err)
	_, err = handle.WriteString("realjpg.jpg")
	assert.NoError(t, err)
	err = handle.Close()
	assert.NoError(t, err)

	ic, err := New(testdir, filepath.Join(testdir, "skipFile"), 0)
	assert.NoError(t, err)

	// this is what will create the error in getFileList
	assert.NoError(t, os.Chmod(filepath.Join(testdir, "skipFile"), 0000))
	assert.NoError(t, acl.Chmod(filepath.Join(testdir, "skipFile"), 0000)) // for windows

	entries, err := ic.getFileList()
	assert.Error(t, err)
	assert.Empty(t, entries)

	assert.NoError(t, acl.Chmod(ic.SkipMapEntry.AbsolutePath, fs.ModePerm)) // for windows
	assert.NoError(t, os.RemoveAll(testdir))
}

func TestHasEOI(t *testing.T) {
	t.Parallel()

	// setup
	var testdir = testimages.MakeTestDir(t)
	var testImage = filepath.Join(testdir, "realjpg.jpg")

	var has, err = hasEOI(testImage)
	assert.NoError(t, err)
	assert.True(t, has)

	has, err = hasEOI("./files_test.go")
	assert.NoError(t, err)
	assert.False(t, has)

	has, err = hasEOI("./doesnotexist")
	assert.Error(t, err)
	assert.False(t, has)

	assert.NoError(t, os.RemoveAll(testdir))
}

func TestWaitTilFileWritesComplete(t *testing.T) {
	t.Parallel()

	// setup
	var testdir = testimages.MakeTestDir(t)
	var testImage = filepath.Join(testdir, "realjpg.jpg")
	var fileAbs, err = filepath.Abs(testImage)
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

	for i := range 1000 {
		if i == 0 {
			eventsIn <- create
		} else {
			eventsIn <- write
		}
	}

	go func() {
		for e := range eventsOut {
			assert.True(t, strings.HasSuffix(e.Entry.AbsolutePath, "realjpg.jpg"))
		}
	}()

	time.Sleep(time.Second * 2)
	close(eventsIn)
	assert.NoError(t, os.RemoveAll(testdir))
}
