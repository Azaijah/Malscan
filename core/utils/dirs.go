package utils

import (
	"os"

	"path/filepath"

	log "github.com/sirupsen/logrus"
)

/*
Author: Liam Hellend
Email: liamhellend@gmail.com

Purpose: Contains functions related creating malscans directories

*/

//GetBaseDir - helper function to get base dir for malscan
func GetBaseDir() string {

	ex, err := os.Executable()
	if err != nil {
		log.Panic(err)
	}

	base := filepath.Join(filepath.Dir(ex), ".malscan")

	return base
}

//GetConfigDir - helper function to get config dir
func GetConfigDir() string {

	return filepath.Join(GetBaseDir(), "config")
}

//GetLogsDir - helper function to get log dir
func GetLogsDir() string {

	return filepath.Join(GetBaseDir(), "logs")
}

//GetPlugDir - helper function to get plugin dir
func GetPlugDir() string {

	return filepath.Join(GetBaseDir(), "plugins")
}

//GetFilestoreDir - helper function to get filestore dir
func GetFilestoreDir() string {

	return filepath.Join(GetBaseDir(), "filestore")
}

//MakeDirs - Responsible for creating malscan dirs is they don't exist already
func MakeDirs() {

	log.Debug("creating malscan directories if they don't exist")

	if _, err := os.Stat(GetBaseDir()); os.IsNotExist(err) {
		os.MkdirAll(GetBaseDir(), 0777)
		log.Debug("creating base directory for malscan")
	}
	if _, err := os.Stat(GetConfigDir()); os.IsNotExist(err) {
		os.MkdirAll(GetConfigDir(), 0777)
		log.Debug("creating config directory for malscan")
	}
	if _, err := os.Stat(GetLogsDir()); os.IsNotExist(err) {
		os.MkdirAll(GetLogsDir(), 0777)
		log.Debug("creating log directory for malscan")
	}
	if _, err := os.Stat(GetPlugDir()); os.IsNotExist(err) {
		os.MkdirAll(GetPlugDir(), 0777)
		log.Debug("creating plugins directory for malscan")
	}
	if _, err := os.Stat(GetFilestoreDir()); os.IsNotExist(err) {
		os.MkdirAll(GetFilestoreDir(), 0777)
		log.Debug("creating filestore directory for malscan")
	}
}
