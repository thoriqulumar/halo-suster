package controller

import (
	"helo-suster/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ImageController struct {
	svc service.ImageService
}

func NewImageController(svc service.ImageService) *ImageController {
	return &ImageController{
		svc: svc,
	}
}

func (ctr *ImageController) PostImage(c echo.Context) error {
	// Read form data including uploaded file
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "params not valid"})
	}

	file := form.File["file"][0]

	urlChan := ctr.svc.UploadImage(file)

	done := make(chan bool)

	go func() {
		defer close(done)
		<-urlChan
	}()

	<-done

	url := <-urlChan
	if url == "" {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "File uploaded sucessfully"})
	// src, err := file.Open()
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Internal server error"})
	// }
	// defer src.Close()

	// fileName := file.Filename
	// if !strings.HasSuffix(strings.ToLower(fileName), ".jpg") && !strings.HasSuffix(strings.ToLower(fileName), ".jpeg") {
	// 	return c.JSON(http.StatusBadRequest, echo.Map{"error": "Only JPG/JPEG files are allowed"})
	// }

	// fileBytes, err := io.ReadAll(src)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Internal server error"})
	// }

	// fileSize := len(fileBytes)
	// if fileSize > 2*1024*1024 || fileSize < 10*1024 {
	// 	return c.JSON(http.StatusBadRequest, echo.Map{"error": "File size must be between 10KB and 2MB"})
	// }

	// // Reset file reader after reading its content
	// src.Seek(0, io.SeekStart)

	// newFileName, err := ctr.svc.SaveImage(src)

	// ctr.s3Svc.UploadImage(fileName, src) // like this
}
