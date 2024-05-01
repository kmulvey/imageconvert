package main

import (
	"bufio"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangeFileName(t *testing.T) {
	t.Parallel()

	var oldFile, err = os.Create("old.log")
	assert.NoError(t, err)
	_, err = oldFile.WriteString("main.go\n")
	assert.NoError(t, err)
	_, err = oldFile.WriteString("main_test.go\n")
	assert.NoError(t, err)
	_, err = oldFile.WriteString("noexist.go\n")
	assert.NoError(t, err)
	_, err = oldFile.WriteString("no-exist.go\n")
	assert.NoError(t, err)
	assert.NoError(t, oldFile.Close())

	assert.NoError(t, cleanLogFile(oldFile.Name(), "new.log"))
	assert.Error(t, cleanLogFile("noexist.log", "new.log"))
	assert.Error(t, cleanLogFile(oldFile.Name(), "./nodir/new.log"))

	newFile, err := os.OpenFile("new.log", os.O_RDONLY, 0755)
	assert.NoError(t, err)
	var fileScanner = bufio.NewScanner(newFile)
	fileScanner.Split(bufio.ScanLines)

	var expectedFileData = map[string]struct{}{
		"main.go":      {},
		"main_test.go": {},
	}
	for fileScanner.Scan() {
		var _, found = expectedFileData[fileScanner.Text()]
		assert.True(t, found)
		delete(expectedFileData, fileScanner.Text())
	}
	assert.Equal(t, 0, len(expectedFileData))
	assert.NoError(t, newFile.Close())

	assert.NoError(t, os.RemoveAll(oldFile.Name()))
	if runtime.GOOS != "windows" { // windows cant do anything right
		assert.NoError(t, os.RemoveAll(newFile.Name()))
	}
}
