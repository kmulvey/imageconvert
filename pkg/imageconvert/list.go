package imageconvert

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// FileInfo is a copy of fs.FileInfo that
// allows us to modify the filename
type FileInfo struct {
	Name    string
	ModTime time.Time
}

// ListFiles lists every file in a directory (recursive) and
// optionally ignores files given in skipMap
func ListFiles(root string) ([]FileInfo, error) {
	var allFiles []FileInfo
	var files, err = ioutil.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("error listing all files in dir: %s, error: %s", root, err.Error())
	}

	suffixRegex, err := regexp.Compile(".*.jpg$|.*.jpeg$|.*.png$|.*.webp$")
	if err != nil {
		return nil, fmt.Errorf("error compiling regex, error: %s", err.Error())
	}

	for _, file := range files {
		var fullPath = filepath.Join(root, file.Name())
		if file.IsDir() {
			recursiveImages, err := ListFiles(fullPath)
			if err != nil {
				return nil, fmt.Errorf("error from recursive call to ListFiles, error: %s", err.Error())
			}
			allFiles = append(allFiles, recursiveImages...)
		} else {
			if suffixRegex.MatchString(strings.ToLower(file.Name())) {
				allFiles = append(allFiles, FileInfo{Name: fullPath, ModTime: file.ModTime()})
			}
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
