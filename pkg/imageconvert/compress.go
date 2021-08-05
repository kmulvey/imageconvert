package imageconvert

import (
	"fmt"
	"os"
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
	var before, err = os.Stat(imagePath)
	HandleErr(fmt.Sprintf("Incorect file name %s", imagePath), err)

	// have to escape the file spaces for the exec call
	var escapedImagePath = EscapeFilePath(imagePath)
	var cmdStr = fmt.Sprintf("jpegoptim -o -m%d %s",
		quality,
		escapedImagePath)

	output, err := exec.Command("bash", "-c", cmdStr).Output()
	HandleErr("Exec", err)

	if strings.Contains(string(output), "skipped.") {
		log.Info(string(output))
	}

	after, err := os.Stat(imagePath)
	HandleErr("stat after", err)

	var afterSize = float64(after.Size())
	var beforeSize = float64(before.Size())
	log.WithFields(log.Fields{
		"file":  imagePath,
		"ratio": ((afterSize - beforeSize) / beforeSize) * 100,
	}).Info("Compress")
}
