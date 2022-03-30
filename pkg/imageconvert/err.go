package imageconvert

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// HandleErr is a convience func to inspect errors
// all errors in this app are fatal
func HandleErr(prefix string, err error) {
	if err != nil {
		log.Fatal(fmt.Errorf("%s: %w", prefix, err))
	}
}
