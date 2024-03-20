package stego

import (
	"math"
)

func boolToInt(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

func toBoolArray[D byte | int](data []D, bitLength int) []bool {
	bits := make([]bool, len(data)*bitLength)
	counter := 0
	for _, part := range data {
		value := int(part)
		binaryCounter := int(math.Pow(2, float64(bitLength-1)))
		for index := 0; index < bitLength; index++ {
			if value-binaryCounter >= 0 {
				bits[counter+index] = true
				value -= binaryCounter
			}
			binaryCounter = binaryCounter / 2
		}
		counter += bitLength
	}
	return bits
}

func fromBoolArray[D byte | int](bits []bool, bitLength int) []D {
	result := make([]D, len(bits)/bitLength)
	for index := 0; index < len(bits); index += bitLength {
		sum := 0
		max := int(math.Pow(2, float64(bitLength-1)))
		for innerIndex := 0; innerIndex < bitLength; innerIndex++ {
			sum += boolToInt(bits[index+innerIndex]) * max
			max = max / 2
		}
		result[index/bitLength] = D(sum)
	}
	return result
}

func makeEven(value uint8) uint8 {
	if value%2 != 0 { // is odd
		if value == 255 {
			// is max value, must remove one (vs adding one) to avoid overflow
			value -= 1
		} else {
			// max odd (that is handled here) is 253, adding one wont cause overflow
			// min odd is 1, adding one wont cause overflow
			value += 1
		}
	}
	return value
}

func makeOdd(value uint8) uint8 {
	if value%2 == 0 { // is even
		// max even is 254, adding one wont cause overflow
		// min even is 0, adding one wont cause overflow
		value += 1
	}
	return value
}
