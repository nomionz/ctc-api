package main

import "github.com/shopspring/decimal"

type Product struct {
	ID     int             `json:"id" gorm:"primary_key"`
	Name   string          `json:"name" gorm:"type:varchar(100);not null"`
	Price  decimal.Decimal `json:"price" gorm:"type:decimal(10,2);not null"`
	Amount int             `json:"amount" gorm:"type:int;not null"`
}

func (p *Product) TableName() string {
	return "products"
}
