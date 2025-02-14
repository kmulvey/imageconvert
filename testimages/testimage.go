package testimages

import (
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kmulvey/goutils"
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

// nolint:gochecknoglobals
var TestCases = []TestCase{
	{InputPath: "test.png", OutputPath: "test.jpg", ImageType: "png", ShouldConvert: true, WouldOverwrite: false, Err: false},
	{InputPath: "fakejpg.jpg", OutputPath: "fakejpg.jpg", ImageType: "png", ShouldConvert: true, WouldOverwrite: true, Err: false},
	{InputPath: "realjpg-portrait.jpg", OutputPath: "realjpg-portrait.jpg", ImageType: "jpeg", ShouldConvert: false, WouldOverwrite: true, Err: false},
	{InputPath: "testwebp.webp", OutputPath: "testwebp.jpg", ImageType: "webp", ShouldConvert: true, WouldOverwrite: false, Err: false},
	{InputPath: "realjpg.jpg", OutputPath: "realjpg.jpg", ImageType: "jpeg", ShouldConvert: false, WouldOverwrite: true, Err: false},
	{InputPath: "realjpg.png", OutputPath: "realjpg.png", ImageType: "png", ShouldConvert: true, WouldOverwrite: true, Err: true},
}

func MakeTestDir() (string, error) {

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	var testdir = filepath.Join(cwd, "testdir_"+goutils.RandomString(5))

	if err := os.MkdirAll(testdir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create test directory: %w", err)
	}

	// copy the embedded images to the testdir
	entries, err := EmbededImages.ReadDir(".")
	if err != nil {
		return "", fmt.Errorf("failed to read embedded images directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			file, err := EmbededImages.Open(entry.Name())
			if err != nil {
				return "", fmt.Errorf("failed to open embedded image %s: %w", entry.Name(), err)
			}
			data, err := io.ReadAll(file)
			if err != nil {
				return "", fmt.Errorf("failed to read embedded image %s: %w", entry.Name(), err)
			}
			if err := os.WriteFile(filepath.Join(testdir, entry.Name()), data, 0600); err != nil {
				return "", fmt.Errorf("failed to write embedded image %s to test directory: %w", entry.Name(), err)
			}
			if err := file.Close(); err != nil {
				return "", fmt.Errorf("failed to close embedded image %s: %w", entry.Name(), err)
			}
		}
	}

	return testdir, nil
}
