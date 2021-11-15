package decimal

import (
	"fmt"
	"strconv"
	"strings"
)

// Float64 returns decimal value approximation as float64.
func (d Number) Float64() float64 {
	f, _ := strconv.ParseFloat(fmt.Sprintf("%de%d", d.val, d.exp), 64)
	// Errors are ignored, for +Inf, -Inf cases, f will be set correctly
	// even if error is returned.
	return f
}

// FromFloat64 creates new decimal value from float64.
func FromFloat64(f float64) (Number, error) {
	parts := strings.Split(strconv.FormatFloat(f, 'e', -1, 64), "e")
	d, err := FromString(parts[0])
	if err != nil {
		return Number{}, err
	}
	exp, err := strconv.Atoi(parts[1])
	if err != nil {
		return Number{}, err
	}
	d.exp += exp
	return d, nil
}
