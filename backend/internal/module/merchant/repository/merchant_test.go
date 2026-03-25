package repository

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupMerchantTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:?_foreign_keys=1")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE merchants (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{"Tokopedia", "Shopee", "Bukalapak", "Grab", "Gojek"} {
		_, err := db.Exec("INSERT INTO merchants(name) VALUES (?)", name)
		if err != nil {
			t.Fatal(err)
		}
	}

	return db
}

func TestMerchant_GetAll(t *testing.T) {
	db := setupMerchantTestDB(t)
	defer db.Close()

	repo := NewMerchantRepo(db)
	merchants, total, err := repo.GetAll(nil, 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}
	if len(merchants) != 5 {
		t.Errorf("expected 5 merchants, got %d", len(merchants))
	}
}

func TestMerchant_GetAll_Pagination(t *testing.T) {
	db := setupMerchantTestDB(t)
	defer db.Close()

	repo := NewMerchantRepo(db)
	merchants, total, err := repo.GetAll(nil, 1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}
	if len(merchants) != 2 {
		t.Errorf("expected 2 merchants on page 1, got %d", len(merchants))
	}

	merchants2, _, err := repo.GetAll(nil, 3, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(merchants2) != 1 {
		t.Errorf("expected 1 merchant on page 3, got %d", len(merchants2))
	}
}

func TestMerchant_GetAll_Search(t *testing.T) {
	db := setupMerchantTestDB(t)
	defer db.Close()

	repo := NewMerchantRepo(db)
	search := "tok"
	merchants, total, err := repo.GetAll(&search, 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
	if len(merchants) != 1 {
		t.Fatalf("expected 1 merchant, got %d", len(merchants))
	}
	if merchants[0].Name != "Tokopedia" {
		t.Errorf("expected Tokopedia, got %s", merchants[0].Name)
	}
}

func TestMerchant_Create(t *testing.T) {
	db := setupMerchantTestDB(t)
	defer db.Close()

	repo := NewMerchantRepo(db)
	m, err := repo.Create("NewMerchant")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Name != "NewMerchant" {
		t.Errorf("expected name NewMerchant, got %s", m.Name)
	}
	if m.ID == 0 {
		t.Error("expected non-zero id")
	}
}

func TestMerchant_GetByID(t *testing.T) {
	db := setupMerchantTestDB(t)
	defer db.Close()

	repo := NewMerchantRepo(db)
	m, err := repo.GetByID(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Name != "Tokopedia" {
		t.Errorf("expected Tokopedia, got %s", m.Name)
	}
}

func TestMerchant_GetByID_NotFound(t *testing.T) {
	db := setupMerchantTestDB(t)
	defer db.Close()

	repo := NewMerchantRepo(db)
	_, err := repo.GetByID(999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMerchant_Update(t *testing.T) {
	db := setupMerchantTestDB(t)
	defer db.Close()

	repo := NewMerchantRepo(db)
	m, err := repo.Update(1, "UpdatedName")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Name != "UpdatedName" {
		t.Errorf("expected UpdatedName, got %s", m.Name)
	}
}

func TestMerchant_Update_NotFound(t *testing.T) {
	db := setupMerchantTestDB(t)
	defer db.Close()

	repo := NewMerchantRepo(db)
	_, err := repo.Update(999, "X")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMerchant_Delete(t *testing.T) {
	db := setupMerchantTestDB(t)
	defer db.Close()

	repo := NewMerchantRepo(db)
	err := repo.Delete(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = repo.GetByID(1)
	if err == nil {
		t.Fatal("expected not found after delete")
	}
}

func TestMerchant_Delete_NotFound(t *testing.T) {
	db := setupMerchantTestDB(t)
	defer db.Close()

	repo := NewMerchantRepo(db)
	err := repo.Delete(999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
