package controller

import (
	"helo-suster/model"
	cerr "helo-suster/pkg/customError"
	"helo-suster/service"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"unicode"

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
		return ctx.JSON(http.StatusBadRequest, model.MedicalGeneralResponse{Message: err.Error()})
	}

	if len(strconv.Itoa(int(patientRequest.IdentityNumber))) != 16 {
		return ctx.JSON(http.StatusBadRequest, model.MedicalGeneralResponse{Message: "identityNumber should be 16 digits"})
	}

	if err := c.validate.Struct(&patientRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.MedicalGeneralResponse{Message: err.Error()})
	}

	_, err := c.service.CreateNewPatient(ctx.Request().Context(), patientRequest)
	if err != nil {
		return ctx.JSON(cerr.GetCode(err), model.MedicalGeneralResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusCreated, model.MedicalGeneralResponse{
		Message: "Patient successfully registered",
	})
}

func (c *MedicalController) PostMedicalReport(ctx echo.Context) error {
	var medicalRequest model.PostMedicalRecordRequest
	if err := ctx.Bind(&medicalRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.MedicalGeneralResponse{Message: err.Error()})
	}

	if len(strconv.Itoa(int(medicalRequest.IdentityNumber))) != 16 {
		return ctx.JSON(http.StatusBadRequest, model.MedicalGeneralResponse{Message: "identityNumber should be 16 digits"})
	}

	if err := c.validate.Struct(&medicalRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.MedicalGeneralResponse{Message: err.Error()})
	}

	mockUser := "8d203c88-9bc2-4838-ac7f-622cc737d614"

	_, err := c.service.CreateNewMedicalRecord(ctx.Request().Context(), medicalRequest, mockUser)
	if err != nil {
		return ctx.JSON(cerr.GetCode(err), model.MedicalGeneralResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusCreated, model.MedicalGeneralResponse{
		Message: "Medical record successfully added",
	})
}

func (c *MedicalController) GetPatient(ctx echo.Context) error {
	value, err := ctx.FormParams()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "params not valid"})
	}

	// query to service
	data, err := c.service.GetAllPatient(ctx.Request().Context(), parseGetPatientParams(value))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, model.MedicalGeneralResponse{
		Message: "success",
		Data:    data,
	})
}

func (c *MedicalController) GetMedicalRecord(ctx echo.Context) error {
	value, err := ctx.FormParams()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "params not valid"})
	}

	// query to service
	data, err := c.service.GetAllMedicalRecord(ctx.Request().Context(), parseGetMedicalRecordParams(value))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, model.MedicalGeneralResponse{
		Message: "success",
		Data:    data,
	})
}





func validatePhoneNumber(fl validator.FieldLevel) bool {
	phoneNumber := fl.Field().String()
	// Regular expression to match phone numbers starting with +62 and having 10 to 15 digits
	regex := `^\+62\d{8,13}$`
	return regexp.MustCompile(regex).MatchString(phoneNumber)
}

func parseGetPatientParams(params url.Values) model.GetPatientParams {
	var result model.GetPatientParams

	for key, values := range params {
		switch key {
		case "identityNumber":
			if len(values[0]) == 16 {
				identityNumber, _ := strconv.Atoi(values[0])
				result.IdentityNumber = &identityNumber
			}
		case "name":
			result.Name = values[0]
		case "phoneNumber":
			result.PhoneNumber = stripNonNumeric(values[0])
		case "limit":
			limit, err := strconv.Atoi(values[0])
			if err == nil {
				result.Limit = limit
			}
		case "offset":
			offset, err := strconv.Atoi(values[0])
			if err == nil {
				result.Offset = offset
			}
		case "createdAt":
			result.CreatedAt = values[0]
		}
	}

	return result
}

func parseGetMedicalRecordParams(params url.Values) model.GetMedicalRecordParams {
	var result model.GetMedicalRecordParams

	for key, values := range params {
		switch key {
		case "identityDetail.identityNumber":
			identityNumber, err := strconv.Atoi(values[0])
			if err == nil {
				if len(values[0]) == 16 {
					result.IdentityNumber = &identityNumber
				}
			}
		case "createdBy.userId":
			result.CreatedByUserId = values[0]
		case "createdBy.nip":
			result.CreatedByNip = stripNonNumeric(values[0])
		case "limit":
			limit, err := strconv.Atoi(values[0])
			if err == nil {
				result.Limit = limit
			}
		case "offset":
			offset, err := strconv.Atoi(values[0])
			if err == nil {
				result.Offset = offset
			}
		case "createdAt":
			result.CreatedAt = values[0]
		}
	}

	return result
}

func stripNonNumeric(input string) string {
	var result []rune
	for _, r := range input {
		if unicode.IsDigit(r) {
			result = append(result, r)
		}
	}
	return string(result)
}
