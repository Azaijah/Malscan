package scp

import (
	"os"

	"malscan/config"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

//Send - Accepts a local filepath and
func Send(localfilepath string, dstip string, remotefilepath string) {

	// Use SSH key authentication from the auth package
	// we ignore the host key in this example, please change this if you use this library
	clientConfig, _ := auth.PrivateKey(config.Values.Alert.ScpUser, config.Values.Alert.ScpPkey, ssh.InsecureIgnoreHostKey())

	clientConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey() //Must add this because of CVE-2017-3204

	// Create a new SCP client
	client := scp.NewClient(dstip, &clientConfig)

	// Connect to the remote server
	err := client.Connect()
	if err != nil {
		log.Error(errors.Wrap(err, "couldn't establish a connection to the remote server"))
		return
	}

	// Open a file
	f, _ := os.Open(localfilepath)

	//Close ssh connection after file has been sent
	defer client.Close()

	// Close the file after it has been copied
	defer f.Close()

	// Finaly, copy the file over
	// Usage: CopyFile(fileReader, remotePath, permission)

	err = client.CopyFile(f, remotefilepath, "0655")

	if err != nil {
		log.Error(errors.Wrap(err, "error while copying file"))
	}

}
