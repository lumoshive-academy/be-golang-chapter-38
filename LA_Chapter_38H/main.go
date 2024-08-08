package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// User model with auto increment ID
type User struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"type:varchar(100);not null"`
	Email     string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password  string `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

func main() {
	// Database connection string
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Jakarta"

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Migrate the schema
	db.AutoMigrate(&User{})

	// Create users
	users := []User{
		{Name: "Budi Santoso", Email: "budi.santoso@example.com", Password: "password"},
		{Name: "Siti Nurhaliza", Email: "siti.nurhaliza@example.com", Password: "password"},
		{Name: "Agus Salim", Email: "agus.salim@example.com", Password: "password"},
	}

	// Insert users into the database
	for _, user := range users {
		db.Create(&user)
		fmt.Printf("User created with ID: %d, Name: %s, Email: %s\n", user.ID, user.Name, user.Email)
	}
}

func UpSertDataUser(db *gorm.DB) {
	// Data user baru atau update data user yang sudah ada
	user1 := User{
		ID:    1, // Jika user dengan ID ini sudah ada, maka data akan diperbarui
		Email: "budi.santoso@example.com",
	}

	// Save akan melakukan upsert
	db.Save(&user1.Email)

	log.Printf("User saved with ID: %d", user1.ID)

	// Data user baru atau update data user yang sudah ada
	user2 := User{
		Name:     "Siti Nurhaliza",
		Email:    "siti.nurhaliza@example.com",
		Password: "password",
	}

	// OnConflict akan melakukan upsert berdasarkan email
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "password"}),
	}).Create(&user2)

	log.Printf("User upserted with Email: %s", user2.Email)
}
