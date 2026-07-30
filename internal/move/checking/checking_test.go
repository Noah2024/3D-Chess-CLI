// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// checking Contians all the business logic for determing the state of checked pieces
// Mostly used inside genMoves
package checking

import (
	"3DC/internal/game/load"
	"3DC/internal/game/new"
	"3DC/internal/move/genMoves"
	"3DC/util/bitutil"
	"3DC/util/testutil"
	"fmt"
	"testing"
)

// type MovementScenario struct {
// 	Name string

// 	Board    map[string]uint32
// 	Position []int

// 	Expected string
// }

// type PieceTestCase struct {
// 	Piece     string
// 	Scenarios []MovementScenario
// }

//Every piece needs to be able to
//Check the enemy king
//Take their own king out of danger
//Not reveal an attack on their king
//Able to deliver checkmate situation

var PiecesUnderTest = []testutil.PieceTestCase{}

// Runs a specific portion of a test case, loads the board state associated with it, then compares its result.
// x,y,z represents the position at which the piece is located
func checkMovement(t *testing.T, piece string, piecesToLoad map[string]uint32, x, y, z int, expected string) {
	allLoadedPieces := new.GenerateSinglePiece(piecesToLoad)

	bs, _ := load.GenerateBoardState(allLoadedPieces, piece)

	loc := bitutil.VecToUint(x, y, z)

	moves := genMoves.MoveMap[piece](bs, loc, x, y, z)

	if got := fmt.Sprint(moves); got != expected {
		t.Errorf("expected %s, got %s", expected, got)
		// t.Errorf("ALL PIECES : expected %064b\n", bs.AllPieces)
		expected2, _ := testutil.BitmapStringToBinary(expected)
		t.Errorf("Human Readable Moves: expected %s\n", expected2)
		t.Errorf("Human Readable Moves: got %064b\n", moves)

	}
}

func AllMovementTest(t *testing.T, tc testutil.PieceTestCase) {
	for _, scenario := range tc.Scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			checkMovement(
				t,
				tc.Piece,
				scenario.Board,
				scenario.Position[0],
				scenario.Position[1],
				scenario.Position[2],
				scenario.Expected,
			)
		})
	}
}

// Runs all movement tests for all pieces
func TestMovement(t *testing.T) {
	for _, testCase := range PiecesUnderTest {
		t.Run(testCase.Piece, func(t *testing.T) {
			AllMovementTest(t, testCase)
		})
	}
	testutil.DumpExpectedMoves(PiecesUnderTest)
}
