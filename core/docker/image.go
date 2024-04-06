package docker

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

/*
Author: Liam Hellend
Email: liamhellend@gmail.com

Purpose: Contains all image related functions

*/

//PullImage - Responsible for pulling the image that is passed into the function
//The image is stored locally on the machine
//DO NOT USE THIS FUNCTION BEFORE REVIEWING
func PullImage(image string) {

	log.Debug("PullImage - Pulling image")

	out, err := cli.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		log.Error(err)
	}

	defer out.Close()

	io.Copy(os.Stdout, out)

}

//GetIntalledImages - This function makes a call to the docker engine
//requesting all installed images, the functon returns a slice of image summary structs
func GetIntalledImages() (installedImages []types.ImageSummary) {

	log.Debug("finding installed images")

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		log.Error(errors.Wrap(err, "error while listing docker images"))
	}

	installedImages = images

	return installedImages
}
