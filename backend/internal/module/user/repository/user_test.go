package repository

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupUserTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:?_foreign_keys=1")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		role TEXT NOT NULL
	)`)
	if err != nil {
		t.Fatal(err)
	}

	users := []struct{ email, hash, role string }{
		{"cs@test.com", "$2a$10$hash1", "cs"},
		{"op@test.com", "$2a$10$hash2", "operation"},
		{"admin@test.com", "$2a$10$hash3", "operation"},
	}
	for _, u := range users {
		_, err := db.Exec("INSERT INTO users(email, password_hash, role) VALUES (?, ?, ?)", u.email, u.hash, u.role)
		if err != nil {
			t.Fatal(err)
		}
	}

	return db
}

func TestUser_GetAll(t *testing.T) {
	db := setupUserTestDB(t)
	defer db.Close()

	repo := NewUserRepo(db)
	users, total, err := repo.GetAll(nil, 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(users) != 3 {
		t.Errorf("expected 3 users, got %d", len(users))
	}
}

func TestUser_GetAll_Pagination(t *testing.T) {
	db := setupUserTestDB(t)
	defer db.Close()

	repo := NewUserRepo(db)
	users, total, err := repo.GetAll(nil, 1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users on page 1, got %d", len(users))
	}

	users2, _, err := repo.GetAll(nil, 2, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users2) != 1 {
		t.Errorf("expected 1 user on page 2, got %d", len(users2))
	}
}

func TestUser_GetAll_Search(t *testing.T) {
	db := setupUserTestDB(t)
	defer db.Close()

	repo := NewUserRepo(db)
	search := "cs"
	users, total, err := repo.GetAll(&search, 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
	if users[0].Email != "cs@test.com" {
		t.Errorf("expected cs@test.com, got %s", users[0].Email)
	}
}

func TestUser_Create(t *testing.T) {
	db := setupUserTestDB(t)
	defer db.Close()

	repo := NewUserRepo(db)
	u, err := repo.Create("new@test.com", "$2a$10$newhash", "cs")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Email != "new@test.com" {
		t.Errorf("expected new@test.com, got %s", u.Email)
	}
	if u.Role != "cs" {
		t.Errorf("expected role cs, got %s", u.Role)
	}
}

func TestUser_GetByID(t *testing.T) {
	db := setupUserTestDB(t)
	defer db.Close()

	repo := NewUserRepo(db)
	u, err := repo.GetByID(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Email != "cs@test.com" {
		t.Errorf("expected cs@test.com, got %s", u.Email)
	}
}

func TestUser_GetByID_NotFound(t *testing.T) {
	db := setupUserTestDB(t)
	defer db.Close()

	repo := NewUserRepo(db)
	_, err := repo.GetByID(999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUser_Update(t *testing.T) {
	db := setupUserTestDB(t)
	defer db.Close()

	repo := NewUserRepo(db)
	u, err := repo.Update(1, "updated@test.com", "operation")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Email != "updated@test.com" {
		t.Errorf("expected updated@test.com, got %s", u.Email)
	}
	if u.Role != "operation" {
		t.Errorf("expected role operation, got %s", u.Role)
	}
}

func TestUser_UpdatePassword(t *testing.T) {
	db := setupUserTestDB(t)
	defer db.Close()

	repo := NewUserRepo(db)
	err := repo.UpdatePassword(1, "$2a$10$newhash")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	u, _ := repo.GetByID(1)
	if u.PasswordHash != "$2a$10$newhash" {
		t.Errorf("expected updated password hash")
	}
}

func TestUser_Delete(t *testing.T) {
	db := setupUserTestDB(t)
	defer db.Close()

	repo := NewUserRepo(db)
	err := repo.Delete(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = repo.GetByID(1)
	if err == nil {
		t.Fatal("expected not found after delete")
	}
}

func TestUser_Delete_NotFound(t *testing.T) {
	db := setupUserTestDB(t)
	defer db.Close()

	repo := NewUserRepo(db)
	err := repo.Delete(999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
