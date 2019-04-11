package decimal

import (
	"database/sql/driver"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

// Number is a type used to store precise decimal values. It is used to
// store data from DECIMAL / NUMERC fields from a database without loosing
// precision. This type also have some helper methods for basic arithmetical
// operation. All except Div operation is executed without a loss of precision.
type Number struct {
	val int64
	exp int
}

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

var logTable = []int64{
	1,
	10,
	100,
	1000,
	10000,
	100000,
	1000000,
	10000000,
	100000000,
	1000000000,
	10000000000,
	100000000000,
	1000000000000,
	10000000000000,
	100000000000000,
	1000000000000000,
	10000000000000000,
	100000000000000000,
	1000000000000000000,
}

// Zero create a new decimal number that is equal to zero.
func Zero() Number {
	return Number{}
}

// New creates a new decimal number having value of val*10^exp.
func New(val int64, exp int) Number {
	return Number{val, exp}
}

// FromInt creates a new instance of decimal number with an integer value and
// zero exponent.
func FromInt(val int) Number {
	return Number{int64(val), 0}
}

// FromString creates a new instance of decimal number by parsing given string.
//
// Accepted format examples: "1234", "123.456". "123.000".
//
// Scientific notation is not supported.
func FromString(str string) (Number, error) {
	d := Number{}
	err := d.UnmarshalText([]byte(str))
	return d, err
}

// Val extracts private field holding current integer value. Real decimal number
// value is val*10^exp.
func (d *Number) Val() int64 {
	return d.val
}

// Exp extracts private field holding current exponent value. Real decimal
// number value is val*10^exp.
func (d *Number) Exp() int {
	return d.exp
}

func (d *Number) denormalize(exp int) {
	if exp >= d.exp {
		return
	}
	log := d.exp - exp
	if log >= len(logTable) {
		panic("decimal.Number: logTable lookup failed, int64 overflow")
	}
	scale := logTable[log]
	d.exp -= log
	d.val *= scale
}

// Round scales decimal value to an integer value with given exponent. On
// exponent scale-down decimal value precision is preserved, on exponent
// scale-up rounding with the given rounding rule is performed.
func (d Number) Round(exp int, rule RoundRule) Number {
	// scale-down case
	if exp <= d.exp {
		d.denormalize(exp)
		return d
	}

	// round value exp > d.exp
	log := exp - d.exp
	if log >= len(logTable) {
		panic("decimal.Number: logTable lookup failed, int64 overflow")
	}
	scale := logTable[log]
	sign := int64(1)
	if d.val < 0 {
		sign = -1
	}
	remainder := sign * (d.val % scale) // remainder is always positive
	d.exp += log
	d.val /= scale

	switch rule {
	case RoundBankers:
		tieVal := 5 * logTable[log-1]
		if remainder > tieVal {
			d.val += sign
			return d
		}
		if remainder == tieVal && d.val%2 != 0 {
			d.val += sign
			return d
		}
		return d
	case RoundMath:
		remainder /= logTable[log-1]
		if remainder >= 5 {
			d.val += sign
		}
		return d
	case RoundFloor:
		if sign < 0 && remainder != 0 {
			d.val += sign
		}
		return d
	case RoundCeil:
		if sign > 0 && remainder != 0 {
			d.val += sign
		}
		return d
	default: // trucate the remainder
		return d
	}
}

// IsZero returns true if value of the decimal is equal to 0. It is a shortcut
// for executing `d.Val() == 0`.
func (d Number) IsZero() bool {
	return d.val == 0
}

// Add implements + operation for decimal numbers
func (d Number) Add(that Number) Number {
	if that.exp < d.exp {
		d.denormalize(that.exp)
	} else {
		that.denormalize(d.exp)
	}
	return Number{d.val + that.val, d.exp}
}

// Sub implements - operation for decimal numbers
func (d Number) Sub(that Number) Number {
	that.val = that.val * -1
	return d.Add(that)
}

// Mul implements * operation for decimal numbers
func (d Number) Mul(that Number) Number {
	return Number{d.val * that.val, d.exp + that.exp}
}

// MulInt calculates d * n value.
func (d Number) MulInt(n int) Number {
	return Number{d.val * int64(n), d.exp}
}

// Neg returnds -d.
func (d Number) Neg() Number {
	return Number{d.val * -1, d.exp}
}

// Rat returns a rational representation of decimal.
func (d Number) Rat() *big.Rat {
	one := big.NewInt(1)
	ten := big.NewInt(10)
	if d.exp <= 0 {
		denom := new(big.Int).Exp(ten, big.NewInt(-int64(d.exp)), nil)
		return new(big.Rat).SetFrac(big.NewInt(d.val), denom)
	}
	mul := new(big.Int).Exp(ten, big.NewInt(int64(d.exp)), nil)
	num := new(big.Int).Mul(big.NewInt(d.val), mul)
	return new(big.Rat).SetFrac(num, one)
}

// ScaledVal scales decimal number to a given exponent and and returns
// internal number integer value. If given exponent is higher then internal
// number exponent this function will loose truncated digits.
//
// Example: number "12.99" with call ScaledVal(-4) would return 129900, with
// call ScaledVal(0) would return 12.
func (d Number) ScaledVal(exp int) int64 {
	if d.exp > exp {
		d.denormalize(exp)
		return d.val
	} else if d.exp < exp {
		for ; d.exp < exp; d.exp++ {
			d.val /= 10
		}
	}
	return d.val
}

// Cmp compares d and n and returns:
//
//   -1 if d <  n
//    0 if d == n
//   +1 if d >  n
//
func (d Number) Cmp(n Number) int {
	if n.exp < d.exp {
		d.denormalize(n.exp)
	} else {
		n.denormalize(d.exp)
	}
	// exponents must be the same for both values now
	switch {
	case d.val < n.val:
		return -1
	case d.val > n.val:
		return 1
	default:
		return 0
	}
}

// String function implements Stringer interface from fmt package.
func (d Number) String() string {
	var sign, str string
	if d.val < 0 {
		sign = "-"
		str = strconv.FormatInt(d.val*-1, 10)
	} else {
		str = strconv.FormatInt(d.val, 10)
	}
	if d.exp == 0 {
		return sign + str
	} else if d.exp > 0 {
		return sign + str + strings.Repeat("0", d.exp)
	}
	// Number have fractional part
	l := len(str)
	pos := l + d.exp
	if pos > 0 {
		return sign + str[:pos] + "." + str[pos:]
	}

	return sign + "0." + strings.Repeat("0", -1*pos) + str
}

// MarshalText function implements TextMarshaler interface from encoding package.
func (d Number) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText implements TextMarshaler interface from encoding package.
func (d *Number) UnmarshalText(data []byte) (err error) {
	parts := strings.SplitN(string(data), ".", 2)
	if len(parts) == 1 {
		d.val, err = strconv.ParseInt(parts[0], 10, 64)
		d.exp = 0
		return err
	}

	if parts[1] == "" {
		return fmt.Errorf("decimal.Number: parsing \"%s\": invalid syntax", data)
	}
	if parts[0] == "" {
		parts[0] = "0"
	}

	if d.val, err = strconv.ParseInt(parts[0]+parts[1], 10, 64); err != nil {
		return err
	}
	d.exp = -1 * len(parts[1])
	return nil
}

// Scan implements Scanner interface database/sql encoding package.
func (d *Number) Scan(src interface{}) (err error) {
	switch src.(type) {
	case []byte:
		err = d.UnmarshalText(src.([]byte))
	default:
		err = fmt.Errorf("decimal.Number: cannot convert %T to decimal.Number, only []byte is supported", src)
	}
	return
}

// Value implements Valuer interface from database/sql/driver package.
func (d Number) Value() (driver.Value, error) {
	return d.MarshalText()
}

// UnmarshalJSON implements Unmarshaler interface from encoding/json package.
func (d *Number) UnmarshalJSON(data []byte) error {
	if len(data) >= 2 && data[0] == '"' && data[len(data)-1] == '"' {
		return d.UnmarshalText(data[1 : len(data)-1])
	}
	return d.UnmarshalText(data)
}

// MarshalJSON implements Marshaler interface from encoding/json package.
func (d Number) MarshalJSON() ([]byte, error) {
	return d.MarshalText()
}
