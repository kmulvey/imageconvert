package imageconvert

import (
	"io/ioutil"
	"path"
	"strings"
)

func ListFiles(root string) (map[string]bool, error) {
	var allFiles = make(map[string]bool)
	var staticBool bool
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return allFiles, err
	}
	for _, file := range files {
		if file.IsDir() {
			var subFiles, err = ListFiles(path.Join(root, file.Name()))
			if err != nil {
				return allFiles, err
			}
			for subFile := range subFiles {
				allFiles[subFile] = staticBool
			}
		} else {
			allFiles[path.Join(root, file.Name())] = staticBool
		}
	}
	return allFiles, nil
}

func FilerPNG(files map[string]bool) {
	for file := range files {
		if !strings.HasSuffix(file, ".png") {
			delete(files, file)
		}
	}
}

func FilerWEBP(files map[string]bool) {
	for file := range files {
		if !strings.HasSuffix(file, ".webp") {
			delete(files, file)
		}
	}
}

func FilerJPG(files map[string]bool) {
	for file := range files {
		if !strings.HasSuffix(file, ".jpg") || !strings.HasSuffix(file, ".jpeg") {
			delete(files, file)
		}
	}
}

func EscapeFilePath(file string) string {
	file = strings.ReplaceAll(file, " ", `\ `)
	return strings.ReplaceAll(file, "'", `\'`)
}
