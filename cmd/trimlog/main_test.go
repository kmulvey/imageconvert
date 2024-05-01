package main

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangeFileName(t *testing.T) {
	t.Parallel()

	var oldFile, err = os.Create("./old.log")
	assert.NoError(t, err)
	oldFile.WriteString("main.go\n")
	oldFile.WriteString("main_test.go\n")
	oldFile.WriteString("noexist.go\n")
	oldFile.WriteString("no-exist.go\n")
	assert.NoError(t, oldFile.Close())

	assert.NoError(t, cleanLogFile(oldFile.Name(), "new.log"))

	newFile, err := os.OpenFile("new.log", os.O_RDONLY, 0755)
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

	assert.Error(t, cleanLogFile("noexist.log", "new.log"))
	assert.Error(t, cleanLogFile(oldFile.Name(), "./nodir/new.log"))

	assert.NoError(t, os.RemoveAll("./old.log"))
	assert.NoError(t, os.RemoveAll("./new.log"))
}
