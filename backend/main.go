package main

import (
	"database/sql"
	"log"
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

	var merchantCnt int
	row = db.QueryRow("SELECT COUNT(1) FROM merchants")
	if err := row.Scan(&merchantCnt); err != nil {
		return err
	}
	if merchantCnt == 0 {
		if err := seedMerchants(db); err != nil {
			return err
		}
	}

	const dbLifetime = time.Minute * 5
	db.SetConnMaxLifetime(dbLifetime)
	return nil
}

func seedMerchants(db *sql.DB) error {
	merchantNames := []string{
		"Tokopedia", "Shopee", "Bukalapak", "Lazada", "Blibli",
		"Grab", "Gojek", "Traveloka", "Tiket.com", "JD.ID",
		"Zalora", "Sociolla", "Bhinneka", "MatahariMall", "Orami",
	}

	for _, name := range merchantNames {
		_, err := db.Exec("INSERT INTO merchants(name) VALUES (?)", name)
		if err != nil {
			return err
		}
	}

	return nil
}
