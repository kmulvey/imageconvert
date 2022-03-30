package imageconvert

import (
	"fmt"
	"os/exec"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// QualityCheck uses imagemagick to determine the quality of the image
// and returns true if the quality is above a given threshold
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

// CompressJPEG uses jpegoptim to compress the image
func CompressJPEG(quality int, imagePath string) {
	// have to escape the file spaces for the exec call
	var escapedImagePath = EscapeFilePath(imagePath)
	var cmdStr = fmt.Sprintf("jpegoptim -p -o -m%d %s",
		quality,
		escapedImagePath)

	output, err := exec.Command("bash", "-c", cmdStr).Output()
	if err != nil {
		log.WithFields(log.Fields{
			"err":        err,
			"cmd output": string(output),
			"file":       imagePath,
		}).Fatal("Compress exec fail")
	}
	HandleErr("Exec", err)

	//if strings.Contains(string(output), "optimized.") {
	log.Info(string(output))
	//}
}
