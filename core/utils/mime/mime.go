package utils

import (
	"malscan/config"
	"malscan/core/utils"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
	log "github.com/sirupsen/logrus"
)

func FileType(filename string) (ftype string) {

	var fileDirOnDisk string

	if config.Values.Env.Filestore == "" {
		fileDirOnDisk = utils.GetFilestoreDir() //If filestore has not be set in the config file then use default filestore
	} else {
		fileDirOnDisk = config.Values.Env.Filestore //If filestore is set use the user provided path in config file
	}

	file := filepath.Join(fileDirOnDisk, filename)

	mime, err := mimetype.DetectFile(file)

	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("failed to detect file type")
	}

	return mime.String()

}
