package main

import (
	"fmt"
	"log"
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
	ID       uint    `gorm:"primaryKey;autoIncrement"`
	Name     string  `gorm:"type:varchar(100);not null"`
	Email    string  `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password string  `gorm:"type:varchar(255);not null"`
	Address  Address `gorm:"embedded;embeddedPrefix:addr_"`
	Timestamps
}

// Product model
type Product struct {
	ID          uint    `gorm:"primaryKey;autoIncrement"`
	Name        string  `gorm:"type:varchar(100);not null"`
	Description string  `gorm:"type:text"`
	Price       float64 `gorm:"not null"`
	Stock       int     `gorm:"not null"`
	Timestamps
}

// Order model
type Order struct {
	ID         uint    `gorm:"primaryKey;autoIncrement"`
	UserID     uint    `gorm:"not null"`
	ProductID  uint    `gorm:"not null"`
	Quantity   int     `gorm:"not null"`
	TotalPrice float64 `gorm:"not null"`
	Timestamps
}

func main() {
	// Database connection details
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Perform auto-migration
	err = db.AutoMigrate(&User{}, &Product{}, &Order{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// Create sample data
	createSampleData(db)

	// Create an order with transaction
	err = createOrderTransaction(db, 1, 1, 2)
	if err != nil {
		log.Fatalf("failed to create order: %v", err)
	}
}

func createSampleData(db *gorm.DB) {
	// Create sample users
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
	}

	// Create sample products
	products := []Product{
		{
			Name:        "Laptop",
			Description: "Laptop high-end dengan spesifikasi tinggi",
			Price:       15000000,
			Stock:       10,
		},
		{
			Name:        "Smartphone",
			Description: "Smartphone dengan kamera berkualitas tinggi",
			Price:       5000000,
			Stock:       20,
		},
	}

	// Insert sample data
	db.Create(&users)
	db.Create(&products)
}

func createOrderTransaction(db *gorm.DB, userID, productID, quantity int) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var product Product
		if err := tx.First(&product, productID).Error; err != nil {
			return err
		}

		if product.Stock < quantity {
			return fmt.Errorf("not enough stock for product %s", product.Name)
		}

		totalPrice := float64(quantity) * product.Price
		order := Order{
			UserID:     uint(userID),
			ProductID:  uint(productID),
			Quantity:   quantity,
			TotalPrice: totalPrice,
		}

		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		product.Stock -= quantity
		if err := tx.Save(&product).Error; err != nil {
			return err
		}

		return nil
	})
}
