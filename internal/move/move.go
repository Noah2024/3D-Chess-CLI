// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// Package move contins the logic for moving pieces. Using the planes from dataplanes.go to generate valid moves then validates them for broken rules
// This is the largest and most compelex file in the project, due to it being the core gameplay element EVERYTHING ELSE is in support of the logic containted here
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

var friendPieces bitmap.Bitmap
var enemyPieces bitmap.Bitmap
var allPieces bitmap.Bitmap
var pieceLoadError error

// BIG NOTE TO SELF 7/13/2026
// uintLoc is ONE indexed, meanwhile the bitmap  is ZERO indexed, hence the fuckywucky -1's everywhere
// So every use of Set, Contains, Remove, etc. needs to be -1 from the uintLoc value

// Helper fucntion to determine offset from intersected pieces during restriction of moves
// When calculating wether a move is restricted or not we need to determine wether or not to include it
// Becuase right and left (positive and negative) sides of a move line are calculated differnetly, we need to know which side were dealing with
// in order to correctly offset it
func teamOffset(foundPieceLocation uint32, isRight bool) uint32 {
	if enemyPieces.Contains(foundPieceLocation) {
		if isRight {
			return foundPieceLocation + 1
		} else { //isLeft
			return foundPieceLocation - 1
		}
	} else {
		if isRight {
			return foundPieceLocation - 1
		} else { //isLeft
			return foundPieceLocation + 1
		}
	}
}

// Takes a given vector of a move and seperates it into right and left halves, before running checks on each half to determine if a piece is intersecting
// If no piece is present nothing happens and that half is jointed with the other, if its an enemy that piece is included, if it is friendly it is not.
// All checks are done using bitmap operations
func restrictMoves(curtPieceUintLoc uint32, moveLine bitmap.Bitmap) bitmap.Bitmap {
	// ==============================================
	// Seperate both directions of attack vector
	// ==============================================

	var leftMask bitmap.Bitmap
	var rightHalf bitmap.Bitmap
	rightHalf.Grow(config.BoardSize - 1) //Grow to normal size
	rightHalf.Ones()
	leftMask.Grow(config.BoardSize - 1)
	leftMask.Ones()
	leftMask.Filter(func(x uint32) bool {
		return x < curtPieceUintLoc+1 //Not sure why this has to be +1 exactly, but otherwise movement lines will be slighly off
	})

	rightHalf.Xor(leftMask) //Leaves ones only on the right half
	rightHalf.And(moveLine) //Cancels out all the movement from the left

	// Calculate left half
	var leftHalf bitmap.Bitmap
	leftHalf.Grow(config.BoardSize - 1) //Filled left half with ones
	leftHalf.Ones()
	leftHalf.Filter(func(x uint32) bool {
		return x < curtPieceUintLoc
	})
	leftHalf.And(moveLine) //Only ones left will be one in both

	// ==============================================
	// Determine friend or foe and mask accordingly
	// ==============================================

	// fmt.Printf("Right Before %064b \n", rightHalf)

	rightHitPieces := rightHalf.Clone(nil)   //bitmap.Bitmap
	rightHitPieces.And(allPieces)            //Contians all the pieces if any in the right half
	foundPerson, rtn := rightHitPieces.Min() //Gets first piece to be in line with attack

	if rtn == true { //If there is an enemy
		var newRight bitmap.Bitmap
		newRight.Grow(config.BoardSize - 1)
		newRight.Ones()
		teamOffset := teamOffset(foundPerson, true) //Determines wether to include or not include the piece itself in possible moves
		newRight.Filter(func(x uint32) bool {
			return x < teamOffset //Team offset handles the
		})

		rightHalf.And(newRight)
	}

	leftHitPieces := leftHalf.Clone(nil)
	leftHitPieces.And(allPieces)                //Contians all the pieces if any in the left half
	foundPersonLeft, rtn := leftHitPieces.Max() //Gets first piece to be in line with attack

	// fmt.Printf("Left Before %064b \n", leftHalf)

	if rtn == true { //If there is an enemy
		var newLeft bitmap.Bitmap
		var newLeftMask bitmap.Bitmap //Supposed to be all zeros
		newLeft.Grow(config.BoardSize - 1)
		newLeft.Ones()
		newLeftMask.Grow(config.BoardSize - 1)
		newLeftMask.Ones()                                  //Can only filter through set bits
		teamOffsetVal := teamOffset(foundPersonLeft, false) //Determines wether to include or not include the piece itself in possible moves

		newLeftMask.Filter(func(x uint32) bool {
			return x > teamOffsetVal
		})
		newLeft.And(newLeftMask) //
		leftHalf.And(newLeft)
	}
	// fmt.Printf("Left After %064b \n", leftHalf)

	// ==============================================
	// Combine into final bitmap and return
	// ==============================================

	rightHalf.Or(leftHalf)
	return rightHalf
}

