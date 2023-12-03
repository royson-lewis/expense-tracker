package configDatabase

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	et_models "api_service/models"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := "root:rl01111998@tcp(mysql:3306)/database?parseTime=true"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	DB.AutoMigrate(
		&et_models.ETAccounts{},
		&et_models.ETTransactionCategories{},
		&et_models.ETTransactions{},
	)

	// Check if the database connection is successful
	currentDB := DB.Migrator().CurrentDatabase()
	if currentDB == "" {
		log.Fatal("Failed to get current database name. Check your connection.")
	} else {
		fmt.Printf("Connected to database: %s\n", currentDB)
	}
}
