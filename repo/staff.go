package repo

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"halo-suster/model"
	"strings"

	"github.com/jmoiron/sqlx"
)

type StaffRepo interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (staff model.Staff, err error)
	GetUserByNIP(ctx context.Context, nip int64) (staff model.Staff, err error)
	InsertStaff(ctx context.Context, staff model.Staff, hashPassword string) error
	GetListUser(ctx context.Context, params model.GetListUserParams) (listUser []model.Staff, err error)
	UpdateUser(ctx context.Context, staff model.Staff) error
	SoftDeleteUser(ctx context.Context, userId uuid.UUID) error
}

type staffRepo struct {
	db *sqlx.DB
}

func NewStaffRepo(db *sqlx.DB) StaffRepo {
	return &staffRepo{db: db}
}

func (r *staffRepo) InsertStaff(ctx context.Context, staff model.Staff, hashPassword string) error {

	query := `INSERT INTO staff (id, nip, name, role, password, "identityCardScanImg", status, "createdAt") VALUES ($1, $2, $3, $4, $5,$6,$7, NOW())`
	_, err := r.db.ExecContext(ctx, query, staff.UserId, staff.NIP, staff.Name, staff.Role, hashPassword, staff.IdentityCardScanImg, staff.Status)
	if err != nil {
		return err
	}
	return nil
}

// GetUserByID retrieves a user by ID
func (r *staffRepo) GetUserByID(ctx context.Context, id uuid.UUID) (model.Staff, error) {
	var staff model.Staff
	query := `SELECT * FROM staff WHERE id = $1`
	err := r.db.GetContext(ctx, &staff, query, id)
	if err != nil {
		return model.Staff{}, err
	}
	return staff, nil
}

// GetUserByNIP retrieves a user by NIP
func (r *staffRepo) GetUserByNIP(ctx context.Context, nip int64) (model.Staff, error) {
	var staff model.Staff
	query := `SELECT "id", nip, name, role, password, "createdAt" FROM staff WHERE nip = $1`
	err := r.db.GetContext(ctx, &staff, query, nip)
	if err != nil {
		return model.Staff{}, err
	}
	return staff, nil
}

func (r *staffRepo) GetListUser(ctx context.Context, params model.GetListUserParams) (listUser []model.Staff, err error) {
	query := `SELECT * FROM staff WHERE true ` + generateGetListUserSQLFilter(params)
	rows, err := r.db.QueryxContext(ctx, query)
	if err != nil {
		return listUser, err
	}
	defer rows.Close()
	for rows.Next() {
		var staff model.Staff
		if err := rows.StructScan(&staff); err != nil {
			return listUser, err
		}
		listUser = append(listUser, staff)
	}
	return listUser, nil
}

func generateGetListUserSQLFilter(params model.GetListUserParams) string {
	var conditions []string

	// Add conditions based on the fields provided
	if params.ID != nil {
		conditions = append(conditions, fmt.Sprintf(`"id" = '%s'`, *params.ID))
	}

	// TODO: explore searchable index
	if params.Name != nil {
		name := *params.Name
		// Append wildcard symbols to allow partial matching
		name = "%" + strings.ToLower(name) + "%"
		conditions = append(conditions, fmt.Sprintf(`lower("name") LIKE '%s'`, name))
	}
	if params.NIP != nil {

		nip := *params.NIP
		nipStr := fmt.Sprintf("%d", nip)
		ns := "%" + fmt.Sprintf("%d", nip) + "%"
		if len(nipStr) > 12 {
			ns = fmt.Sprintf("%d", nip)
		}
		conditions = append(conditions, fmt.Sprintf(`"nip"::text LIKE '%s'`, ns))
	}
	if params.Role != nil {
		conditions = append(conditions, fmt.Sprintf(`"role" = '%s'`, *params.Role))
	}
	if params.Status != nil {
		conditions = append(conditions, fmt.Sprintf(`"status" = '%s'`, *params.Status))
	} else {
		conditions = append(conditions, fmt.Sprintf(`"status" = 'active'`))
	}

	// Combine conditions with AND
	filter := strings.Join(conditions, " AND ")
	if filter != "" {
		// add and clause in the front
		filter = "AND " + filter
	}

	orderByClause := ""
	// set default sort
	if params.Sort.CreatedAt == nil {
		defaultSort := "desc"
		params.Sort.CreatedAt = &defaultSort
	}
	orderByClause = fmt.Sprintf(` ORDER BY "createdAt" %s`, *params.Sort.CreatedAt)
	filter += orderByClause

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

func (r *staffRepo) UpdateUser(ctx context.Context, staff model.Staff) error {
	// UpdateStaff updates a staff member's details in the database
	query := `
		UPDATE staff
		SET nip = :nip,
		    name = :name,
		    role = :role,
		    "identityCardScanImg" = :identityCardScanImg,
		    status = :status,
		    password = :password
		WHERE id = :id
	`

	_, err := r.db.NamedExecContext(ctx, query, staff)
	return err
}

// SoftDeleteUser performs a soft delete by updating the status to 'deleted'
func (r *staffRepo) SoftDeleteUser(ctx context.Context, userId uuid.UUID) error {
	query := `UPDATE staff SET status = 'deleted' WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, userId)
	return err
}
