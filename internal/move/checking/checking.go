// checking Contians all the business logic for determing the state of checked pieces
// Mostly used inside genMoves
package checking

import (
	"3DC/config"
	"3DC/internal/game/load"
	"3DC/internal/move/genMoves"
	"3DC/util/bitutil"

	"github.com/kelindar/bitmap"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func inBounds(v int) bool {
	return 0 <= v && v <= 7
}

// Determines if two positions are related to one another
func distance(loc1 uint32, loc2 uint32) int {
	x1, y1, z1 := bitutil.UintToVec(loc1)
	x2, y2, z2 := bitutil.UintToVec(loc2)

	dx := x2 - x1
	dy := y2 - y1
	dz := z2 - z1

	return abs(dx) + abs(dy) + abs(dz)
}

// Determines if the given piece will put its king in danger by moving.
// If so, it returns the only legal moves it can make, if not then it reuturns and empty bitmap
func KingInDanger(BoardState load.BoardState, loc uint32) (bitmap.Bitmap, bool) {

	// I will admit this is signficantly less elegant than my move generation pipeline
	// But it took me so long to do even this, I couldn't think of anything better
	// AND I can always come back to this code after the v1 release (though it is conceptually quite simple)

	//King can never protect itself out of check
	//Also without this there is a really big recursion loop (obviously)
	if BoardState.ReferencePiece == "♚" || BoardState.ReferencePiece == "♔" {
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
			(dx == 0 && abs(dy) == abs(dz)) ||
			(abs(dx) == abs(dy) && abs(dy) == abs(dz))

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
	r := []rune(BoardState.ReferencePiece)[0]
	badQueen := "♛"
	if r > 9817 {
		lookingFor -= 6
		badQueen = "♕"
	}

	var protectingMoves bitmap.Bitmap
	protectingMoves.Grow(config.BoardSize - 1)

	var PinningEnemies bitmap.Bitmap
	PinningEnemies.Grow(config.BoardSize - 1)
	if BoardState.AllIndividualPieces[badQueen] != nil {
		PinningEnemies.Or(BoardState.AllIndividualPieces[badQueen])
	}

	if BoardState.AllIndividualPieces[string(lookingFor)] != nil {
		PinningEnemies.Or(BoardState.AllIndividualPieces[string(lookingFor)])
	}

	// ==========================================
	// Step through the vector opposite to the king we are in line with
	// If we find a friendly we exit immediatley
	// If we find an enemy we return all the moves it took for us to get there, including that piece
	// ==========================================

	x := x1
	y := y1
	z := z1

	for inBounds(x+stepX) && inBounds(y+stepY) && inBounds(z+stepZ) {
		x += stepX
		y += stepY
		z += stepZ

		uloc := bitutil.VecToUint(x, y, z)

		if uloc == BoardState.FriendKingLoc {
			break
		}

		protectingMoves.Set(uloc)
	}

	for inBounds(x1-stepX) && inBounds(y1-stepY) && inBounds(z1-stepZ) {
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
// Returns isCheck, isCheckMate, allowedKingMoves, allowedProtectingMoves
func IsKingInCheck(BoardState load.BoardState) (bool, bool, bitmap.Bitmap, bitmap.Bitmap) {
	var validKingMoves bitmap.Bitmap
	var enemyAttackingMoves bitmap.Bitmap
	var savingKingMoves bitmap.Bitmap
	savingKingMoves.Grow(config.BoardSize - 1)
	var inCheck bool
	var directionOfAttack = 0 //If there is more than one direction of attack then
	//There is no need to check for friends that can take the king out of check

	// var inCheckMate bool

	// ======================
	// Swaping enemy and friendly (so that the rest of move generation works)
	// Because this happens in sequnce with the rest of the movement command
	// This should not affect anything else because it is swtiched back after
	// ======================
	tmp := BoardState.EnemyPieces
	BoardState.EnemyPieces = BoardState.FriendPieces
	BoardState.FriendPieces = tmp

	//The swaped pieces ensure that this refernce to enemies... I think (look theres a lot of swtiching happening here)
	start, end := '♔', '♙'
	if []rune(BoardState.ReferencePiece)[0] <= '♙' {
		start, end = '♚', '♟'
	}

	//Calcuting king moves to determine if there are valid moves avilable
	x, y, z := bitutil.UintToVec(BoardState.FriendKingLoc)
	validKingMoves = genMoves.MoveMap[string(start)](BoardState, BoardState.FriendKingLoc, x, y, z)

	//Temporarily removing the king from its own bitmap so that when calculating enemy moves we prevent the king from moving INTO potentially dangerous sapces
	BoardState.AllPieces.Remove(BoardState.FriendKingLoc)
	BoardState.SwapPawn = true

	var kingCantMove bitmap.Bitmap
	kingCantMove.Grow(config.BoardSize - 1)

	//This WILL Be sped up with go rountines after all the logic is verified
	for vis, bm := range BoardState.AllIndividualPieces {
		visAsRune := []rune(vis)[0]
		if visAsRune >= start && visAsRune <= end { //Computing types of enemy moves

			// pieceWaitGroup.Go(func() {

			bm.Range(func(curtLoc uint32) {

				//Get this pieces move
				x, y, z := bitutil.UintToVec(curtLoc)
				attackLine := genMoves.MoveMap[vis](BoardState, curtLoc, x, y, z)

				//Remove these moves from possible king moves
				kingCantMove.Or(attackLine)

				//Determine if king is in check from this piece
				if attackLine.Contains(BoardState.FriendKingLoc) { //Attack contains the king
					if directionOfAttack < 1 { //Only need to do this the first time
						inCheck = true
						enemyAttackingMoves.Or(attackLine)

						//Calculating the same move type as if it was coming form the king
						//ANDing them togther, then limiting their output to the given distance between
						//The two pieces to prevent false positives
						x, y, z := bitutil.UintToVec(BoardState.FriendKingLoc)
						limitedAttackingMoves := genMoves.MoveMap[vis](BoardState, BoardState.FriendKingLoc, x, y, z)
						enemyAttackingMoves.And(limitedAttackingMoves)

						maxDist := distance(curtLoc, BoardState.FriendKingLoc)
						enemyAttackingMoves.Filter(func(x uint32) bool {
							if distance(x, BoardState.FriendKingLoc) < maxDist {
								return true
							} else {
								return false
							}
						})
						enemyAttackingMoves.Set(curtLoc)

						directionOfAttack += 1
					} else {
						enemyAttackingMoves.Clear()
						directionOfAttack += 1
					}
				}
			})
		}
	}

	//Actually removing illegal moves from kings valid moves
	// fmt.Printf("kingCantMove : %064b\n", kingCantMove)
	tmpKing := kingCantMove.Clone(nil)
	tmpKing.And(validKingMoves)
	// fmt.Printf("King Location : %064b\n", tmpKing)
	validKingMoves.Xor(tmpKing)

	//Adding the king back to his correct palce
	BoardState.AllPieces.Set(BoardState.FriendKingLoc)
	BoardState.SwapPawn = false

	//Swaping the teams back so the rest of the move calculating can be done as normal
	tmp2 := BoardState.EnemyPieces
	BoardState.EnemyPieces = BoardState.FriendPieces
	BoardState.FriendPieces = tmp2

	if !inCheck { //If the king is not in check there is no need to test the other pieces
		return false, false, validKingMoves, savingKingMoves
	}

	//Swap them back so we can loop over friends again
	if start == '♔' {
		start, end = '♚', '♟'
	} else {
		start, end = '♔', '♙'
	}

	//If the king is attacked from more than one direction it dosn't matter if any piece can move, becuase there is no move
	//from another piece which could take the king out of check
	// fmt.Println("Direction of Attack: ", directionOfAttack)
	if directionOfAttack <= 1 {
		// return inCheck, false, validKingMoves, bitmap.Bitmap{}
		// fmt.Printf("Enemy Attacking Moves %064b", enemyAttackingMoves)

		//After computing if the king is in check we need to see if there are any pieces which can take it out of check
		//Otherwise there is no way to tell if were in a state of checkmate
		for vis, bm := range BoardState.AllIndividualPieces {
			if vis == "♚" || vis == "♔" { //AGAIN, the king cannont save itself
				continue
			}
			visAsRune := []rune(vis)[0]
			if visAsRune >= start && visAsRune <= end {
				// x, y, z := bitutil.UintToVec(BoardState.FriendKingLoc)
				//Computing friendlies moves
				bm.Range(func(X uint32) {
					x, y, z := bitutil.UintToVec(X)
					friendAttempt := genMoves.MoveMap[vis](BoardState, X, x, y, z)

					friendAttempt.And(enemyAttackingMoves)
					_, present := friendAttempt.Max()
					if present { //If this freind intersects with any of the attacking moves
						savingKingMoves = friendAttempt
					}
				})
			}
		}
	}

	//By now we know that the king is in check
	//So we know need to test if either the king can move, AND/OR if it can be taken out of check
	_, canMove := validKingMoves.Max()
	_, canBeSaved := savingKingMoves.Max()

	//If neither then its checkmate
	if !canMove && !canBeSaved {
		return true, true, validKingMoves, savingKingMoves
	}

	return true, false, validKingMoves, savingKingMoves
}
