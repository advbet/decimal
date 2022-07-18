package decimal

import (
	"math"
	"math/big"

	newDecimal "github.com/shopspring/decimal"
)

type Number = newDecimal.Decimal

// RoundRule is enum type for specifying rounding algorithm when decimal number
// is scaled with loss of precision. List of supported rounding rules are listed
// in `Round*` consts.
type RoundRule int

// List of supported rounding rules
const (
	RoundTruncate RoundRule = iota // Directed rounding towards zero
	RoundFloor                     // Directed rounding towards positive infinity
	RoundCeil                      // Directed rounding towards negative infinity
	RoundMath                      // Round to nearest, on tie round away from zero
	RoundBankers                   // Round to nearest, on tie round to even number
)

func init() {
	newDecimal.MarshalJSONWithoutQuotes = true
}

// NewFromRat this only works for small Rats, big.Rats won't necessarily fit into a Float64,
// so you'll be losing information here.
// https://github.com/shopspring/decimal/pull/4/files
func NewFromRat(value *big.Rat, exp int) newDecimal.Decimal {
	if value.IsInt() {
		f, _ := value.Float64()
		return newDecimal.New(int64(math.Floor(f)), 0)
	}

	floatValue, _ := value.Float64()
	return newDecimal.NewFromFloatWithExponent(floatValue, int32(exp))
}

// Round scales decimal value to an integer value with given exponent. On
// exponent scale-down decimal value precision is preserved, on exponent
// scale-up rounding with the given rounding rule is performed.
func Round(value newDecimal.Decimal, exp int, rule RoundRule) newDecimal.Decimal {
	// scale-down case
	if exp <= int(value.Exponent()) {
		return Rescale(value, int32(exp))
	}

	switch rule {
	case RoundBankers:
		return Rescale(value.RoundBank(-1*int32(exp)), int32(exp))
	case RoundMath:
		return Rescale(value.Round(-1*int32(exp)), int32(exp))
	case RoundFloor:
		return Rescale(value.RoundFloor(-1*int32(exp)), int32(exp))
	case RoundCeil:
		return Rescale(value.RoundCeil(-1*int32(exp)), int32(exp))
	default: // truncate the remainder
		return Rescale(value, int32(exp))
	}
}

// MulInt calculates d * n value.
func MulInt(value newDecimal.Decimal, n int) newDecimal.Decimal {
	d := newDecimal.NewFromInt(int64(n))
	return value.Mul(d)
}

// ScaledVal scales decimal number to a given exponent and returns
// internal number integer value. If given exponent is higher than internal
// number exponent this function will lose truncated digits.
//
// Example: number "12.99" with call ScaledVal(-4) would return 129900, with
// call ScaledVal(0) would return 12.
func ScaledVal(d newDecimal.Decimal, exp int) int64 {
	return Rescale(d, int32(exp)).CoefficientInt64()
}

// Rescale copied from `shopspring/decimal`
func Rescale(d newDecimal.Decimal, exp int32) newDecimal.Decimal {
	if d.Exponent() == exp {
		return d
	}

	// NOTE(vadim): must convert exps to float64 before - to prevent overflow
	diff := math.Abs(float64(exp) - float64(d.Exponent()))
	value := new(big.Int).Set(d.Coefficient())

	expScale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(diff)), nil)
	if exp > d.Exponent() {
		value = value.Quo(value, expScale)
	} else if exp < d.Exponent() {
		value = value.Mul(value, expScale)
	}

	return newDecimal.NewFromBigInt(value, exp)
}
