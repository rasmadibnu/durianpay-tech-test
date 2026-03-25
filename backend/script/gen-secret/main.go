package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var KeyLength = 32

func main() {
	db, err := sql.Open("sqlite3", "dashboard.db?_foreign_keys=1")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := initDB(db); err != nil {
		log.Fatal(err)
	}

	key := make([]byte, KeyLength)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}
	secret := base64.StdEncoding.EncodeToString(key)

	envFile := ".env"
	content, _ := os.ReadFile(envFile)
	lines := bytes.Split(content, []byte("\n"))

	var found bool
	for i, line := range lines {
		if bytes.HasPrefix(line, []byte("JWT_SECRET=")) {
			lines[i] = []byte("JWT_SECRET=" + secret)
			found = true
			break
		}
	}
	if !found {
		lines = append(lines, []byte("JWT_SECRET="+secret))
	}

	f, err := os.Create(envFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, line := range lines {
		if len(strings.TrimSpace(string(line))) > 0 {
			fmt.Fprintln(w, string(line))
		}
	}
	w.Flush()

	fmt.Println("✅ New JWT_SECRET saved to .env")
}

func initDB(db *sql.DB) error {
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

	// Seed users
	var userCnt int
	if err := db.QueryRow("SELECT COUNT(1) FROM users").Scan(&userCnt); err != nil {
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
	if err := seedMerchant(db); err != nil {
		return err
	}

	const dbLifetime = time.Minute * 5
	db.SetConnMaxLifetime(dbLifetime)
	return nil
}

func seedMerchant(db *sql.DB) error {
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
	return nil
}
