package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
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

	_, err = db.Exec(`CREATE TABLE payments (
		id TEXT PRIMARY KEY,
		merchant_id INTEGER NOT NULL,
		amount TEXT NOT NULL,
		status TEXT NOT NULL CHECK(status IN ('completed', 'processing', 'failed')),
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (merchant_id) REFERENCES merchants(id)
	)`)
	if err != nil {
		t.Fatal(err)
	}

	// Seed test merchants
	testMerchants := []string{"Tokopedia", "Shopee", "Grab", "Gojek", "Traveloka"}
	merchantIDs := make(map[string]int64)
	for _, name := range testMerchants {
		res, err := db.Exec("INSERT INTO merchants(name) VALUES (?)", name)
		if err != nil {
			t.Fatal(err)
		}
		id, _ := res.LastInsertId()
		merchantIDs[name] = id
	}

	// Seed test data
	now := time.Now()
	testPayments := []struct {
		id       string
		merchant string
		amount   string
		status   string
		createdAt time.Time
	}{
		{"PAY-0001", "Tokopedia", "50000.00", "completed", now.Add(-24 * time.Hour)},
		{"PAY-0002", "Shopee", "75000.00", "processing", now.Add(-12 * time.Hour)},
		{"PAY-0003", "Grab", "25000.00", "failed", now.Add(-6 * time.Hour)},
		{"PAY-0004", "Gojek", "100000.00", "completed", now.Add(-1 * time.Hour)},
		{"PAY-0005", "Traveloka", "200000.00", "completed", now},
	}

	for _, p := range testPayments {
		_, err := db.Exec("INSERT INTO payments(id, merchant_id, amount, status, created_at) VALUES (?, ?, ?, ?, ?)",
			p.id, merchantIDs[p.merchant], p.amount, p.status, p.createdAt)
		if err != nil {
			t.Fatal(err)
		}
	}

	return db
}

func TestGetPayments_All(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPaymentRepo(db)
	payments, _, err := repo.GetPayments(nil, nil, nil, nil, 1, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payments) != 5 {
		t.Errorf("expected 5 payments, got %d", len(payments))
	}
}

func TestGetPayments_FilterByStatus(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPaymentRepo(db)

	status := "completed"
	payments, _, err := repo.GetPayments(&status, nil, nil, nil, 1, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payments) != 3 {
		t.Errorf("expected 3 completed payments, got %d", len(payments))
	}
	for _, p := range payments {
		if p.Status != "completed" {
			t.Errorf("expected status 'completed', got '%s'", p.Status)
		}
	}
}

func TestGetPayments_FilterByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPaymentRepo(db)

	id := "PAY-0002"
	payments, _, err := repo.GetPayments(nil, &id, nil, nil, 1, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payments) != 1 {
		t.Fatalf("expected 1 payment, got %d", len(payments))
	}
	if payments[0].ID != "PAY-0002" {
		t.Errorf("expected id PAY-0002, got %s", payments[0].ID)
	}
	if payments[0].MerchantName != "Shopee" {
		t.Errorf("expected merchant Shopee, got %s", payments[0].MerchantName)
	}
}

func TestGetPayments_Sort(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPaymentRepo(db)

	sort := "status"
	payments, _, err := repo.GetPayments(nil, nil, &sort, nil, 1, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payments) != 5 {
		t.Fatalf("expected 5 payments, got %d", len(payments))
	}
	// ASC by status: completed, completed, completed, failed, processing
	if payments[0].Status != "completed" {
		t.Errorf("expected first status completed, got %s", payments[0].Status)
	}
	if payments[4].Status != "processing" {
		t.Errorf("expected last status processing, got %s", payments[4].Status)
	}
}

func TestGetPayments_SortDesc(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPaymentRepo(db)

	sort := "-status"
	payments, _, err := repo.GetPayments(nil, nil, &sort, nil, 1, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payments) != 5 {
		t.Fatalf("expected 5 payments, got %d", len(payments))
	}
	// DESC by status: processing, failed, completed, completed, completed
	if payments[0].Status != "processing" {
		t.Errorf("expected first status processing, got %s", payments[0].Status)
	}
}

func TestGetPayments_InvalidSortField(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPaymentRepo(db)

	// Invalid sort field should be ignored, falls back to default order
	sort := "invalid_field"
	payments, _, err := repo.GetPayments(nil, nil, &sort, nil, 1, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payments) != 5 {
		t.Errorf("expected 5 payments, got %d", len(payments))
	}
}

func TestGetPayments_EmptyResult(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPaymentRepo(db)

	id := "PAY-9999"
	payments, _, err := repo.GetPayments(nil, &id, nil, nil, 1, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payments) != 0 {
		t.Errorf("expected 0 payments, got %d", len(payments))
	}
}

