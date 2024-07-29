package helper

import (
	"bytes"
	"image"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-pdf/fpdf"
	"github.com/urfave/cli/v2"
)

func CreatePdf(images [][]byte, outputPath string, mangaName string, chapter float64) {
	var imageType = "jpg"

	var filePath = filepath.Join(outputPath, mangaName)

	// Create directory if not exists
	CreateDirectoryIfNotExists(filePath)

	// Create new pdf document
	var document = fpdf.New("P", "mm", "A4", "")

	for index, img := range images {
		pngImage, _, _ := image.DecodeConfig(bytes.NewReader(img))

		document.AddPageFormat("P", fpdf.SizeType{
			Wd: float64(pngImage.Width),
			Ht: float64(pngImage.Height),
		})

		var imageKey = "image" + strconv.Itoa(index)
		document.RegisterImageOptionsReader(imageKey, fpdf.ImageOptions{ImageType: imageType, ReadDpi: true}, bytes.NewReader(img))
		document.ImageOptions(imageKey, 0, 0, float64(pngImage.Width), float64(pngImage.Height), false, fpdf.ImageOptions{ImageType: imageType, ReadDpi: true}, 0, "")
	}

	var err = document.OutputFileAndClose(filepath.Join(filePath ,"chapter_" + strconv.FormatFloat(chapter, 'f', -1, 64) + ".pdf"))
	if err != nil {
		cli.Exit("Error: Failed to write chapter to outputpath", 1)
	}
}

func CreateDirectoryIfNotExists(dirctoryName string) error {
    var err = os.MkdirAll(dirctoryName, os.ModeDir)

    if err == nil || os.IsExist(err) {
        return nil
    } else {
        return err
    }
}