package stego

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"image"
)

func EmbedData(key string, plainText []byte, img *image.RGBA) error {
	rect := img.Bounds()
	salt := fmt.Sprintf("%d", rect.Max.X*rect.Max.Y)
	derivedKey := deriveKey(key, salt)
	random := predictableRandom(derivedKey)

	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return err
	}

	iv := make([]byte, 12)
	_, err = random.Read(iv)
	if err != nil {
		return err
	}

	additionalData := make([]byte, 64)
	_, err = random.Read(additionalData)
	if err != nil {
		return err
	}

	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	cipherText := aesGcm.Seal(nil, iv, plainText, additionalData)

	instructionFactory := NewInstructionFactory(random, rect.Max.X, rect.Max.Y)

	// reserve 3 bytes worth of instructions
	lengthInstructions := make([]*Instruction, 9*8)
	for index := 0; index < 9*8; index++ {
		lengthInstructions[index] = instructionFactory.New()
	}

	// encode cipher text
	var lastInstruction *Instruction
	bits := bytesToBoolArray(cipherText, 8)
	for _, bit := range bits {
		lastInstruction = instructionFactory.New()
		lastInstruction.Write(bit, img)
	}

	// turn last generated instruction into message that fits into lengthInstructions
	lengthBits := intArrayToBoolArray([]int{
		lastInstruction.X,
		lastInstruction.Y,
		lastInstruction.Channel,
	}, 24)

	// encode lengthBits
	for index, instruction := range lengthInstructions {
		instruction.Write(lengthBits[index], img)
	}

	return nil
}

func ReadEmbededData(key string, img *image.RGBA) ([]byte, error) {
	rect := img.Bounds()
	salt := fmt.Sprintf("%d", rect.Max.X*rect.Max.Y)
	derivedKey := deriveKey(key, salt)
	random := predictableRandom(derivedKey)

	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, 12)
	_, err = random.Read(iv)
	if err != nil {
		return nil, err
	}

	additionalData := make([]byte, 64)
	_, err = random.Read(additionalData)
	if err != nil {
		return nil, err
	}

	instructionFactory := NewInstructionFactory(random, rect.Max.X, rect.Max.Y)

	lengthBits := make([]bool, 9*8)
	for index := 0; index < 9*8; index++ {
		instruction := instructionFactory.New()
		lengthBits[index] = instruction.Read(img)
	}
	lengthBytes := boolArrayToIntArray(lengthBits, 24)
	lastInstruction := &Instruction{
		X:       lengthBytes[0],
		Y:       lengthBytes[1],
		Channel: lengthBytes[2],
	}

	instructions := []*Instruction{}
	for {
		instruction := instructionFactory.New()
		instructions = append(instructions, instruction)
		if lastInstruction.X == instruction.X && lastInstruction.Y == instruction.Y && lastInstruction.Channel == instruction.Channel {
			break
		}
	}

	bits := make([]bool, len(instructions))
	for index := len(instructions) - 1; index != -1; index-- {
		instruction := instructions[index]
		bits[index] = instruction.Read(img)
	}

	cipherText := boolArrayToBytes(bits, 8)

	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plainText, err := aesGcm.Open(nil, iv, cipherText, additionalData)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}
