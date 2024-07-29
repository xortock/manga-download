package handlers

import (
	"github.com/urfave/cli"
	"github.com/xortock/mangafire-download/flags"
	"github.com/xortock/mangafire-download/services"
)


type ICliHandler interface {
}

type CliHandler struct {
	mangaService services.MangaService
}


func NewCliHandler() *CliHandler {
	return &CliHandler{
		mangaService: *services.NewMangaService(),
	}
}

func (handler *CliHandler) Handle(context *cli.Context) error {
	var mangaCode = context.String(flags.CODE)
	var mangaName = context.String(flags.NAME)
	var outputPath = context.String(flags.OUTPUTPATH)
	// assert that output path end with \
	handler.mangaService.Download(mangaName, mangaCode, outputPath)

	return cli.NewExitError("download completed", 0)
}