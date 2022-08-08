package decimal

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	newDecimal "github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNumberAddSub(t *testing.T) {
	tests := []struct {
		a   string
		b   string
		sum string
	}{
		{"56.78", "123400", "123456.78"},
		{"123400", "56.78", "123456.78"},
		{"-12.34", "12.34", "0.00"},
		{"12.34", "-12.34", "0.00"},
		{"56.78", "0.22", "57.00"},
		{"0", "0", "0"},
		{"0", "0.00", "0.00"},
		{"0.00", "0", "0.00"},
	}

	for _, test := range tests {
		a, err := newDecimal.NewFromString(test.a)
		assert.NoError(t, err)
		b, err := newDecimal.NewFromString(test.b)
		assert.NoError(t, err)
		sum, err := newDecimal.NewFromString(test.sum)
		assert.NoError(t, err)

		sum1 := a.Add(b)
		sum2 := a.Add(b).Add(b).Sub(b)

		// Zero should never be compared with == or != directly, please use decimal.Equal or decimal.Cmp instead.
		if sum.IsZero() {
			assert.Equal(t, true, sum.Equal(sum1)) // sum = a + b
			assert.Equal(t, true, sum.Equal(sum2)) // sum = a + b + b - b
		} else {
			assert.Equal(t, sum, sum1) // sum = a + b
			assert.Equal(t, sum, sum2) // sum = a + b + b -b
		}
	}
}

func TestNumberMul(t *testing.T) {
	a := newDecimal.New(12, 2)
	b := newDecimal.New(12, -2)

	assert.Equal(t, newDecimal.New(144, 0), a.Mul(b))
}

func TestNumberScaledInt64(t *testing.T) {
	assert.Equal(t, int64(12), ScaledVal(newDecimal.New(1234, -2), 0))
	assert.Equal(t, int64(1234), ScaledVal(newDecimal.New(1234, -2), -2))
	assert.Equal(t, int64(123400), ScaledVal(newDecimal.New(1234, -2), -4))
}

func TestNumberValExp(t *testing.T) {
	a := newDecimal.New(1, 2)
	assert.Equal(t, int64(1), a.CoefficientInt64())
	assert.Equal(t, int32(2), a.Exponent())
}

func TestNumberString(t *testing.T) {
	assert.Equal(t, "123400", newDecimal.New(1234, 2).String())
	assert.Equal(t, "12340", newDecimal.New(1234, 1).String())
	assert.Equal(t, "1234", newDecimal.New(1234, 0).String())
	assert.Equal(t, "123.4", newDecimal.New(1234, -1).String())
	assert.Equal(t, "12.34", newDecimal.New(1234, -2).String())
	assert.Equal(t, "1.234", newDecimal.New(1234, -3).String())
	assert.Equal(t, "0.1234", newDecimal.New(1234, -4).String())
	assert.Equal(t, "0.001234", newDecimal.New(1234, -6).String())

	assert.Equal(t, "0", newDecimal.New(0, 0).String())

	assert.Equal(t, "-123400", newDecimal.New(-1234, 2).String())
	assert.Equal(t, "-12340", newDecimal.New(-1234, 1).String())
	assert.Equal(t, "-1234", newDecimal.New(-1234, 0).String())
	assert.Equal(t, "-123.4", newDecimal.New(-1234, -1).String())
	assert.Equal(t, "-12.34", newDecimal.New(-1234, -2).String())
	assert.Equal(t, "-1.234", newDecimal.New(-1234, -3).String())
	assert.Equal(t, "-0.1234", newDecimal.New(-1234, -4).String())
	assert.Equal(t, "-0.001234", newDecimal.New(-1234, -6).String())
}

func TestNumberMarshalText(t *testing.T) {
	res, err := newDecimal.New(-1234, -6).MarshalText()
	assert.Equal(t, []byte("-0.001234"), res)
	assert.Nil(t, err)
}

