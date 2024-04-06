package alert

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"malscan/config"
	"malscan/core/scp"
	"malscan/core/utils"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

/*
Author: Liam Hellend
Email: liamhellend@gmail.com

Purpose: Generate alerts for malscan

*/

//Generate - Responsible for generating alerts, accepts any struct and concurrently converts it to json and appends output to a file.

type alertStruct struct {
	Filename string `json:"Filename"`
	Result   string `json:"Result"`
}

//Generate - Responsible for generating an alert, accepts in bytes, parses and outputs the alert to json file
//Should only be called for av type plugins
func Generate(buf []byte, filename *string) {

	log.Debugf("malware detected: generating an alert for:%s", *filename)

	var fullHostname string
	if config.Values.Alert.DynamicRemoteHost == true {
		fullHostname = utils.ParseInstance(*filename) + "." + config.Values.Alert.RemoteHost + ":" + config.Values.Alert.RemotePort
	} else {
		fullHostname = config.Values.Alert.RemoteHost + ":" + config.Values.Alert.RemotePort
	}

	var avresult map[string]map[string]interface{}
	if err := json.Unmarshal(buf, &avresult); err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("unmarshaling av detection")
	}

	var alert alertStruct
	var ok bool
	alert.Filename = *filename
	alert.Result, ok = avresult["analysis"]["result"].(string)
	if !ok {
		log.WithFields(log.Fields{"err": ok}).Fatal("no result")
	}

	hdjson, err := json.Marshal(alert)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("marshaling")
	}

	var remotePath string
	var localPath string
	remotePath = config.Values.Alert.RemotePath
	if config.Values.Alert.LocalPath != "" {
		localPath = config.Values.Alert.LocalPath
	} else {
		localPath = filepath.Join(utils.GetLogsDir(), "alert.log")
	}

	f, err := os.OpenFile(localPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("opening alert file")
	}
	defer f.Close()

	nwritten, err := f.Write(hdjson)
	if err != nil {
		log.Error(errors.Wrap(err, "Error writing alert too alert file"))
	}
	log.Debugf("wrote %d characters to alert.log", nwritten)

	_, err = f.WriteString("\n")
	if err != nil {
		log.Error(errors.Wrap(err, "Error writing newline too alert file"))
	}

	time.Sleep(time.Millisecond * 100)
	log.Debugf("sending alert to remote location:%s:%s", fullHostname, remotePath)
	scp.Send(localPath, fullHostname, remotePath)

}
