// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// Hardcoded creation of a default board game state
package new

import (
	"3DC/internal/game/save"
	"3DC/util/logger"
	"fmt"

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
func NewCommand() {

	save.SaveGame(GenerateSinglePiece("♖", 5, 5, 5), config.CurrentGame)
	// save.SaveGame(DefaultStartState(), config.CurrentGame)
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

	for i := 1; i <= 8; i++ {
		whitePawn.Set(VecToUint(i, 3, 2))
	}

	var whiteKnight bitmap.Bitmap
	whiteKnight.Grow(BoardSize - 1)
	whiteKnight.Set(VecToUint(2, 3, 1))
	whiteKnight.Set(VecToUint(7, 3, 1))

	var whiteBishop bitmap.Bitmap
	whiteBishop.Grow(BoardSize - 1)
	whiteBishop.Set(VecToUint(3, 3, 1))
	whiteBishop.Set(VecToUint(6, 3, 1))

	var whiteRook bitmap.Bitmap
	whiteRook.Grow(BoardSize - 1)
	whiteRook.Set(VecToUint(1, 3, 1))
	whiteRook.Set(VecToUint(8, 3, 1))

	var whiteQueen bitmap.Bitmap
	whiteQueen.Grow(BoardSize - 1)
	whiteQueen.Set(VecToUint(4, 3, 1))

	var whiteKing bitmap.Bitmap
	whiteKing.Grow(BoardSize - 1)
	whiteKing.Set(VecToUint(5, 3, 1))

	//=============================
	// Defining Black Pieces Bitmaps
	//=============================

	var blackPawn bitmap.Bitmap
	blackPawn.Grow(BoardSize - 1)
	for i := 1; i <= 8; i++ {
		blackPawn.Set(VecToUint(i, 3, 7))
	}

	var blackKnight bitmap.Bitmap
	blackKnight.Grow(BoardSize - 1)
	blackKnight.Set(VecToUint(2, 3, 8))
	blackKnight.Set(VecToUint(7, 3, 8))

	var blackBishop bitmap.Bitmap
	blackBishop.Grow(BoardSize - 1)
	blackBishop.Set(VecToUint(3, 3, 8))
	blackBishop.Set(VecToUint(6, 3, 8))

	var blackRook bitmap.Bitmap
	blackRook.Grow(BoardSize - 1)
	blackRook.Set(VecToUint(1, 3, 8)) ///Why?
	blackRook.Set(VecToUint(8, 3, 8))

	var blackQueen bitmap.Bitmap
	blackQueen.Grow(BoardSize - 1)
	blackQueen.Set(VecToUint(4, 3, 8))

	var blackKing bitmap.Bitmap
	blackKing.Grow(BoardSize - 1)
	blackKing.Set(VecToUint(5, 3, 8))

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

// Generates a board with a single piece at a given position.
// Used extensivly in testing the genMoves unit tests
func GenerateSinglePiece(vis string, x int, y int, z int) map[string]bitmap.Bitmap {
	var singlePiece bitmap.Bitmap
	singlePiece.Grow(config.BoardSize - 1)
	loc := bitutil.VecToUint(x, y, z)
	singlePiece.Set(loc)
	logger.Debug(fmt.Sprintf("Init single piece test board for '%s' setup", vis))
	fullMap := map[string]bitmap.Bitmap{
		"♙": bitmap.Bitmap{},
		"♘": bitmap.Bitmap{},
		"♗": bitmap.Bitmap{},
		"♖": singlePiece,
		"♕": bitmap.Bitmap{},
		"♔": bitmap.Bitmap{},
		"♟": bitmap.Bitmap{},
		"♞": bitmap.Bitmap{},
		"♝": bitmap.Bitmap{},
		"♜": bitmap.Bitmap{},
		"♛": bitmap.Bitmap{},
		"♚": bitmap.Bitmap{},
	}
	return fullMap
}

// A single pawn on either side of the board in its starting position
func LonePawn() map[string]bitmap.Bitmap {
	logger.Debug("Running New Command Now")
	var whitePawn bitmap.Bitmap

	whitePawn.Grow(BoardSize - 1)

	for i := 1; i <= 8; i++ {
		whitePawn.Set(VecToUint(i, 3, 2))
	}

	// var whiteKing bitmap.Bitmap
	// whiteKing.Grow(BoardSize - 1)
	// whiteKing.Set(VecToUint(5, 3, 1))

	//=============================
	// Defining Black Pieces Bitmaps
	//=============================

	var blackPawn bitmap.Bitmap
	blackPawn.Grow(BoardSize - 1)
	for i := 1; i <= 8; i++ {
		blackPawn.Set(VecToUint(i, 3, 7))
	}

	// var blackKing bitmap.Bitmap
	// blackKing.Grow(BoardSize - 1)
	// blackKing.Set(VecToUint(5, 3, 8))

	logger.Debug("Init new LonePawn setup")
	fullMap := map[string]bitmap.Bitmap{
		"♙": blackPawn,
		"♟": whitePawn,
	}
	return fullMap
}