func TestNumberUnmarshalText(t *testing.T) {
	tests := []struct {
		str   string
		d     Number
		valid bool
	}{
		{"-12340", newDecimal.New(-12340, 0), true},
		{"-1234", newDecimal.New(-1234, 0), true},
		{"-123.4", newDecimal.New(-1234, -1), true},
		{"-12.34", newDecimal.New(-1234, -2), true},
		{"-1.234", newDecimal.New(-1234, -3), true},
		{"-0.1234", newDecimal.New(-1234, -4), true},
		{"-0.01234", newDecimal.New(-1234, -5), true},
		{"-0.001234", newDecimal.New(-1234, -6), true},
		{"-0.0012340", newDecimal.New(-12340, -7), true},
		{"-0.0000000", newDecimal.New(0, -7), true},
		{"-00", newDecimal.New(0, 0), true},
		{"-0", newDecimal.New(0, 0), true},
		{"0", newDecimal.New(0, 0), true},
		{"00", newDecimal.New(0, 0), true},
		{"0.0000000", newDecimal.New(0, -7), true},
		{"0.0012340", newDecimal.New(12340, -7), true},
		{"0.001234", newDecimal.New(1234, -6), true},
		{"0.01234", newDecimal.New(1234, -5), true},
		{"0.1234", newDecimal.New(1234, -4), true},
		{"1.234", newDecimal.New(1234, -3), true},
		{"12.34", newDecimal.New(1234, -2), true},
		{"123.4", newDecimal.New(1234, -1), true},
		{"1234", newDecimal.New(1234, 0), true},
		{"12340", newDecimal.New(12340, 0), true},
		{".2", newDecimal.New(2, -1), true},
		{".0", newDecimal.New(0, -1), true},
		{"-.4", newDecimal.New(-4, -1), true},

		{".-2", newDecimal.New(-2, -2), true},
		{"1.", newDecimal.New(1, 0), true},

		// this is a side effect of using strconv.ParseInt
		{"+1", newDecimal.New(1, 0), true},
		{"+1.2", newDecimal.New(12, -1), true},

		{" 1", Number{}, false},
		{"1 ", Number{}, false},
		{"1 2", Number{}, false},
		{"1,2", Number{}, false},
		{"1.+2", Number{}, false},
		{"--1", Number{}, false},
		{".", Number{}, false},
		{"1.-2", Number{}, false},
		{"1.-", Number{}, false},
		{"a1", Number{}, false},
		{"1.a2", Number{}, false},
		{"a3.9", Number{}, false},
	}

	for _, test := range tests {
		d := Number{}
		err := d.UnmarshalText([]byte(test.str))
		if test.valid {
			assert.NoError(t, err)
			assert.Equal(t, test.d, d)
		} else {
			assert.Error(t, err)
			assert.Equal(t, test.d, d, fmt.Sprintf("\"%s\" -> %s", test.str, d))
		}
	}
}

func TestNumberScan(t *testing.T) {
	var a Number
	assert.Nil(t, a.Scan([]byte("0.015")))
	assert.Equal(t, newDecimal.New(15, -3), a)

	assert.NotNil(t, a.Scan("Strings are not supported"))
}

// FIXME: this needs discussion, output changed:
//    expected: []uint8([]byte{0x31, 0x32, 0x2e, 0x33})
//    actual  : string("12.3")
func TestNumberValue(t *testing.T) {
	val, err := newDecimal.New(123, -1).Value()
	assert.Nil(t, err)
	assert.Equal(t, []byte("12.3"), val)
}

func TestNumberUnmarshalJSON(t *testing.T) {
	var data struct {
		Num Number `json:"num"`
	}
	err := json.Unmarshal([]byte(`{"num": 123.456}`), &data)
	expected := newDecimal.New(123456, -3)
	assert.NoError(t, err, "unmarshalling should not fail")
	assert.Equal(t, expected, data.Num)
}

func TestNumberUnmarshalJSONString(t *testing.T) {
	var data struct {
		Num Number `json:"num"`
	}
	err := json.Unmarshal([]byte(`{"num": "123.456"}`), &data)
	expected := newDecimal.New(123456, -3)
	assert.NoError(t, err, "unmarshalling should not fail")
	assert.Equal(t, expected, data.Num)
}

func TestNumberMarshalJSON(t *testing.T) {
	data := struct {
		Num Number `json:"num"`
	}{
		newDecimal.New(123456, -3),
	}
	blob, err := json.Marshal(&data)
	assert.NoError(t, err, "json marshaling should not fail")
	assert.Equal(t, `{"num":123.456}`, string(blob))
}

func TestNumberCmp(t *testing.T) {
	tests := []struct {
		x        Number
		y        Number
		expected int
	}{{
		x:        newDecimal.New(0, 0),
		y:        newDecimal.New(0, 0),
		expected: 0,
	}, {
		x:        newDecimal.New(5, 0),
		y:        newDecimal.New(2, 0),
		expected: 1,
	}, {
		x:        newDecimal.New(2, 0),
		y:        newDecimal.New(5, 0),
		expected: -1,
	}, {
		x:        newDecimal.New(10, -1),
		y:        newDecimal.New(1, 0),
		expected: 0,
	}, {
		x:        newDecimal.New(50, -1),
		y:        newDecimal.New(2, 0),
		expected: 1,
	}, {
		x:        newDecimal.New(1, 0),
		y:        newDecimal.New(10, -1),
		expected: 0,
	}, {
		x:        newDecimal.New(2, 0),
		y:        newDecimal.New(50, -1),
		expected: -1,
	}}

	for _, test := range tests {
		actual := test.x.Cmp(test.y)
		assert.Equal(t, test.expected, actual, fmt.Sprintf("%s cmp %s", test.x, test.y))
	}
}

