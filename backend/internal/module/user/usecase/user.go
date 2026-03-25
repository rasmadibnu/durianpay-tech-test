package usecase

import (
	"github.com/durianpay/fullstack-boilerplate/internal/entity"
	"github.com/durianpay/fullstack-boilerplate/internal/module/user/repository"
	"golang.org/x/crypto/bcrypt"
)

var validRoles = map[string]bool{
	"cs":        true,
	"operation": true,
}

type UserUsecase interface {
	GetAll(search *string, page, limit int) (entity.PaginatedResult[entity.User], error)
	GetByID(id int) (*entity.User, error)
	Create(email, password, role string) (*entity.User, error)
	Update(id int, email, role string) (*entity.User, error)
	UpdatePassword(id int, password string) error
	Delete(id int) error
}

type UserUC struct {
	repo repository.UserCRUDRepository
}

func NewUserUsecase(repo repository.UserCRUDRepository) *UserUC {
	return &UserUC{repo: repo}
}

func (u *UserUC) GetAll(search *string, page, limit int) (entity.PaginatedResult[entity.User], error) {
	users, total, err := u.repo.GetAll(search, page, limit)
	if err != nil {
		return entity.PaginatedResult[entity.User]{}, err
	}
	return entity.NewPaginatedResult(users, page, limit, total), nil
}

func (u *UserUC) GetByID(id int) (*entity.User, error) {
	return u.repo.GetByID(id)
}

func (u *UserUC) Create(email, password, role string) (*entity.User, error) {
	if email == "" {
		return nil, entity.ErrorBadRequest("email is required")
	}
	if password == "" {
		return nil, entity.ErrorBadRequest("password is required")
	}
	if !validRoles[role] {
		return nil, entity.ErrorBadRequest("role must be one of: cs, operation")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, entity.WrapError(err, entity.ErrorCodeInternal, "failed to hash password")
	}

	return u.repo.Create(email, string(hash), role)
}

func (u *UserUC) Update(id int, email, role string) (*entity.User, error) {
	if email == "" {
		return nil, entity.ErrorBadRequest("email is required")
	}
	if !validRoles[role] {
		return nil, entity.ErrorBadRequest("role must be one of: cs, operation")
	}
	return u.repo.Update(id, email, role)
}

func (u *UserUC) UpdatePassword(id int, password string) error {
	if password == "" {
		return entity.ErrorBadRequest("password is required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return entity.WrapError(err, entity.ErrorCodeInternal, "failed to hash password")
	}

	return u.repo.UpdatePassword(id, string(hash))
}

func (u *UserUC) Delete(id int) error {
	return u.repo.Delete(id)
}
