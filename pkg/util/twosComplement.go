package util

// https://en.wikipedia.org/wiki/Two%27s_complement

func TwosComplementConvertTo(bytes []byte) {
	flipBits(bytes)
	addOneIgnoringOverflow(bytes)
}

func TwosComplementConvertFrom(bytes []byte) {
	subtractOneIgnoringOverflow(bytes)
	flipBits(bytes)
}

func addOneIgnoringOverflow(bytes []byte) {
	for i := len(bytes) - 1; i >= 0; i-- {
		if bytes[i] < 255 {
			bytes[i] += 1
			return
		}
		bytes[i] = 0
		if i == 0 {
			// Overflow
		}
	}
}

func subtractOneIgnoringOverflow(bytes []byte) {
	for i := len(bytes) - 1; i >= 0; i-- {
		if bytes[i] > 0 {
			bytes[i] -= 1
			return
		}
		bytes[i] = 255
		if i == 0 {
			// Overflow
		}
	}
}

func flipBits(bytes []byte) {
	for i := 0; i < len(bytes); i++ {
		bytes[i] = ^bytes[i]
	}
}
