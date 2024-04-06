package plugins

import (
	"encoding/json"
	"strings"
	"sync"
	"time"

	malscanconfig "malscan/config"
	"malscan/core/alert"
	"malscan/core/docker"
	"malscan/core/utils"
	hash "malscan/core/utils/hash"
	"malscan/elastic"
	"malscan/structs"

	log "github.com/sirupsen/logrus"
)

/*
Author: Liam Hellend
Email: liamhellend@gmail.com

Purpose: Contains functions related running malscan plugins

*/

//RunPlugin - Responsible for running a single plugin passed into the function against a file
func RunPlugin(filename *string, image string) {
	log.Debug("RunPlugin - running plugin")
	docker.RunContainerOnFile(image, filename)
}

//RunPluginUpdate - Responsible for running a update on a single plugin
func (pconfig PluginConfig) RunPluginUpdate(plugName string) (msg string) {
	log.Debug("RunPlugin - running plugin update")

	enabled := pconfig.GetEnabledPlugins()

	found := false

	for _, plugin := range enabled {

		if plugin.Name == plugName {
			found = true
			if plugin.Updatable == true {
				status := docker.RunContainerUpdate(plugin.Image)
				msg = plugin.Name + " Update: " + status
				break
			} else {
				log.Debug("The plugin specified cannot be updated")
				msg = "The plugin you have specified cannot be updated"
				break
			}
		}

	}

	if found != true {
		log.Debug("Could not update: " + plugName + ", either the plugin does not exist or it is not enabled")
		msg = "Could not update: " + plugName + ", either the plugin does not exist or it is not enabled"
	}

	//docker.Prune() (NOT SAFE TO USE) - containers now removed invidually in container.go

	return msg
}

//RunPluginUpdateAll - Responsible for attempting to update all enabled plugins
func (pconfig PluginConfig) RunPluginUpdateAll() (statusAndTime map[string]map[string]string) {

	log.Debug("RunPlugin - running plugin update for all plugins")

	statusAndTime = make(map[string]map[string]string)
	status := make(map[string]string)
	var msg string

	enabled := pconfig.GetEnabledPlugins()

	for _, plugin := range enabled {

		if plugin.Updatable == true {

			msg = docker.RunContainerUpdate(plugin.Image)

			status[plugin.Name] = msg

			log.Debug(plugin.Name + " Update: " + msg)

		}
	}

	//docker.Prune() (NOT SAFE TO USE) - containers now removed invidually in container.go

	statusAndTime[time.Now().Format(time.RFC3339)] = status

	return statusAndTime
}

