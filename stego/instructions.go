package stego

import (
	"fmt"
	"image"
	"math/rand"
)

type InstructionFactory struct {
	random *rand.Rand
	maxX   int
	maxY   int
	used   map[string]bool
}

func NewInstructionFactory(random *rand.Rand, maxX int, maxY int) *InstructionFactory {
	return &InstructionFactory{
		random: random,
		maxX:   maxX,
		maxY:   maxY,
		used:   map[string]bool{},
	}
}

func (i *InstructionFactory) New() *Instruction {
	x := i.random.Intn(i.maxX)
	y := i.random.Intn(i.maxY)
	channel := i.random.Intn(3)
	uniqueKey := fmt.Sprintf("%d,%d,%d", x, y, channel)
	if i.used[uniqueKey] {
		for {
			x = i.random.Intn(i.maxX)
			y = i.random.Intn(i.maxY)
			channel = i.random.Intn(3)
			newUniqueKey := fmt.Sprintf("%d,%d,%d", x, y, channel)
			if !i.used[newUniqueKey] {
				i.used[newUniqueKey] = true
				break
			}
		}
	} else {
		i.used[uniqueKey] = true
	}
	return &Instruction{X: x, Y: y, Channel: channel}
}

type Instruction struct {
	X       int
	Y       int
	Channel int
}

func (i *Instruction) Read(img *image.RGBA) bool {
	rgba := img.RGBAAt(i.X, i.Y)

	var channel uint8

	switch i.Channel {
	case 0: // red
		channel = rgba.R
	case 1: // green
		channel = rgba.G
	case 2: // blue
		channel = rgba.B
	}

	return channel%2 == 0
}

func (i *Instruction) Write(bit bool, img *image.RGBA) {
	rgba := img.RGBAAt(i.X, i.Y)

	switch i.Channel {
	case 0: // red
		if bit { // bit is "1" so needs to be even
			if rgba.R%2 != 0 {
				if rgba.R == 255 { // isnt even
					rgba.R -= 1
				} else {
					rgba.R += 1
				}
			}
		} else { // bit is "0" so needs to be odd
			if rgba.R%2 == 0 { // is even
				rgba.R += 1
			}
		}
	case 1: // green
		if bit { // bit is "1" so needs to be even
			if rgba.G%2 != 0 { // isnt even
				if rgba.G == 255 {
					rgba.G -= 1
				} else {
					rgba.G += 1
				}
			}
		} else { // bit is "0" so needs to be odd
			if rgba.G%2 == 0 { // is even
				rgba.G += 1
			}
		}
	case 2: // blue
		if bit { // bit is "1" so needs to be even
			if rgba.B%2 != 0 { // isnt even
				if rgba.B == 255 {
					rgba.B -= 1
				} else {
					rgba.B += 1
				}
			}
		} else { // bit is "0" so needs to be odd
			if rgba.B%2 == 0 { // is even
				rgba.B += 1
			}
		}
	}

	img.SetRGBA(i.X, i.Y, rgba)
}
