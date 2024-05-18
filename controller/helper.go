package controller

import (
	"github.com/google/uuid"
	"halo-suster/model"
	"net/url"
	"strconv"
)

func parseGetUserParams(params url.Values) model.GetListUserParams {
	var result model.GetListUserParams

	for key, values := range params {
		switch key {
		case "userId":
			id, err := uuid.Parse(values[0])
			if err == nil {
				result.ID = &id
			}
		case "limit":
			limit, err := strconv.Atoi(values[0])
			if err == nil {
				result.Limit = &limit
			}
		case "offset":
			offset, err := strconv.Atoi(values[0])
			if err == nil {
				result.Offset = &offset
			}
		case "name":
			result.Name = &values[0]
		case "nip":
			nip, err := strconv.ParseInt(values[0], 10, 64)
			if err == nil {
				result.NIP = &nip
			}
		case "role":
			temp := model.Role(values[0])
			result.Role = &temp
		case "status":
			temp := model.Status(values[0])
			result.Status = &temp
		case "createdAt":
			result.Sort.CreatedAt = &values[0]
		}

	}

	return result
}
