package controller

import (
	"helo-suster/model"
	"helo-suster/service"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/labstack/echo/v4"
)

type StaffController struct {
	svc service.StaffService
}

func NewStaffController(svc service.StaffService) *StaffController {
	return &StaffController{svc: svc}
}

func (ctr *StaffController) GetStaff(c echo.Context) error {
	var request model.GetStaffRequest

	params := c.QueryParams()

	decoder := schema.NewDecoder()

	err := decoder.Decode(&request, params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid params"})
	}

	data, err := ctr.svc.GetStaff(c.Request().Context(), request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// compose response
	return c.JSON(http.StatusOK, model.GetStaffResponse{
		Message: "success",
		Data:    data,
	})
}
