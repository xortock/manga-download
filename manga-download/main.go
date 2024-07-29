package main

import (
	"os"

	"github.com/urfave/cli"
	"github.com/xortock/mangafire-download/flags"
	"github.com/xortock/mangafire-download/handlers"
)

var version = "1.0.0.0"

func main() {
	var cliHandler = handlers.NewCliHandler()
	
	var app = &cli.App{
		Name: "semanticloq",
		Version: version,
		Authors: []cli.Author{
			{
				Name:  "xortock",
				Email: "bgmaduro@gmail.com",
			},
		},
		Copyright:       "(C) 2024 xortock",
		HideHelp:        true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     flags.CODE,
				Usage:    "the mangafire manga code",
				Required: true,
			},
			&cli.StringFlag{
				Name:     flags.NAME,
				Usage:    "the folder name where the manga chapters should be stored",
				Required: true,
			},
			&cli.StringFlag{
				Name:     flags.OUTPUTPATH,
				Usage:    "the output path for the manga folder ",
				Required: true,
			},
		},
		Action: cliHandler.Handle,
	}

	var _ = app.Run(os.Args)
}
