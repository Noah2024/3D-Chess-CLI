package bitutil

import (
	"3DC/config"
	"3DC/util/logger"
	"fmt"
)

const (
	BoardSize = config.BoardSize
	LayerSize = config.LayerSize
	LineSize  = config.LineSize
	SpaceSize = config.SpaceSize
)

func VecToUint(x, y, z int) uint32 {
	return uint32(x + (y-1)*int(LayerSize) + (z-1)*int(LineSize))
}

// Decodes uint32 position into integer x,y,z position
func UintToVec(space uint32) (int, int, int) {
	if space < 1 || space > 512 {
		logger.Error(fmt.Sprintf("uint32 %d out of range for board size %d ", space, BoardSize))
		panic("See above error")
	}

	// index = x + y*8 + z*64 essentially decoding this
	//Step by step removing the largest term at a time
	space-- // convert to 0-based
	y := space / LayerSize
	space %= LayerSize
	z := space / LineSize
	x := space % LineSize

	return int(x + 1), int(y + 1), int(z + 1)
}