func TestFromInt(t *testing.T) {
	assert.Equal(t, newDecimal.New(5, 0), newDecimal.NewFromInt(5))
	assert.Equal(t, newDecimal.New(0, 0), newDecimal.NewFromInt(0))
	assert.Equal(t, newDecimal.New(-5, 0), newDecimal.NewFromInt(-5))
}

func TestFromRat(t *testing.T) {
	tests := []struct {
		rat      *big.Rat
		exp      int
		expected Number
	}{
		{
			rat:      big.NewRat(1234, 100),
			exp:      -2,
			expected: newDecimal.New(1234, -2),
		},
		{
			rat:      big.NewRat(1234, 100),
			exp:      -1,
			expected: newDecimal.New(123, -1),
		},
		{
			rat:      big.NewRat(1234, 100),
			exp:      0,
			expected: newDecimal.New(12, 0),
		},
		{
			rat:      big.NewRat(1234, 100),
			exp:      1,
			expected: newDecimal.New(1, 1),
		},
		{
			rat:      big.NewRat(1000000000, 3),
			exp:      0,
			expected: newDecimal.New(333333333, 0),
		},
		{
			rat:      big.NewRat(1000000000, 3),
			exp:      -7,
			expected: newDecimal.New(3333333333333333, -7),
		},
		{
			rat:      big.NewRat(1000000000, 3),
			exp:      -8,
			expected: newDecimal.New(33333333333333331, -8),
		},
		{
			rat:      big.NewRat(1000000000, 3),
			exp:      -9,
			expected: newDecimal.New(333333333333333313, -9),
		},
		{
			rat:      big.NewRat(1000000000, 3),
			exp:      -10,
			expected: newDecimal.New(3333333333333333135, -10),
		},
		{
			rat:      big.NewRat(1000000000, 3),
			exp:      -11,
			expected: newDecimal.NewFromFloatWithExponent(333333333.33333331347, -11),
		},
		{
			rat:      big.NewRat(1000000000, 3),
			exp:      -12,
			expected: newDecimal.NewFromFloatWithExponent(333333333.33333331347, -12),
		},
		{
			rat:      big.NewRat(1000000000, 3),
			exp:      -13,
			expected: newDecimal.NewFromFloatWithExponent(333333333.33333331347, -13),
		},
	}

	for _, tt := range tests {
		actual := NewFromRat(tt.rat, tt.exp)
		assert.Equalf(t, tt.expected, actual, "%s (%d) expected %s, got %s (%d * 10^%d)", tt.rat, tt.exp, tt.expected, actual, actual.CoefficientInt64(), actual.Exponent())
	}
}

func TestNumberMulInt(t *testing.T) {
	tests := []struct {
		x        Number
		y        int
		expected Number
	}{{
		x:        newDecimal.New(0, 0),
		y:        5,
		expected: newDecimal.New(0, 0),
	}, {
		x:        newDecimal.New(1, 0),
		y:        5,
		expected: newDecimal.New(5, 0),
	}, {
		// Assert exponent is not normalized
		x:        newDecimal.New(2, -1),
		y:        5,
		expected: newDecimal.New(10, -1),
	}, {
		// Assert exponent is not normalized
		x:        newDecimal.New(2, -1),
		y:        0,
		expected: newDecimal.New(0, -1),
	}}

	for _, test := range tests {
		actual := MulInt(test.x, test.y)
		assert.Equal(t, test.expected, actual, fmt.Sprintf("%s * %d", test.x, test.y))
	}
}

func TestNumberIsZero(t *testing.T) {
	assert.True(t, newDecimal.Zero.IsZero())
	assert.True(t, newDecimal.New(0, -1).IsZero())
	assert.True(t, newDecimal.New(0, 0).IsZero())
	assert.True(t, newDecimal.New(0, 1).IsZero())
	assert.False(t, newDecimal.New(1, 0).IsZero())
}

