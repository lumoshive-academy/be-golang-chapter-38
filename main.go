package main

import (
	"fmt"
	"log"
	"os"
	"time"
	
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	ID       uint    `gorm:"primaryKey;autoIncrement"`
	Name     string  `gorm:"type:varchar(100);not null"`
	Email    string  `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password string  `gorm:"type:varchar(255);not null"`
	Address  Address `gorm:"embedded;embeddedPrefix:addr_"`
	Timestamps
}

func main() {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)

	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	// Initialize the database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}

	// Perform auto-migration
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("failed to migrate database")
	}

	fmt.Printf("success migration database")

	// Manual migration: add new column 'PhoneNumber' to 'users' table
	err = db.Exec("ALTER TABLE users ADD COLUMN phone_number VARCHAR(15)").Error
	if err != nil {
		panic("failed to add phone_number column to users table")
	}

	// Create a slice of users with Indonesian data
	users := []User{
		{
			Name:     "Budi Santoso",
			Email:    "budi.santoso@example.com",
			Password: "rahasia123",
			Address: Address{
				Street:  "Jl. Merdeka No. 123",
				City:    "Jakarta",
				State:   "DKI Jakarta",
				ZipCode: "10110",
			},
		},
		{
			Name:     "Siti Aminah",
			Email:    "siti.aminah@example.com",
			Password: "rahasia456",
			Address: Address{
				Street:  "Jl. Sudirman No. 45",
				City:    "Bandung",
				State:   "Jawa Barat",
				ZipCode: "40235",
			},
		},
		{
			Name:     "Agus Setiawan",
			Email:    "agus.setiawan@example.com",
			Password: "rahasia789",
			Address: Address{
				Street:  "Jl. Diponegoro No. 78",
				City:    "Surabaya",
				State:   "Jawa Timur",
				ZipCode: "60265",
			},
		},
	}

	// Perform batch insert
	result := db.Create(&users)

	// Check for errors
	if result.Error != nil {
		panic("failed to insert batch records")
	}

	// 	// Create a new user with Indonesian data
	// 	user := User{
	// 		Name:     "Budi Santoso",
	// 		Email:    "budi.santoso@example.com",
	// 		Password: "rahasia123",
	// 		Address: Address{
	// 			Street:  "Jl. Merdeka No. 123",
	// 			City:    "Jakarta",
	// 			State:   "DKI Jakarta",
	// 			ZipCode: "10110",
	// 		},
	// 	}

	// 	// Save the user to the database
	// 	db.Create(&user)
}
