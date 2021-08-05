package imageconvert

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func QualityCheck(maxQuality int, file string) bool {
	file = EscapeFilePath(file)
	cmd := fmt.Sprintf("identify -format %s %s", "'%Q'", file)
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		HandleErr(fmt.Sprintf("Incorect file name %s", file), err)
	}
	imageQuality, err := strconv.ParseInt(string(out), 10, 0)
	HandleErr("parse quality int", err)

	return int64(maxQuality) >= imageQuality
}

func CompressJPEG(quality int, imagePath string) {
	// have to escape the file spaces for the exec call
	var escapedImagePath = EscapeFilePath(imagePath)
	var cmdStr = fmt.Sprintf("jpegoptim -p -o -t -m%d %s",
		quality,
		escapedImagePath)

	output, err := exec.Command("bash", "-c", cmdStr).Output()
	if err != nil {
		log.Fatal("Compress exec fail: ", err.Error(), ", ", string(output))
	}
	HandleErr("Exec", err)

	if strings.Contains(string(output), "optimized.") {
		log.Info(string(output))
	}
}
