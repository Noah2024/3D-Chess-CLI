// checking Contians all the business logic for determing the state of checked pieces
// Mostly used inside genMoves
package checking

import (
	"3DC/config"
	"3DC/internal/game/load"
	"3DC/util/bitutil"

	"github.com/kelindar/bitmap"
)

// var EnemyPieces bitmap.Bitmap
// var AllPieces bitmap.Bitmap
// var PieceLoadError error
// var BlackPawns bitmap.Bitmap //Used to determine direction of pawns move dynamically at runtime
// var wg sync.WaitGroup

// var tmp = genMoves.

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func inBounds(v int) bool {
	return 0 <= v && v <= 8
}

// Determines if the given piece will put its king in danger by moving.
// If so, it returns the only legal moves it can make, if not then it reuturns and empty bitmap
func KingInDanger(BoardState load.BoardState, loc uint32) (bitmap.Bitmap, bool) {

	// I will admit this is signficantly less elegant than my move generation pipeline
	// But it took me so long to do even this, I couldn't think of anything better
	// I can always come back to this code after the v1 release

	//King can never protect itself out of check
	//Also without this there is a really big recursion loop
	if BoardState.PieceInProcess == "♚" || BoardState.PieceInProcess == "♔" {
		return bitmap.Bitmap{}, false
	}

	// ==========================================
	//Compute the vector connecting this piece to the king
	// ==========================================

	x1, y1, z1 := bitutil.UintToVec(loc)
	x2, y2, z2 := bitutil.UintToVec(BoardState.FriendKingLoc)

	////===========Some AI used below

	dx := x2 - x1
	dy := y2 - y1
	dz := z2 - z1

	stepX := 0
	if dx > 0 {
		stepX = 1
	} else if dx < 0 {
		stepX = -1
	}

	stepY := 0
	if dy > 0 {
		stepY = 1
	} else if dy < 0 {
		stepY = -1
	}

	stepZ := 0
	if dz > 0 {
		stepZ = 1
	} else if dz < 0 {
		stepZ = -1
	}

	rookLine :=
		(dx == 0 && dy == 0) ||
			(dx == 0 && dz == 0) ||
			(dy == 0 && dz == 0)

	bishopLine :=
		(dz == 0 && abs(dx) == abs(dy)) ||
			(dy == 0 && abs(dx) == abs(dz)) ||
			(dx == 0 && abs(dy) == abs(dz))

	// ONLY check if we are in line with the king
	if !(rookLine || bishopLine) {
		return bitmap.Bitmap{}, false
	}

	// ==========================================
	// Determine what pieces to check for on our given vector
	// ==========================================
	var lookingFor rune

	//First assuming we are playing as black
	if rookLine {
		lookingFor = '♜'
	} else {
		lookingFor = '♝'
	}

	//=========== Some AI Used above

	//Then determining if we are black or white
	r := []rune(BoardState.PieceInProcess)[0]
	badQueen := "♛"
	if r > 9817 {
		lookingFor -= 6
		badQueen = "♕"
	}

	var protectingMoves bitmap.Bitmap
	protectingMoves.Grow(config.BoardSize - 1)

	var PinningEnemies bitmap.Bitmap
	PinningEnemies.Grow(config.BoardSize - 1)
	PinningEnemies.Or(BoardState.AllIndividualPieces[badQueen])
	PinningEnemies.Or(BoardState.AllIndividualPieces[string(lookingFor)])

	// ==========================================
	// Step through the vector opposite to the king we are in line with
	// If we find a friendly we exit immediatley
	// If we find an enemy we return all the moves it took for us to get there, including that piece
	// ==========================================

	for inBounds(x1+stepX) && inBounds(y1+stepY) && inBounds(z1+stepZ) {
		x1 -= stepX
		y1 -= stepY
		z1 -= stepZ

		uloc := bitutil.VecToUint(x1, y1, z1)

		if BoardState.FriendPieces.Contains(uloc) {
			return protectingMoves, false
		}

		protectingMoves.Set(uloc)
		if PinningEnemies.Contains(uloc) {
			return protectingMoves, true
		}
	}
	return bitmap.Bitmap{0}, false
}

// Runs all movement functions from the enemies perspective to determine if friendly king is in check.
// If king in check it returns a list of all the moves which other peices could make to move the king out of check
// func IsKingInCheck() (bool, bitmap.Bitmap, bitmap.Bitmap, bitmap.Bitmap) {

// }
