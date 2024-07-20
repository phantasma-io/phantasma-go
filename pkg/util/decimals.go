package util

import (
	"strings"
	"unicode"
)

// Internal utils

func isInteger(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func trimWholePrefix(s string, prefix string) string {
	for strings.HasPrefix(s, prefix) {
		s = s[len(prefix):]
	}
	return s
}

func trimWholeSuffix(s string, suffix string) string {
	for strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func stringIsZeroOrEmptyBigint(number string) bool {
	for _, c := range number {
		if c != '0' {
			return false
		}
	}

	return true
}

// Public utils
// BigintStringToDecimalString converts big integer number to decimal number, both serialized as a string.
func BigintStringToDecimalString(number string, decimals int) string {
	if stringIsZeroOrEmptyBigint(number) {
		return "0"
	}

	if decimals == 0 {
		return number
	}

	if len(number) <= decimals {
		return "0." + strings.Repeat("0", decimals-len(number)) + trimWholeSuffix(number, "0")
	}

	integerPart := number[:len(number)-decimals]
	if integerPart == "" {
		integerPart = "0"
	}

	fractionalPart := number[len(number)-decimals:]

	if stringIsZeroOrEmptyBigint(fractionalPart) {
		return integerPart
	}
	return integerPart + "." + trimWholeSuffix(fractionalPart, "0")
}

// BigintStringToDecimalStringEx converts decimal number to big integer number, both serialized as a string.
func BigintStringFromDecimalStringEx(number string, decimals int, separator string, panicIfRoundingNeeded bool) string {
	if stringIsZeroOrEmptyBigint(number) {
		return "0"
	}

	if !strings.Contains(number, separator) {
		// No fractional part found, we need to put zeroes instead
		return number + strings.Repeat("0", decimals)
	}

	split := strings.SplitN(number, separator, 2)

	integerPart := split[0]
	fractionalPart := split[1]

	if decimals == 0 {
		// Nothing to do, only to check if passed number is correct
		if !isInteger(number) && !stringIsZeroOrEmptyBigint(fractionalPart) {
			panic("Non-integer number passed but decimals set to 0")
		}
		return number
	}

	if len(fractionalPart) < decimals {
		// We need to add more zeroes to fractional part
		fractionalPart = fractionalPart + strings.Repeat("0", decimals-len(fractionalPart))
	} else if len(fractionalPart) > decimals {
		if stringIsZeroOrEmptyBigint(fractionalPart[decimals:]) {
			// We can safely drop zeroes
			fractionalPart = fractionalPart[:decimals]
		} else {
			// Fractional part is too big and requires rounding
			if panicIfRoundingNeeded {
				panic("Fractional part is too big and requires rounding")
			} else {
				// We have to round
				fractionalPart = fractionalPart[:decimals]
			}
		}
	}

	result := trimWholePrefix(integerPart+fractionalPart, "0")
	if result == "" { // We killed all zeroes, we need one
		result = "0"
	}

	return result
}

// BigintStringToDecimalString converts decimal number to big integer number, both serialized as a string.
func BigintStringFromDecimalString(number string, decimals int) string {
	return BigintStringFromDecimalStringEx(number, decimals, ".", true)
}
