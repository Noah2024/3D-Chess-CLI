// Package move contins the logic for moving pieces. Using the planes from dataplanes.go to generate valid moves then validates them for broken rules
// This is the largest and most compelx file in the project, due to it being the core gameplay element
package move

import (
	"3DC/config"
	"3DC/internal/game/load"
	"3DC/internal/game/save"
	"3DC/util/bitutil"
	"3DC/util/dataplane"
	"3DC/util/logger"
	"3DC/util/metadata"

	"fmt"

	"github.com/kelindar/bitmap"
)

//To Do
// - Need to go through a better log regiment for all the moves
// - Need to build out a test suite to ensure that moves don't break with new updates

// BIG NOTE TO SELF 7/13/2026
// uintLoc is ONE indexed, meanwhile the bitmap  is ZERO indexed, hence the fuckywucky -1's everywhere
// So every use of Set, Contains, Remove, etc. needs to be -1 from the uintLoc value

// rookMove contains the bitwise operations necessary to generate all possible moves for a rook piece
// it takes x y and z integer cooridnates and outputs a size 511 bitmap all ones of which represent possible moves
// inputs: x, y, z int | outputs: bitmap.Bitmap
func rookMove(x int, y int, z int) bitmap.Bitmap { // Will parallelize with go rountine

	logger.Debug(fmt.Sprintf("Generating all possible rook moves from :x: %d, y: %d, z: %d", x, y, z))
	forward := dataplane.XPlane[x-1].Clone(nil) //-2
	forward.And(dataplane.ZPlane[z-1])

	sideToSide := dataplane.YPlane[y-1].Clone(nil)
	sideToSide.And(dataplane.XPlane[x-1]) //-2

	upAndDown := dataplane.YPlane[y-1].Clone(nil)
	upAndDown.And(dataplane.ZPlane[z-1])

	forward.Or(upAndDown)
	forward.Or(sideToSide)
	// fmt.Printf("Intersection %064b\n", forward)
	return forward
}

// moveMap matches a pieces visual representation to the function that generates all possible moves for that piece
// inputs: string | outputs: function(int, int, int) bitmap.Bitmap
var moveMap = map[string]func(int, int, int) bitmap.Bitmap{
	// "♙": blackPawn,
	// "♘": blackKnight,
	// "♗": blackBishop,
	"♖": rookMove,
	// "♕": blackQueen,
	// "♔": blackKing,
	// "♟": whitePawn,
	// "♞": whiteKnight,
	// "♝": whiteBishop,
	// "♜": whiteRook,
	// "♛": whiteQueen,
	// "♚": whiteKing,
}

//Taking input from
//X               Z               Y
//a b c d e f g - 1 2 3 4 5 6 7 8 - A B C D E F G

// parseLoc turns the user friendly notation input by the user (e.g., "a1A") to a uint32 index which represents that location in the bitmap
// inputs: string | outputs: uint32, x, y, z
func parseLoc(loc string) (uint32, int, int, int) {

	if len(loc) != 3 {
		logger.Error(fmt.Sprintf("Could not parse location '%v' - invalid length of string", loc))
	}
	x, z, y := int(loc[0]-'a'+1), int(loc[1]-'1'+1), int(loc[2]-'A'+1)

	return bitutil.VecToUint(x, y, z), x, y, z //bitutil.VecToUint(x, y, z)
}

// Determines peice type
// inputs uint32 location | outputs: string, bitmap.Bitmap (bitmap )
func pieceType(loc uint32) (string, bitmap.Bitmap) {

	allPieces, _ := load.LoadGame(config.CurrentGame)

	//Contains is simd vectorized, I don't feel the need to optimize this search
	for meta, bm := range allPieces {

		if bm.Contains(loc) {
			// logger.Info(meta)
			return meta, bm
		}

	}
	return "", bitmap.Bitmap{}
}

// Move is the standard function to move one piece to another location.
// It takes the user friendly notation of the from and to locations, parses them into uint32 locations, and then checks if the move is valid for that piece type.
// If it is valid, it updates the bitmaps for both pieces and saves the game state.
// inputs: from string, to string | outputs: none
func MoveCommand(from string, to string) {
	uLocFrom, fX, fY, fZ := parseLoc(from)
	// fmt.Println("FROM")
	logger.Debug(fmt.Sprintf("Move called from %v to %v", from, to))

	logger.Debug(fmt.Sprintf("uLoc: %d | x: %d | y: %d | z: %d", uLocFrom, fX, fY, fZ))

	visFrom, bmFrom := pieceType(uLocFrom)
	if visFrom == "" {
		logger.Error(fmt.Sprintf("Could not find piece at location %v", from))
		return
	}

	// logger.Info(fmt.Sprintf("%v", parseLoc(to)))
	uintLocTo, _, _, _ := parseLoc(to)
	// fmt.Println("TO")
	// fmt.Printf("uLoc: %d | x: %d | y: %d | z: %d \n", uintLocTo, tX, tY, tZ)
	visTo, bmTo := pieceType(uintLocTo)

	//visFrom encodes the type of piece, and thus the move function we use to generate all possible moves
	moveFunction := moveMap[visFrom]

	if moveFunction == nil {
		logger.Error(fmt.Sprintf("Unknown piece [%v]", visFrom))
	}

	allMoves := moveFunction(fX, fY, fZ)

	// uintLoc is a 1-based board coordinate.
	// bitmap.Bitmap uses 0-based bit indices.
	// Move generation already outputs bitmap indices, but
	// external board locations need conversion before use in bitmap.
	if !(allMoves.Contains(uintLocTo - 1)) {
		fmt.Printf("From %d | To %d \n", uLocFrom, uintLocTo)
		logger.Error(fmt.Sprintf("Piece %v cannot move in that way", visFrom))
	}

	//Now, the fun part is that for some ungodly reason, the usage of Remove and Set here DON'T require that -1 offset
	//Why? you may ask, no fucking clue.

	//Updates bitmap of piece being moved - does not validate if move is legal
	atomicMove(uLocFrom, uintLocTo, visTo, visFrom, bmFrom, bmTo)

	logger.Info("Piece Moved Successfully!")

}

// atomicMove instantly moves a piece from one location to another without any validation or checks.
// It is only used in practice from within a validated move funciton, and should only be used for debugging elsehwere
// So many variables are needed becuase no state is stored in the compiled binary itself, and thus the piece must be updated here for changes to take effect.
// inputs: from string, to string | outputs: none
func atomicMove(uintLocFrom uint32, uintLocTo uint32, visTo string, visFrom string, bmFrom bitmap.Bitmap, bmTo bitmap.Bitmap) {

	bmFrom.Remove(uintLocFrom)
	bmFrom.Set(uintLocTo)

	//Updates bitmap (if it exists) of piece being taken
	bmTo.Remove(uintLocTo)

	// logger.Warn(fmt.Sprintf("%v", visFrom))
	// logger.Warn(fmt.Sprintf("%v", visTo))

	metadata.CreateSaveMetaData(config.CurrentGame)
	save.SavePieceType(visFrom, bmFrom)

	if visTo != "" {
		save.SavePieceType(visTo, bmTo)
	}
}
