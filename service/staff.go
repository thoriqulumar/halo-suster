package service

import (
	"context"
	"errors"
	"fmt"
	"helo-suster/config"
	"helo-suster/model"
	"helo-suster/pkg/crypto"
	"helo-suster/repo"

	"github.com/google/uuid"
)

type StaffService interface {
	Register(newStaff model.Staff) (model.StaffWithToken, error)
	GetStaff(ctx context.Context, param model.GetStaffRequest) ([]model.Staff, error)
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

func (s *staffSvc) Register(newStaff model.Staff) (model.StaffWithToken, error) {
	// Validasi data yang diperlukan
	if newStaff.NIP == 0 || newStaff.Name == "" || newStaff.Password == "" {
		return model.StaffWithToken{}, errors.New("nip, name, and password are required")
	}
	role, err := DetermineRoleByNIP(newStaff.NIP)
	if err != nil {
		return model.StaffWithToken{}, err
	}
	newStaff.Role = role

	// Hash password
	hashedPassword, err := crypto.GenerateHashedPassword(newStaff.Password, s.cfg.BcryptSalt)
	if err != nil {
		return model.StaffWithToken{}, err
	}
	newStaff.Password = hashedPassword

	id := uuid.New()
	newStaff.UserId = id

	// Simpan ke database
	err = s.repo.InsertStaff(newStaff, hashedPassword)
	if err != nil {
		return model.StaffWithToken{}, err
	}

	// Generate token
	token, err := crypto.GenerateToken(id, newStaff.Name, s.cfg.JWTSecret)
	if err != nil {
		return model.StaffWithToken{}, err

	}
	serviceResponse := model.StaffWithToken{
		UserId:      id.String(),
		AccessToken: token,
	}
	return serviceResponse, err

}

func (s *staffSvc) GetStaff(ctx context.Context, param model.GetStaffRequest) ([]model.Staff, error) {
	staffList, err := s.repo.GetStaff(ctx, param)
	if staffList == nil || err != nil {
		staffList = []model.Staff{}
	}
	return staffList, err
}

func DetermineRoleByNIP(nip int64) (string, error) {
	nipStr := fmt.Sprintf("%013d", nip)
	if len(nipStr) != 13 || nipStr[:3] != "615" {
		return "", errors.New("invalid NIP format")
	}

	switch nipStr[3] {
	case '1':
		return "it", nil
	case '2':
		return "nurse", nil
	default:
		return "", errors.New("invalid role digit")
	}
}
