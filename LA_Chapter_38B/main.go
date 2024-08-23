package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID       uint
	Name     string
	Email    string
	Password string
}

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Contoh penggunaan Exec
	execExamples(db)

	// raw select by id
	rawSQLSelectByID(db, 1)

	// select all data
	rawSQLSelectAll(db)

}

func execExamples(db *gorm.DB) {
	result := db.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", "Lumoshive", "lumoshive@example.com", "12345")
	if result.Error != nil {
		log.Fatalf("Failed to insert user: %v", result.Error)
	}
	log.Printf("Rows affected: %d", result.RowsAffected)
}

func rawSQLSelectByID(db *gorm.DB, id int) {
	var user User
	// Menjalankan query SELECT dengan parameter
	result := db.Raw("SELECT id, name, email, password FROM users WHERE id = ?", id).Scan(&user)
	if result.Error != nil {
		log.Fatalf("Failed to query user: %v", result.Error)
	}
	log.Printf("User: %+v", user)
}

func rawSQLSelectAll(db *gorm.DB) {
	rows, err := db.Raw("SELECT id, name, email, password FROM users").Rows()
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := db.ScanRows(rows, &user); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		log.Printf("User: %+v", user)
	}
}