//RunEnabledConcurrent - Responsible for running all plugins against a file, except each plugin is run concurrently
func (pconfig PluginConfig) RunEnabledConcurrent(filename string) {

	log.Debug("running enabled plugins concurrently")

	fileReport := structs.FullFileReport{}

	//Set default/static values
	if malscanconfig.Values.Alert.DynamicRemoteHost == true {
		fileReport.File.Name = strings.Replace(filename, utils.ParseInstance(filename), "", -1)
	} else {
		fileReport.File.Name = filename
	}

	fileReport.File.Sha1, _ = hash.GenerateFileSha1(&filename)
	fileReport.File.Md5, _ = hash.GenerateFileMd5(&filename)
	fileReport.File.Malware.Infected = false

	//Initialize maps
	fileReport.File.Malware.Analyzers.RawAnalysis.AntiVirus = make(map[string]json.RawMessage)
	fileReport.File.Malware.Analyzers.RawAnalysis.Enricher = make(map[string]json.RawMessage)

	var pluginsUsed []string     //Stores plugins used in the analysis
	var pluginsDetected []string //Stores plugins that detected malware

	detected := false //Used to know if there has been a detection

	var mutex = &sync.Mutex{} //Used so only one av can generate an alert

	enabledAV := pconfig.GetEnabledDectionPlugins() //Stores enabled av plugins

	var wgAV sync.WaitGroup //Used to wait for all av plugins to finish before posting results to es
	var wgER sync.WaitGroup //Used to wait for all er plugins to finish before posting results to es

	wgAV.Add(len(enabledAV)) //Tell waitgroup how many av plugins to wait for (er is done in the run er plugins function)

	for _, plugin := range enabledAV { //Range over all enabled av plugins

		pluginsUsed = append(pluginsUsed, plugin.Name) //Append used plugin

		go func(plugin Plugin, detected *bool) { //Detach running of plugin as goroutine

			dockerOutput := docker.RunContainerOnFile(plugin.Image, &filename) //Run plugin container

			//If there has been no detection test if the current av has detected malware
			mutex.Lock()
			if *detected == false {
				if testInfected(dockerOutput, &plugin.Name) == true {
					go alert.Generate(dockerOutput, &filename) //If there is an av detection generate an alert detached as goroutine
					*detected = true
					fileReport.File.Malware.Infected = true

					go func() { pconfig.RunEnricherPlugins(&filename, &fileReport, &pluginsUsed, &wgER) }() //Run er plugins detached as goroutine

				}

			}
			mutex.Unlock()

			//If there plugin deteced malware add to the list of plugins that detected malware
			if testInfected(dockerOutput, &plugin.Name) == true {
				pluginsDetected = append(pluginsDetected, plugin.Name)
			}

			//Parse raw results
			parsedResult := parseResult(dockerOutput, plugin.Name)
			fileReport.File.Malware.Analyzers.RawAnalysis.AntiVirus[plugin.Name] = parsedResult

			//Parse and add results to list of av results
			analysisResult := parseAnalysisResult(dockerOutput, &plugin.Name)
			if analysisResult != "" {
				fileReport.File.Malware.Results = append(fileReport.File.Malware.Results, analysisResult)
			}

			wgAV.Done() //Inform wait group this plugin has finished

		}(plugin, &detected)
	}

	//Wait for plugins to finish
	wgAV.Wait()
	time.Sleep(time.Second * 5)
	wgER.Wait()

	//Set plugins that detected malware
	fileReport.File.Malware.Analyzers.Names = pluginsDetected

	//Set tags
	fileReport.File.Tags = append(fileReport.File.Tags, malscanconfig.Values.Env.Client)
	fileReport.File.Tags = append(fileReport.File.Tags, malscanconfig.Values.Env.Site)
	fileReport.File.Tags = append(fileReport.File.Tags, malscanconfig.Values.Env.Network)
	if malscanconfig.Values.Alert.DynamicRemoteHost == true {
		fileReport.File.Tags = append(fileReport.File.Tags, "sen"+string(utils.ParseInstance(filename)[7]))
	}

	fileReport.File.Tags = append(fileReport.File.Tags, "malscan")

	//Set timestamp of scan
	fileReport.File.Date = time.Now().Format(time.RFC3339)

	log.Infof("analyzed:%s:with:%s:infected:%t", filename, strings.Join(pluginsUsed, ","), detected)

	//docker.Prune() //Clean docker system (NOT SAFE TO USE) - containers now removed invidually in container.go
	// no but seriously using this could make a lot of people mad

	if malscanconfig.Values.Elasticsearch.Enabled == true {
		elastic.Index(fileReport, &filename) //Post es results
	}

}

//RunEnabledConcurrentWithChannel - Responsible for running all plugins against a file, except each plugin is run concurrently
//Includes a channel for recieving a signal when the function has completed
func (pconfig PluginConfig) RunEnabledConcurrentWithChannel(filename string, done chan bool) {

	log.Debug("Running enabled plugins concurrently")

	pconfig.RunEnabledConcurrent(filename)

	done <- true

}

func (pconfig PluginConfig) RunEnricherPlugins(filename *string, fileReport *structs.FullFileReport, pluginsUsed *[]string, wgER *sync.WaitGroup) {

	enabledER := pconfig.GetEnabledEnrichmentPlugins()

	wgER.Add(len(enabledER))

	for _, plugin := range enabledER {

		*pluginsUsed = append(*pluginsUsed, plugin.Name)

		go func(plugin Plugin) {

			dockerOutput := docker.RunContainerOnFile(plugin.Image, filename)

			parsedResult := parseResult(dockerOutput, plugin.Name)
			fileReport.File.Malware.Analyzers.RawAnalysis.Enricher[plugin.Name] = parsedResult

			wgER.Done()

		}(plugin)
	}
}

