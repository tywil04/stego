package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"stego/stego"
)

// func predictableRand(key []byte) *rand.Rand {
// 	hasher := fnv.New64()
// 	hasher.Write(key)
// 	return rand.New(rand.NewSource(int64(hasher.Sum64())))
// }

// func deriveKey(key string, salt string) []byte {
// 	const (
// 		time    = 2
// 		memory  = 64 * 1024
// 		threads = 4
// 		keyLen  = 32 // for aes 256 bit
// 	)
// 	return argon2.IDKey([]byte(key), []byte(salt), time, memory, threads, keyLen)
// }

// func boolToInt(b bool) int {
// 	if b {
// 		return 1
// 	} else {
// 		return 0
// 	}
// }

func readImage(filePath string) (*image.RGBA, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	rgba, ok := img.(*image.RGBA)
	if !ok {
		b := img.Bounds()
		rgba = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(rgba, rgba.Bounds(), img, b.Min, draw.Src)
	}

	return rgba, nil
}

func writeImage(filePath string, img *image.RGBA) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return err
	}

	return nil
}

// func embedData(key string, plainText []byte, img *image.RGBA) error {
// 	rect := img.Bounds()
// 	salt := fmt.Sprintf("%d", rect.Max.X*rect.Max.Y)
// 	derivedKey := deriveKey(key, salt)
// 	random := predictableRand(derivedKey)

// 	block, err := aes.NewCipher(derivedKey)
// 	if err != nil {
// 		return err
// 	}

// 	iv := make([]byte, 12)
// 	_, err = random.Read(iv)
// 	if err != nil {
// 		return err
// 	}

// 	additionalData := make([]byte, 64)
// 	_, err = random.Read(additionalData)
// 	if err != nil {
// 		return err
// 	}

// 	aesGcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		return err
// 	}

// 	cipherText := aesGcm.Seal(nil, iv, plainText, additionalData)

// 	used := map[string]bool{}
// 	bits := make([]bool, len(cipherText)*8)
// 	counter := 0
// 	for _, part := range cipherText {
// 		value := int(part)
// 		binaryCounter := 128
// 		for index := 0; index < 8; index++ {
// 			if value-binaryCounter >= 0 {
// 				bits[counter+index] = true
// 				value -= binaryCounter
// 			}
// 			binaryCounter = binaryCounter / 2
// 		}
// 		counter += 8
// 	}

// 	for _, bit := range bits {
// 		x := random.Intn(rect.Max.X)
// 		y := random.Intn(rect.Max.Y)
// 		channel := random.Intn(3)
// 		uniqueKey := fmt.Sprintf("%d,%d,%d", x, y, channel)
// 		if used[uniqueKey] {
// 			for {
// 				x = random.Intn(rect.Max.X)
// 				y = random.Intn(rect.Max.Y)
// 				channel = random.Intn(3)
// 				newUniqueKey := fmt.Sprintf("%d,%d,%d", x, y, channel)
// 				if !used[newUniqueKey] {
// 					used[newUniqueKey] = true
// 					break
// 				}
// 			}
// 		} else {
// 			used[uniqueKey] = true
// 		}

// 		rgba := img.RGBAAt(x, y)

// 		switch channel {
// 		case 0: // red
// 			if bit { // bit is "1" so needs to be even
// 				if rgba.R%2 != 0 {
// 					if rgba.R == 255 { // isnt even
// 						rgba.R -= 1
// 					} else {
// 						rgba.R += 1
// 					}
// 				}
// 			} else { // bit is "0" so needs to be odd
// 				if rgba.R%2 == 0 { // is even
// 					rgba.R += 1
// 				}
// 			}
// 		case 1: // green
// 			if bit { // bit is "1" so needs to be even
// 				if rgba.G%2 != 0 { // isnt even
// 					if rgba.G == 255 {
// 						rgba.G -= 1
// 					} else {
// 						rgba.G += 1
// 					}
// 				}
// 			} else { // bit is "0" so needs to be odd
// 				if rgba.G%2 == 0 { // is even
// 					rgba.G += 1
// 				}
// 			}
// 		case 2: // blue
// 			if bit { // bit is "1" so needs to be even
// 				if rgba.B%2 != 0 { // isnt even
// 					if rgba.B == 255 {
// 						rgba.B -= 1
// 					} else {
// 						rgba.B += 1
// 					}
// 				}
// 			} else { // bit is "0" so needs to be odd
// 				if rgba.B%2 == 0 { // is even
// 					rgba.B += 1
// 				}
// 			}
// 		}

