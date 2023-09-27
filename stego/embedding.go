package stego

import (
	"fmt"
	"image"
)

func EmbedData(key string, plainText []byte, img *image.RGBA) error {
	rect := img.Bounds()
	salt := fmt.Sprintf("%d", rect.Max.X*rect.Max.Y)
	derivedKey := deriveKey(key, salt)
	random := predictableRandom(derivedKey)

	iv := make([]byte, 12)
	_, err := random.Read(iv)
	if err != nil {
		return err
	}

	additionalData := make([]byte, 64)
	_, err = random.Read(additionalData)
	if err != nil {
		return err
	}

	cipherText, err := encrypt(derivedKey, iv, additionalData, plainText)
	if err != nil {
		return err
	}

	instructionFactory := NewInstructionFactory(random, rect.Max.X, rect.Max.Y)

	// 3 bytes for X, 3 bytes for Y
	lastXYInstructions := make([]*Instruction, 4*8)
	for index := 0; index < len(lastXYInstructions); index++ {
		lastXYInstructions[index] = instructionFactory.New()
	}

	// 2 bits for channel
	lastChannelInstructions := make([]*Instruction, 2)
	for index := 0; index < len(lastChannelInstructions); index++ {
		lastChannelInstructions[index] = instructionFactory.New()
	}

	// encode cipher text
	var lastInstruction *Instruction
	bits := toBoolArray[byte](cipherText, 8)
	for _, bit := range bits {
		lastInstruction = instructionFactory.New()
		lastInstruction.Write(bit, img)
	}

	// turn last generated instruction into message that fits into lengthInstructions
	lastXYBits := toBoolArray[int]([]int{
		lastInstruction.X,
		lastInstruction.Y,
	}, 16)

	lastChannelBits := toBoolArray[int]([]int{
		lastInstruction.Channel,
	}, 2)

	// write
	for index, instruction := range lastXYInstructions {
		instruction.Write(lastXYBits[index], img)
	}

	for index, instruction := range lastChannelInstructions {
		instruction.Write(lastChannelBits[index], img)
	}

	for index := 0; index < instructionFactory.remaining; index++ {
		bit := random.Intn(2)
		instructionFactory.New().Write(bit == 1, img)
	}

	return nil
}

func ReadEmbededData(key string, img *image.RGBA) ([]byte, error) {
	rect := img.Bounds()
	salt := fmt.Sprintf("%d", rect.Max.X*rect.Max.Y)
	derivedKey := deriveKey(key, salt)
	random := predictableRandom(derivedKey)

	iv := make([]byte, 12)
	_, err := random.Read(iv)
	if err != nil {
		return nil, err
	}

	additionalData := make([]byte, 64)
	_, err = random.Read(additionalData)
	if err != nil {
		return nil, err
	}

	instructionFactory := NewInstructionFactory(random, rect.Max.X, rect.Max.Y)

	lastXYBits := make([]bool, 4*8)
	for index := 0; index < len(lastXYBits); index++ {
		lastXYBits[index] = instructionFactory.New().Read(img)
	}

	lastChannelBits := make([]bool, 2)
	for index := 0; index < len(lastChannelBits); index++ {
		lastChannelBits[index] = instructionFactory.New().Read(img)
	}

	lastXYBytes := fromBoolArray[int](lastXYBits, 16)
	lastChannelBytes := fromBoolArray[int](lastChannelBits, 2)

	lastInstruction := &Instruction{
		X:       lastXYBytes[0],
		Y:       lastXYBytes[1],
		Channel: lastChannelBytes[0],
	}

	instructions := []*Instruction{}
	for {
		instruction := instructionFactory.New()
		instructions = append(instructions, instruction)
		if lastInstruction.X == instruction.X &&
			lastInstruction.Y == instruction.Y &&
			lastInstruction.Channel == instruction.Channel {
			break
		}
	}

	bits := make([]bool, len(instructions))
	for index := len(bits) - 1; index != -1; index-- {
		instruction := instructions[index]
		bits[index] = instruction.Read(img)
	}

	cipherText := fromBoolArray[byte](bits, 8)
	plainText, err := decrypt(derivedKey, iv, additionalData, cipherText)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}
