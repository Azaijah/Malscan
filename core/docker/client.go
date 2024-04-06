package docker

import (
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

/*
Author: Liam Hellend
Email: liamhellend@gmail.com

Purpose: Initialize docker API

*/

var cli *client.Client //Docker client used to make calls down to the docker engine
var err error

//init - Responsible for initializing a new docker client
//The docker client is used to interact with the docker daemon
func init() {

	log.Debug("initializing docker client")

	cli, err = client.NewEnvClient()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("creating docker client")
	}
}