// 		img.SetRGBA(x, y, rgba)
// 	}

// 	return nil
// }

// func readEmbededData(key string, img *image.RGBA) ([]byte, error) {
// 	length := 21

// 	rect := img.Bounds()
// 	salt := fmt.Sprintf("%d", rect.Max.X*rect.Max.Y)
// 	derivedKey := deriveKey(key, salt)
// 	random := predictableRand(derivedKey)

// 	block, err := aes.NewCipher(derivedKey)
// 	if err != nil {
// 		return nil, err
// 	}

// 	iv := make([]byte, 12)
// 	_, err = random.Read(iv)
// 	if err != nil {
// 		return nil, err
// 	}

// 	additionalData := make([]byte, 64)
// 	_, err = random.Read(additionalData)
// 	if err != nil {
// 		return nil, err
// 	}

// 	used := map[string]bool{}
// 	instructions := make([][]int, length*8)
// 	for index := 0; index < length*8; index++ {
// 		x := random.Intn(rect.Max.X)
// 		y := random.Intn(rect.Max.Y)
// 		channel := random.Intn(3)
// 		uniqueKey := fmt.Sprintf("%d,%d,%d", x, y, channel)
// 		if used[uniqueKey] {
// 			for {
// 				x = random.Intn(rect.Max.X)
// 				y = random.Intn(rect.Max.Y)
// 				channel = random.Intn(3)
// 				newUniqueKey := fmt.Sprintf("%d,%d,%d", x, y, channel)
// 				if !used[newUniqueKey] {
// 					used[newUniqueKey] = true
// 					break
// 				}
// 			}
// 		} else {
// 			used[uniqueKey] = true
// 		}

// 		instructions[index] = []int{x, y, channel}
// 	}

// 	bits := make([]bool, length*8)
// 	for index := len(instructions) - 1; index != -1; index-- {
// 		instruction := instructions[index]
// 		rgba := img.RGBAAt(instruction[0], instruction[1])

// 		switch instruction[2] {
// 		case 0: // red
// 			bits[index] = rgba.R%2 == 0
// 		case 1: // green
// 			bits[index] = rgba.G%2 == 0
// 		case 2: // blue
// 			bits[index] = rgba.B%2 == 0
// 		}
// 	}

// 	cipherText := make([]byte, length)
// 	for index := 0; index < len(bits); index += 8 {
// 		sum := (boolToInt(bits[index]) * 128) +
// 			(boolToInt(bits[index+1]) * 64) +
// 			(boolToInt(bits[index+2]) * 32) +
// 			(boolToInt(bits[index+3]) * 16) +
// 			(boolToInt(bits[index+4]) * 8) +
// 			(boolToInt(bits[index+5]) * 4) +
// 			(boolToInt(bits[index+6]) * 2) +
// 			(boolToInt(bits[index+7]) * 1)
// 		cipherText[index/8] = byte(sum)
// 	}

// 	aesGcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		return nil, err
// 	}

// 	plainText, err := aesGcm.Open(nil, iv, cipherText, additionalData)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return plainText, nil
// }

func main() {
	img, _ := readImage("./test.png")

	err := stego.EmbedData("supersecretkey", []byte("you have found me! look at my amazing skills"), img)
	if err != nil {
		log.Fatal(err)
	}

	writeImage("./output.png", img)

	fmt.Println("written")

	result, err := stego.ReadEmbededData("supersecretkey", img)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(result))
}
