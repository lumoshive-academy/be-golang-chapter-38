package main

import (
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
	ID         uint    `gorm:"primaryKey;autoIncrement"`
	Name       string  `gorm:"type:varchar(100);not null"`
	Email      string  `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password   string  `gorm:"type:varchar(255);not null"`
	Address    Address `gorm:"embedded;embeddedPrefix:addr_"` // prefix field name, example user_street
	Timestamps         // embed struct timestamps
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
		Logger: newLogger, // config logger
	})
	if err != nil {
		panic("failed to connect database")
	}

	// create a new user
	createDataUser(db)

	// create a slice of user
	createWithSliceDataUser(db)
}

func createDataUser(db *gorm.DB) {
	// Create a new user
	user := User{
		Name:     "Budi Santoso",
		Email:    "budi.santoso@example.com",
		Password: "rahasia123",
		Address: Address{
			Street:  "Jl. Merdeka No. 123",
			City:    "Jakarta",
			State:   "DKI Jakarta",
			ZipCode: "10110",
		},
	}

	// Save the user to the database
	db.Create(&user)
}

func createWithSliceDataUser(db *gorm.DB) {
	// Create a slice of users
	users := []User{
		{
			Name:     "Rudi Budiman",
			Email:    "rudi.budiman@example.com",
			Password: "rahasia343",
			Address: Address{
				Street:  "Jl. kedoya raya",
				City:    "Jakarta",
				State:   "DKI Jakarta",
				ZipCode: "10110",
			},
		}, {
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
}