func TestNumberRound(t *testing.T) {
	tests := []struct {
		rule   RoundRule
		num    string
		exp    int
		result string
	}{
		{RoundTruncate, "1.23", -4, "1.2300"}, // Scale down
		{RoundTruncate, "1.23", -2, "1.23"},   // Noop
		// Truncate
		{RoundTruncate, "123.45", -1, "123.4"},
		{RoundTruncate, "-123.45", -1, "-123.4"},
		// Floor
		{RoundFloor, "123.4500", -2, "123.45"},
		{RoundFloor, "123.45", -1, "123.4"},
		{RoundFloor, "-123.45", -1, "-123.5"},
		//// Ceil
		{RoundCeil, "123.4500", -2, "123.45"},
		{RoundCeil, "123.45", -1, "123.5"},
		{RoundCeil, "-123.45", -1, "-123.4"},
		// Math rounding, positive numbers
		{RoundMath, "0.45", -1, "0.5"},
		{RoundMath, "-0.45", -1, "-0.5"},
		// Bankers rounding, positive numbers, round-up
		{RoundBankers, "0.349999", -1, "0.3"},
		{RoundBankers, "0.350000", -1, "0.4"},
		{RoundBankers, "0.350001", -1, "0.4"},
		// Bankers rounding, positive numbers, round-down
		{RoundBankers, "0.449999", -1, "0.4"},
		{RoundBankers, "0.450000", -1, "0.4"},
		{RoundBankers, "0.450001", -1, "0.5"},
		// Bankers rounding, negative numbers, round-up
		{RoundBankers, "-0.349999", -1, "-0.3"},
		{RoundBankers, "-0.350000", -1, "-0.4"},
		{RoundBankers, "-0.350001", -1, "-0.4"},
		// Bankers rounding, negative numbers, round-down
		{RoundBankers, "-0.449999", -1, "-0.4"},
		{RoundBankers, "-0.450000", -1, "-0.4"},
		{RoundBankers, "-0.450001", -1, "-0.5"},
		// Bankers rounding to zero
		{RoundBankers, "0.500000", 0, "0"},
		{RoundBankers, "-0.500000", 0, "0"},
	}

	var err error
	var num, res Number

	for _, test := range tests {
		err = num.UnmarshalText([]byte(test.num))
		assert.NoError(t, err)
		err = res.UnmarshalText([]byte(test.result))
		assert.NoError(t, err)

		result := Round(num, test.exp, test.rule)

		// Zero should never be compared with == or != directly, please use decimal.Equal or decimal.Cmp instead.
		if result.IsZero() {
			assert.Equal(t, true, result.Equal(res), fmt.Sprintf("%s round(%d, %d)", test.num, test.exp, test.rule))
		} else {
			assert.Equal(t, res, result, fmt.Sprintf("%s round(%d, %d)", test.num, test.exp, test.rule))
		}
	}
}

func TestDecimalNeg(t *testing.T) {
	tests := []struct {
		n        Number
		expected Number
	}{
		{
			n:        newDecimal.Zero,
			expected: newDecimal.Zero,
		},
		{
			n:        newDecimal.NewFromInt(1),
			expected: newDecimal.NewFromInt(-1),
		},
		{
			n:        newDecimal.NewFromInt(-1),
			expected: newDecimal.NewFromInt(1),
		},
		{
			n:        newDecimal.New(13, 2),
			expected: newDecimal.New(-13, 2),
		},
		{
			n:        newDecimal.New(13, -2),
			expected: newDecimal.New(-13, -2),
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.n.Neg())
	}
}

func TestDecimalRat(t *testing.T) {
	tests := []struct {
		n Number
		r *big.Rat
	}{
		{
			n: newDecimal.New(1, 0),
			r: big.NewRat(1, 1),
		},
		{
			n: newDecimal.New(5, -1),
			r: big.NewRat(1, 2),
		},
		{
			n: newDecimal.New(5, 1),
			r: big.NewRat(50, 1),
		},
	}
	for _, test := range tests {
		a := test.n.Rat()
		assert.Equal(t, test.r.Cmp(a), 0)
	}
}

func BenchmarkNumberScanRoundMarshal(b *testing.B) {
	var d Number
	for i := 0; i < b.N; i++ {
		_ = d.Scan([]byte("123456789.123456789"))
		d = Round(d, -2, RoundMath)
		_, _ = d.MarshalText()
	}
}

func BenchmarkExternalNumberScanRoundMarshal(b *testing.B) {
	var d newDecimal.Decimal
	for i := 0; i < b.N; i++ {
		_ = d.Scan([]byte("123456789.123456789"))
		d = d.Round(2)
		_, _ = d.MarshalText()
	}
}
