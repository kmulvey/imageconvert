package main

import (
	"os"
	"strings"
	"testing"

	"github.com/kmulvey/path"
	"github.com/stretchr/testify/assert"
)

func TestValidCharacter(t *testing.T) {
	t.Parallel()

	assert.True(t, validCharacter('a'))
	assert.True(t, validCharacter('2'))
	assert.True(t, validCharacter('R'))
	assert.True(t, validCharacter('-'))
	assert.True(t, validCharacter('_'))
	assert.False(t, validCharacter('%'))
	assert.False(t, validCharacter(':'))
	assert.False(t, validCharacter('世'))
}

func TestChangeFileName(t *testing.T) {
	t.Parallel()

	var entry = path.Entry{AbsolutePath: "/here/goodname.jpg"}
	var newName, changed = changeFileName(entry)
	assert.False(t, changed)
	assert.Equal(t, "goodname", newName)

	entry = path.Entry{AbsolutePath: "/here/bad&name.jpg"}
	newName, changed = changeFileName(entry)
	assert.True(t, changed)
	assert.True(t, strings.HasPrefix(newName, "bad"))
	assert.True(t, strings.HasSuffix(newName, "name"))
}

func TestRenameFiles(t *testing.T) {
	t.Parallel()

	assert.NoError(t, os.Mkdir("testdir", os.ModePerm))

	var goodFile, err = os.Create("./testdir/goodname.jpg")
	assert.NoError(t, err)
	assert.NoError(t, renameFiles("testdir"))

	badFile, err := os.Create("./testdir/bad$name.jpg")
	assert.NoError(t, err)
	assert.NoError(t, renameFiles("testdir"))

	// already exists collision
	_, err = os.Create("./testdir/bad$name.jpg")
	assert.NoError(t, err)
	assert.NoError(t, renameFiles("testdir"))

	assert.Error(t, renameFiles("noexistdir"))

	assert.NoError(t, goodFile.Close())
	assert.NoError(t, badFile.Close())
	assert.NoError(t, os.RemoveAll("testdir"))
}

func TestMoveFile(t *testing.T) {
	t.Parallel()
	assert.Error(t, moveFile("/does/not/exist/", "testdir/goodname.jpg"))
	assert.Error(t, moveFile("main.go", "/does/not/exist/"))
}
