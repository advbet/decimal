package decimal

import (
	newDecimal "github.com/shopspring/decimal"
)

// FromFloat64 creates new decimal value from float64.
func FromFloat64(f float64) (Number, error) {
	return newDecimal.NewFromFloat(f), nil
}
