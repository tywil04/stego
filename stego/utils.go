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

func bytesToBoolArray(data []byte, bitLength int) []bool {
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

func intArrayToBoolArray(data []int, bitLength int) []bool {
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

func boolArrayToBytes(bits []bool, bitLength int) []byte {
	result := make([]byte, len(bits)/bitLength)
	for index := 0; index < len(bits); index += bitLength {
		sum := 0
		max := int(math.Pow(2, float64(bitLength-1)))
		for innerIndex := 0; innerIndex < bitLength; innerIndex++ {
			sum += boolToInt(bits[index+innerIndex]) * max
			max = max / 2
		}
		result[index/bitLength] = byte(sum)
	}
	return result
}

func boolArrayToIntArray(bits []bool, bitLength int) []int {
	result := make([]int, len(bits)/bitLength)
	for index := 0; index < len(bits); index += bitLength {
		sum := 0
		max := int(math.Pow(2, float64(bitLength-1)))
		for innerIndex := 0; innerIndex < bitLength; innerIndex++ {
			sum += boolToInt(bits[index+innerIndex]) * max
			max = max / 2
		}
		result[index/bitLength] = sum
	}
	return result
}
