package repo

import (
	"helo-suster/model"

	"github.com/jmoiron/sqlx"
)

type StaffRepo interface {
	InsertStaff(staff model.Staff, hashPassword string) error
}

type staffRepo struct {
	db *sqlx.DB
}

func NewStaffRepo(db *sqlx.DB) StaffRepo {
	return &staffRepo{db: db}
}

func (r *staffRepo) InsertStaff(staff model.Staff, hashPassword string) error {

	query := `INSERT INTO staff (id, nip, name, role, password, "createdAt") VALUES ($1, $2, $3, $4, $5, NOW())`
	_, err := r.db.Exec(query, staff.UserId, staff.NIP, staff.Name, staff.Role, string(hashPassword))
	if err != nil {
		return err
	}
	return nil
}
