package logger

import (
	"os"
	"path/filepath"

	"malscan/config"
	"malscan/core/utils"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

/*
Author: Liam Hellend
Email: liamhellend@gmail.com

Purpose: Contains functions related to different log settings

*/

//DebugLog - Logs everything to the console
//This function can be called to override the log settings in the malscan config
func DebugLog() {
	logrus.Debug("initializing - debug logging")
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)
}

//DevLog - Logs info, erros and fatal errors to the console
func DevLog() {
	logrus.Debug("initializing - development logging")
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout)
}

//ProdLog - only logs errors and fatal errors
//This function can be called to override the log settings in the malscan config
func ProdLog() {

	var logpath string

	if config.Values.Logging.Filename == "" {
		logpath = filepath.Join(utils.GetLogsDir(), "malscan.log")

	} else {
		logpath = config.Values.Logging.Filename
	}

	logrus.Debug("initializing - production logging")
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(&lumberjack.Logger{
		Filename:   logpath,
		MaxSize:    config.Values.Logging.MaxSize,
		MaxBackups: config.Values.Logging.MaxBackups,
		MaxAge:     config.Values.Logging.MaxAge,
		Compress:   config.Values.Logging.Compress,
	})

}

func Load() {

	switch config.Values.Env.Runtime {
	case "debug":
		DebugLog()
	case "dev":
		DevLog()
	case "prod":
		ProdLog()
	default:
		log.Errorf("failed to intialize:%s:logging", config.Values.Env.Runtime)
	}

}
