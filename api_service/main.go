package main

import (
	"errors"
	"net/http"
	"strconv"

	database "api_service/config/database"
	databaseHelpers "api_service/helpers/database"
	et_models "api_service/models"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccountWithCalcBal struct {
	Account       et_models.ETAccounts
	CalculatedBal float64
}

func setupRouter() *gin.Engine {
	// Initialize the DB Connection
	database.InitDB()

	r := gin.Default()

	var DB = database.DB

	r.GET("/transactions", func(c *gin.Context) {
		var transactions []et_models.ETTransactions
		if err := DB.Preload(clause.Associations).Find(&transactions).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, transactions)
	})

	r.POST("/transactions", func(c *gin.Context) {
		var jsonInput struct {
			Description               string  `json:"description" binding:"required"`
			Amount                    float64 `json:"amount" binding:"required"`
			ExpenseType               string  `json:"expense_type"`
			IsPaid                    uint8   `json:"is_paid" binding:"required"`
			ETTransactionCategoriesID uint    `json:"category_id"`
			ETAccountsID              uint    `json:"account_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&jsonInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newTransaction := et_models.ETTransactions{
			Description:               jsonInput.Description,
			Amount:                    jsonInput.Amount,
			ExpenseType:               jsonInput.ExpenseType,
			IsPaid:                    jsonInput.IsPaid,
			ETTransactionCategoriesID: jsonInput.ETTransactionCategoriesID,
			ETAccountsID:              jsonInput.ETAccountsID,
		}

		if err := DB.Create(&newTransaction).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Transaction added successfully", "transaction": newTransaction})
	})

	r.GET("/transactions/categories", func(c *gin.Context) {
		var categories []et_models.ETTransactionCategories
		if err := DB.Find(&categories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, categories)
	})

	r.POST("/transactions/categories", func(c *gin.Context) {
		var jsonInput struct {
			Name string `json:"name" binding:"required"`
			Type string `json:"type" binding:"required"`
		}

		if err := c.ShouldBindJSON(&jsonInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newCategory := et_models.ETTransactionCategories{
			Name: jsonInput.Name,
			Type: jsonInput.Type,
		}

		if err := DB.Create(&newCategory).Error; err != nil {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok {
				switch mysqlErr.Number {
				case 1062: // MySQL error number for duplicate entry
					if databaseHelpers.IsDuplicateKeyForField(mysqlErr, "name") {
						c.JSON(http.StatusConflict, gin.H{"error": "Name must be unique"})
						return
					}
				}
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Category added successfully", "category": newCategory})
	})

	r.GET("/accounts", func(c *gin.Context) {
		var accounts []et_models.ETAccounts
		if err := DB.Preload("Transactions.TransactionCategory").Preload("Transactions.Account").Find(&accounts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		type AccountWithCalculatedBal struct {
			Account       et_models.ETAccounts
			CalculatedBal float64
		}

		accountsWithCalcBal := make([]AccountWithCalculatedBal, len(accounts))
		for i, account := range accounts {
			var totalBal float64 = 0
			for _, transaction := range account.Transactions {
				totalBal = totalBal + transaction.Amount
			}
			accountsWithCalcBal[i].CalculatedBal = totalBal
			accountsWithCalcBal[i].Account = account
		}

		c.JSON(http.StatusOK, accountsWithCalcBal)
	})

	r.GET("/accounts/:id", func(c *gin.Context) {
		accountID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
			return
		}

		type AccountWithCalculatedBal struct {
			Account       et_models.ETAccounts
			CalculatedBal float64
		}

		var data AccountWithCalculatedBal
		if err := DB.Preload("Transactions").First(&data.Account, accountID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var totalBal float64
		for _, item := range data.Account.Transactions {
			totalBal = totalBal + item.Amount
		}

		c.JSON(http.StatusOK, AccountWithCalculatedBal{
			Account:       data.Account,
			CalculatedBal: totalBal,
		})
		return
	})

	r.POST("/accounts", func(c *gin.Context) {
		var jsonInput struct {
			Name          string  `json:"name" binding:"required"`
			Description   string  `json:"description"`
			ActualBalance float64 `json:"actual_balance" binding:"min=0"`
		}

		if err := c.ShouldBindJSON(&jsonInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newAccount := et_models.ETAccounts{
			Name:          jsonInput.Name,
			Description:   jsonInput.Description,
			ActualBalance: jsonInput.ActualBalance,
		}

		if err := DB.Create(&newAccount).Error; err != nil {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok {
				switch mysqlErr.Number {
				case 1062: // MySQL error number for duplicate entry
					if databaseHelpers.IsDuplicateKeyForField(mysqlErr, "name") {
						c.JSON(http.StatusConflict, gin.H{"error": "Name must be unique"})
						return
					}
				}
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Account added successfully", "account": newAccount})
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	// authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
	// 	"foo":  "bar", // user:foo password:bar
	// 	"manu": "123", // user:manu password:123
	// }))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	// authorized.POST("admin", func(c *gin.Context) {
	// 	user := c.MustGet(gin.AuthUserKey).(string)

	// 	// Parse JSON
	// 	var json struct {
	// 		Value string `json:"value" binding:"required"`
	// 	}

	// 	if c.Bind(&json) == nil {
	// 		db[user] = json.Value
	// 		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	// 	}
	// })

	return r
}

func main() {
	r := setupRouter()
	r.Run(":9000")
}
