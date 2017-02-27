package decimal

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
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
		a, err := FromString(test.a)
		assert.NoError(t, err)
		b, err := FromString(test.b)
		assert.NoError(t, err)
		sum, err := FromString(test.sum)
		assert.NoError(t, err)

		assert.Equal(t, sum, a.Add(b))               // sum = a + b
		assert.Equal(t, sum, a.Add(b).Add(b).Sub(b)) // sum = a + b + b -b
	}
}

func TestNumberMul(t *testing.T) {
	a := Number{12, 2}
	b := Number{12, -2}

	assert.Equal(t, Number{144, 0}, a.Mul(b))
}

func TestNumberScaledInt64(t *testing.T) {
	assert.Equal(t, int64(12), Number{1234, -2}.ScaledVal(0))
	assert.Equal(t, int64(1234), Number{1234, -2}.ScaledVal(-2))
	assert.Equal(t, int64(123400), Number{1234, -2}.ScaledVal(-4))
}

func TestNumberValExp(t *testing.T) {
	a := Number{1, 2}
	assert.Equal(t, int64(1), a.Val())
	assert.Equal(t, int(2), a.Exp())
}

func TestNumberString(t *testing.T) {
	assert.Equal(t, "123400", Number{1234, 2}.String())
	assert.Equal(t, "12340", Number{1234, 1}.String())
	assert.Equal(t, "1234", Number{1234, 0}.String())
	assert.Equal(t, "123.4", Number{1234, -1}.String())
	assert.Equal(t, "12.34", Number{1234, -2}.String())
	assert.Equal(t, "1.234", Number{1234, -3}.String())
	assert.Equal(t, "0.1234", Number{1234, -4}.String())
	assert.Equal(t, "0.001234", Number{1234, -6}.String())

	assert.Equal(t, "0", Number{0, 0}.String())

	assert.Equal(t, "-123400", Number{-1234, 2}.String())
	assert.Equal(t, "-12340", Number{-1234, 1}.String())
	assert.Equal(t, "-1234", Number{-1234, 0}.String())
	assert.Equal(t, "-123.4", Number{-1234, -1}.String())
	assert.Equal(t, "-12.34", Number{-1234, -2}.String())
	assert.Equal(t, "-1.234", Number{-1234, -3}.String())
	assert.Equal(t, "-0.1234", Number{-1234, -4}.String())
	assert.Equal(t, "-0.001234", Number{-1234, -6}.String())
}

func TestNumberMarshalText(t *testing.T) {
	res, err := Number{-1234, -6}.MarshalText()
	assert.Equal(t, []byte("-0.001234"), res)
	assert.Nil(t, err)
}

func TestNumberUnmarshalText(t *testing.T) {
	tests := []struct {
		str   string
		d     Number
		valid bool
	}{
		{"-12340", Number{-12340, 0}, true},
		{"-1234", Number{-1234, 0}, true},
		{"-123.4", Number{-1234, -1}, true},
		{"-12.34", Number{-1234, -2}, true},
		{"-1.234", Number{-1234, -3}, true},
		{"-0.1234", Number{-1234, -4}, true},
		{"-0.01234", Number{-1234, -5}, true},
		{"-0.001234", Number{-1234, -6}, true},
		{"-0.0012340", Number{-12340, -7}, true},
		{"-0.0000000", Number{0, -7}, true},
		{"-00", Number{0, 0}, true},
		{"-0", Number{0, 0}, true},
		{"0", Number{0, 0}, true},
		{"00", Number{0, 0}, true},
		{"0.0000000", Number{0, -7}, true},
		{"0.0012340", Number{12340, -7}, true},
		{"0.001234", Number{1234, -6}, true},
		{"0.01234", Number{1234, -5}, true},
		{"0.1234", Number{1234, -4}, true},
		{"1.234", Number{1234, -3}, true},
		{"12.34", Number{1234, -2}, true},
		{"123.4", Number{1234, -1}, true},
		{"1234", Number{1234, 0}, true},
		{"12340", Number{12340, 0}, true},

		// this is a side effect of using strconv.ParseInt
		{"+1", Number{1, 0}, true},
		{"+1.2", Number{12, -1}, true},

		{" 1", Number{}, false},
		{"1 ", Number{}, false},
		{"1 2", Number{}, false},
		{"1,2", Number{}, false},
		{"1.+2", Number{}, false},
		{"--1", Number{}, false},
		{".-2", Number{}, false},
		{".2", Number{}, false},
		{".0", Number{}, false},
		{".", Number{}, false},
		{"1.-2", Number{}, false},
		{"1.-", Number{}, false},
		{"1.", Number{}, false},
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
	assert.Equal(t, Number{15, -3}, a)

	assert.NotNil(t, a.Scan("Strings are not supported"))
}

func TestNumberValue(t *testing.T) {
	val, err := Number{123, -1}.Value()
	assert.Nil(t, err)
	assert.Equal(t, []byte("12.3"), val)
}

func TestNumberUnmarshalJSON(t *testing.T) {
	var data struct {
		Num Number `json:"num"`
	}
	err := json.Unmarshal([]byte(`{"num": 123.456}`), &data)
	expected := Number{
		val: 123456,
		exp: -3,
	}
	assert.NoError(t, err, "unmarshaling should not fail")
	assert.Equal(t, expected, data.Num)
}

func TestNumberUnmarshalJSONString(t *testing.T) {
	var data struct {
		Num Number `json:"num"`
	}
	err := json.Unmarshal([]byte(`{"num": "123.456"}`), &data)
	expected := Number{
		val: 123456,
		exp: -3,
	}
	assert.NoError(t, err, "unmarshaling should not fail")
	assert.Equal(t, expected, data.Num)
}

func TestNumberMarshalJSON(t *testing.T) {
	data := struct {
		Num Number `json:"num"`
	}{Number{
		val: 123456,
		exp: -3,
	}}
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
		x:        New(0, 0),
		y:        New(0, 0),
		expected: 0,
	}, {
		x:        New(5, 0),
		y:        New(2, 0),
		expected: 1,
	}, {
		x:        New(2, 0),
		y:        New(5, 0),
		expected: -1,
	}, {
		x:        New(10, -1),
		y:        New(1, 0),
		expected: 0,
	}, {
		x:        New(50, -1),
		y:        New(2, 0),
		expected: 1,
	}, {
		x:        New(1, 0),
		y:        New(10, -1),
		expected: 0,
	}, {
		x:        New(2, 0),
		y:        New(50, -1),
		expected: -1,
	}}

	for _, test := range tests {
		actual := test.x.Cmp(test.y)
		assert.Equal(t, test.expected, actual, fmt.Sprintf("%s cmp %s", test.x, test.y))
	}
}

