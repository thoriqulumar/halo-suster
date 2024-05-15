package controller

import (
	"helo-suster/model"
	cerr "helo-suster/pkg/customError"
	"helo-suster/service"
	"net/http"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type MedicalController struct {
	service  service.MedicalService
	validate *validator.Validate
}

func NewMedicalController(service service.MedicalService, validate *validator.Validate) *MedicalController {
	_ = validate.RegisterValidation("phone_number", validatePhoneNumber)

	return &MedicalController{
		service:  service,
		validate: validate,
	}
}

func (c *MedicalController) PostPatient(ctx echo.Context) error {
	var patientRequest model.PostPatientRequest
	if err := ctx.Bind(&patientRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest,  model.MedicalGeneralResponse{Message: err.Error()})
	}

	if err := c.validate.Struct(&patientRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest,  model.MedicalGeneralResponse{Message: err.Error()})
	}
	
	_, err := c.service.CreateNewPatient(ctx.Request().Context(), patientRequest)
	if err != nil {
		return ctx.JSON(cerr.GetCode(err), model.MedicalGeneralResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusCreated, model.MedicalGeneralResponse{
		Message: "Customer registered successfully",
	})
}


func validatePhoneNumber(fl validator.FieldLevel) bool {
	phoneNumber := fl.Field().String()
	// Regular expression to match phone numbers starting with +62 and having 10 to 15 digits
	regex := `^\+62\d{8,13}$`
	return regexp.MustCompile(regex).MatchString(phoneNumber)
}