package et_models

import "gorm.io/gorm"

type ETAccounts struct {
    gorm.Model
    Name string
    ActualBalance int
    Description string
}