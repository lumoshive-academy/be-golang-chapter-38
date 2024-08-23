package main

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Address model to be embedded
type Address struct {
	Street  string
	City    string
	State   string
	ZipCode string
}

// Timestamps model to be embedded
type Timestamps struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// User model with embedded Address and Timestamps
type User struct {
	ID         uint    `gorm:"primaryKey;autoIncrement"`
	Name       string  `gorm:"type:varchar(100);not null"`
	Email      string  `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password   string  `gorm:"type:varchar(255);not null"`
	Address    Address `gorm:"embedded;embeddedPrefix:addr_"` // prefix field name, example user_street
	Timestamps         // embed struct timestamps
}

func main() {

	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	// Initialize the database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Perform auto-migration
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("failed to migrate database")
	}

	// Manual migration: add new column 'PhoneNumber' to 'users' table
	err = db.Exec("ALTER TABLE users ADD COLUMN phone_number VARCHAR(15)").Error
	if err != nil {
		panic("failed to add phone_number column to users table")
	}

	fmt.Printf("success migration database")

}
