package scan

import (
	"time"

	"malscan/config"
	"malscan/core/utils"
	file "malscan/core/utils/file"
	pconfig "malscan/plugins"

	mime "malscan/core/utils/mime"

	"github.com/pkg/errors"
	"github.com/radovskyb/watcher"
	log "github.com/sirupsen/logrus"
)

/*
Author: Liam Hellend
Email: liamhellend@gmail.com

Purpose: Contains functions related to file scanning

*/

//Mode1 - Responsible for watching the malscan filestore
//when a file enters the filestore the file is scanned using malscan plugins
func Mode1() {

	log.Debug("malscan is ready to start scanning files in mode-1 ... waiting for files")

	plugins := pconfig.PluginConfig{}
	plugins = plugins.Load()

	w := watcher.New()

	// SetMaxEvents to 1 to allow at most 1 event's to be received
	// on the Event channel per watching cycle.
	// If SetMaxEvents is not set, the default is to send all events.
	//w.SetMaxEvents(1)

	// Only notify create events
	w.FilterOps(watcher.Create)

	go func() {
		for {
			select {
			case event := <-w.Event:
				log.Infof("file type:%s", mime.FileType(event.Name()))
				plugins.RunEnabled(event.Name())
				filename := event.Name()
				file.Remove(&filename)
			case err := <-w.Error:
				log.Error(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch this folder for changes.

	var folderToWatch string

	if config.Values.Env.Filestore == "" {
		folderToWatch = utils.GetFilestoreDir() //If filestore has not be set in the config file then use default filestore
	} else {
		folderToWatch = config.Values.Env.Filestore //If filestore is set use the user provided path in config file
	}

	if err := w.Add(folderToWatch); err != nil {
		log.Fatal(err)
	}

	// Print a list of all of the files and folders currently
	// being watched and their paths.
	//for path, f := range w.WatchedFiles() {
	//log.Debug("Watching: ", f.Name(), " at: ", path)
	//}

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatal(err)
	}

}

//Mode2 - Responsible for watching the malscan filestore
//when a file enters the filestore the file is scanned using malscan plugins
func Mode2() {

	log.Debug("malscan is ready to start scanning files in mode-2 ... waiting for files")

	plugins := pconfig.PluginConfig{}
	plugins = plugins.Load()

	w := watcher.New()

	// SetMaxEvents to 1 to allow at most 1 event's to be received
	// on the Event channel per watching cycle.
	// If SetMaxEvents is not set, the default is to send all events.
	//w.SetMaxEvents(1)

	// Only notify create events
	w.FilterOps(watcher.Create)

	go func() {
		for {
			select {
			case event := <-w.Event:
				log.Infof("file type:%s", mime.FileType(event.Name()))
				plugins.RunEnabledConcurrent(event.Name())
				filename := event.Name()
				file.Remove(&filename)
			case err := <-w.Error:
				log.Error(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch this folder for changes.

	var folderToWatch string

	if config.Values.Env.Filestore == "" {
		folderToWatch = utils.GetFilestoreDir() //If filestore has not be set in the config file then use default filestore
	} else {
		folderToWatch = config.Values.Env.Filestore //If filestore is set use the user provided path in config file
	}

	if err := w.Add(folderToWatch); err != nil {
		log.Fatal(err)
	}

	// Print a list of all of the files and folders currently
	// being watched and their paths.
	//for path, f := range w.WatchedFiles() {
	//log.Debug("Watching: ", f.Name(), " at: ", path)
	//}

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatal(err)
	}

}

//Mode3 - Responsible for watching the malscan filestore
//when a file enters the filestore the file is scanned using malscan plugins
func Mode3() {

	log.Debug("malscan is ready to start scanning files in mode-3 ... waiting for files")

	plugins := pconfig.PluginConfig{}

	plugins = plugins.Load()

	w := watcher.New()

	// SetMaxEvents to 1 to allow at most 1 event's to be received
	// on the Event channel per watching cycle.
	// If SetMaxEvents is not set, the default is to send all events.
	//w.SetMaxEvents(1)

	// Only notify create events
	w.FilterOps(watcher.Create)

	filesProcessingCount := 0

	go func() {
		for {
			select {
			case event := <-w.Event:
				log.Infof("file type:%s", mime.FileType(event.Name()))
				done := make(chan bool)
				go func() {
					filename := event.Name()
					filesProcessingCount++
					for filesProcessingCount > config.Values.Env.MaxFileProc {
						log.Debugf("max files:%d:processing", filesProcessingCount)
						time.Sleep(time.Second * 5)
					}
					go plugins.RunEnabledConcurrentWithChannel(event.Name(), done)
					<-done
					file.Remove(&filename)
					filesProcessingCount--
				}()
			case err := <-w.Error:
				log.Error(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch this folder for changes.

	var folderToWatch string

	if config.Values.Env.Filestore == "" {
		folderToWatch = utils.GetFilestoreDir() //If filestore has not be set in the config file then use default filestore
	} else {
		folderToWatch = config.Values.Env.Filestore //If filestore is set use the user provided path in config file
	}

	if err := w.Add(folderToWatch); err != nil {
		log.Fatal(errors.Wrap(err, "error while opening filestore folder"))
	}

	// Print a list of all of the files and folders currently
	// being watched and their paths.
	for path, f := range w.WatchedFiles() {
		log.Debug("Watching: ", f.Name(), " at: ", path)
	}

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatal(errors.Wrap(err, "error while starting filestore watcher"))
	}

}
