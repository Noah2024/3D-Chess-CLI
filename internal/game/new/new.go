// Hardcoded creation of a default board game state
package new

import (
	"3DC/util/logger"
	"fmt"

	"github.com/kelindar/bitmap"
)

// Defining size and shape of board
// Stored in Uints right now to make Uint -> Vec easier
// BUT it may be benificial later to store them as ints
// And to make Vec -> Uint easier
const BoardSize uint32 = 512
const LayerSize uint32 = 64
const LineSize uint32 = 8
const SpaceSize uint32 = 1

// Encodes integer x,y,z positions into uint32 encoded position
func VecToUint(x, y, z int) uint32 {
	return uint32(x + (y-1)*int(LineSize) + (z-1)*int(LayerSize))
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
	z := space / LayerSize
	space %= LayerSize
	y := space / LineSize
	x := space % LineSize

	return int(x + 1), int(y + 1), int(z + 1)
}

// Generates new board witha hardcoded start state as a series of len 512 bitmaps
func NewCommand() {
	logger.Warn("Running New Command Now")
	var whitePawn bitmap.Bitmap
	whitePawn.Grow(BoardSize)
	for i := 1; i <= 8; i++ {
		whitePawn.Set(VecToUint(i, 3, 2))
	}

	var whiteKnight bitmap.Bitmap
	whiteKnight.Grow(BoardSize)
	whiteKnight.Set(VecToUint(2, 3, 1))
	whiteKnight.Set(VecToUint(7, 3, 1))

	var whiteBishop bitmap.Bitmap
	whiteBishop.Grow(BoardSize)
	whiteBishop.Set(VecToUint(3, 3, 1))
	whiteBishop.Set(VecToUint(6, 3, 1))

	var whiteRook bitmap.Bitmap
	whiteRook.Grow(BoardSize)
	whiteRook.Set(VecToUint(1, 3, 1))
	whiteRook.Set(VecToUint(8, 3, 1))

	var whiteQueen bitmap.Bitmap
	whiteQueen.Grow(BoardSize)
	whiteQueen.Set(VecToUint(4, 3, 1))

	var whiteKing bitmap.Bitmap
	whiteKing.Grow(BoardSize)
	whiteKing.Set(VecToUint(5, 3, 1))

	//=============================
	// Defining Black Pieces Bitmaps
	//=============================

	var blackPawn bitmap.Bitmap
	blackPawn.Grow(BoardSize)
	for i := 1; i <= 8; i++ {
		blackPawn.Set(VecToUint(i, 3, 7))
	}

	var blackKnight bitmap.Bitmap
	blackKnight.Grow(BoardSize)
	blackKnight.Set(VecToUint(2, 3, 8))
	blackKnight.Set(VecToUint(7, 3, 8))

	var blackBishop bitmap.Bitmap
	blackBishop.Grow(BoardSize)
	blackBishop.Set(VecToUint(3, 3, 8))
	blackBishop.Set(VecToUint(6, 3, 8))

	var blackRook bitmap.Bitmap
	blackRook.Grow(BoardSize)
	blackRook.Set(VecToUint(1, 3, 8))
	blackRook.Set(VecToUint(8, 3, 8))

	var blackQueen bitmap.Bitmap
	blackQueen.Grow(BoardSize)
	blackQueen.Set(VecToUint(4, 3, 8))

	var blackKing bitmap.Bitmap
	blackKing.Grow(BoardSize)
	blackQueen.Set(VecToUint(5, 3, 8))

	logger.Debug("Init new game setup")
}
