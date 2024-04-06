package main

import (
	"os"

	"malscan/config"
	mlog "malscan/core/logger"
	"malscan/core/scan"
	"malscan/core/utils"
	"malscan/system"

	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli"
)

/*
Author: Liam Hellend
Email: liamhellend@gmail.com

Purpose: entrypoint for malscan

*/

func main() {

	cli.AppHelpTemplate = utils.AppHelpTemplate
	app := cli.NewApp()

	app.Name = "Malscan"
	app.Version = "1.0.0"
	app.Usage = "Malscan"
	app.Flags = []cli.Flag{}
	app.Commands = []cli.Command{
		{
			Name:    "mode-1",
			Aliases: []string{"m1"},
			Usage:   "files and antivirus plugins are ran one at a time",
			Action: func(c *cli.Context) error {
				//config.Load()
				scan.Mode1()
				return nil
			},
		},
		{
			Name:    "mode-2",
			Aliases: []string{"m2"},
			Usage:   "files are ran one at a time, antivirus plugins are ran concurrently",
			Action: func(c *cli.Context) error {
				//config.Load()
				scan.Mode2()
				return nil
			},
		},
		{
			Name:    "mode-3",
			Aliases: []string{"m3"},
			Usage:   "files are ran concurrently (limit set in config), antivirus plugins are ran concurrently",
			Action: func(c *cli.Context) error {
				//config.Load()
				scan.Mode3()
				return nil
			},
		},
	}
	system.SetCPUCores()
	utils.MakeDirs()
	config.Load()
	mlog.Load()
	err := app.Run(os.Args)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("oh no failed to even start the app")
	}
}
