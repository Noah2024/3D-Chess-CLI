package genMoves

import (
	"3DC/internal/new"
	"testing"

	"github.com/kelindar/bitmap"
)

type pieceTestCase struct {
	piece string

	movementCenter bitmap.Bitmap
	movementCorner1  bitmap.Bitmap
	movementCorner2 bitmap.Bitmap
	movementEdge1 bitmap.Bitmap
	movementEdge2 bitmap.Bitmap

	stopAtNearestPiece bitmap.Bitmap
	multipleBlockers bitmap.Bitmap

	preventTakingFriends bitmap.Bitmap
	allowTakingEnemies bitmap.Bitmap
}

type pawnTestCase struct {
	piece string 
	
	movemenCenter
	movementWhiteStart
	movementBlackStart

	attackBlack
	attackWhite

	//Not implmented features yet, do not need to test for

	// doubleStartWhite
	// doubleStartBlack

	// promotionWhite
	// promotionBlack
}

var PiecesUnderTest = []pieceTestCase {
	pieceTestCase{
		piece: "♖ ",
		movementCenter: ,
	}
}

func GeneralTest(testCase pieceTestCase){
	new.GenerateSinglePiece()
}

func AllMovementTests(t *testing.T){

}