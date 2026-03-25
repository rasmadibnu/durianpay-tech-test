package usecase

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
	"github.com/durianpay/fullstack-boilerplate/internal/module/payment/repository"
)

var validStatuses = map[string]bool{
	"completed":  true,
	"processing": true,
	"failed":     true,
}

type PaymentUsecase interface {
	GetPayments(status, id, sort, search *string, page, limit int) (entity.PaginatedResult[entity.Payment], error)
	CreatePayment(merchantID int, amount, status string) (*entity.Payment, error)
	UpdatePayment(id string, merchantID int, amount, status string) (*entity.Payment, error)
	UpdatePaymentStatus(id, status string) error
	DeletePayment(id string) error
}

type PaymentUC struct {
	repo repository.PaymentRepository
}

func NewPaymentUsecase(repo repository.PaymentRepository) *PaymentUC {
	return &PaymentUC{repo: repo}
}

func (u *PaymentUC) GetPayments(status, id, sort, search *string, page, limit int) (entity.PaginatedResult[entity.Payment], error) {
	payments, total, err := u.repo.GetPayments(status, id, sort, search, page, limit)
	if err != nil {
		return entity.PaginatedResult[entity.Payment]{}, err
	}
	return entity.NewPaginatedResult(payments, page, limit, total), nil
}

func (u *PaymentUC) CreatePayment(merchantID int, amount, status string) (*entity.Payment, error) {
	if merchantID <= 0 {
		return nil, entity.ErrorBadRequest("merchant_id is required")
	}
	if amount == "" {
		return nil, entity.ErrorBadRequest("amount is required")
	}
	if !validStatuses[status] {
		return nil, entity.ErrorBadRequest("status must be one of: completed, processing, failed")
	}

	payment := &entity.Payment{
		ID:         fmt.Sprintf("PAY-%s", uuid.New().String()[:8]),
		MerchantID: merchantID,
		Amount:     amount,
		Status:     status,
		CreatedAt:  time.Now(),
	}
	if err := u.repo.CreatePayment(payment); err != nil {
		return nil, err
	}
	return payment, nil
}

func (u *PaymentUC) UpdatePayment(id string, merchantID int, amount, status string) (*entity.Payment, error) {
	if merchantID <= 0 {
		return nil, entity.ErrorBadRequest("merchant_id is required")
	}
	if amount == "" {
		return nil, entity.ErrorBadRequest("amount is required")
	}
	if !validStatuses[status] {
		return nil, entity.ErrorBadRequest("status must be one of: completed, processing, failed")
	}

	payment := &entity.Payment{
		ID:         id,
		MerchantID: merchantID,
		Amount:     amount,
		Status:     status,
	}
	if err := u.repo.UpdatePayment(payment); err != nil {
		return nil, err
	}
	return payment, nil
}

func (u *PaymentUC) UpdatePaymentStatus(id, status string) error {
	if !validStatuses[status] {
		return entity.ErrorBadRequest("status must be one of: completed, processing, failed")
	}
	return u.repo.UpdatePaymentStatus(id, status)
}

func (u *PaymentUC) DeletePayment(id string) error {
	return u.repo.DeletePayment(id)
}
