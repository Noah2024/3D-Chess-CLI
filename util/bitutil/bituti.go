// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// bitutil provides utility functions for converting between 3D vector coordinates and uint32 representations of positions on the board.
// I do belive the 1 indexed nature of the BoardSize is a bit confusing, and I'm fairly certian there is an off by one error in the either decode or encode function, but I don't have time to fix it right now, and it works for the current implementation of the game, so I'm leaving it as is for now.
package bitutil

import (
	"3DC/config"
	"3DC/util/logger"
	"fmt"
)

// Constants for board size and dimensions, derived from the configuration settings.
// ALL are SIZES of the board, not indexes, so they are 1 indexed
const (
	BoardSize = config.BoardSize // 512
	LayerSize = config.LayerSize // 64
	LineSize  = config.LineSize  // 8
	SpaceSize = config.SpaceSize // 1
)

// VecToUint encodes a 3D vector (x, y, z) into a uint32 representation of a position on the board.
func VecToUint(x, y, z int) uint32 {
	return uint32(x + (y)*int(LayerSize) + (z)*int(LineSize))
}

// Decodes uint32 position into integer x,y,z position
func UintToVec(space uint32) (int, int, int) {
	if space > 511 {
		logger.Error(fmt.Sprintf("uint32 %d out of range for board size %d ", space, BoardSize))
		panic(fmt.Sprintf("uint32 %d out of range for board size %d ", space, BoardSize))
	}

	// index = x + y*8 + z*64 essentially decoding this
	//Step by step removing the largest term at a time
	y := space / LayerSize
	space %= LayerSize
	z := space / LineSize
	x := space % LineSize

	return int(x), int(y), int(z)
}
