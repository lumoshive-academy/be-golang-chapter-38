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
	dsn := "host=localhost user=postgres password=mysecretpassword dbname=testdb port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Auto-migrate the User model
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	//insert data user
	InsertDataUser(db)

	// implementation advance query
	AdvanceQuery(db)

	// implementation query non model default
	QueryNonModelDefault(db)

	// implementation update data all column
	UpdateDataUserAllColumn(db)

	// implementation update data one column
	UpdateDataUserOneColumn(db)

	// implementation update data multiple column
	UpdateDataUserMultipleColumn(db)
}

func InsertDataUser(db *gorm.DB) {
	// Menambahkan 5 data pengguna
	users := []User{
		{Name: "Budi Santoso", Email: "budi.santoso@example.com"},
		{Name: "Siti Rahmawati", Email: "siti.rahmawati@example.com"},
		{Name: "Andi Wijaya", Email: "andi.wijaya@example.com"},
		{Name: "Dewi Kartika", Email: "dewi.kartika@example.com"},
		{Name: "Joko Prasetyo", Email: "joko.prasetyo@example.com"},
	}

	// Melakukan batch insert
	db.Create(&users)
	fmt.Println("success insert data user")
}

func AdvanceQuery(db *gorm.DB) {
	// Contoh menggunakan Where, Or, Not
	var users1 []User
	db.Where("name = ?", "Budi Santoso").Or("email = ?", "budi.santoso@example.com").Find(&users1)
	fmt.Printf("Users found with Where and Or: %+v\n", users1)

	// Contoh menggunakan Not
	var users2 []User
	db.Where("name != ?", "Budi Santoso").Find(&users2)
	fmt.Printf("Users found with Not: %+v\n", users2)

	// Contoh menggunakan Select
	var users3 []User
	db.Select("name, email").Find(&users3)
	fmt.Printf("Users found with Select: %+v\n", users3)

	// Contoh menggunakan Struct
	var users4 []User
	db.Where(User{Name: "Budi Santoso", Email: "budi.santoso@example.com"}).Find(&users4)
	fmt.Printf("Users found with Struct: %+v\n", users4)

	// Contoh menggunakan Map Condition
	var users5 []User
	conditions := map[string]interface{}{
		"name":  "Budi Santoso",
		"email": "budi.santoso@example.com",
	}
	db.Where(conditions).Find(&users5)
	fmt.Printf("Users found with Map Condition: %+v\n", users5)

	// Contoh menggunakan Order
	var users6 []User
	db.Order("created_at desc").Find(&users6)
	fmt.Printf("Users found with Order: %+v\n", users6)

	// Contoh menggunakan Limit dan Offset
	var users7 []User
	db.Offset(0).Limit(5).Order("created_at desc").Find(&users7)
	fmt.Printf("Users found with Limit and Offset: %+v\n", users7)
}

func QueryNonModelDefault(db *gorm.DB) {
	type UserResponse struct {
		Name  string `gorm:"type:varchar(100);not null"`
		Email string `gorm:"type:varchar(100);uniqueIndex;not null"`
	}

	// Mengambil objek pertama
	var firstUser UserResponse
	if err := db.Model(&User{}).First(&firstUser).Error; err != nil {
		log.Fatalf("failed to find first user: %v", err)
	}
	fmt.Printf("First user: %+v\n", firstUser)

	// Mengambil objek terakhir
	var lastUser UserResponse
	if err := db.Model(&User{}).Last(&lastUser).Error; err != nil {
		log.Fatalf("failed to find last user: %v", err)
	}
	fmt.Printf("Last user: %+v\n", lastUser)

	// Mengambil objek acak
	var randomUser UserResponse
	if err := db.Model(&User{}).Take(&randomUser).Error; err != nil {
		log.Fatalf("failed to find random user: %v", err)
	}
	fmt.Printf("Random user: %+v\n", randomUser)

	// Mengambil semua object objek
	var allUser []UserResponse
	if err := db.Model(&User{}).Select("name", "email").Find(&allUser).Error; err != nil {
		log.Fatalf("failed to find random user: %v", err)
	}
	fmt.Printf("Random user: %+v\n", randomUser)
}

func UpdateDataUserAllColumn(db *gorm.DB) {
	var user User
	db.First(&user, "email = ?", "budi.santoso@example.com")

	user.Name = "Budi Santoso Baru"
	db.Save(&user)
}

func UpdateDataUserOneColumn(db *gorm.DB) {
	db.Model(&User{}).Where("email = ?", "budi.santoso@example.com").Update("name", "Budi Santoso Baru")
}

func UpdateDataUserMultipleColumn(db *gorm.DB) {
	db.Model(&User{}).Where("email = ?", "budi.santoso@example.com").Updates(User{Name: "Budi Santoso Baru", Email: "budi.baru@gmail.com"})
}
