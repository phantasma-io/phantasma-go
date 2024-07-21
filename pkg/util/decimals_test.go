package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrimWholePrefix(t *testing.T) {
	assert.Equal(t, "1", trimWholePrefix("01", "0"))
	assert.Equal(t, "1000000001", trimWholePrefix("001000000001", "0"))
	assert.Equal(t, "", trimWholePrefix("00000000", "0"))
	assert.Equal(t, "", trimWholePrefix("0", "0"))
	assert.Equal(t, "", trimWholePrefix("", "0"))
	assert.Equal(t, "100000000", trimWholePrefix("100000000", "0"))
	assert.Equal(t, "1000000001", trimWholePrefix("1000000001", "0"))
}

func TestTrimWholeSuffix(t *testing.T) {
	assert.Equal(t, "1", trimWholeSuffix("100000000", "0"))
	assert.Equal(t, "1000000001", trimWholeSuffix("1000000001", "0"))
	assert.Equal(t, "", trimWholeSuffix("00000000", "0"))
	assert.Equal(t, "", trimWholeSuffix("0", "0"))
	assert.Equal(t, "", trimWholeSuffix("", "0"))
}

func TestStringIsZeroOrEmptyBigint(t *testing.T) {
	assert.Equal(t, false, stringIsZeroOrEmptyBigint("100000000"))
	assert.Equal(t, true, stringIsZeroOrEmptyBigint("00000000"))
	assert.Equal(t, true, stringIsZeroOrEmptyBigint("0"))
	assert.Equal(t, true, stringIsZeroOrEmptyBigint(""))
}

func TestConvertDecimalsEx(t *testing.T) {
	assert.Equal(t, "0.0000000001", ConvertDecimalsEx("1", 10, "."))
	assert.Equal(t, "0.000009", ConvertDecimalsEx("90000", 10, "."))
	assert.Equal(t, "0.01", ConvertDecimalsEx("100000000", 10, "."))
	assert.Equal(t, "0.1", ConvertDecimalsEx("1000000000", 10, "."))
	assert.Equal(t, "0", ConvertDecimalsEx("000000000", 10, "."))
	assert.Equal(t, "0", ConvertDecimalsEx("0000000000", 10, "."))
	assert.Equal(t, "0", ConvertDecimalsEx("00000000000", 10, "."))
	assert.Equal(t, "1", ConvertDecimalsEx("10000000000", 10, "."))
	assert.Equal(t, "1.1", ConvertDecimalsEx("11000000000", 10, "."))
	assert.Equal(t, "10", ConvertDecimalsEx("100000000000", 10, "."))
	assert.Equal(t, "11", ConvertDecimalsEx("110000000000", 10, "."))
	assert.Equal(t, "1.19", ConvertDecimalsEx("11900000000", 10, "."))
	assert.Equal(t, "1.019", ConvertDecimalsEx("10190000000", 10, "."))
}

func TestConvertDecimalsBackEx(t *testing.T) {
	assert.Panics(t, func() { ConvertDecimalsBackEx("0.0000000001", 0, ".", false) })
	assert.Equal(t, "0", ConvertDecimalsBackEx("0.0000000000", 10, ".", true))
	assert.Equal(t, "1", ConvertDecimalsBackEx("0.0000000001", 10, ".", true))
	assert.Equal(t, "1", ConvertDecimalsBackEx("0.00000000010", 10, ".", true))
	assert.Equal(t, "1", ConvertDecimalsBackEx("0.00000000011", 10, ".", false))
	assert.Equal(t, "90000", ConvertDecimalsBackEx("0.000009", 10, ".", true))
	assert.Equal(t, "100000000", ConvertDecimalsBackEx("0.01", 10, ".", true))
	assert.Equal(t, "1000000000", ConvertDecimalsBackEx("0.1", 10, ".", true))
	assert.Equal(t, "1", ConvertDecimalsBackEx("0.1", 1, ".", true))
	assert.Equal(t, "10", ConvertDecimalsBackEx("0.1", 2, ".", true))
	assert.Equal(t, "0", ConvertDecimalsBackEx("0", 0, ".", true))
	assert.Equal(t, "0", ConvertDecimalsBackEx("0", 10, ".", true))
	assert.Equal(t, "0", ConvertDecimalsBackEx("0.0", 10, ".", true))
	assert.Equal(t, "0", ConvertDecimalsBackEx("0,0", 10, ",", true))
	assert.Equal(t, "1", ConvertDecimalsBackEx("1", 0, ".", true))
	assert.Equal(t, "9999", ConvertDecimalsBackEx("9999", 0, ".", true))
	assert.Equal(t, "10000000000", ConvertDecimalsBackEx("1", 10, ".", true))
	assert.Equal(t, "11000000000", ConvertDecimalsBackEx("1.1", 10, ".", true))
	assert.Equal(t, "100000000000", ConvertDecimalsBackEx("10", 10, ".", true))
	assert.Equal(t, "110000000000", ConvertDecimalsBackEx("11", 10, ".", true))
	assert.Equal(t, "11900000000", ConvertDecimalsBackEx("1.19", 10, ".", true))
	assert.Equal(t, "10190000000", ConvertDecimalsBackEx("1.019", 10, ".", true))
	assert.Equal(t, "10190000000", ConvertDecimalsBackEx("000000001.019", 10, ".", true))
	assert.Equal(t, "10190000000", ConvertDecimalsBackEx("000000001.01900000000000000000", 10, ".", true))
	assert.Equal(t, "10190000000", ConvertDecimalsBackEx("000000001.019000000000000000001", 10, ".", false))
}

func TestConvertDecimalsBack(t *testing.T) {
	assert.Equal(t, "1", ConvertDecimalsBack("0.0000000001", 10).String())
	assert.Equal(t, "90000", ConvertDecimalsBack("0.000009", 10).String())
	assert.Equal(t, "100000000", ConvertDecimalsBack("0.01", 10).String())
	assert.Equal(t, "1000000000", ConvertDecimalsBack("0.1", 10).String())
	assert.NotEqual(t, "1000000000", ConvertDecimalsBack("0,1", 10).String())
	assert.Equal(t, "0", ConvertDecimalsBack("0", 10).String())
	assert.Equal(t, "0", ConvertDecimalsBack("0.0", 10).String())
	assert.Equal(t, "10000000000", ConvertDecimalsBack("1", 10).String())
	assert.Equal(t, "11000000000", ConvertDecimalsBack("1.1", 10).String())
	assert.Equal(t, "100000000000", ConvertDecimalsBack("10", 10).String())
	assert.Equal(t, "110000000000", ConvertDecimalsBack("11", 10).String())
	assert.Equal(t, "11900000000", ConvertDecimalsBack("1.19", 10).String())
	assert.Equal(t, "10190000000", ConvertDecimalsBack("1.019", 10).String())
}
