package service

import (
	"context"
	"helo-suster/model"
	"helo-suster/repo"
)

type StaffService interface {
	GetStaff(ctx context.Context, param model.GetStaffRequest) ([]model.Staff, error)
}

type staffService struct {
	repo repo.StaffRepo
}

func NewStaffService(repo repo.StaffRepo) StaffService {
	return &staffService{
		repo: repo,
	}
}

func (s *staffService) GetStaff(ctx context.Context, param model.GetStaffRequest) ([]model.Staff, error) {
	staffList, err := s.repo.GetStaff(ctx, param)
	if staffList == nil || err != nil {
		staffList = []model.Staff{}
	}
	return staffList, err
}
