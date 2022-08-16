package imageconvert

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// FileInfo is a copy of fs.FileInfo that
// allows us to modify the filename
type FileInfo struct {
	ModTime time.Time
	Name    string
}

// ListAllFiles un-globs input as well as recursively list all
// files in the given input
func ListAllFiles(inputPath string) ([]FileInfo, error) {
	var allFiles []FileInfo
	var suffixRegex = regexp.MustCompile(".*.jpg$|.*.jpeg$|.*.png$|.*.webp$")

	// expand ~ paths
	if strings.Contains(inputPath, "~") {
		user, err := user.Current()
		if err != nil {
			return nil, fmt.Errorf("error getting current user, error: %s", err.Error())
		}
		inputPath = filepath.Join(user.HomeDir, strings.ReplaceAll(inputPath, "~", ""))
	}

	// try un-globing the input
	files, err := filepath.Glob(inputPath)
	if err != nil {
		return nil, fmt.Errorf("could not glob files: %w", err)
	}

	// go through the glob output and expand all dirs
	for _, file := range files {
		stat, err := os.Stat(file)
		if err != nil {
			return nil, fmt.Errorf("could not stat file: %s, err: %w", file, err)
		}

		if stat.IsDir() {
			dirFiles, err := ListDirFiles(file)
			if err != nil {
				return nil, fmt.Errorf("could not list files in dir: %s, err: %w", file, err)
			}
			allFiles = append(allFiles, dirFiles...)
		} else if suffixRegex.MatchString(strings.ToLower(file)) {
			allFiles = append(allFiles, FileInfo{Name: file, ModTime: stat.ModTime()})
		}
	}

	return allFiles, nil
}

// ListFiles lists every file in a directory (recursive) and
// optionally ignores files given in skipMap
func ListDirFiles(root string) ([]FileInfo, error) {
	var allFiles []FileInfo
	var files, err = os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("error listing all files in dir: %s, error: %s", root, err.Error())
	}

	var suffixRegex = regexp.MustCompile(".*.jpg$|.*.jpeg$|.*.png$|.*.webp$")

	for _, file := range files {
		var fullPath = filepath.Join(root, file.Name())
		if file.IsDir() {
			recursiveImages, err := ListDirFiles(fullPath)
			if err != nil {
				return nil, fmt.Errorf("error from recursive call to ListFiles, error: %s", err.Error())
			}
			allFiles = append(allFiles, recursiveImages...)
		} else if suffixRegex.MatchString(strings.ToLower(file.Name())) {
			var info, err = file.Info()
			if err != nil {
				return nil, fmt.Errorf("could not get FileInfo for file: %s, error: %s", file.Name(), err.Error())
			}
			allFiles = append(allFiles, FileInfo{Name: fullPath, ModTime: info.ModTime()})
		}
	}
	return allFiles, nil
}

// FileInfoToString converts a slice of fs.FileInfo to a slice
// of just the files names joined with a given root directory
func FileInfoToString(files []FileInfo) []string {
	var fileNames = make([]string, len(files))
	for i, file := range files {
		fileNames[i] = file.Name
	}
	return fileNames
}

// FilterFilesByDate removes files from the slice if they were modified
// before the modifiedSince
func FilterFilesByDate(files []FileInfo, modifiedSince time.Time) []FileInfo {
	for i := len(files) - 1; i >= 0; i-- {
		if files[i].ModTime.Before(modifiedSince) {
			files = remove(files, i)
		}
	}
	return files
}

// FilterFilesBySkipMap removes files from the map that are also in the skipMap
func FilterFilesBySkipMap(files []FileInfo, skipMap map[string]bool) []FileInfo {
	for i := len(files) - 1; i >= 0; i-- {
		if _, has := skipMap[files[i].Name]; has {
			files = remove(files, i)
		}
	}
	return files
}

// FilterPNG filters a slice of files to return only pngs
func FilerPNG(files []string) []string {
	var filtered []string
	for _, file := range files {
		if strings.HasSuffix(strings.ToLower(file), ".png") {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

// FilterWEBP filters a slice of files to return only webps
func FilerWEBP(files []string) []string {
	var filtered []string
	for _, file := range files {
		if strings.HasSuffix(strings.ToLower(file), ".webp") {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

// FilterJPG filters a slice of files to return only jpgs
func FilerJPG(files []string) []string {
	var filtered []string
	for _, file := range files {
		if strings.HasSuffix(strings.ToLower(file), ".jpg") {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

// FilterJPEG filters a slice of files to return only jpegs
func FilerJPEG(files []string) []string {
	var filtered []string
	for _, file := range files {
		if strings.HasSuffix(strings.ToLower(file), ".jpeg") {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

// EscapeFilePath escapes spaces in the filepath used for an exec() call
func EscapeFilePath(file string) string {
	var r = strings.NewReplacer(" ", `\ `, "(", `\(`, ")", `\)`, "'", `\'`, "&", `\&`, "@", `\@`)
	return r.Replace(file)
}

func remove(slice []FileInfo, s int) []FileInfo {
	return append(slice[:s], slice[s+1:]...)
}
