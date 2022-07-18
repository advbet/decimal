package decimal

import (
	"fmt"
	"math"
	"testing"

	newDecimal "github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestFloat64(t *testing.T) {
	tests := []struct {
		n Number
		f float64
	}{
		{newDecimal.New(0, 0), 0.0},
		{newDecimal.New(1, 0), 1.0},
		{newDecimal.New(1, 1), 10.0},
		{newDecimal.New(10, 0), 10.0},
		{newDecimal.New(1, -1), 0.1},
		{newDecimal.New(-1, 0), -1.0},
		{newDecimal.New(-1, 1), -10.0},
		{newDecimal.New(-10, 0), -10.0},
		{newDecimal.New(-1, -1), -0.1},
		{newDecimal.New(123456, -3), 123.456},
		{newDecimal.New(17976931348623157, 292), math.MaxFloat64},
		{newDecimal.New(-17976931348623157, 292), -math.MaxFloat64},
		{newDecimal.New(5, -324), math.SmallestNonzeroFloat64},
		{newDecimal.New(4, -324), math.SmallestNonzeroFloat64},
		{newDecimal.New(3, -324), math.SmallestNonzeroFloat64},
		{newDecimal.New(2, -324), 0.0},
		{newDecimal.New(1, 309), math.Inf(1)},
		{newDecimal.New(-1, 309), math.Inf(-1)},
	}

	for _, test := range tests {
		f := test.n.InexactFloat64()
		assert.Equal(t, test.f, f, fmt.Sprintf("Number(%s).Float64() = %g", test.n, test.f))
	}
}

func TestFromFloat64(t *testing.T) {
	tests := []struct {
		f   float64
		n   Number
		err bool
	}{
		{0.0, newDecimal.New(0, 0), false},
		{1.0, newDecimal.New(1, 0), false},
		{10.0, newDecimal.New(1, 1), false},
		{0.1, newDecimal.New(1, -1), false},
		{-1.0, newDecimal.New(-1, 0), false},
		{-10.0, newDecimal.New(-1, 1), false},
		{-0.1, newDecimal.New(-1, -1), false},
		{123.456, newDecimal.New(123456, -3), false},
		{math.MaxFloat64, newDecimal.New(17976931348623157, 292), false},
		{-math.MaxFloat64, newDecimal.New(-17976931348623157, 292), false},
		{math.SmallestNonzeroFloat64, newDecimal.New(5, -324), false},
		// FIXME: ?
		{math.Inf(1), Number{}, true},
		{math.Inf(-1), Number{}, true},
		{math.NaN(), Number{}, true},
	}

	for _, test := range tests {
		n := newDecimal.NewFromFloat(test.f)
		assert.Equal(t, test.n, n, fmt.Sprintf("FromFloat(%g), expected %s", test.f, test.n))
	}
}