// generateRookMoves contains the bitwise operations necessary to generate all possible moves for a rook piece
// it takes x y and z integer cooridnates and outputs a size 511 bitmap all ones of which represent possible moves
// inputs: x, y, z int | outputs: bitmap.Bitmap
func generateRookMoves(loc uint32, x int, y int, z int) bitmap.Bitmap { // Will parallelize with go rountine
	//Note to self, OK SO, the bitmaps when storing values store an ENTIRE BYTE at a time

	// logger.Debug(fmt.Sprintf("Generating all possible rook moves from :x: %d, y: %d, z: %d", x, y, z))
	forward := dataplane.XPlane[x-1].Clone(nil) //-2
	forward.And(dataplane.ZPlane[z-1])
	forward = restrictMoves(loc, forward) //x

	sideToSide := dataplane.YPlane[y-1].Clone(nil)
	sideToSide.And(dataplane.XPlane[x-1])       //-2
	sideToSide = restrictMoves(loc, sideToSide) //z

	upAndDown := dataplane.YPlane[y-1].Clone(nil)
	upAndDown.And(dataplane.ZPlane[z-1])
	upAndDown = restrictMoves(loc, upAndDown) //y

	forward.Or(upAndDown)
	forward.Or(sideToSide)
	// fmt.Printf("All Pieces %064b\n", allPieces) //For Debug
	// fmt.Printf("All Allowed Moves %064b\n", forward) //For Debug
	return forward
}

// moveMap matches a pieces visual representation to the function that generates all possible moves for that piece
// inputs: string | outputs: function(int, int, int) bitmap.Bitmap
var moveMap = map[string]func(uint32, int, int, int) bitmap.Bitmap{
	// "♙": blackPawn,
	// "♘": blackKnight,
	// "♗": blackBishop,
	"♖": generateRookMoves,
	// "♕": blackQueen,
	// "♔": blackKing,
	// "♟": whitePawn,
	// "♞": whiteKnight,
	// "♝": whiteBishop,
	"♜": generateRookMoves,
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
	// fmt.Printf("All Pieces Alternate %064b\n", allPieces)

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
	logger.Debug(fmt.Sprintf("Move called from %v to %v", from, to))

	visFrom, bmFrom := pieceType(uLocFrom)
	if visFrom == "" {
		logger.Error(fmt.Sprintf("Could not find piece at location %v", from))
		return
	}

	friendPieces, enemyPieces, allPieces, pieceLoadError = load.GetFriendsAndEnemies(config.CurrentGame, visFrom)

	if pieceLoadError != nil {
		logger.Error(fmt.Sprintf("Error in determing pieces team %v", pieceLoadError))
		return
	}
	uintLocTo, _, _, _ := parseLoc(to)
	// fmt.Println("TO")
	// fmt.Printf("uLoc: %d | x: %d | y: %d | z: %d \n", uintLocTo, tX, tY, tZ)
	visTo, bmTo := pieceType(uintLocTo)

	//visFrom encodes the type of piece, and thus the move function we use to generate all possible moves
	moveFunction := moveMap[visFrom]

	if moveFunction == nil {
		logger.Error(fmt.Sprintf("Unknown piece [%v]", visFrom))
		return
	}

	allMoves := moveFunction(uLocFrom, fX, fY, fZ)

	//Kept for eventual need at debug
	// fmt.Printf("Piece Moving %s: ", visFrom)
	// fmt.Printf("Piece Being Taken %s ", visTo)

	if !(allMoves.Contains(uintLocTo)) {
		logger.Error(fmt.Sprintf("Piece %v cannot move in that way", visFrom))
		return
	}

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

	metadata.CreateSaveMetaData(config.CurrentGame)
	save.SavePieceType(visFrom, bmFrom)

	if visTo != "" {
		save.SavePieceType(visTo, bmTo)
	}
}
