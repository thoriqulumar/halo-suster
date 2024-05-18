package controller

import (
	"helo-suster/model"
	"helo-suster/service"
	"net/http"

	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/labstack/echo/v4"
)

type StaffController struct {
	svc      service.StaffService
	validate *validator.Validate
}

func NewStaffController(svc service.StaffService, validate *validator.Validate) *StaffController {
	validate.RegisterValidation("nipValidator", validateNIP)
	return &StaffController{
		svc:      svc,
		validate: validate,
	}
}

func (c *StaffController) Register(ctx echo.Context) error {
	var newStaffReq model.RegisterStaffRequest

	newStaff := model.Staff{
		Name:     newStaffReq.Name,
		NIP:      newStaffReq.NIP,
		Password: newStaffReq.Password,
	}

	serviceRes, err := c.svc.Register(newStaff)
	if err != nil {
		switch err.Error() {
		case "User already exist":
			return ctx.JSON(http.StatusConflict, err.Error())
		default:
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	registerStaffResponse := model.RegisterStaffResponse{
		Message: "User registered successfully",
		Data: model.StaffWithToken{
			UserId: serviceRes.UserId,
			//NIP:         newStaff.NIP,
			Name:        newStaff.Name,
			AccessToken: serviceRes.AccessToken,
		},
	}

	return ctx.JSON(http.StatusCreated, registerStaffResponse)
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

func validateNIP(fl validator.FieldLevel) bool {
	nip := fmt.Sprintf("%d", fl.Field().Int())
	currentYear := time.Now().Year()

	// Validasi format secara keseluruhan
	regexPattern := fmt.Sprintf(`^615[12]20[0-2][0-9](0[1-9]|1[0-2])[0-9]{3}$`)
	match, _ := regexp.MatchString(regexPattern, nip)
	if !match {
		return false
	}

	// Validasi tahun (digit ke-5 sampai ke-8) harus antara 2000 dan tahun saat ini
	year := nip[3:7]
	yearInt, err := strconv.Atoi(year)
	if err != nil || yearInt < 2000 || yearInt > currentYear {
		return false
	}

	// Validasi bulan (digit ke-9 dan ke-10) harus antara 01 dan 12
	month := nip[7:9]
	monthInt, err := strconv.Atoi(month)
	if err != nil || monthInt < 1 || monthInt > 12 {
		return false
	}

	// Tidak perlu validasi untuk tiga digit acak terakhir karena sudah dijamin oleh regex

	return true
}
