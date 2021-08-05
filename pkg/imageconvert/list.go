package imageconvert

import (
	"io/ioutil"
	"path"
	"strings"
)

func ListFiles(root string) []string {
	var allFiles []string
	files, err := ioutil.ReadDir(root)
	HandleErr("readdir", err)

	for _, file := range files {
		if file.IsDir() {
			var subFiles = ListFiles(path.Join(root, file.Name()))

			for _, subFile := range subFiles {
				allFiles = append(allFiles, subFile)
			}
		} else {
			allFiles = append(allFiles, path.Join(root, file.Name()))
		}
	}
	return allFiles
}

func FilerPNG(files []string) []string {
	var filtered []string
	for _, file := range files {
		if strings.HasSuffix(file, ".png") {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

func FilerWEBP(files []string) []string {
	var filtered []string
	for _, file := range files {
		if strings.HasSuffix(file, ".webp") {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

func FilerJPG(files []string) []string {
	var filtered []string
	for _, file := range files {
		if strings.HasSuffix(file, ".jpg") || strings.HasSuffix(file, ".jpeg") {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

func EscapeFilePath(file string) string {
	file = strings.ReplaceAll(file, " ", `\ `)
	return strings.ReplaceAll(file, "'", `\'`)
}