func TestFromInt(t *testing.T) {
	assert.Equal(t, Number{5, 0}, FromInt(5))
	assert.Equal(t, Number{0, 0}, FromInt(0))
	assert.Equal(t, Number{-5, 0}, FromInt(-5))
}

func TestNew(t *testing.T) {
	assert.Equal(t, Number{5, 0}, New(5, 0))
	assert.Equal(t, Number{0, 0}, New(0, 0))
	assert.Equal(t, Number{-5, 0}, New(-5, 0))

	// Assert value is not normalized
	assert.Equal(t, Number{50, -1}, New(50, -1))
	assert.Equal(t, Number{50, 0}, New(50, 0))
	assert.Equal(t, Number{50, 1}, New(50, 1))
}

func TestNumberMulInt(t *testing.T) {
	tests := []struct {
		x        Number
		y        int
		expected Number
	}{{
		x:        Number{0, 0},
		y:        5,
		expected: Number{0, 0},
	}, {
		x:        Number{1, 0},
		y:        5,
		expected: Number{5, 0},
	}, {
		// Assert exponent is not normalized
		x:        Number{2, -1},
		y:        5,
		expected: Number{10, -1},
	}, {
		// Assert exponent is not normalized
		x:        Number{2, -1},
		y:        0,
		expected: Number{0, -1},
	}}

	for _, test := range tests {
		actual := test.x.MulInt(test.y)
		assert.Equal(t, test.expected, actual, fmt.Sprintf("%s * %d", test.x, test.y))
	}
}

func TestNumberIsZero(t *testing.T) {
	assert.True(t, Number{0, -1}.IsZero())
	assert.True(t, Number{0, 0}.IsZero())
	assert.True(t, Number{0, 1}.IsZero())
	assert.False(t, Number{1, 0}.IsZero())
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
		// Ceil
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

		result := num.Round(test.exp, test.rule)
		assert.Equal(t, res, result, fmt.Sprintf("%s round(%d, %d)", test.num, test.exp, test.rule))
	}
}

func BenchmarkNumberScanRoundMarshal(b *testing.B) {
	var d Number
	for i := 0; i < b.N; i++ {
		_ = d.Scan([]byte("123456789.123456789"))
		d = d.Round(-2, RoundMath)
		_, _ = d.MarshalText()
	}
}

func BenchmarkExternalNumberScanRoundMarshal(b *testing.B) {
	var d decimal.Decimal
	for i := 0; i < b.N; i++ {
		_ = d.Scan([]byte("123456789.123456789"))
		d = d.Round(2)
		_, _ = d.MarshalText()
	}
}

func BenchmarkNumberDenormalize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		n := Number{val: 123456, exp: 0}
		n.denormalize(-8)
	}
}
