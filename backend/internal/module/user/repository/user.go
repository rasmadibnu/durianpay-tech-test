package repository

import (
	"database/sql"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
)

type UserCRUDRepository interface {
	GetAll(search *string, page, limit int) ([]entity.User, int, error)
	GetByID(id int) (*entity.User, error)
	Create(email, passwordHash, role string) (*entity.User, error)
	Update(id int, email, role string) (*entity.User, error)
	UpdatePassword(id int, passwordHash string) error
	Delete(id int) error
}

type UserCRUD struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserCRUD {
	return &UserCRUD{db: db}
}

func (r *UserCRUD) GetAll(search *string, page, limit int) ([]entity.User, int, error) {
	whereClause := ""
	var args []interface{}
	if search != nil && *search != "" {
		whereClause = " WHERE email LIKE ? OR role LIKE ?"
		q := "%" + *search + "%"
		args = append(args, q, q)
	}

	var total int
	if err := r.db.QueryRow("SELECT COUNT(*) FROM users"+whereClause, args...).Scan(&total); err != nil {
		return nil, 0, entity.WrapError(err, entity.ErrorCodeInternal, "failed to count users")
	}

	offset := (page - 1) * limit
	rows, err := r.db.Query("SELECT id, email, password_hash, role FROM users"+whereClause+" ORDER BY id ASC LIMIT ? OFFSET ?", append(args, limit, offset)...)
	if err != nil {
		return nil, 0, entity.WrapError(err, entity.ErrorCodeInternal, "failed to query users")
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var u entity.User
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role); err != nil {
			return nil, 0, entity.WrapError(err, entity.ErrorCodeInternal, "failed to scan user")
		}
		users = append(users, u)
	}
	return users, total, nil
}

func (r *UserCRUD) GetByID(id int) (*entity.User, error) {
	row := r.db.QueryRow("SELECT id, email, password_hash, role FROM users WHERE id = ?", id)
	var u entity.User
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role); err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrorNotFound("user not found")
		}
		return nil, entity.WrapError(err, entity.ErrorCodeInternal, "failed to get user")
	}
	return &u, nil
}

func (r *UserCRUD) Create(email, passwordHash, role string) (*entity.User, error) {
	res, err := r.db.Exec("INSERT INTO users(email, password_hash, role) VALUES (?, ?, ?)", email, passwordHash, role)
	if err != nil {
		return nil, entity.WrapError(err, entity.ErrorCodeInternal, "failed to create user")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, entity.WrapError(err, entity.ErrorCodeInternal, "failed to get user id")
	}
	return r.GetByID(int(id))
}

func (r *UserCRUD) Update(id int, email, role string) (*entity.User, error) {
	res, err := r.db.Exec("UPDATE users SET email = ?, role = ? WHERE id = ?", email, role, id)
	if err != nil {
		return nil, entity.WrapError(err, entity.ErrorCodeInternal, "failed to update user")
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, entity.WrapError(err, entity.ErrorCodeInternal, "failed to check rows affected")
	}
	if affected == 0 {
		return nil, entity.ErrorNotFound("user not found")
	}
	return r.GetByID(id)
}

func (r *UserCRUD) UpdatePassword(id int, passwordHash string) error {
	res, err := r.db.Exec("UPDATE users SET password_hash = ? WHERE id = ?", passwordHash, id)
	if err != nil {
		return entity.WrapError(err, entity.ErrorCodeInternal, "failed to update password")
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return entity.WrapError(err, entity.ErrorCodeInternal, "failed to check rows affected")
	}
	if affected == 0 {
		return entity.ErrorNotFound("user not found")
	}
	return nil
}

func (r *UserCRUD) Delete(id int) error {
	res, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return entity.WrapError(err, entity.ErrorCodeInternal, "failed to delete user")
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return entity.WrapError(err, entity.ErrorCodeInternal, "failed to check rows affected")
	}
	if affected == 0 {
		return entity.ErrorNotFound("user not found")
	}
	return nil
}
