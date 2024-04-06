package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"malscan/config"
	"malscan/core/utils"
)

//GenerateFileSha256 - accepts file in filestore dir and returns shar256
func GenerateFileSha256(filename *string) (result string, err error) {

	var filePath string

	if config.Values.Env.Filestore == "" {
		filePath = filepath.Join(utils.GetFilestoreDir(), *filename)
	} else {
		filePath = filepath.Join(config.Values.Env.Filestore, *filename)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return
	}

	result = hex.EncodeToString(hash.Sum(nil))
	return
}

//GenerateFileMd5 - accepts file in filestore dir and returns md5
func GenerateFileMd5(filename *string) (result string, err error) {

	var filePath string

	if config.Values.Env.Filestore == "" {
		filePath = filepath.Join(utils.GetFilestoreDir(), *filename)
	} else {
		filePath = filepath.Join(config.Values.Env.Filestore, *filename)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return
	}

	result = hex.EncodeToString(hash.Sum(nil))
	return
}

//GenerateFileSha1 - accepts file in filestore dir and returns shar1
func GenerateFileSha1(filename *string) (result string, err error) {

	var filePath string

	if config.Values.Env.Filestore == "" {
		filePath = filepath.Join(utils.GetFilestoreDir(), *filename)
	} else {
		filePath = filepath.Join(config.Values.Env.Filestore, *filename)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	hash := sha1.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return
	}

	result = hex.EncodeToString(hash.Sum(nil))
	return
}
