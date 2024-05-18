package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"halo-suster/model"
	"halo-suster/pkg/crypto"
	cerr "halo-suster/pkg/customErr"
	"net/http"
)

func (s *staffSvc) RegisterNurse(ctx context.Context, newStaff model.Staff) (model.StaffWithToken, error) {
	role := getRoleFromNIP(newStaff.NIP)
	if role != model.RoleNurse {
		return model.StaffWithToken{}, cerr.New(http.StatusBadRequest, "nip not valid")
	}
	return s.Register(ctx, newStaff)
}

func (s *staffSvc) LoginNurse(ctx context.Context, staff model.Staff) (model.StaffWithToken, error) {
	role := getRoleFromNIP(staff.NIP)
	if role != model.RoleNurse {
		return model.StaffWithToken{}, cerr.New(http.StatusBadRequest, "nip not valid")
	}
	return s.Login(ctx, staff)
}

func (s *staffSvc) UpdateNurse(ctx context.Context, staff model.UpdateNurseRequest) error {
	//TODO: utilise TX
	user, err := s.repo.GetUserByID(ctx, staff.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return cerr.New(http.StatusNotFound, "Nurse not found")
		}
		return cerr.New(http.StatusInternalServerError, err.Error())
	}
	role := getRoleFromNIP(user.NIP)
	if role != model.RoleNurse {
		return cerr.New(http.StatusNotFound, "nip not for nurse")
	}

	existUserWithNip, err := s.repo.GetUserByNIP(ctx, staff.NIP)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return cerr.New(http.StatusInternalServerError, err.Error())
	}
	// self NIP
	if existUserWithNip.UserId != user.UserId && !errors.Is(err, sql.ErrNoRows) {
		return cerr.New(http.StatusNotFound, "NIP already exists")
	}

	// filled value
	user.NIP = staff.NIP
	user.Name = staff.Name
	return s.repo.UpdateUser(ctx, user)
}

func (s *staffSvc) DeleteNurse(ctx context.Context, userId uuid.UUID) error {
	user, err := s.repo.GetUserByID(ctx, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return cerr.New(http.StatusNotFound, "Nurse not found")
		}
		return cerr.New(http.StatusInternalServerError, err.Error())
	}
	role := getRoleFromNIP(user.NIP)
	if role != model.RoleNurse {
		return cerr.New(http.StatusNotFound, "nip not for nurse")
	}
	return s.repo.SoftDeleteUser(ctx, userId)
}

func (s *staffSvc) GrantAccessNurse(ctx context.Context, userId uuid.UUID, password string) error {
	user, err := s.repo.GetUserByID(ctx, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return cerr.New(http.StatusNotFound, "Nurse not found")
		}
		return cerr.New(http.StatusInternalServerError, err.Error())
	}
	role := getRoleFromNIP(user.NIP)
	if role != model.RoleNurse {
		return cerr.New(http.StatusNotFound, "nip not nurse")
	}
	// Hash password
	hashedPassword, err := crypto.GenerateHashedPassword(password, s.cfg.BcryptSalt)
	if err != nil {
		return cerr.New(http.StatusInternalServerError, err.Error())
	}
	user.Password = hashedPassword
	return s.repo.UpdateUser(ctx, user)
}