func TestParseSortParam(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"created_at", 1},
		{"-created_at", 1},
		{"amount,created_at", 2},
		{"-amount,-created_at", 2},
		{"invalid", 0},
		{"", 0},
		{"amount,invalid,created_at", 2},
	}

	for _, tt := range tests {
		result := parseSortParam(tt.input)
		if len(result) != tt.expected {
			t.Errorf("parseSortParam(%q): expected %d clauses, got %d", tt.input, tt.expected, len(result))
		}
	}
}

func TestGetPayments_Pagination(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPaymentRepo(db)
	payments, total, err := repo.GetPayments(nil, nil, nil, nil, 1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}
	if len(payments) != 2 {
		t.Errorf("expected 2 payments on page 1, got %d", len(payments))
	}

	payments2, total2, err := repo.GetPayments(nil, nil, nil, nil, 3, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total2 != 5 {
		t.Errorf("expected total 5, got %d", total2)
	}
	if len(payments2) != 1 {
		t.Errorf("expected 1 payment on page 3, got %d", len(payments2))
	}
}

func TestGetPayments_Search(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPaymentRepo(db)
	search := "Shopee"
	payments, total, err := repo.GetPayments(nil, nil, nil, &search, 1, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
	if len(payments) != 1 {
		t.Fatalf("expected 1 payment, got %d", len(payments))
	}
	if payments[0].MerchantName != "Shopee" {
		t.Errorf("expected merchant Shopee, got %s", payments[0].MerchantName)
	}
}

func TestGetPayments_SearchByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPaymentRepo(db)
	search := "PAY-0003"
	payments, total, err := repo.GetPayments(nil, nil, nil, &search, 1, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
	if len(payments) != 1 {
		t.Fatalf("expected 1 payment, got %d", len(payments))
	}
	if payments[0].ID != "PAY-0003" {
		t.Errorf("expected PAY-0003, got %s", payments[0].ID)
	}
}

func TestGetPayments_SearchNoMatch(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPaymentRepo(db)
	search := "nonexistent"
	payments, total, err := repo.GetPayments(nil, nil, nil, &search, 1, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 0 {
		t.Errorf("expected total 0, got %d", total)
	}
	if len(payments) != 0 {
		t.Errorf("expected 0 payments, got %d", len(payments))
	}
}

func TestCreatePayment(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPaymentRepo(db)
	p := &entity.Payment{
		ID:         "PAY-NEW1",
		MerchantID: 1,
		Amount:     "99999.00",
		Status:     "processing",
		CreatedAt:  time.Now(),
	}
	err := repo.CreatePayment(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	payments, total, _ := repo.GetPayments(nil, nil, nil, nil, 1, 100)
	if total != 6 {
		t.Errorf("expected total 6 after create, got %d", total)
	}
	_ = payments
}

func TestUpdatePayment(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPaymentRepo(db)
	p := &entity.Payment{
		ID:         "PAY-0001",
		MerchantID: 2,
		Amount:     "77777.00",
		Status:     "failed",
	}
	err := repo.UpdatePayment(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	id := "PAY-0001"
	payments, _, _ := repo.GetPayments(nil, &id, nil, nil, 1, 100)
	if len(payments) != 1 {
		t.Fatalf("expected 1 payment, got %d", len(payments))
	}
	if payments[0].Amount != "77777.00" {
		t.Errorf("expected amount 77777.00, got %s", payments[0].Amount)
	}
	if payments[0].Status != "failed" {
		t.Errorf("expected status failed, got %s", payments[0].Status)
	}
}

func TestUpdatePayment_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPaymentRepo(db)
	err := repo.UpdatePayment(&entity.Payment{ID: "PAY-XXXX", MerchantID: 1, Amount: "1", Status: "completed"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdatePaymentStatus(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPaymentRepo(db)
	err := repo.UpdatePaymentStatus("PAY-0002", "completed")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	id := "PAY-0002"
	payments, _, _ := repo.GetPayments(nil, &id, nil, nil, 1, 100)
	if payments[0].Status != "completed" {
		t.Errorf("expected status completed, got %s", payments[0].Status)
	}
}

func TestDeletePayment(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPaymentRepo(db)
	err := repo.DeletePayment("PAY-0001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, total, _ := repo.GetPayments(nil, nil, nil, nil, 1, 100)
	if total != 4 {
		t.Errorf("expected total 4 after delete, got %d", total)
	}
}

func TestDeletePayment_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPaymentRepo(db)
	err := repo.DeletePayment("PAY-XXXX")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
