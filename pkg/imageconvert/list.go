package imageconvert

import (
	"io/ioutil"
	"path"
	"regexp"
	"strings"
)

func ListFiles(root string, skipMap map[string]bool) []string {
	var allFiles []string
	files, err := ioutil.ReadDir(root)
	HandleErr("readdir", err)
	suffixRegex, err := regexp.Compile(".*.jpg$|.*.jpeg$|.*.png$|.*.webp$")
	HandleErr("regex", err)

	for _, file := range files {
		if file.IsDir() {
			var subFiles = ListFiles(path.Join(root, file.Name()), skipMap)

			for _, subFile := range subFiles {
				allFiles = append(allFiles, subFile)
			}
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

func FilerPNG(files []string) []string {
	var filtered []string
	for _, file := range files {
		if strings.HasSuffix(strings.ToLower(file), ".png") {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

func FilerWEBP(files []string) []string {
	var filtered []string
	for _, file := range files {
		if strings.HasSuffix(strings.ToLower(file), ".webp") {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

func FilerJPG(files []string) []string {
	var filtered []string
	for _, file := range files {
		if strings.HasSuffix(strings.ToLower(file), ".jpg") || strings.HasSuffix(strings.ToLower(file), ".jpeg") {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

func EscapeFilePath(file string) string {
	var r = strings.NewReplacer(" ", `\ `, "(", `\(`, ")", `\)`, "'", `\'`, "&", `\&`, "@", `\@`)
	return r.Replace(file)
}
