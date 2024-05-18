package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"halo-suster/config"
	"halo-suster/model"
	"halo-suster/pkg/crypto"
	cerr "halo-suster/pkg/customErr"
	"halo-suster/repo"
	"net/http"
	"time"
)

type StaffService interface {
	GetUser(ctx context.Context, param model.GetListUserParams) ([]model.Staff, error)
	RegisterIT(ctx context.Context, newStaff model.Staff) (model.StaffWithToken, error)
	Login(ctx context.Context, staff model.Staff) (model.StaffWithToken, error)
	RegisterNurse(ctx context.Context, newStaff model.Staff) (model.StaffWithToken, error)
	LoginNurse(ctx context.Context, staff model.Staff) (model.StaffWithToken, error)
	UpdateNurse(ctx context.Context, staff model.UpdateNurseRequest) error
	DeleteNurse(ctx context.Context, userId uuid.UUID) error
	GrantAccessNurse(ctx context.Context, userId uuid.UUID, password string) error
}

type staffSvc struct {
	cfg  *config.Config
	repo repo.StaffRepo
}

func NewStaffService(cfg *config.Config, r repo.StaffRepo) StaffService {
	return &staffSvc{
		cfg:  cfg,
		repo: r,
	}
}

func (s *staffSvc) RegisterIT(ctx context.Context, newStaff model.Staff) (model.StaffWithToken, error) {
	role := getRoleFromNIP(newStaff.NIP)
	if role != model.RoleIt {
		return model.StaffWithToken{}, cerr.New(http.StatusBadRequest, "nip not valid")
	}
	return s.Register(ctx, newStaff)
}

func (s *staffSvc) Register(ctx context.Context, newStaff model.Staff) (model.StaffWithToken, error) {
	staff, err := s.repo.GetUserByNIP(ctx, newStaff.NIP)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return model.StaffWithToken{}, cerr.New(http.StatusInternalServerError, err.Error())
		}
	}
	if staff.NIP > 0 {
		return model.StaffWithToken{}, cerr.New(http.StatusConflict, "nip exist")
	}

	// Hash password
	hashedPassword, err := crypto.GenerateHashedPassword(newStaff.Password, s.cfg.BcryptSalt)
	if err != nil {
		return model.StaffWithToken{}, err
	}
	newStaff.Password = hashedPassword

	id := uuid.New()
	newStaff.UserId = id
	newStaff.CreatedAt = time.Now()
	newStaff.Status = model.StatusActive

	// save to database
	err = s.repo.InsertStaff(ctx, newStaff, hashedPassword)
	if err != nil {
		return model.StaffWithToken{}, cerr.New(http.StatusInternalServerError, err.Error())
	}

	// Generate token
	token, err := crypto.GenerateToken(newStaff, s.cfg.JWTSecret)
	if err != nil {
		return model.StaffWithToken{}, cerr.New(http.StatusInternalServerError, err.Error())

	}
	serviceResponse := model.StaffWithToken{
		UserId:      id.String(),
		AccessToken: token,
	}
	return serviceResponse, nil

}

func (s *staffSvc) Login(ctx context.Context, staff model.Staff) (model.StaffWithToken, error) {
	user, err := s.repo.GetUserByNIP(ctx, staff.NIP)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.StaffWithToken{}, cerr.New(http.StatusNotFound, "user not found")
		}
		return model.StaffWithToken{}, cerr.New(http.StatusInternalServerError, err.Error())
	}
	err = crypto.VerifyPassword(staff.Password, user.Password)
	if err != nil {
		return model.StaffWithToken{}, cerr.New(http.StatusBadRequest, "Invalid password")
	}
	token, err := crypto.GenerateToken(user, s.cfg.JWTSecret)
	if err != nil {
		return model.StaffWithToken{}, cerr.New(http.StatusBadRequest, err.Error())
	}

	serviceResponse := model.StaffWithToken{
		UserId:      user.UserId.String(),
		Name:        user.Name,
		NIP:         user.NIP,
		AccessToken: token,
	}

	return serviceResponse, nil
}

func (s *staffSvc) GetUser(ctx context.Context, param model.GetListUserParams) ([]model.Staff, error) {
	// generate filter query from param
	// do request
	listUser, err := s.repo.GetListUser(ctx, param)
	if listUser == nil || err != nil {
		listUser = []model.Staff{}
	}
	return listUser, err
}
