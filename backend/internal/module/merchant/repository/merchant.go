package repository

import (
	"database/sql"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
)

type MerchantRepository interface {
	Create(name string) (*entity.Merchant, error)
	GetAll(search *string, page, limit int) ([]entity.Merchant, int, error)
	GetByID(id int) (*entity.Merchant, error)
	Update(id int, name string) (*entity.Merchant, error)
	Delete(id int) error
}

type Merchant struct {
	db *sql.DB
}

func NewMerchantRepo(db *sql.DB) *Merchant {
	return &Merchant{db: db}
}

func (r *Merchant) Create(name string) (*entity.Merchant, error) {
	res, err := r.db.Exec("INSERT INTO merchants(name) VALUES (?)", name)
	if err != nil {
		return nil, entity.WrapError(err, entity.ErrorCodeInternal, "failed to create merchant")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, entity.WrapError(err, entity.ErrorCodeInternal, "failed to get merchant id")
	}
	return r.GetByID(int(id))
}

func (r *Merchant) GetAll(search *string, page, limit int) ([]entity.Merchant, int, error) {
	whereClause := ""
	var args []interface{}
	if search != nil && *search != "" {
		whereClause = " WHERE name LIKE ?"
		args = append(args, "%"+*search+"%")
	}

	var total int
	if err := r.db.QueryRow("SELECT COUNT(*) FROM merchants"+whereClause, args...).Scan(&total); err != nil {
		return nil, 0, entity.WrapError(err, entity.ErrorCodeInternal, "failed to count merchants")
	}

	offset := (page - 1) * limit
	rows, err := r.db.Query("SELECT id, name, created_at, updated_at FROM merchants"+whereClause+" ORDER BY created_at DESC LIMIT ? OFFSET ?", append(args, limit, offset)...)
	if err != nil {
		return nil, 0, entity.WrapError(err, entity.ErrorCodeInternal, "failed to query merchants")
	}
	defer rows.Close()

	var merchants []entity.Merchant
	for rows.Next() {
		var m entity.Merchant
		if err := rows.Scan(&m.ID, &m.Name, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, 0, entity.WrapError(err, entity.ErrorCodeInternal, "failed to scan merchant")
		}
		merchants = append(merchants, m)
	}
	return merchants, total, nil
}

func (r *Merchant) GetByID(id int) (*entity.Merchant, error) {
	row := r.db.QueryRow("SELECT id, name, created_at, updated_at FROM merchants WHERE id = ?", id)
	var m entity.Merchant
	if err := row.Scan(&m.ID, &m.Name, &m.CreatedAt, &m.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrorNotFound("merchant not found")
		}
		return nil, entity.WrapError(err, entity.ErrorCodeInternal, "failed to get merchant")
	}
	return &m, nil
}

func (r *Merchant) Update(id int, name string) (*entity.Merchant, error) {
	res, err := r.db.Exec("UPDATE merchants SET name = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", name, id)
	if err != nil {
		return nil, entity.WrapError(err, entity.ErrorCodeInternal, "failed to update merchant")
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, entity.WrapError(err, entity.ErrorCodeInternal, "failed to check rows affected")
	}
	if affected == 0 {
		return nil, entity.ErrorNotFound("merchant not found")
	}
	return r.GetByID(id)
}

func (r *Merchant) Delete(id int) error {
	res, err := r.db.Exec("DELETE FROM merchants WHERE id = ?", id)
	if err != nil {
		return entity.WrapError(err, entity.ErrorCodeInternal, "failed to delete merchant")
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return entity.WrapError(err, entity.ErrorCodeInternal, "failed to check rows affected")
	}
	if affected == 0 {
		return entity.ErrorNotFound("merchant not found")
	}
	return nil
}
