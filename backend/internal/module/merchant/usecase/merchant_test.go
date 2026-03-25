package usecase

import (
	"errors"
	"testing"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
)

type merchantRepoStub struct {
	createName string
	updateID   int
	updateName string
	deleteID   int

	createRes *entity.Merchant
	updateRes *entity.Merchant
	getByIDRes *entity.Merchant
	listRes   []entity.Merchant
	total     int
	err       error
}

func (r *merchantRepoStub) Create(name string) (*entity.Merchant, error) {
	r.createName = name
	return r.createRes, r.err
}

func (r *merchantRepoStub) GetAll(search *string, page, limit int) ([]entity.Merchant, int, error) {
	return r.listRes, r.total, r.err
}

func (r *merchantRepoStub) GetByID(id int) (*entity.Merchant, error) {
	return r.getByIDRes, r.err
}

func (r *merchantRepoStub) Update(id int, name string) (*entity.Merchant, error) {
	r.updateID = id
	r.updateName = name
	return r.updateRes, r.err
}

func (r *merchantRepoStub) Delete(id int) error {
	r.deleteID = id
	return r.err
}

func TestMerchantUC_Create_RequiresName(t *testing.T) {
	uc := NewMerchantUsecase(&merchantRepoStub{})

	merchant, err := uc.Create("")
	if merchant != nil {
		t.Fatalf("expected nil merchant, got %+v", merchant)
	}

	var appErr *entity.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %v", err)
	}
	if appErr.Code != entity.ErrorCodeBadRequest {
		t.Fatalf("expected bad_request, got %s", appErr.Code)
	}
}

func TestMerchantUC_Create_DelegatesToRepository(t *testing.T) {
	repo := &merchantRepoStub{
		createRes: &entity.Merchant{ID: 1, Name: "Tokopedia"},
	}
	uc := NewMerchantUsecase(repo)

	merchant, err := uc.Create("Tokopedia")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.createName != "Tokopedia" {
		t.Fatalf("expected create name Tokopedia, got %s", repo.createName)
	}
	if merchant == nil || merchant.Name != "Tokopedia" {
		t.Fatalf("unexpected merchant result: %+v", merchant)
	}
}

func TestMerchantUC_GetAll_ReturnsPaginatedResult(t *testing.T) {
	repo := &merchantRepoStub{
		listRes: []entity.Merchant{{ID: 1, Name: "Tokopedia"}, {ID: 2, Name: "Shopee"}},
		total:   5,
	}
	uc := NewMerchantUsecase(repo)

	result, err := uc.GetAll(nil, 2, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 5 || result.Page != 2 || result.TotalPages != 3 {
		t.Fatalf("unexpected pagination result: %+v", result)
	}
	if len(result.Data) != 2 {
		t.Fatalf("expected 2 merchants, got %d", len(result.Data))
	}
}

func TestMerchantUC_Update_ValidatesAndDelegates(t *testing.T) {
	repo := &merchantRepoStub{
		updateRes: &entity.Merchant{ID: 9, Name: "Updated"},
	}
	uc := NewMerchantUsecase(repo)

	_, err := uc.Update(9, "")
	var appErr *entity.AppError
	if !errors.As(err, &appErr) || appErr.Code != entity.ErrorCodeBadRequest {
		t.Fatalf("expected bad_request for empty name, got %v", err)
	}

	merchant, err := uc.Update(9, "Updated")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.updateID != 9 || repo.updateName != "Updated" {
		t.Fatalf("unexpected update call: id=%d name=%s", repo.updateID, repo.updateName)
	}
	if merchant == nil || merchant.Name != "Updated" {
		t.Fatalf("unexpected merchant result: %+v", merchant)
	}
}

func TestMerchantUC_Delete_Delegates(t *testing.T) {
	repo := &merchantRepoStub{}
	uc := NewMerchantUsecase(repo)

	if err := uc.Delete(7); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.deleteID != 7 {
		t.Fatalf("expected delete id 7, got %d", repo.deleteID)
	}
}
