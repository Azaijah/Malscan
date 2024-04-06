package system

import (
	"malscan/config"
	"runtime"

	log "github.com/sirupsen/logrus"
)

//SetCPUCores - Responsible for setting cpu cores as specified in the malscan config
func SetCPUCores() {

	log.Debug("Setting cpu cores to: ", config.Values.Env.CPUcores)

	runtime.GOMAXPROCS(config.Values.Env.CPUcores)

}

//SetCPUCoresOverride - Responsible for setting cpu cores as passed into the function
func SetCPUCoresOverride(cores int) {

	log.Debug("Setting cpu cores to: ", cores)

	runtime.GOMAXPROCS(cores)

}
