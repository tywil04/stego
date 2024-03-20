package stego

import (
	"image"
	"math/rand"
)

type InstructionFactory struct {
	random    *rand.Rand
	maxX      int
	maxY      int
	remaining int
	used      map[[3]int]bool
}

func NewInstructionFactory(random *rand.Rand, maxX int, maxY int) *InstructionFactory {
	return &InstructionFactory{
		random:    random,
		maxX:      maxX,
		maxY:      maxY,
		remaining: maxX * maxY * 3, // Remaining number of instructions
		used:      map[[3]int]bool{},
	}
}

func (i *InstructionFactory) New() *Instruction {
	x := i.random.Intn(i.maxX)
	y := i.random.Intn(i.maxY)
	channel := i.random.Intn(3)
	uniqueKey := [3]int{x, y, channel}

	if i.used[uniqueKey] {
		for {
			x = i.random.Intn(i.maxX)
			y = i.random.Intn(i.maxY)
			channel = i.random.Intn(3)
			newUniqueKey := [3]int{x, y, channel}

			if !i.used[newUniqueKey] {
				i.used[newUniqueKey] = true
				break
			}
		}
	} else {
		i.used[uniqueKey] = true
	}

	i.remaining -= 1
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

	var bitFunc func(uint8) uint8
	if bit {
		bitFunc = makeEven
	} else {
		bitFunc = makeOdd
	}

	switch i.Channel {
	case 0: // red
		rgba.R = bitFunc(rgba.R)
	case 1: // green
		rgba.G = bitFunc(rgba.G)
	case 2: // blue
		rgba.B = bitFunc(rgba.B)
	}

	img.SetRGBA(i.X, i.Y, rgba)
}
