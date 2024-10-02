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

type DecimalsTestData struct {
	stringWithDecimalsRef string
	stringWithDecimals    string
	stringNoDecimalsRef   string
	stringNoDecimals      string
	decimals              int
	separator             string
	panicIfRoundingNeeded bool
}

var decimalsTestData []DecimalsTestData = []DecimalsTestData{
	{"0", "0", "0", "000000000", 10, ".", true},
	{"0", "-0", "0", "-000000000", 10, ".", true},
	{"0", "0.0000000000", "0", "0000000000", 10, ".", true},
	{"0", "-0.0000000000", "0", "-0000000000", 10, ".", true},
	{"0", "0", "0", "00000000000", 0, ".", true},
	{"0", "-0", "0", "-00000000000", 0, ".", true},
	{"0", "0", "0", "00000000000", 10, ".", true},
	{"0", "-0", "0", "-00000000000", 10, ".", true},
	{"0", "0.0", "0", "00000000000", 10, ".", true},
	{"0", "-0.0", "0", "-00000000000", 10, ".", true},
	{"0", "0,0", "0", "00000000000", 10, ",", true},
	{"0", "-0,0", "0", "-00000000000", 10, ",", true},

	{"0.0000000001", "00.00000000010", "1", "1", 10, ".", true},
	{"-0.0000000001", "-00.00000000010", "-1", "-01", 10, ".", true},
	{"0.0000000001", "0.00000000011", "1", "1", 10, ".", false},
	{"-0.0000000001", "-00000.00000000011", "-1", "-00001", 10, ".", false},
	{"0.000009", "0.000009", "90000", "90000", 10, ".", true},
	{"-0.000009", "-0.000009", "-90000", "-90000", 10, ".", true},
	{"0.01", "0.01", "100000000", "100000000", 10, ".", true},
	{"-0.01", "-0.01", "-100000000", "-100000000", 10, ".", true},
	{"0.1", "0.10000000000", "1000000000", "1000000000", 10, ".", true},
	{"-0.1", "-0.1", "-1000000000", "-1000000000", 10, ".", true},
	{"0.1", "0.1", "10", "10", 2, ".", true},
	{"-0.1", "-0.1", "-10", "-10", 2, ".", true},

	{"1", "1", "10", "10", 1, ".", true},
	{"-1", "-1", "-10", "-10", 1, ".", true},
	{"1", "1", "1", "1", 0, ".", true},
	{"-1", "-1", "-1", "-1", 0, ".", true},
	{"1", "1", "10000000000", "010000000000", 10, ".", true},
	{"-1", "-01", "-10000000000", "-010000000000", 10, ".", true},
	{"1.1", "1.1", "11000000000", "11000000000", 10, ".", true},
	{"-1.1", "-1.1", "-11000000000", "-11000000000", 10, ".", true},
	{"1.19", "1.19", "11900000000", "11900000000", 10, ".", true},
	{"-1.19", "-1.19", "-11900000000", "-11900000000", 10, ".", true},
	{"1.019", "000000001.019", "10190000000", "10190000000", 10, ".", true},
	{"1.019", "000000001.01900000000000000000", "10190000000", "10190000000", 10, ".", true},
	{"1.019", "000000001.019000000000000000001", "10190000000", "10190000000", 10, ".", false},

	{"10", "10", "100000000000", "100000000000", 10, ".", true},
	{"-10", "-10", "-100000000000", "-100000000000", 10, ".", true},
	{"11", "11", "110000000000", "110000000000", 10, ".", true},
	{"-11", "-11", "-110000000000", "-110000000000", 10, ".", true},

	{"9999", "09999", "9999", "09999", 0, ".", true},
	{"-9999", "-009999", "-9999", "-009999", 0, ".", true},
}

func TestConvertDecimalsEx(t *testing.T) {
	for _, d := range decimalsTestData {
		assert.Equal(t, d.stringWithDecimalsRef, ConvertDecimalsEx(d.stringNoDecimals, d.decimals, d.separator))
	}
}

func TestConvertDecimalsBackEx(t *testing.T) {
	assert.Panics(t, func() { ConvertDecimalsBackEx("0.0000000001", 0, ".", false) })

	for _, d := range decimalsTestData {
		assert.Equal(t, d.stringNoDecimalsRef, ConvertDecimalsBackEx(d.stringWithDecimals, d.decimals, d.separator, d.panicIfRoundingNeeded))
	}
}

func TestConvertDecimalsBack(t *testing.T) {
	for _, d := range decimalsTestData {
		if d.separator != "," && d.panicIfRoundingNeeded != false {
			assert.Equal(t, d.stringNoDecimalsRef, ConvertDecimalsBack(d.stringWithDecimals, d.decimals).String())
		}
	}
}
