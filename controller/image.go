package controller

import (
	"fmt"
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

	url := <-urlChan
	fmt.Println("url in ctrl", url)
	if url == "" {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "File uploaded sucessfully"})
}
