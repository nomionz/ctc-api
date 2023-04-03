package models

import "github.com/shopspring/decimal"

// form is used to bind the request body to the Product struct
type Product struct {
	ID     int             `json:"id,omitempty" gorm:"primary_key;auto_increment;not_null"`
	Name   string          `json:"name" gorm:"type:varchar(100);not null" form:"name" binding:"required"`
	Price  decimal.Decimal `json:"price" gorm:"type:decimal(10,2);not null" form:"price" binding:"required"`
	Amount int             `json:"amount" gorm:"type:int;not null" form:"amount" binding:"required"`
}

func (p *Product) TableName() string {
	return "products"
}
