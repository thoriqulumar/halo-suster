package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"halo-suster/model"
	cerr "halo-suster/pkg/customErr"
	"halo-suster/service"
	"net/http"
)

type StaffController struct {
	svc      service.StaffService
	validate *validator.Validate
}

func NewStaffController(svc service.StaffService, validate *validator.Validate) *StaffController {
	_ = validate.RegisterValidation("nip", validateNIP)
	_ = validate.RegisterValidation("customURL", customURL)
	return &StaffController{
		svc:      svc,
		validate: validate,
	}
}

func (c *StaffController) RegisterIT(ctx echo.Context) error {
	var newStaffReq model.RegisterStaffRequest
	if err := ctx.Bind(&newStaffReq); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.GeneralResponse{Message: err.Error()})
	}

	err := c.validate.Struct(newStaffReq)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, model.GeneralResponse{Message: err.Error()})
	}
	newStaff := model.Staff{
		NIP:                 newStaffReq.NIP,
		Name:                newStaffReq.Name,
		Role:                model.RoleIt,
		IdentityCardScanImg: "",
		Password:            newStaffReq.Password,
	}

	serviceRes, err := c.svc.RegisterIT(ctx.Request().Context(), newStaff)
	if err != nil {
		return ctx.JSON(cerr.GetCode(err), model.GeneralResponse{Message: err.Error()})
	}

	registerStaffResponse := model.RegisterStaffResponse{
		Message: "User registered successfully",
		Data: model.StaffWithToken{
			UserId:      serviceRes.UserId,
			NIP:         newStaff.NIP,
			Name:        newStaff.Name,
			AccessToken: serviceRes.AccessToken,
		},
	}

	return ctx.JSON(http.StatusCreated, registerStaffResponse)
}

func (c *StaffController) LoginIT(ctx echo.Context) error {
	var staffReq model.LoginStaffRequest
	if err := ctx.Bind(&staffReq); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.GeneralResponse{Message: err.Error()})
	}

	err := c.validate.Struct(staffReq)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, model.GeneralResponse{Message: err.Error()})
	}
	staff := model.Staff{
		NIP:      staffReq.NIP,
		Password: staffReq.Password,
	}
	serviceRes, err := c.svc.Login(ctx.Request().Context(), staff)
	if err != nil {
		return ctx.JSON(cerr.GetCode(err), model.GeneralResponse{Message: err.Error()})
	}

	registerStaffResponse := model.RegisterStaffResponse{
		Message: "User registered successfully",
		Data: model.StaffWithToken{
			UserId:      serviceRes.UserId,
			Name:        serviceRes.Name,
			NIP:         serviceRes.NIP,
			AccessToken: serviceRes.AccessToken,
		},
	}

	return ctx.JSON(http.StatusOK, registerStaffResponse)
}

func (ctr *StaffController) GetUser(c echo.Context) error {
	// parse param to model
	value, err := c.FormParams()
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "params not valid"})
	}
	value.Add("status", string(model.StatusActive))
	// query to service
	data, err := ctr.svc.GetUser(c.Request().Context(), parseGetUserParams(value))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// compose response
	return c.JSON(http.StatusOK, model.GeneralResponse{
		Message: "success",
		Data:    data,
	})
}
