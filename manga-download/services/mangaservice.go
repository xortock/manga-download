package services

import (
	"slices"

	"github.com/schollz/progressbar/v3"
	"github.com/xortock/mangafire-download/clients"
	"github.com/xortock/mangafire-download/helper"
)

type IMangaService interface {
}

type MangaService struct {
	mangClient clients.MangaClient
}

func NewMangaService() *MangaService {
	return &MangaService{
		mangClient: *clients.NewMangaClient(),
	}
}

func (service *MangaService) Download(mangaName string, mangaCode string, outputpath string) {
	var chapters = service.mangClient.GetChapters(mangaCode)
	slices.Reverse(chapters)

	var progressBar = progressbar.Default(int64(len(chapters)))

	for _, chapter := range chapters {

		var uris = service.mangClient.GetChapterUris(chapter.Id)

		var images [][]byte
		for _, uri := range uris {
			var image = service.mangClient.GetChapterImages(uri)
			images = append(images, image)
		}

		helper.CreatePdf(images, outputpath, mangaName, chapter.Number)
		progressBar.Add(1)
	}
}
