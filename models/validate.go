package models

import (
	"fmt"

	"github.com/shopspring/decimal"
)

func Validate(prod *Product) error {
	if prod.Name == "" {
		return fmt.Errorf("name is empty")
	}
	if prod.Price.LessThanOrEqual(decimal.NewFromFloat(0.0)) {
		return fmt.Errorf("price is less than or equal to 0")
	}
	if prod.Amount <= 0 {
		return fmt.Errorf("amount is less than or equal to 0")
	}

	return nil
}
