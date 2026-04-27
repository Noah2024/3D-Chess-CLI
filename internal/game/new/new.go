// Hardcoded creation of a default board game state
package new

import (
	"3DC/util/logger"
	"fmt"

	"3DC/util/bitutil"

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

// Generates new board with a hardcoded start state as a series of len 512 bitmaps
func NewCommand() {
	logger.Warn("Running New Command Now")
	var whitePawn bitmap.Bitmap

	// Debugging
	// vec := VecToUint(1, 3, 7)
	// logger.Debug(fmt.Sprintf("%d", vec))
	// x, y, z := UintToVec(vec)
	// logger.Debug(fmt.Sprintf("x: %d, y: %d, z:%d", x, y, z))

	whitePawn.Grow(BoardSize)
	logger.Warn(fmt.Sprintf("Presave len: %d", len(whitePawn)))

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
		logger.Debug(fmt.Sprintf("%d", VecToUint(i, 3, 7)))
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
	fullMap := map[string]bitmap.Bitmap{
		"white_pawn":   whitePawn,
		"white_knight": whiteKnight,
		"white_bishop": whiteBishop,
		"white_rook":   whiteRook,
		"white_queen":  whiteQueen,
		"white_king":   whiteKing,
		"black_pawn":   blackPawn,
		"black_knight": blackKnight,
		"black_bishop": blackBishop,
		"black_rook":   blackRook,
		"black_queen":  blackQueen,
		"black_king":   blackKing,
	}
	bitutil.SaveGame(fullMap, "data/output")
	out, _ := bitutil.LoadGame("data/output")
	for _, V := range out {
		fmt.Printf("bit map for piece index %v", V)
	}

	fmt.Print(out)
	fmt.Print("Hello")
}
