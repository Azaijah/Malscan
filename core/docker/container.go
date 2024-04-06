package docker

import (
	"bytes"
	"path/filepath"

	"malscan/config"
	"malscan/core/utils"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

/*
Author: Liam Hellend
Email: liamhellend@gmail.com

Purpose: Contains all container related functions

*/

//RunContainerOnFile - Used to run whatever image is passed into the function as a container against the file
//passed onto the function, returns the results of the scan as a slice of bytes
func RunContainerOnFile(image string, fileToScanName *string) []byte {

	log.Debugf("running container:%s:against file:%s", image, *fileToScanName)

	var fileDirOnDisk string
	command := filepath.Join("/malware", *fileToScanName)
	if config.Values.Env.Filestore == "" {
		fileDirOnDisk = utils.GetFilestoreDir() //If filestore has not be set in the config file then use default filestore
	} else {
		fileDirOnDisk = config.Values.Env.Filestore //If filestore is set use the user provided path in config file
	}

	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image: image,
		Cmd:   []string{command},
		Tty:   true,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: fileDirOnDisk,
				Target: "/malware",
			},
		},
	}, nil, &v1.Platform{Architecture: "amd64", OS: "linux"}, "")
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Errorf("creating container for:%s", image)
	}

	if err := cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		log.WithFields(log.Fields{"err": err}).Errorf("starting container for:%s", image)
	}

	cli.ContainerWait(context.Background(), resp.ID, container.WaitConditionNotRunning)

	/*if err != nil {
		log.Error(errors.Wrap(err, "waiting for container, image: "+image))
	}
	*/

	out, err := cli.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		Follow:     true,
	})
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Errorf("getting logs for:%s", image)
	}

	buf := new(bytes.Buffer)

	buf.ReadFrom(out)

	out.Close()

	err = cli.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{
		Force: true,
	})

	if err != nil {
		log.WithFields(log.Fields{"err": err}).Errorf("removing container:%s", image)
	}

	log.Debugf("finished running container:%s:against file:%s", image, *fileToScanName)

	return (buf.Bytes())
}

//RunContainerUpdate - Responsible for accepting an image, spawning the container and running the update command for that image/plugin.
//Once the update command has fully run in the container the container is then commited back to an image.
func RunContainerUpdate(image string) (outcome string) {

	log.Debugf("running container:%s:update", image)

	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image: image,
		Cmd:   []string{"update"},
		Tty:   true,
	}, nil, nil, &v1.Platform{Architecture: "amd64", OS: "linux"}, "")
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Errorf("creating container for:%s", image)
	}

	if err := cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		log.WithFields(log.Fields{"err": err}).Errorf("starting container for:%s", image)
	}

	cli.ContainerWait(context.Background(), resp.ID, container.WaitConditionNotRunning)

	/*if err != nil {
		log.Error(err)
	}
	*/

	image = filepath.Join("docker.io/", image)

	com, err := cli.ContainerCommit(context.Background(), resp.ID, types.ContainerCommitOptions{Reference: image})
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Errorf("committing container for:%s", image)
	} else {
		log.Debugf("committed updated container:%s:ID:%s ", image, com.ID)
	}

	out, err := cli.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		Follow:     true,
	})
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Errorf("getting logs for:%s", image)
	}

	buf := new(bytes.Buffer)

	buf.ReadFrom(out)

	out.Close()

	err = cli.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{
		Force: true,
	})

	if err != nil {
		log.WithFields(log.Fields{"err": err}).Errorf("removing container:%s", image)
	}

	if string(buf.Bytes()[0]) != "1" {
		return "success"
	}

	log.Debugf("finished running container:%s:update", image)

	return "failure"

}
