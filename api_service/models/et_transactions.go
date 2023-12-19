package et_models

import "gorm.io/gorm"

type ETTransactions struct {
	gorm.Model
	Description               string
	Amount                    float64
	ExpenseType               string
	IsPaid                    uint8
	TransactionCategory       ETTransactionCategories `gorm:"ForeignKey:ETTransactionCategoriesID"`
	ETTransactionCategoriesID uint
	Account                   ETAccounts `gorm:"ForeignKey:ETAccountsID"`
	ETAccountsID              uint
}

type ETTransactionCategories struct {
	gorm.Model
	Name         string `gorm:"unique"`
	Type         string
	Transactions []ETTransactions
}
