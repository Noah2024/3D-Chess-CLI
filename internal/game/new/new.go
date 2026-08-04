// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// Contains the command to create and save a new board state
// the intial boardstate is Hardcoded to ensure that 'game new' command always creates the same start state
package new

import (
	"3DC/internal/game/save"
	"3DC/util/logger"

	"3DC/config"
	"3DC/util/bitutil"

	"github.com/kelindar/bitmap"
)

const (
	BoardSize = config.BoardSize
	LayerSize = config.LayerSize
	LineSize  = config.LineSize
	SpaceSize = config.SpaceSize
)

var VecToUint = bitutil.VecToUint
var UintToVec = bitutil.UintToVec

// Generates new board with a hardcoded start state as a series of len 64 bitmaps
// DESTRUCTIVE: Overwrites any existing game in the CurrentGame folder without prompting for confirmation.
func NewCommand() {
	save.SaveGame(DefaultStartState(), config.CurrentGame)
	// save.SaveGame(GenerateSinglePiece(map[string]uint32{
	// 	"♟": 292,
	// 	"♗": 235,
	// }), config.CurrentGame)

}

// For now I will store possible game states to start from here, for NOW
// Until I can think of a better place to put them
//As a note all valid board states REQUIRE a king on either side
//This is a requiremnet for the checking validation system

// Default starting state for board
func DefaultStartState() map[string]bitmap.Bitmap {
	logger.Debug("Running New Command Now")
	var whitePawn bitmap.Bitmap

	whitePawn.Grow(BoardSize - 1)
	// logger.Debug(fmt.Sprintf("Presave len: %d", len(whitePawn)))

	for i := 0; i <= 7; i++ {
		whitePawn.Set(VecToUint(i, 2, 1))
	}

	var whiteKnight bitmap.Bitmap
	whiteKnight.Grow(BoardSize - 1)
	whiteKnight.Set(VecToUint(1, 2, 0))
	whiteKnight.Set(VecToUint(6, 2, 0))

	var whiteBishop bitmap.Bitmap
	whiteBishop.Grow(BoardSize - 1)
	whiteBishop.Set(VecToUint(2, 2, 0))
	whiteBishop.Set(VecToUint(5, 2, 0))

	var whiteRook bitmap.Bitmap
	whiteRook.Grow(BoardSize - 1)
	whiteRook.Set(VecToUint(0, 2, 0))
	whiteRook.Set(VecToUint(7, 2, 0))

	var whiteQueen bitmap.Bitmap
	whiteQueen.Grow(BoardSize - 1)
	whiteQueen.Set(VecToUint(3, 2, 0))

	var whiteKing bitmap.Bitmap
	whiteKing.Grow(BoardSize - 1)
	whiteKing.Set(VecToUint(4, 2, 0))

	//=============================
	// Defining Black Pieces Bitmaps
	//=============================

	var blackPawn bitmap.Bitmap
	blackPawn.Grow(BoardSize - 1)
	for i := 0; i <= 7; i++ {
		blackPawn.Set(VecToUint(i, 2, 6))
	}

	var blackKnight bitmap.Bitmap
	blackKnight.Grow(BoardSize - 1)
	blackKnight.Set(VecToUint(1, 2, 7))
	blackKnight.Set(VecToUint(6, 2, 7))

	var blackBishop bitmap.Bitmap
	blackBishop.Grow(BoardSize - 1)
	blackBishop.Set(VecToUint(2, 2, 7))
	blackBishop.Set(VecToUint(5, 2, 7))

	var blackRook bitmap.Bitmap
	blackRook.Grow(BoardSize - 1)
	blackRook.Set(VecToUint(0, 2, 7)) ///Why?
	blackRook.Set(VecToUint(7, 2, 7))

	var blackQueen bitmap.Bitmap
	blackQueen.Grow(BoardSize - 1)
	blackQueen.Set(VecToUint(4, 2, 7))

	var blackKing bitmap.Bitmap
	blackKing.Grow(BoardSize - 1)
	blackKing.Set(VecToUint(3, 2, 7))

	logger.Debug("Init new game setup")
	fullMap := map[string]bitmap.Bitmap{
		"♙": blackPawn,
		"♘": blackKnight,
		"♗": blackBishop,
		"♖": blackRook,
		"♕": blackQueen,
		"♔": blackKing,
		"♟": whitePawn,
		"♞": whiteKnight,
		"♝": whiteBishop,
		"♜": whiteRook,
		"♛": whiteQueen,
		"♚": whiteKing,
	}
	return fullMap
}

// Generates a board with pieces at the specified uint32 index locations
func GenerateSinglePiece(allPieces map[string]uint32) map[string]bitmap.Bitmap {
	//FullMap starts as a completly empty (but initalized) array
	fullMap := map[string]bitmap.Bitmap{
		"♙": bitmap.Bitmap{0, 0, 0, 0, 0, 0, 0, 0},
		"♘": bitmap.Bitmap{0, 0, 0, 0, 0, 0, 0, 0},
		"♗": bitmap.Bitmap{0, 0, 0, 0, 0, 0, 0, 0},
		"♖": bitmap.Bitmap{0, 0, 0, 0, 0, 0, 0, 0},
		"♕": bitmap.Bitmap{0, 0, 0, 0, 0, 0, 0, 0},
		"♔": bitmap.Bitmap{0, 0, 0, 0, 0, 0, 0, 0},
		"♟": bitmap.Bitmap{0, 0, 0, 0, 0, 0, 0, 0},
		"♞": bitmap.Bitmap{0, 0, 0, 0, 0, 0, 0, 0},
		"♝": bitmap.Bitmap{0, 0, 0, 0, 0, 0, 0, 0},
		"♜": bitmap.Bitmap{0, 0, 0, 0, 0, 0, 0, 0},
		"♛": bitmap.Bitmap{0, 0, 0, 0, 0, 0, 0, 0},
		"♚": bitmap.Bitmap{0, 0, 0, 0, 0, 0, 0, 0},
	}
	for piece, loc := range allPieces {
		tmp := fullMap[piece].Clone(nil)
		tmp.Set(loc)
		fullMap[piece] = tmp

	}
	// logger.Debug(fmt.Sprintf("Init custom inline board setup finsihed '%s' ", allPieces))

	return fullMap
}
