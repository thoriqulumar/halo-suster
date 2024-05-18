package controller

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"halo-suster/model"
	cerr "halo-suster/pkg/customErr"
	"net/http"
)

func (c *StaffController) RegisterNurse(ctx echo.Context) error {
	var req model.RegisterNurseRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.GeneralResponse{Message: err.Error()})
	}

	err := c.validate.Struct(req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, model.GeneralResponse{Message: err.Error()})
	}
	newStaff := model.Staff{
		NIP:                 req.NIP,
		Name:                req.Name,
		Role:                model.RoleNurse,
		IdentityCardScanImg: req.IdentityCardScanImg,
	}

	serviceRes, err := c.svc.RegisterNurse(ctx.Request().Context(), newStaff)
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

func (c *StaffController) LoginNurse(ctx echo.Context) error {
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
	serviceRes, err := c.svc.LoginNurse(ctx.Request().Context(), staff)
	if err != nil {
		return ctx.JSON(cerr.GetCode(err), model.GeneralResponse{Message: err.Error()})
	}

	registerStaffResponse := model.RegisterStaffResponse{
		Message: "User login successfully",
		Data: model.StaffWithToken{
			UserId:      serviceRes.UserId,
			Name:        serviceRes.Name,
			NIP:         serviceRes.NIP,
			AccessToken: serviceRes.AccessToken,
		},
	}

	return ctx.JSON(http.StatusOK, registerStaffResponse)
}

func (c *StaffController) UpdateNurse(ctx echo.Context) error {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, echo.Map{"error": "Invalid ID format"})
	}

	var payload model.UpdateNurseRequest
	if err := ctx.Bind(&payload); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.GeneralResponse{Message: err.Error()})
	}

	payload.ID = id
	err = c.validate.Struct(payload)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, model.GeneralResponse{Message: err.Error()})
	}

	err = c.svc.UpdateNurse(ctx.Request().Context(), payload)
	if err != nil {
		return ctx.JSON(cerr.GetCode(err), model.GeneralResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, model.GeneralResponse{})
}

func (c *StaffController) DeleteNurse(ctx echo.Context) error {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, echo.Map{"error": "Invalid ID format"})
	}

	err = c.svc.DeleteNurse(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(cerr.GetCode(err), model.GeneralResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, model.GeneralResponse{})
}
func (c *StaffController) GrantAccessNurse(ctx echo.Context) error {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, echo.Map{"error": "Invalid ID format"})
	}
	var payload model.GrantAccessNurseRequest
	if err := ctx.Bind(&payload); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.GeneralResponse{Message: err.Error()})
	}

	err = c.validate.Struct(payload)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, model.GeneralResponse{Message: err.Error()})
	}

	err = c.svc.GrantAccessNurse(ctx.Request().Context(), id, payload.Password)
	if err != nil {
		return ctx.JSON(cerr.GetCode(err), model.GeneralResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, model.GeneralResponse{})
}
