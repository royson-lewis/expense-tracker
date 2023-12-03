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
	dsn := "admin:Paniyoor0111@tcp(rl-portfolio-prod.c2zyjnkyryxg.ap-south-1.rds.amazonaws.com:3306)/et"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
        log.Fatal(err)
    }

	DB.AutoMigrate(
		&et_models.ETTransactions{},
		&et_models.ETTransactionCategories{},
		&et_models.ETAccounts{},
	)

	// Check if the database connection is successful
	currentDB := DB.Migrator().CurrentDatabase()
	if currentDB == "" {
		log.Fatal("Failed to get current database name. Check your connection.")
	} else {
		fmt.Printf("Connected to database: %s\n", currentDB)
	}
}