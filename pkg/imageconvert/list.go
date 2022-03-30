package imageconvert

import (
	"io/ioutil"
	"path"
	"regexp"
	"strings"
)

// ListFiles lists every file in a directory (recursive) and
// optionally ignores files given in skipMap
func ListFiles(root string, skipMap map[string]bool) []string {
	var allFiles []string
	files, err := ioutil.ReadDir(root)
	HandleErr("readdir", err)
	suffixRegex, err := regexp.Compile(".*.jpg$|.*.jpeg$|.*.png$|.*.webp$")
	HandleErr("regex", err)

	for _, file := range files {
		if file.IsDir() {
			allFiles = append(allFiles, ListFiles(path.Join(root, file.Name()), skipMap)...)
		} else {
			if _, exists := skipMap[file.Name()]; !exists { // we dont process images that have already been processed
				if suffixRegex.MatchString(strings.ToLower(file.Name())) {
					allFiles = append(allFiles, path.Join(root, file.Name()))
				}
			}
		}
	}
	return allFiles
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
		if strings.HasSuffix(strings.ToLower(file), ".jpg") || strings.HasSuffix(strings.ToLower(file), ".jpeg") {
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
