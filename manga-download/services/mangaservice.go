package services

import (
	"archive/zip"
	"compress/flate"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/schollz/progressbar/v3"
	"github.com/xortock/mangafire-download/clients"
	"github.com/xortock/mangafire-download/constants"
	"github.com/xortock/mangafire-download/helper"
	"github.com/xortock/mangafire-download/models"
	"github.com/xortock/mangafire-download/styles"
)

type IMangaService interface {
	Download(mangaName string, mangaCode string, outputpath string, fileType string, division string) error
}

type MangaService struct {
	mangClient clients.IMangaClient
}

func NewMangaService() *MangaService {
	return &MangaService{
		mangClient: clients.NewMangaClient(),
	}
}

func (service *MangaService) Download(mangaName string, mangaCode string, outputpath string, fileType string, divisionType string) error {

	// Get division based on input
	var divisions, errGetDivision = service.mangClient.GetDivisions(divisionType, mangaCode)
	if errGetDivision != nil {
		return errGetDivision
	}

	fmt.Println(styles.RenderSuccess(divisionType + "s found:" + strconv.Itoa(len(divisions))))
	var progressBar = progressbar.Default(int64(len(divisions)))

	// Create file type based on input
	var errCreateFileType = service.CreateFileType(fileType, mangaName, mangaCode, outputpath, divisionType, divisions, progressBar)
	if errCreateFileType != nil {
		return errCreateFileType
	}

	fmt.Println(styles.RenderSuccess("download completed"))
	return nil
}

func (service *MangaService) CreateFileType(fileType string, mangaName string, mangaCode string, outputpath string, divisionType string, divisions []models.Division, progressBar *progressbar.ProgressBar) error {
	switch fileType {
	case constants.FILE_TYPE_ZIP:
		var errCreatePdf = service.CreatePDFs(mangaName, mangaCode, outputpath, divisionType, divisions, progressBar)
		if errCreatePdf != nil {
			return errCreatePdf
		}
	case constants.FILE_TYPE_CBZ:
		var errCreateCbz = service.CreateCBZ(mangaName, mangaCode, outputpath, divisionType, divisions, progressBar)
		if errCreateCbz != nil {
			return errCreateCbz
		}
	default:
		return errors.New(fileType + " file type not supported")
	}
	return nil
}

func (service *MangaService) CreatePDFs(mangaName string, mangaCode string, outputpath string, divisionType string, divisions []models.Division, progressBar *progressbar.ProgressBar) error {
	for _, division := range divisions {

		var uris, err = service.mangClient.GetDivisionsUris(divisionType, division.Id)
		if err != nil {
			return err
		}

		var images [][]byte
		for _, uri := range uris {
			// download jpeg images as byte array
			var image, err = service.mangClient.GetDivisionImages(uri)
			if err != nil {
				return err
			}
			images = append(images, image)
		}
		// create pdf from images
		helper.CreatePdf(images, outputpath, mangaName, divisionType, division.Number)
		progressBar.Add(1)
	}
	return nil
}

func (service *MangaService) CreateCBZ(mangaName string, mangaCode string, outputpath string, divisionType string, divisions []models.Division, progressBar *progressbar.ProgressBar) error {
	// create zip file
	zipFile, err := os.Create(filepath.Join(outputpath, mangaName+".cbz"))
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// create zip writer
	var zipWriter = zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Register a custom Deflate compressor to compress zip contents
	zipWriter.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})

	// start page index at 1
	var pageCount = 1
	for _, division := range divisions {

		var uris, err = service.mangClient.GetDivisionsUris(divisionType, division.Id)
		if err != nil {
			return err
		}

		for _, uri := range uris {
			// download jpeg images as byte array in jpg format
			var jpgImage, errGetChapter = service.mangClient.GetDivisionImages(uri)
			if errGetChapter != nil {
				return errGetChapter
			}

			// create file inside zip
			var fileWriter, err = zipWriter.Create(mangaName + "_chapter_" + strconv.FormatFloat(division.Number, 'f', -1, 64) + "_page_" + strconv.Itoa(pageCount) + ".jpg")
			if err != nil {
				return err
			}

			// write bytes to file inside zip
			var _, errWriting = fileWriter.Write(jpgImage)
			if errWriting != nil {
				return errWriting
			}

			pageCount++
		}

		progressBar.Add(1)
	}
	return nil
}
