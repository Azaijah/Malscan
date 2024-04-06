package utils

import (
	"os"
	"path/filepath"

	"malscan/config"
	"malscan/core/utils"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

//Remove - Removes specified file
func Remove(filename *string) {

	if config.Values.Env.Filestore != "" {
		err := os.Remove(filepath.Join(config.Values.Env.Filestore, *filename))
		if err != nil {
			log.Error(errors.Wrap(err, "error while trying to scanned file: "+*filename))
		}
	} else {
		err := os.Remove(filepath.Join(utils.GetFilestoreDir(), *filename))
		if err != nil {
			log.Error(errors.Wrap(err, "error while trying to scanned file: "+*filename))
		}
	}
}
