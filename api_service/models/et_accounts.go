package et_models

import "gorm.io/gorm"

type ETAccounts struct {
	gorm.Model
	Name          string  `json:"name" binding:"required" gorm:"unique"`
	ActualBalance float64 `json:"actual_balance" binding:"min=0"`
	Description   string  `json:"description"`
	Transactions  []ETTransactions
}
