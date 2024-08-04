package handlers

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/xortock/mangafire-download/constants"
	"github.com/xortock/mangafire-download/services"
	"github.com/xortock/mangafire-download/styles"
	"github.com/xortock/mangafire-download/validators"
)

type ICliHandler interface {
	Handle(context *cli.Context) error
}

type CliHandler struct {
	mangaService services.IMangaService
}

func NewCliHandler() *CliHandler {
	return &CliHandler{
		mangaService: services.NewMangaService(),
	}
}

func (handler *CliHandler) Handle(context *cli.Context) error {
	var mangaCode = context.String(constants.FLAG_CODE)
	var mangaName = context.String(constants.FLAG_NAME)
	var outputPath = context.String(constants.FLAG_OUTPUTPATH)

	var fileType = context.String(constants.FLAG_TYPE)
	var errTypeFlag = validators.ValidateTypeFlag(validators.TypeFlag{Value: fileType})
	if errTypeFlag != nil {
		return cli.Exit(styles.RenderFailed(fmt.Errorf("--"+constants.FLAG_DIVISION + " %w", errTypeFlag).Error()), 1)
	}

	var divisionType = context.String(constants.FLAG_DIVISION)
	var errDivisionFlag = validators.ValidateDivisionFlag(validators.DivisionFlag{Value: divisionType})
	if errDivisionFlag != nil {
		return cli.Exit(styles.RenderFailed(fmt.Errorf("--"+constants.FLAG_DIVISION + " %w", errDivisionFlag).Error()), 1)
	}

	var err = handler.mangaService.Download(mangaName, mangaCode, outputPath, fileType, divisionType)
	if err != nil {
		return cli.Exit(styles.RenderFailed(err.Error()), 1)
	}

	return nil
}
