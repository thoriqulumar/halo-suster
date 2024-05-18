package repo

import (
	"context"
	"fmt"
	"helo-suster/model"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

type StaffRepo interface {
	GetStaff(ctx context.Context, param model.GetStaffRequest) (staffList []model.Staff, err error)
}

type staffRepo struct {
	db *sqlx.DB
}

func NewStaffRepo(db *sqlx.DB) StaffRepo {
	return &staffRepo{
		db: db,
	}
}

func (r *staffRepo) GetStaff(ctx context.Context, param model.GetStaffRequest) (staffList []model.Staff, err error) {
	query := `SELECT * FROM staff WHERE 1=1 ` + generateGetStaffSQLFilter(param)

	rows, err := r.db.QueryxContext(ctx, query)
	if err != nil {
		return staffList, err
	}
	defer rows.Close()
	for rows.Next() {
		var staff model.Staff
		if err := rows.StructScan(&staff); err != nil {
			return staffList, err
		}
		staffList = append(staffList, staff)
	}
	return staffList, nil
}

func generateGetStaffSQLFilter(params model.GetStaffRequest) string {
	var conditions []string

	v := reflect.ValueOf(params)
	t := reflect.TypeOf(params)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		typeField := t.Field(i)
		dbTag := typeField.Tag.Get("schema")

		// Skip limit and offset fields
		if dbTag == "limit" || dbTag == "offset" || dbTag == "createdAt" {
			continue
		}

		if !field.IsNil() {
			switch field.Kind() {
			case reflect.Ptr:
				value := field.Elem()
				switch value.Kind() {
				case reflect.String:
					if dbTag == "name" {
						// Special handling for partial matching
						name := "%" + strings.ToLower(value.String()) + "%"
						conditions = append(conditions, fmt.Sprintf(`lower("%s") LIKE '%s'`, dbTag, name))
					} else {
						conditions = append(conditions, fmt.Sprintf(`"%s" = '%s'`, dbTag, value.String()))
					}
				case reflect.Int32, reflect.Int:
					if dbTag == "nip" {
						// Convert nip to string and handle partial matching
						nipStr := fmt.Sprintf("%%%d%%", value.Int())
						conditions = append(conditions, fmt.Sprintf(`"%s" LIKE '%s'`, dbTag, nipStr))
					} else {
						conditions = append(conditions, fmt.Sprintf(`"%s" = %d`, dbTag, value.Int()))
					}
				}
			}
		}
	}

	// Combine conditions with AND
	filter := strings.Join(conditions, " AND ")
	if filter != "" {
		filter = "WHERE " + filter
	}

	orderByClause := ""
	if params.CreatedAt != nil {
		orderByClause = fmt.Sprintf(` ORDER BY "createdAt" %s`, *params.CreatedAt)
	}
	if orderByClause != "" {
		filter += " " + orderByClause
	}

	// Add additional clauses such as LIMIT and OFFSET
	if params.Limit != nil {
		filter += fmt.Sprintf(" LIMIT %d", *params.Limit)
	} else {
		filter += " LIMIT 5"
	}
	if params.Offset != nil {
		filter += fmt.Sprintf(" OFFSET %d", *params.Offset)
	} else {
		filter += " OFFSET 0"
	}

	return filter
}
