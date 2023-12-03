package et_models

import "gorm.io/gorm"

type ETTransactions struct {
    gorm.Model
    Description string
    Amount int
    ExpenseType string
    IsPaid uint8
	CategoryId uint
}

type ETTransactionCategories struct {
    gorm.Model
    Name string
    Type string
    ETTransactions []ETTransactions
}