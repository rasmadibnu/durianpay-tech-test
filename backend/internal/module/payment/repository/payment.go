package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
)

type PaymentRepository interface {
	GetPayments(status, id, sort, search *string, page, limit int) ([]entity.Payment, int, error)
	CreatePayment(payment *entity.Payment) error
	UpdatePayment(payment *entity.Payment) error
	UpdatePaymentStatus(id string, status string) error
	DeletePayment(id string) error
}

type Payment struct {
	db *sql.DB
}

func NewPaymentRepo(db *sql.DB) *Payment {
	return &Payment{db: db}
}

func (r *Payment) GetPayments(status, id, sort, search *string, page, limit int) ([]entity.Payment, int, error) {
	baseQuery := "FROM payments p LEFT JOIN merchants m ON p.merchant_id = m.id"
	var conditions []string
	var args []interface{}

	if status != nil && *status != "" {
		conditions = append(conditions, "p.status = ?")
		args = append(args, *status)
	}
	if id != nil && *id != "" {
		conditions = append(conditions, "p.id = ?")
		args = append(args, *id)
	}
	if search != nil && *search != "" {
		conditions = append(conditions, "(p.id LIKE ? OR m.name LIKE ? OR p.amount LIKE ?)")
		q := "%" + *search + "%"
		args = append(args, q, q, q)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	var total int
	countQuery := "SELECT COUNT(*) " + baseQuery + whereClause
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, entity.WrapError(err, entity.ErrorCodeInternal, "failed to count payments")
	}

	// Select with pagination
	query := "SELECT p.id, p.merchant_id, m.name, p.amount, p.status, p.created_at " + baseQuery + whereClause

	if sort != nil && *sort != "" {
		orderClauses := parseSortParam(*sort)
		if len(orderClauses) > 0 {
			query += " ORDER BY " + strings.Join(orderClauses, ", ")
		}
	} else {
		query += " ORDER BY p.created_at DESC"
	}

	offset := (page - 1) * limit
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, entity.WrapError(err, entity.ErrorCodeInternal, "failed to query payments")
	}
	defer rows.Close()

	var payments []entity.Payment
	for rows.Next() {
		var p entity.Payment
		if err := rows.Scan(&p.ID, &p.MerchantID, &p.MerchantName, &p.Amount, &p.Status, &p.CreatedAt); err != nil {
			return nil, 0, entity.WrapError(err, entity.ErrorCodeInternal, "failed to scan payment")
		}
		payments = append(payments, p)
	}
	return payments, total, nil
}

func (r *Payment) CreatePayment(payment *entity.Payment) error {
	_, err := r.db.Exec(
		"INSERT INTO payments(id, merchant_id, amount, status, created_at) VALUES (?, ?, ?, ?, ?)",
		payment.ID, payment.MerchantID, payment.Amount, payment.Status, payment.CreatedAt,
	)
	if err != nil {
		return entity.WrapError(err, entity.ErrorCodeInternal, "failed to create payment")
	}
	return nil
}

func (r *Payment) UpdatePayment(payment *entity.Payment) error {
	res, err := r.db.Exec(
		"UPDATE payments SET merchant_id = ?, amount = ?, status = ? WHERE id = ?",
		payment.MerchantID, payment.Amount, payment.Status, payment.ID,
	)
	if err != nil {
		return entity.WrapError(err, entity.ErrorCodeInternal, "failed to update payment")
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return entity.WrapError(err, entity.ErrorCodeInternal, "failed to check rows affected")
	}
	if affected == 0 {
		return entity.ErrorNotFound("payment not found")
	}
	return nil
}

func (r *Payment) UpdatePaymentStatus(id string, status string) error {
	res, err := r.db.Exec("UPDATE payments SET status = ? WHERE id = ?", status, id)
	if err != nil {
		return entity.WrapError(err, entity.ErrorCodeInternal, "failed to update payment status")
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return entity.WrapError(err, entity.ErrorCodeInternal, "failed to check rows affected")
	}
	if affected == 0 {
		return entity.ErrorNotFound("payment not found")
	}
	return nil
}

func (r *Payment) DeletePayment(id string) error {
	res, err := r.db.Exec("DELETE FROM payments WHERE id = ?", id)
	if err != nil {
		return entity.WrapError(err, entity.ErrorCodeInternal, "failed to delete payment")
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return entity.WrapError(err, entity.ErrorCodeInternal, "failed to check rows affected")
	}
	if affected == 0 {
		return entity.ErrorNotFound("payment not found")
	}
	return nil
}

var allowedSortFields = map[string]bool{
	"id":          true,
	"merchant_id": true,
	"amount":      true,
	"status":      true,
	"created_at":  true,
}

func parseSortParam(sort string) []string {
	var clauses []string
	fields := strings.Split(sort, ",")
	for _, f := range fields {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		dir := "ASC"
		if strings.HasPrefix(f, "-") {
			dir = "DESC"
			f = f[1:]
		}
		if allowedSortFields[f] {
			clauses = append(clauses, fmt.Sprintf("%s %s", f, dir))
		}
	}
	return clauses
}
