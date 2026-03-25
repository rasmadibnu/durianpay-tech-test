package usecase

import (
	"errors"
	"testing"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
	"golang.org/x/crypto/bcrypt"
)

type userRepoStub struct {
	createEmail        string
	createPasswordHash string
	createRole         string
	updateID           int
	updateEmail        string
	updateRole         string
	updatePasswordID   int
	updatePasswordHash string
	deleteID           int

	createRes *entity.User
	updateRes *entity.User
	getByIDRes *entity.User
	listRes   []entity.User
	total     int
	err       error
}

func (r *userRepoStub) GetAll(search *string, page, limit int) ([]entity.User, int, error) {
	return r.listRes, r.total, r.err
}

func (r *userRepoStub) GetByID(id int) (*entity.User, error) {
	return r.getByIDRes, r.err
}

func (r *userRepoStub) Create(email, passwordHash, role string) (*entity.User, error) {
	r.createEmail = email
	r.createPasswordHash = passwordHash
	r.createRole = role
	return r.createRes, r.err
}

func (r *userRepoStub) Update(id int, email, role string) (*entity.User, error) {
	r.updateID = id
	r.updateEmail = email
	r.updateRole = role
	return r.updateRes, r.err
}

func (r *userRepoStub) UpdatePassword(id int, passwordHash string) error {
	r.updatePasswordID = id
	r.updatePasswordHash = passwordHash
	return r.err
}

func (r *userRepoStub) Delete(id int) error {
	r.deleteID = id
	return r.err
}

func TestUserUC_Create_ValidatesInputs(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		pass    string
		role    string
		wantMsg string
	}{
		{"missing email", "", "password", "cs", "email is required"},
		{"missing password", "user@test.com", "", "cs", "password is required"},
		{"invalid role", "user@test.com", "password", "admin", "role must be one of: cs, operation"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewUserUsecase(&userRepoStub{})
			user, err := uc.Create(tt.email, tt.pass, tt.role)
			if user != nil {
				t.Fatalf("expected nil user, got %+v", user)
			}

			var appErr *entity.AppError
			if !errors.As(err, &appErr) {
				t.Fatalf("expected AppError, got %v", err)
			}
			if appErr.Message != tt.wantMsg {
				t.Fatalf("expected %q, got %q", tt.wantMsg, appErr.Message)
			}
		})
	}
}

func TestUserUC_Create_HashesPasswordBeforePersisting(t *testing.T) {
	repo := &userRepoStub{
		createRes: &entity.User{ID: "1", Email: "user@test.com", Role: "cs"},
	}
	uc := NewUserUsecase(repo)

	user, err := uc.Create("user@test.com", "password", "cs")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.createEmail != "user@test.com" || repo.createRole != "cs" {
		t.Fatalf("unexpected create call: email=%s role=%s", repo.createEmail, repo.createRole)
	}
	if repo.createPasswordHash == "" || repo.createPasswordHash == "password" {
		t.Fatal("expected hashed password to be stored")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(repo.createPasswordHash), []byte("password")); err != nil {
		t.Fatalf("expected valid password hash, got %v", err)
	}
	if user == nil || user.Email != "user@test.com" {
		t.Fatalf("unexpected user result: %+v", user)
	}
}

func TestUserUC_Update_ValidatesAndDelegates(t *testing.T) {
	repo := &userRepoStub{
		updateRes: &entity.User{ID: "1", Email: "new@test.com", Role: "operation"},
	}
	uc := NewUserUsecase(repo)

	_, err := uc.Update(1, "", "cs")
	var appErr *entity.AppError
	if !errors.As(err, &appErr) || appErr.Message != "email is required" {
		t.Fatalf("expected email validation error, got %v", err)
	}

	_, err = uc.Update(1, "new@test.com", "admin")
	if !errors.As(err, &appErr) || appErr.Message != "role must be one of: cs, operation" {
		t.Fatalf("expected role validation error, got %v", err)
	}

	user, err := uc.Update(1, "new@test.com", "operation")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.updateID != 1 || repo.updateEmail != "new@test.com" || repo.updateRole != "operation" {
		t.Fatalf("unexpected update call: %+v", repo)
	}
	if user == nil || user.Role != "operation" {
		t.Fatalf("unexpected user result: %+v", user)
	}
}

func TestUserUC_UpdatePassword_ValidatesAndHashes(t *testing.T) {
	repo := &userRepoStub{}
	uc := NewUserUsecase(repo)

	err := uc.UpdatePassword(5, "")
	var appErr *entity.AppError
	if !errors.As(err, &appErr) || appErr.Message != "password is required" {
		t.Fatalf("expected password validation error, got %v", err)
	}

	if err := uc.UpdatePassword(5, "new-password"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.updatePasswordID != 5 {
		t.Fatalf("expected update password id 5, got %d", repo.updatePasswordID)
	}
	if repo.updatePasswordHash == "" || repo.updatePasswordHash == "new-password" {
		t.Fatal("expected hashed password to be stored")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(repo.updatePasswordHash), []byte("new-password")); err != nil {
		t.Fatalf("expected valid password hash, got %v", err)
	}
}

func TestUserUC_Delete_Delegates(t *testing.T) {
	repo := &userRepoStub{}
	uc := NewUserUsecase(repo)

	if err := uc.Delete(3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.deleteID != 3 {
		t.Fatalf("expected delete id 3, got %d", repo.deleteID)
	}
}
