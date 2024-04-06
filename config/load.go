package config

/*
Author: Liam Hellend
Email: liamhellend@gmail.com

Purpose: Contains data structures and functions related to loading malscans configuration

*/

import (
	"io/ioutil"
	"path/filepath"

	"malscan/core/utils"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

const (
	configFile = "config.toml"
)

//Values is the complete configuration variable that malscans configuration is loaded into
var Values tomlConfig

//tomlConfig is the complete malscan configuration
type tomlConfig struct {
	Env           env
	Logging       logging
	Alert         alert
	Elasticsearch elasticsearch
}

type env struct {
	Runtime     string `toml:"runtime"`
	Filestore   string `toml:"filestore"`
	CPUcores    int    `toml:"cpu_cores"`
	MaxFileProc int    `toml:"max_file_proc"`
	Client      string `toml:"client"`
	Site        string `toml:"site"`
	Network     string `toml:"network"`
}

type logging struct {
	Filename   string `toml:"filename"`
	MaxSize    int    `toml:"max_size"`
	MaxBackups int    `toml:"max_backups"`
	MaxAge     int    `toml:"max_age"`
	Compress   bool   `toml:"compress"`
}

type alert struct {
	LocalPath         string `toml:"local_path"`
	RemotePath        string `toml:"remote_path"`
	DynamicRemoteHost bool   `toml:"dynamic_remote_host"`
	RemoteHost        string `toml:"remote_host"`
	RemotePort        string `toml:"remote_port"`
	ScpPkey           string `toml:"scp_pkey"`
	ScpUser           string `toml:"scp_user"`
}

type elasticsearch struct {
	Enabled  bool   `toml:"enabled"`
	TLS      bool   `toml:"tls"`
	URL1     string `toml:"url1"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Cert     string `toml:"es_cert"`
	Key      string `toml:"es_key"`
	Ca       string `toml:"es_ca"`
}

func Load() {

	configPath := filepath.Join(utils.GetConfigDir(), configFile)

	log.Debug(configPath)

	data, err := ioutil.ReadFile(configPath)

	if err != nil {
		data, err = Asset(configFile)
		ioutil.WriteFile(configPath, data, 0644)
		log.WithFields(log.Fields{"err": err}).Warnf("Failed to read %s from disk", configFile)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Fatalf("Failed to read %s", configFile)
		}
	}
	err = toml.Unmarshal(data, &Values)

	log.Debug("config", Values)

	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatalf("Failed to unmarshal %s", configFile)
	}
}
