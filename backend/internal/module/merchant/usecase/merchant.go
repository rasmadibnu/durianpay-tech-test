package usecase

import (
	"github.com/durianpay/fullstack-boilerplate/internal/entity"
	"github.com/durianpay/fullstack-boilerplate/internal/module/merchant/repository"
)

type MerchantUsecase interface {
	Create(name string) (*entity.Merchant, error)
	GetAll() ([]entity.Merchant, error)
	GetByID(id int) (*entity.Merchant, error)
	Update(id int, name string) (*entity.Merchant, error)
	Delete(id int) error
}

type MerchantUC struct {
	repo repository.MerchantRepository
}

func NewMerchantUsecase(repo repository.MerchantRepository) *MerchantUC {
	return &MerchantUC{repo: repo}
}

func (u *MerchantUC) Create(name string) (*entity.Merchant, error) {
	if name == "" {
		return nil, entity.ErrorBadRequest("merchant name is required")
	}
	return u.repo.Create(name)
}

func (u *MerchantUC) GetAll() ([]entity.Merchant, error) {
	return u.repo.GetAll()
}

func (u *MerchantUC) GetByID(id int) (*entity.Merchant, error) {
	return u.repo.GetByID(id)
}

func (u *MerchantUC) Update(id int, name string) (*entity.Merchant, error) {
	if name == "" {
		return nil, entity.ErrorBadRequest("merchant name is required")
	}
	return u.repo.Update(id, name)
}

func (u *MerchantUC) Delete(id int) error {
	return u.repo.Delete(id)
}
