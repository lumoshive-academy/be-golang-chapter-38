package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"type:varchar(100);not null"`
	Email     string `gorm:"type:varchar(100);uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Auto-migrate the User model
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// Add users within a transaction
	err = createUsersTransaction(db)
	if err != nil {
		log.Fatalf("failed to create users: %v", err)
	}

	// implementation query GORM
	err = QueryImplem(db)
	if err != nil {
		log.Fatalf("failed query: %v", err)
	}
}

func createUsersTransaction(db *gorm.DB) error {
	users := []User{
		{
			Name:  "Budi Santoso",
			Email: "budi.santoso@example.com",
		},
		{
			Name:  "Siti Aminah",
			Email: "siti.aminah@example.com",
		},
		{
			Name:  "Andi Wijaya",
			Email: "andi.wijaya@example.com",
		},
	}

	// Using db.Transaction to ensure atomicity
	return db.Transaction(func(tx *gorm.DB) error {
		for _, user := range users {

			if err := tx.Create(&user).Error; err != nil {
				return err // Rollback the transaction in case of an error
			}
		}
		return nil // Commit the transaction
	})
}

func QueryImplem(db *gorm.DB) error {
	var user User
	// Querying a single object using First
	if err := db.First(&user, 1).Error; err != nil {
		log.Fatalf("failed to find user: %v", err)
	}
	fmt.Printf("User found using First: %+v\n", user)

	// Querying a single object using Take
	if err := db.Take(&user).Error; err != nil {
		log.Fatalf("failed to take user: %v", err)
	}
	fmt.Printf("User found using Take: %+v\n", user)

	// Querying a single object using Last
	if err := db.Last(&user).Error; err != nil {
		log.Fatalf("failed to find last user: %v", err)
	}
	fmt.Printf("User found using Last: %+v\n", user)

	// Querying with condition
	if err := db.Where("email = ?", "budi.santoso@example.com").First(&user).Error; err != nil {
		log.Fatalf("failed to find user: %v", err)
	}
	fmt.Printf("User found: %+v\n", user)

	// Querying all objects
	var users []User
	if err := db.Find(&users).Error; err != nil {
		log.Fatalf("failed to find users: %v", err)
	}
	fmt.Printf("Users found: %+v\n", users)

	// Querying with IN condition using Find
	if err := db.Find(&users, "email IN ?", []string{"budi.santoso@example.com", "siti.aminah@example.com"}).Error; err != nil {
		log.Fatalf("failed to find users: %v", err)
	}
	fmt.Printf("Users found: %+v\n", users)

	return nil
}
