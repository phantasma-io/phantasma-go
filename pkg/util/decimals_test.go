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

func TestBigintStringToDecimalString(t *testing.T) {
	assert.Equal(t, "0.0000000001", BigintStringToDecimalString("1", 10))
	assert.Equal(t, "0.000009", BigintStringToDecimalString("90000", 10))
	assert.Equal(t, "0.01", BigintStringToDecimalString("100000000", 10))
	assert.Equal(t, "0.1", BigintStringToDecimalString("1000000000", 10))
	assert.Equal(t, "0", BigintStringToDecimalString("000000000", 10))
	assert.Equal(t, "0", BigintStringToDecimalString("0000000000", 10))
	assert.Equal(t, "0", BigintStringToDecimalString("00000000000", 10))
	assert.Equal(t, "1", BigintStringToDecimalString("10000000000", 10))
	assert.Equal(t, "1.1", BigintStringToDecimalString("11000000000", 10))
	assert.Equal(t, "10", BigintStringToDecimalString("100000000000", 10))
	assert.Equal(t, "11", BigintStringToDecimalString("110000000000", 10))
	assert.Equal(t, "1.19", BigintStringToDecimalString("11900000000", 10))
	assert.Equal(t, "1.019", BigintStringToDecimalString("10190000000", 10))
}

func TestBigintStringFromDecimalStringEx(t *testing.T) {
	assert.Panics(t, func() { BigintStringFromDecimalStringEx("0.0000000001", 0, ".", false) })
	assert.Equal(t, "0", BigintStringFromDecimalStringEx("0.0000000000", 10, ".", true))
	assert.Equal(t, "1", BigintStringFromDecimalStringEx("0.0000000001", 10, ".", true))
	assert.Equal(t, "1", BigintStringFromDecimalStringEx("0.00000000010", 10, ".", true))
	assert.Equal(t, "1", BigintStringFromDecimalStringEx("0.00000000011", 10, ".", false))
	assert.Equal(t, "90000", BigintStringFromDecimalStringEx("0.000009", 10, ".", true))
	assert.Equal(t, "100000000", BigintStringFromDecimalStringEx("0.01", 10, ".", true))
	assert.Equal(t, "1000000000", BigintStringFromDecimalStringEx("0.1", 10, ".", true))
	assert.Equal(t, "1", BigintStringFromDecimalStringEx("0.1", 1, ".", true))
	assert.Equal(t, "10", BigintStringFromDecimalStringEx("0.1", 2, ".", true))
	assert.Equal(t, "0", BigintStringFromDecimalStringEx("0", 0, ".", true))
	assert.Equal(t, "0", BigintStringFromDecimalStringEx("0", 10, ".", true))
	assert.Equal(t, "0", BigintStringFromDecimalStringEx("0.0", 10, ".", true))
	assert.Equal(t, "0", BigintStringFromDecimalStringEx("0,0", 10, ",", true))
	assert.Equal(t, "1", BigintStringFromDecimalStringEx("1", 0, ".", true))
	assert.Equal(t, "9999", BigintStringFromDecimalStringEx("9999", 0, ".", true))
	assert.Equal(t, "10000000000", BigintStringFromDecimalStringEx("1", 10, ".", true))
	assert.Equal(t, "11000000000", BigintStringFromDecimalStringEx("1.1", 10, ".", true))
	assert.Equal(t, "100000000000", BigintStringFromDecimalStringEx("10", 10, ".", true))
	assert.Equal(t, "110000000000", BigintStringFromDecimalStringEx("11", 10, ".", true))
	assert.Equal(t, "11900000000", BigintStringFromDecimalStringEx("1.19", 10, ".", true))
	assert.Equal(t, "10190000000", BigintStringFromDecimalStringEx("1.019", 10, ".", true))
	assert.Equal(t, "10190000000", BigintStringFromDecimalStringEx("000000001.019", 10, ".", true))
	assert.Equal(t, "10190000000", BigintStringFromDecimalStringEx("000000001.01900000000000000000", 10, ".", true))
	assert.Equal(t, "10190000000", BigintStringFromDecimalStringEx("000000001.019000000000000000001", 10, ".", false))
}

func TestBigintStringFromDecimalString(t *testing.T) {
	assert.Equal(t, "1", BigintStringFromDecimalString("0.0000000001", 10))
	assert.Equal(t, "90000", BigintStringFromDecimalString("0.000009", 10))
	assert.Equal(t, "100000000", BigintStringFromDecimalString("0.01", 10))
	assert.Equal(t, "1000000000", BigintStringFromDecimalString("0.1", 10))
	assert.Equal(t, "0", BigintStringFromDecimalString("0", 10))
	assert.Equal(t, "0", BigintStringFromDecimalString("0.0", 10))
	assert.NotEqual(t, "0", BigintStringFromDecimalString("0,0", 10))
	assert.Equal(t, "10000000000", BigintStringFromDecimalString("1", 10))
	assert.Equal(t, "11000000000", BigintStringFromDecimalString("1.1", 10))
	assert.Equal(t, "100000000000", BigintStringFromDecimalString("10", 10))
	assert.Equal(t, "110000000000", BigintStringFromDecimalString("11", 10))
	assert.Equal(t, "11900000000", BigintStringFromDecimalString("1.19", 10))
	assert.Equal(t, "10190000000", BigintStringFromDecimalString("1.019", 10))
}
