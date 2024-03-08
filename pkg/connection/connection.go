package connection

import (
	"fmt"
	"time"

	"github.com/Bappy60/BookStore_in_Go/pkg/config"
	"github.com/Bappy60/BookStore_in_Go/pkg/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func Connect() {
	time.Sleep(1000 * time.Millisecond)
	config := config.GConfig
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUser, config.DBPass, config.DBHost, config.DBPort, config.DbName)

	d, err := gorm.Open("mysql", connectionString)
	if err != nil {
		panic(err.Error())
	}
	DB = d
}
func ConnectToDBWithRetry() {
	// Set the maximum number of retries and retry interval
	maxRetries := 5
	retryInterval := 3 * time.Second

	for i := 1; i <= maxRetries; i++ {
		// Add a delay before each retry
		time.Sleep(retryInterval)

		config := config.GConfig
		connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.DBUser, config.DBPass, config.DBHost, config.DBPort, config.DbName)

		// Attempt to connect to the database
		d, err := gorm.Open("mysql", connectionString)
		if err == nil {
			// Connection successful
			fmt.Println("Database Connected")
			DB = d
			return
		}
		// Log the error, and retry
		fmt.Printf("Error connecting to DB: %v. Retrying...\n", err)

	}
	fmt.Printf("Failed to connect DB")

}

func GetDB() *gorm.DB {
	return DB
}

func Initialize() *gorm.DB {
	ConnectToDBWithRetry()
	db := GetDB()
	db.AutoMigrate(&models.Book{}, &models.Author{})
	return db
}
