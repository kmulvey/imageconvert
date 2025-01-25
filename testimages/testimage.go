package testimages

import (
	"embed"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/kmulvey/goutils"
	"github.com/stretchr/testify/assert"
)

//go:embed *.jpg
//go:embed *.png
//go:embed *.webp
var EmbededImages embed.FS

type TestCase struct {
	InputPath      string
	OutputPath     string
	ImageType      string
	ShouldConvert  bool
	WouldOverwrite bool
	Err            bool
}

var TestCases = []TestCase{
	{InputPath: "test.png", OutputPath: "test.jpg", ImageType: "png", ShouldConvert: true, WouldOverwrite: false, Err: false},
	{InputPath: "fakejpg.jpg", OutputPath: "fakejpg.jpg", ImageType: "png", ShouldConvert: true, WouldOverwrite: true, Err: false},
	{InputPath: "realjpg-portrait.jpg", OutputPath: "realjpg-portrait.jpg", ImageType: "jpeg", ShouldConvert: false, WouldOverwrite: true, Err: false},
	{InputPath: "testwebp.webp", OutputPath: "testwebp.jpg", ImageType: "webp", ShouldConvert: true, WouldOverwrite: false, Err: false},
	{InputPath: "realjpg.jpg", OutputPath: "realjpg.jpg", ImageType: "jpeg", ShouldConvert: false, WouldOverwrite: true, Err: false},
	{InputPath: "realjpg.png", OutputPath: "realjpg.png", ImageType: "png", ShouldConvert: true, WouldOverwrite: true, Err: true},
}

func MakeTestDir(t *testing.T) string {

	cwd, err := os.Getwd()
	assert.NoError(t, err)

	var testdir = "testdir_" + goutils.RandomString(5)
	testdir = filepath.Join(cwd, testdir)

	assert.NoError(t, os.MkdirAll(testdir, os.ModePerm))

	// copy the embedded images to the testdir
	entries, err := EmbededImages.ReadDir(".")
	assert.NoError(t, err)

	for _, entry := range entries {
		if !entry.IsDir() {
			file, err := EmbededImages.Open(entry.Name())
			assert.NoError(t, err)

			data, err := io.ReadAll(file)
			assert.NoError(t, err)

			assert.NoError(t, os.WriteFile(filepath.Join(testdir, entry.Name()), data, os.ModePerm))
			assert.NoError(t, file.Close())
		}
	}
	return testdir
}
