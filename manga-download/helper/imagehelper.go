package helper

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"path/filepath"
	"strconv"

	"github.com/go-pdf/fpdf"
	"github.com/pkg/errors"
)

func CreatePdf(images [][]byte, outputPath string, mangaName string, divisionType string, division float64) error {
	var imageType = "jpg"

	var filePath = filepath.Join(outputPath, mangaName)

	// Create directory if not exists
	var errCreateDir = CreateDirectoryIfNotExists(filePath)
	if errCreateDir != nil {
		return errCreateDir
	}

	// Create new pdf document
	var document = fpdf.New("P", "mm", "A4", "")

	for index, img := range images {
		pngImage, _, _ := image.DecodeConfig(bytes.NewReader(img))

		document.AddPageFormat("P", fpdf.SizeType{
			Wd: float64(pngImage.Width),
			Ht: float64(pngImage.Height),
		})

		var imageKey = "image" + strconv.Itoa(index)
		// Set pdf image size and options
		document.RegisterImageOptionsReader(imageKey, fpdf.ImageOptions{ImageType: imageType, ReadDpi: true}, bytes.NewReader(img))
		document.ImageOptions(imageKey, 0, 0, float64(pngImage.Width), float64(pngImage.Height), false, fpdf.ImageOptions{ImageType: imageType, ReadDpi: true}, 0, "")
	}

	var err = document.OutputFileAndClose(filepath.Join(filePath , divisionType + "_" + strconv.FormatFloat(division, 'f', -1, 64) + ".pdf"))
	if err != nil {
		return err
	}

	return nil
}


func JpgToPng(jpgBytes []byte) ([]byte, error) {
	var img, errDecode = jpeg.Decode(bytes.NewReader(jpgBytes))
	if errDecode != nil {
		return nil, errors.Wrap(errDecode, "unable to decode jpeg")
	}

	buffer := new(bytes.Buffer)
	var err = png.Encode(buffer, img); 
	if err != nil {
		return nil, errors.Wrap(err, "unable to encode png")
	}

	return buffer.Bytes(), nil
}