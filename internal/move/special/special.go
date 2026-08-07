// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// Contains logic behind certian special moves, including enPessent and Castling
package special

import (
	"3DC/internal/game/load"
	"3DC/internal/game/save"
	"3DC/util/bitutil"
	"3DC/util/dataplane"
	"fmt"

	"github.com/kelindar/bitmap"
)

func CouldMakeEnPessent(canEnPessentFrom bitmap.Bitmap, uintloc uint32, x int, enemyEnPessentable int) bool {
	if canEnPessentFrom.Contains(uintloc) {
		if x == enemyEnPessentable {
			return true
		}
	}
	return false
}

// Checks state of move just made and determines if the move was 'enPessantable' if so
// It updates the value stored in metadata. It also clears the metadata for expired Enpessent moves
func UpdateEnPessent(BoardState *load.BoardState, uLocFrom uint32, uintLocTo uint32, tX int) {

	var DoubleMovePlane bitmap.Bitmap
	var EnPessentPlane bitmap.Bitmap
	var CanEnPessentFrom bitmap.Bitmap
	var enemyVis string
	var update *uint8

	//Rewrite description of function
	//Determine if EnPessent Removes enemy piece - Test case one
	//Determine if double move triggers EnPessent - test case two

	if !BoardState.Team { //If White
		enemyVis = "♟"
		DoubleMovePlane = dataplane.WhiteDoubleMovePlane
		CanEnPessentFrom = dataplane.BlackEnPessentPlane
		EnPessentPlane = dataplane.WhiteEnPessentPlane
		update = &BoardState.Meta.WhiteEnPessent
		BoardState.Meta.BlackEnPessent = 9 //Weather we used it or not this enPessable move has expired
	} else { //If Black
		enemyVis = "♙"
		DoubleMovePlane = dataplane.BlackDoubleMovePlane
		CanEnPessentFrom = dataplane.WhiteEnPessentPlane
		EnPessentPlane = dataplane.BlackEnPessentPlane
		BoardState.Meta.WhiteEnPessent = 9
		update = &BoardState.Meta.BlackEnPessent
	}

	if BoardState.ReferencePiece != "♙" && BoardState.ReferencePiece != "♟" {
		return //IF were not looking at a pawn return immediatley
	}

	//Consdiering were looking at pawn lets check if it's made an enPessentMoveitself
	attackX, _, _ := bitutil.UintToVec(uintLocTo) //Used to help determine if enPessent move was ade
	if CanEnPessentFrom.Contains(uLocFrom) && attackX == tX {
		_, fromY, fromZ := bitutil.UintToVec(uLocFrom)
		enemyLoc := bitutil.VecToUint(tX, fromY, fromZ)
		enemyPieceBm := BoardState.AllIndividualPieces[enemyVis]
		enemyPieceBm.Remove(enemyLoc)
		save.SavePieceType(enemyVis, enemyPieceBm)
		fmt.Println()
	}

	from := DoubleMovePlane.Contains(uLocFrom)
	to := EnPessentPlane.Contains(uintLocTo)
	if from && to {
		//WE'VE MADE A DOUBLE MOVE AND MAY BE ENPESSANTED
		*update = uint8(tX)
	}
}
