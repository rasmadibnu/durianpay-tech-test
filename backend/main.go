package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/durianpay/fullstack-boilerplate/internal/api"
	"github.com/durianpay/fullstack-boilerplate/internal/config"
	ah "github.com/durianpay/fullstack-boilerplate/internal/module/auth/handler"
	ar "github.com/durianpay/fullstack-boilerplate/internal/module/auth/repository"
	au "github.com/durianpay/fullstack-boilerplate/internal/module/auth/usecase"
	mh "github.com/durianpay/fullstack-boilerplate/internal/module/merchant/handler"
	mr "github.com/durianpay/fullstack-boilerplate/internal/module/merchant/repository"
	mu "github.com/durianpay/fullstack-boilerplate/internal/module/merchant/usecase"
	srv "github.com/durianpay/fullstack-boilerplate/internal/service/http"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	_ = godotenv.Load()

	db, err := sql.Open("sqlite3", "dashboard.db?_foreign_keys=1")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := initDB(db); err != nil {
		log.Fatal(err)
	}

	JwtExpiredDuration, err := time.ParseDuration(config.JwtExpired)
	if err != nil {
		panic(err)
	}

	userRepo := ar.NewUserRepo(db)
	merchantRepo := mr.NewMerchantRepo(db)

	authUC := au.NewAuthUsecase(userRepo, config.JwtSecret, JwtExpiredDuration)
	merchantUC := mu.NewMerchantUsecase(merchantRepo)

	authH := ah.NewAuthHandler(authUC)
	merchantH := mh.NewMerchantHandler(merchantUC)

	apiHandler := &api.APIHandler{
		Auth:     authH,
		Merchant: merchantH,
	}

	server := srv.NewServer(apiHandler, config.OpenapiYamlLocation)

	addr := config.HttpAddress
	log.Printf("starting server on %s", addr)
	server.Start(addr)
}

func initDB(db *sql.DB) error {
	// create tables if not exists
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
		  id INTEGER PRIMARY KEY AUTOINCREMENT,
		  email TEXT NOT NULL UNIQUE,
		  password_hash TEXT NOT NULL,
		  role TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS merchants (
		  id INTEGER PRIMARY KEY AUTOINCREMENT,
		  name TEXT NOT NULL UNIQUE,
		  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS payments (
		  id TEXT PRIMARY KEY,
		  merchant_id INTEGER NOT NULL,
		  amount TEXT NOT NULL,
		  status TEXT NOT NULL CHECK(status IN ('completed', 'processing', 'failed')),
		  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		  FOREIGN KEY (merchant_id) REFERENCES merchants(id)
		);`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return err
		}
	}
	// seed admin user if not exists
	var userCnt int
	row := db.QueryRow("SELECT COUNT(1) FROM users")
	if err := row.Scan(&userCnt); err != nil {
		return err
	}
	if userCnt == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		if _, err := db.Exec("INSERT INTO users(email, password_hash, role) VALUES (?, ?, ?)", "cs@test.com", string(hash), "cs"); err != nil {
			return err
		}
		if _, err := db.Exec("INSERT INTO users(email, password_hash, role) VALUES (?, ?, ?)", "operation@test.com", string(hash), "operation"); err != nil {
			return err
		}
	}

	// Seed payments
	var payCnt int
	if err := db.QueryRow("SELECT COUNT(1) FROM payments").Scan(&payCnt); err != nil {
		return err
	}

	if payCnt == 0 {
		if err := seedMerchantsAndPayments(db); err != nil {
			return err
		}
	}

	const dbLifetime = time.Minute * 5
	db.SetConnMaxLifetime(dbLifetime)
	return nil
}

func seedMerchantsAndPayments(db *sql.DB) error {
	merchantNames := []string{
		"Tokopedia", "Shopee", "Bukalapak", "Lazada", "Blibli",
		"Grab", "Gojek", "Traveloka", "Tiket.com", "JD.ID",
		"Zalora", "Sociolla", "Bhinneka", "MatahariMall", "Orami",
	}

	merchantIDs := make([]int64, len(merchantNames))
	for i, name := range merchantNames {
		res, err := db.Exec("INSERT INTO merchants(name) VALUES (?)", name)
		if err != nil {
			return err
		}
		merchantIDs[i], err = res.LastInsertId()
		if err != nil {
			return err
		}
	}

	statuses := []string{"completed", "processing", "failed"}
	rng := rand.New(rand.NewSource(42))
	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 1; i <= 50; i++ {
		id := fmt.Sprintf("PAY-%04d", i)
		merchantID := merchantIDs[rng.Intn(len(merchantIDs))]
		status := statuses[rng.Intn(len(statuses))]
		amount := fmt.Sprintf("%.2f", 10000+rng.Float64()*990000)
		createdAt := baseTime.Add(time.Duration(rng.Intn(365*24)) * time.Hour)

		_, err := db.Exec(
			"INSERT INTO payments(id, merchant_id, amount, status, created_at) VALUES (?, ?, ?, ?, ?)",
			id, merchantID, amount, status, createdAt,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
