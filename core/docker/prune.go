package docker

import (
	"context"

	"github.com/docker/docker/api/types/filters"
	log "github.com/sirupsen/logrus"
)

/*
Author: Liam Hellend
Email: liamhellend@gmail.com

Purpose: Contains functions related to cleaning the docker system

*/

//Prune - This function simply simply makes a call to the docker engine
//requesting that all left over container data is cleaned up
//BE VERY CAREFUL WHEN USING THIS FUNCTION (ALL CONTAINERS GET MARKED FOR REMOVAL)
func Prune() {

	log.Debug("cleaning up containers")

	_, err := cli.ContainersPrune(context.Background(), filters.NewArgs())
	if err != nil {
		log.Error(err)
	}
}
