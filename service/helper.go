package service

import (
	"halo-suster/model"
	"strconv"
)

func getRoleFromNIP(nip int64) model.Role {
	nipStr := strconv.FormatInt(nip, 10)

	switch nipStr[:3] {
	case "615":
		return model.RoleIt
	case "303":
		return model.RoleNurse
	}
	return model.RoleUnknown
}
