package decimal

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloat64(t *testing.T) {
	tests := []struct {
		n Number
		f float64
	}{
		{Number{0, 0}, 0.0},
		{Number{1, 0}, 1.0},
		{Number{1, 1}, 10.0},
		{Number{10, 0}, 10.0},
		{Number{1, -1}, 0.1},
		{Number{-1, 0}, -1.0},
		{Number{-1, 1}, -10.0},
		{Number{-10, 0}, -10.0},
		{Number{-1, -1}, -0.1},
		{Number{123456, -3}, 123.456},
		{Number{17976931348623157, 292}, math.MaxFloat64},
		{Number{-17976931348623157, 292}, -math.MaxFloat64},
		{Number{5, -324}, math.SmallestNonzeroFloat64},
		{Number{4, -324}, math.SmallestNonzeroFloat64},
		{Number{3, -324}, math.SmallestNonzeroFloat64},
		{Number{2, -324}, 0.0},
		{Number{1, 309}, math.Inf(1)},
		{Number{-1, 309}, math.Inf(-1)},
	}

	for _, test := range tests {
		f := test.n.Float64()
		assert.Equal(t, test.f, f, fmt.Sprintf("Number(%s).Float64() = %g", test.n, test.f))
	}
}

func TestFromFloat64(t *testing.T) {
	tests := []struct {
		f   float64
		n   Number
		err bool
	}{
		{0.0, Number{0, 0}, false},
		{1.0, Number{1, 0}, false},
		{10.0, Number{1, 1}, false},
		{0.1, Number{1, -1}, false},
		{-1.0, Number{-1, 0}, false},
		{-10.0, Number{-1, 1}, false},
		{-0.1, Number{-1, -1}, false},
		{123.456, Number{123456, -3}, false},
		{math.MaxFloat64, Number{17976931348623157, 292}, false},
		{-math.MaxFloat64, Number{-17976931348623157, 292}, false},
		{math.SmallestNonzeroFloat64, Number{5, -324}, false},
		{math.Inf(1), Number{}, true},
		{math.Inf(-1), Number{}, true},
		{math.NaN(), Number{}, true},
	}

	for _, test := range tests {
		n, err := FromFloat64(test.f)
		if test.err {
			assert.Error(t, err, fmt.Sprintf("FromFloat(%g), expected error", test.f))
			continue
		}
		assert.NoError(t, err)
		assert.Equal(t, test.n, n, fmt.Sprintf("FromFloat(%g), expected %s", test.f, test.n))
	}
}
