package main

import (
	"os"

	"github.com/urfave/cli/v2"
	"github.com/xortock/mangafire-download/constants"
	"github.com/xortock/mangafire-download/handlers"
)

var version = "development"

func main() {
	var cliHandler = handlers.NewCliHandler()

	var app = &cli.App{
		Name:        "manga-dl",
		Usage:       "download manga in PDF or CBZ format",
		Version:     version,
		Description: "manga-dl is a cli tool build to download manga in PDF or CBZ format from mangafire.to",
		Authors: []*cli.Author{
			{
				Name:  "xortock",
				Email: "bgmaduro@gmail.com",
			},
		},
		Copyright: "(C) 2024 xortock",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     constants.FLAG_CODE,
				Aliases:  []string{"c"},
				Usage:    "the mangafire manga code",
				Required: true,
			},
			&cli.StringFlag{
				Name:     constants.FLAG_NAME,
				Aliases:  []string{"n"},
				Usage:    "the folder name where the manga chapters should be stored",
				Required: true,
			},
			&cli.StringFlag{
				Name:     constants.FLAG_OUTPUTPATH,
				Aliases:  []string{"o"},
				Usage:    "the output path for the manga folder",
				Required: true,
			},
			&cli.StringFlag{
				Name:     constants.FLAG_TYPE,
				Aliases:  []string{"t"},
				Value:    constants.FILE_TYPE_ZIP,
				Usage:    "the type in which the manga should be downloaded " + constants.FILE_TYPE_ZIP + " or " + constants.FILE_TYPE_CBZ,
				Required: false,
			},
			&cli.StringFlag{
				Name:     constants.FLAG_DIVISION,
				Aliases:  []string{"d"},
				Value:    constants.DIVISION_CHAPTER,
				Usage:    "the division type in which the manga should be download (does not make a difference in case of file type cbz) " + constants.DIVISION_CHAPTER + " or " + constants.DVISION_VOLUME,
				Required: false,
			},
		},
		Action: cliHandler.Handle,
	}

	var _ = app.Run(os.Args)
}