//RunEnabled - Responsible for running all plugins against a file, except each plugin is run concurrently
func (pconfig PluginConfig) RunEnabled(filename string) {

	log.Debug("running enabled plugins")

	fileReport := structs.FullFileReport{}

	//Set default/static values
	if malscanconfig.Values.Alert.DynamicRemoteHost == true {
		fileReport.File.Name = strings.Replace(filename, utils.ParseInstance(filename), "", -1)
	} else {
		fileReport.File.Name = filename
	}

	fileReport.File.Sha1, _ = hash.GenerateFileSha1(&filename)
	fileReport.File.Md5, _ = hash.GenerateFileMd5(&filename)
	fileReport.File.Malware.Infected = false

	//Initialize maps
	fileReport.File.Malware.Analyzers.RawAnalysis.AntiVirus = make(map[string]json.RawMessage)
	fileReport.File.Malware.Analyzers.RawAnalysis.Enricher = make(map[string]json.RawMessage)

	var pluginsUsed []string     //Stores plugins used in the analysis
	var pluginsDetected []string //Stores plugins that detected malware

	detected := false //Used to know if there has been a detection

	enabledAV := pconfig.GetEnabledDectionPlugins() //Stores enabled av plugins

	var wgER sync.WaitGroup //Used to wait for all er plugins to finish before posting results to es

	for _, plugin := range enabledAV { //Range over all enabled av plugins

		pluginsUsed = append(pluginsUsed, plugin.Name) //Append used plugin

		dockerOutput := docker.RunContainerOnFile(plugin.Image, &filename) //Run plugin container

		//If there has been no detection test if the current av has detected malware

		if detected == false {
			if testInfected(dockerOutput, &plugin.Name) == true {
				go alert.Generate(dockerOutput, &filename) //If there is an av detection generate an alert detached as goroutine
				detected = true
				fileReport.File.Malware.Infected = true

				go func() { pconfig.RunEnricherPlugins(&filename, &fileReport, &pluginsUsed, &wgER) }() //Run er plugins detached as goroutine

			}

		}

		//If a plugin deteced malware add to the list of plugins that detected malware
		if testInfected(dockerOutput, &plugin.Name) == true {
			pluginsDetected = append(pluginsDetected, plugin.Name)
		}

		//Parse raw results
		parsedResult := parseResult(dockerOutput, plugin.Name)
		fileReport.File.Malware.Analyzers.RawAnalysis.AntiVirus[plugin.Name] = parsedResult

		//Parse and add results to list of av results
		analysisResult := parseAnalysisResult(dockerOutput, &plugin.Name)
		if analysisResult != "" {
			fileReport.File.Malware.Results = append(fileReport.File.Malware.Results, analysisResult)
		}

	}

	time.Sleep(time.Second * 5)
	wgER.Wait()

	//Set plugins that detected malware
	fileReport.File.Malware.Analyzers.Names = pluginsDetected

	//Set tags
	fileReport.File.Tags = append(fileReport.File.Tags, malscanconfig.Values.Env.Client)
	fileReport.File.Tags = append(fileReport.File.Tags, malscanconfig.Values.Env.Site)
	fileReport.File.Tags = append(fileReport.File.Tags, malscanconfig.Values.Env.Network)
	if malscanconfig.Values.Alert.DynamicRemoteHost == true {
		fileReport.File.Tags = append(fileReport.File.Tags, "sen"+string(utils.ParseInstance(filename)[7]))
	}

	fileReport.File.Tags = append(fileReport.File.Tags, "malscan")

	//Set timestamp of scan
	fileReport.File.Date = time.Now().Format(time.RFC3339)

	log.Infof("analyzed:%s:with:%s:infected:%t", filename, strings.Join(pluginsUsed, ","), detected)

	//docker.Prune() //Clean docker system (NOT SAFE TO USE) - containers now removed invidually in container.go
	// no but seriously using this could make a lot of people mad

	if malscanconfig.Values.Elasticsearch.Enabled == true {
		elastic.Index(fileReport, &filename) //Post es results
	}

	b, err := json.Marshal(fileReport)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatalf("failed to read %s", pluginFile)
	}
	log.Info(string(b))

}
