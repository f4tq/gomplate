package funcs

import (
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseNum(t *testing.T) {
	i, f, _ := parseNum("42")
	assert.Equal(t, int64(42), i)
	assert.Equal(t, int64(0), f)

	i, f, _ = parseNum(42)
	assert.Equal(t, int64(42), i)
	assert.Equal(t, int64(0), f)

	i, f, _ = parseNum(big.NewInt(42))
	assert.Equal(t, int64(42), i)
	assert.Equal(t, int64(0), f)

	i, f, _ = parseNum(big.NewFloat(42.0))
	assert.Equal(t, int64(42), i)
	assert.Equal(t, int64(0), f)

	i, f, _ = parseNum(uint64(math.MaxInt64))
	assert.Equal(t, int64(uint64(math.MaxInt64)), i)
	assert.Equal(t, int64(0), f)

	i, f, _ = parseNum("9223372036854775807.9223372036854775807")
	assert.Equal(t, int64(9223372036854775807), i)
	assert.Equal(t, int64(9223372036854775807), f)

	_, _, err := parseNum("bogus.9223372036854775807")
	assert.Error(t, err)

	_, _, err = parseNum("bogus")
	assert.Error(t, err)

	_, _, err = parseNum("1.2.3")
	assert.Error(t, err)

	_, _, err = parseNum(1.1)
	assert.Error(t, err)

	i, f, err = parseNum(nil)
	assert.Zero(t, i)
	assert.Zero(t, f)
	assert.NoError(t, err)
}
