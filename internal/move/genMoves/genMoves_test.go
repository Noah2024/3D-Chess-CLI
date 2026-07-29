package genMoves_test

import (
	"3DC/internal/game/load"
	"3DC/internal/game/new"
	"3DC/internal/move/genMoves"
	"3DC/util/bitutil"
	"fmt"
	"testing"
)

type pieceTestCase struct {
	piece string
	// boardState map[string]uint32

	movementCenter  string
	movementCorner1 string
	movementCorner2 string
	movementEdge1   string
	movementEdge2   string

	nearestBoardState  map[string]uint32
	stopAtNearestPiece string

	friendBoardState     map[string]uint32
	preventTakingFriends string

	enemyBoardState    map[string]uint32
	allowTakingEnemies string
}

// type pawnTestCase struct {
// 	piece string

// 	movemenCenter
// 	movementWhiteStart
// 	movementBlackStart

// 	attackBlack
// 	attackWhite

// 	//Not implmented features yet, do not need to test for

// 	// doubleStartWhite
// 	// doubleStartBlack

// 	// promotionWhite
// 	// promotionBlack
// }

var PiecesUnderTest = []pieceTestCase{
	pieceTestCase{
		piece: "♖",

		//Uses an empty board to ensure there is nothing that could interfere with movement
		movementCenter:  "[68719476736 68719476736 68719476736 68719476736 1157443723186933776 68719476736 68719476736 68719476736]",                                                         //5,5,5
		movementCorner1: "[72340172838076926 1 1 1 1 1 1 1]",                                                                                                                                 //1,1,1
		movementCorner2: "[9223372036854775808 9223372036854775808 9223372036854775808 9223372036854775808 9223372036854775808 9223372036854775808 9223372036854775808 9187484529235886208]", //8,8,8
		movementEdge1:   "",
		movementEdge2:   "",

		//The combination of each of these will check all 3 cardinal directions
		//to ensure a piece cannont move through, but more importantly that it
		nearestBoardState:  map[string]uint32{"♖": 292, "♗": 276},
		stopAtNearestPiece: "[68719476736 68719476736 68719476736 68719476736 1157443723185881088 68719476736 68719476736 68719476736]",

		friendBoardState:     map[string]uint32{"♖": 292, "♗": 164},
		preventTakingFriends: "[0 0 0 68719476736 1157443723186933776 68719476736 68719476736 68719476736]",

		enemyBoardState:    map[string]uint32{"♖": 292, "♟": 290},
		allowTakingEnemies: "[68719476736 68719476736 68719476736 68719476736 1157443710302031888 68719476736 68719476736 68719476736]",
	},
}

// Runs a specific portion of a test case, loads the board state associated with it, then compares its result.
// x,y,z represents the position at which the piece is located
func checkMovement(t *testing.T, piece string, piecesToLoad map[string]uint32, x, y, z int, expected string) {
	allLoadedPieces := new.GenerateSinglePiece(piecesToLoad)

	bs, _ := load.GenerateBoardState(allLoadedPieces, piece)

	loc := bitutil.VecToUint(x, y, z)

	moves := genMoves.MoveMap[piece](bs, loc, x, y, z)

	if got := fmt.Sprint(moves); got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

// Takes a given test case and runs each test in its own subtest
func AllMovementTest(t *testing.T, testCase pieceTestCase) {
	//General Movement
	t.Run("movementCenter", func(t *testing.T) {
		checkMovement(t, testCase.piece, map[string]uint32{}, 5, 5, 5, testCase.movementCenter)

	})
	t.Run("movementCorner1", func(t *testing.T) {
		checkMovement(t, testCase.piece, map[string]uint32{}, 1, 1, 1, testCase.movementCorner1)

	})
	t.Run("movementCorner2", func(t *testing.T) {
		checkMovement(t, testCase.piece, map[string]uint32{}, 8, 8, 8, testCase.movementCorner2)
	})

	t.Run("stopAtNearestPiece", func(t *testing.T) {
		checkMovement(t, testCase.piece, testCase.nearestBoardState, 5, 5, 5, testCase.stopAtNearestPiece)
	})

	t.Run("preventTakingFriends", func(t *testing.T) {
		checkMovement(t, testCase.piece, testCase.friendBoardState, 5, 5, 5, testCase.preventTakingFriends)
	})

	t.Run("allowTakingEnemies", func(t *testing.T) {
		checkMovement(t, testCase.piece, testCase.enemyBoardState, 5, 5, 5, testCase.allowTakingEnemies)
	})

}

// Runs all movement tests for all pieces
func TestMovement(t *testing.T) {
	for _, testCase := range PiecesUnderTest {
		t.Run(testCase.piece, func(t *testing.T) {
			AllMovementTest(t, testCase)
		})
	}
}
