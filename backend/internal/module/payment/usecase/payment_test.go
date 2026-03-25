package usecase

import (
	"errors"
	"testing"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
)

type paymentRepoStub struct {
	createArg       *entity.Payment
	updateArg       *entity.Payment
	updateStatusID  string
	updateStatusVal string
	deleteID        string

	createErr       error
	updateErr       error
	updateStatusErr error
	deleteErr       error
}

func (r *paymentRepoStub) GetPayments(status, id, sort, search *string, page, limit int) ([]entity.Payment, int, error) {
	return nil, 0, nil
}

func (r *paymentRepoStub) CreatePayment(payment *entity.Payment) error {
	r.createArg = payment
	return r.createErr
}

func (r *paymentRepoStub) UpdatePayment(payment *entity.Payment) error {
	r.updateArg = payment
	return r.updateErr
}

func (r *paymentRepoStub) UpdatePaymentStatus(id string, status string) error {
	r.updateStatusID = id
	r.updateStatusVal = status
	return r.updateStatusErr
}

func (r *paymentRepoStub) DeletePayment(id string) error {
	r.deleteID = id
	return r.deleteErr
}

func TestPaymentUC_CreatePayment_ValidatesRequiredFields(t *testing.T) {
	tests := []struct {
		name      string
		merchant  int
		amount    string
		status    string
		wantCode  entity.Code
		wantMsg   string
	}{
		{
			name:     "missing merchant id",
			merchant: 0,
			amount:   "100.00",
			status:   "completed",
			wantCode: entity.ErrorCodeBadRequest,
			wantMsg:  "merchant_id is required",
		},
		{
			name:     "missing amount",
			merchant: 12,
			amount:   "",
			status:   "completed",
			wantCode: entity.ErrorCodeBadRequest,
			wantMsg:  "amount is required",
		},
		{
			name:     "invalid status",
			merchant: 12,
			amount:   "100.00",
			status:   "pending",
			wantCode: entity.ErrorCodeBadRequest,
			wantMsg:  "status must be one of: completed, processing, failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewPaymentUsecase(&paymentRepoStub{})

			payment, err := uc.CreatePayment(tt.merchant, tt.amount, tt.status)
			if payment != nil {
				t.Fatalf("expected no payment, got %+v", payment)
			}

			var appErr *entity.AppError
			if !errors.As(err, &appErr) {
				t.Fatalf("expected AppError, got %v", err)
			}
			if appErr.Code != tt.wantCode {
				t.Fatalf("expected code %s, got %s", tt.wantCode, appErr.Code)
			}
			if appErr.Message != tt.wantMsg {
				t.Fatalf("expected message %q, got %q", tt.wantMsg, appErr.Message)
			}
		})
	}
}

func TestPaymentUC_CreatePayment_PersistsGeneratedPayment(t *testing.T) {
	repo := &paymentRepoStub{}
	uc := NewPaymentUsecase(repo)

	payment, err := uc.CreatePayment(9, "25000.00", "processing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if payment == nil {
		t.Fatal("expected payment to be returned")
	}
	if repo.createArg == nil {
		t.Fatal("expected repository CreatePayment to be called")
	}
	if payment != repo.createArg {
		t.Fatal("expected returned payment to match repository input")
	}
	if payment.MerchantID != 9 {
		t.Fatalf("expected merchant id 9, got %d", payment.MerchantID)
	}
	if payment.Amount != "25000.00" {
		t.Fatalf("expected amount 25000.00, got %s", payment.Amount)
	}
	if payment.Status != "processing" {
		t.Fatalf("expected status processing, got %s", payment.Status)
	}
	if len(payment.ID) != len("PAY-12345678") || payment.ID[:4] != "PAY-" {
		t.Fatalf("expected generated payment id, got %s", payment.ID)
	}
	if payment.CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to be set")
	}
}

func TestPaymentUC_UpdatePayment_DelegatesValidatedPayload(t *testing.T) {
	repo := &paymentRepoStub{}
	uc := NewPaymentUsecase(repo)

	payment, err := uc.UpdatePayment("PAY-0001", 4, "75000.00", "completed")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if payment == nil {
		t.Fatal("expected payment to be returned")
	}
	if repo.updateArg == nil {
		t.Fatal("expected repository UpdatePayment to be called")
	}
	if repo.updateArg.ID != "PAY-0001" || repo.updateArg.MerchantID != 4 || repo.updateArg.Amount != "75000.00" || repo.updateArg.Status != "completed" {
		t.Fatalf("unexpected update payload: %+v", repo.updateArg)
	}
}

func TestPaymentUC_UpdatePaymentStatus_RejectsInvalidStatus(t *testing.T) {
	repo := &paymentRepoStub{}
	uc := NewPaymentUsecase(repo)

	err := uc.UpdatePaymentStatus("PAY-0001", "queued")

	var appErr *entity.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %v", err)
	}
	if appErr.Code != entity.ErrorCodeBadRequest {
		t.Fatalf("expected bad_request, got %s", appErr.Code)
	}
	if repo.updateStatusID != "" || repo.updateStatusVal != "" {
		t.Fatal("expected repository UpdatePaymentStatus not to be called")
	}
}

func TestPaymentUC_DeletePayment_DelegatesToRepository(t *testing.T) {
	repo := &paymentRepoStub{}
	uc := NewPaymentUsecase(repo)

	if err := uc.DeletePayment("PAY-0010"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.deleteID != "PAY-0010" {
		t.Fatalf("expected delete id PAY-0010, got %s", repo.deleteID)
	}
}
