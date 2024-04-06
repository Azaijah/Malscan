package plugins

import (
	"fmt"
	"io/ioutil"
	"malscan/core/docker"
	"malscan/core/utils"
	"path/filepath"
	"strings"

	toml "github.com/pelletier/go-toml"
	log "github.com/sirupsen/logrus"
)

const (
	pluginFile = "plugins.toml"

	dectection = "av"
	enrichment = "er"
)

type Plugin struct {
	Enabled     bool   `toml:"enabled"`
	Name        string `toml:"name"`
	Description string `toml:"description"`
	Category    string `toml:"category"`
	Image       string `toml:"image"`
	Repository  string `toml:"repository"`
	Updatable   bool   `toml:"updatable"`
	Mime        string `toml:"mime"`
}

type PluginConfig struct {
	Plugins []Plugin `toml:"plugin"`
}

func (pconfig PluginConfig) PrintAll() {
	for _, plugin := range pconfig.Plugins {
		fmt.Println(plugin.Name)
	}
}

func (pconfig PluginConfig) Load() PluginConfig {

	pconfigPath := filepath.Join(utils.GetPlugDir(), pluginFile)

	data, err := ioutil.ReadFile(pconfigPath)

	if err != nil {
		data, err = Asset(pluginFile)
		ioutil.WriteFile(pconfigPath, data, 0644)
		log.WithFields(log.Fields{"err": err}).Warnf("Failed to read %s from disk", pluginFile)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Fatalf("Failed to read %s", pluginFile)
		}
	}

	err = toml.Unmarshal(data, &pconfig)

	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatalf("Failed to unmarshal %s", pluginFile)
	}

	return pconfig

}

//EnablePlugin - Responsible for enabling a plugin
func (pconfig PluginConfig) EnablePlugin(name string) {

	for _, plugin := range pconfig.Plugins {

		if plugin.Name == name {
			plugin.Enabled = true
		}
	}
}

//DisablePlugin - Responsible for disabling a plugin
func (pconfig PluginConfig) DisablePlugin(name string) {

	for _, plugin := range pconfig.Plugins {

		if plugin.Name == name {
			plugin.Enabled = false
		}
	}
}

func (pconfig PluginConfig) GetEnabledPlugins() (enabled []Plugin) {

	for _, plugin := range pconfig.Plugins {
		if plugin.Enabled == true {
			enabled = append(enabled, plugin)
		}
	}
	return enabled
}

func (pconfig PluginConfig) GetEnabledInstalledPlugins() (enabledInstalled []Plugin) {

	installed := docker.GetIntalledImages()

	for _, plugin := range pconfig.GetEnabledPlugins() {
		for _, installed := range installed {
			if strings.Contains(installed.RepoTags[0], plugin.Image) {
				enabledInstalled = append(enabledInstalled, plugin)
			}
		}

	}
	return enabledInstalled
}

func (pconfig PluginConfig) GetEnabledDectionPlugins() (enabled []Plugin) {

	for _, plugin := range pconfig.GetEnabledInstalledPlugins() {
		if plugin.Category == dectection {
			enabled = append(enabled, plugin)
		}
	}
	return enabled
}

func (pconfig PluginConfig) GetEnabledEnrichmentPlugins() (enabled []Plugin) {

	for _, plugin := range pconfig.GetEnabledInstalledPlugins() {
		if plugin.Category == enrichment {
			enabled = append(enabled, plugin)
		}
	}
	return enabled
}

/*
func InstallEnabledPlugins() {

	enabled := getEnabledPlugins()
	installed := getInstalledPlugins()

	if installed != nil { //Install all enabled not yet installed
		for _, en := range enabled {
			for _, instld := range installed {
				if en.Image != instld.Image {
					docker.PullImage(en.Image)
					break
				}
			}
		}

	} else { //If nothing is installed install all enabled
		for _, en := range enabled {
			docker.PullImage(en.Image)
		}

	}

}
*/
