// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// genMoves_test contains tests for the genMoves package, which is responsible for generating all possible moves for a given piece on the board.
package genMoves_test

import (
	"3DC/internal/game/load"
	"3DC/internal/game/new"
	"3DC/internal/move/genMoves"
	"3DC/util/bitutil"
	"3DC/util/testutil"
	"fmt"
	"testing"
)

type movementScenario struct {
	Name string

	Board    map[string]uint32
	Position []int

	Expected string
}

type pieceTestCase struct {
	Piece     string
	Scenarios []movementScenario
}

//As it stands many movement scenarious are exactly the same
//So there is room to create common Scenarios which could be better

var PiecesUnderTest = []pieceTestCase{
	{
		Piece: "♖",
		Scenarios: []movementScenario{
			{
				Name:     "movementCenter",
				Board:    map[string]uint32{},
				Position: []int{4, 4, 4},
				Expected: "[68719476736 68719476736 68719476736 68719476736 1157443723186933776 68719476736 68719476736 68719476736]",
			},
			{
				Name:     "movementCorner1",
				Board:    map[string]uint32{},
				Position: []int{0, 0, 0},
				Expected: "[72340172838076926 1 1 1 1 1 1 1]",
			},
			{
				Name:     "movementCorner2",
				Board:    map[string]uint32{},
				Position: []int{7, 7, 7},
				Expected: "[9223372036854775808 9223372036854775808 9223372036854775808 9223372036854775808 9223372036854775808 9223372036854775808 9223372036854775808 9187484529235886208]",
			},
			{
				Name:     "movementEdge1",
				Board:    map[string]uint32{},
				Position: []int{0, 4, 4},
				Expected: "[4294967296 4294967296 4294967296 4294967296 72341259464802561 4294967296 4294967296 4294967296]",
			},
			{
				Name:     "movementEdge2",
				Board:    map[string]uint32{},
				Position: []int{4, 0, 4},
				Expected: "[1157443723186933776 68719476736 68719476736 68719476736 68719476736 68719476736 68719476736 68719476736]",
			},
			{
				Name: "stopAtNearestPiece",
				Board: map[string]uint32{
					"♖": 292,
					"♗": 276,
				},
				Position: []int{4, 4, 4},
				Expected: "[68719476736 68719476736 68719476736 68719476736 1157443723185881088 68719476736 68719476736 68719476736]",
			},
			{
				Name: "preventTakingFriends",
				Board: map[string]uint32{
					"♖": 292,
					"♗": 164,
				},
				Position: []int{4, 4, 4},
				Expected: "[0 0 0 68719476736 1157443723186933776 68719476736 68719476736 68719476736]",
			},
			{
				Name: "allowTakingEnemies",
				Board: map[string]uint32{
					"♖": 292,
					"♟": 290,
				},
				Position: []int{4, 4, 4},
				Expected: "[68719476736 68719476736 68719476736 68719476736 1157443710302031888 68719476736 68719476736 68719476736]",
			},
		},
	},
	pieceTestCase{
		Piece: "♗",
		Scenarios: []movementScenario{
			{
				Name:     "movementCenter",
				Board:    map[string]uint32{},
				Position: []int{4, 4, 4},
				Expected: "[1 9367487224930664960 19140298420781056 43981136199680 9386671504487645697 43981136199680 19140298420781056 9367487224930664960]",
			},
			{
				Name:     "movementCorner1",
				Board:    map[string]uint32{},
				Position: []int{0, 0, 0},
				Expected: "[9241421688590303744 512 262144 134217728 68719476736 35184372088832 18014398509481984 9223372036854775808]",
			},
			{
				Name:     "movementCorner2",
				Board:    map[string]uint32{},
				Position: []int{7, 7, 7},
				Expected: "[1 512 262144 134217728 68719476736 35184372088832 18014398509481984 18049651735527937]",
			},
			{
				Name:     "movementEdge1",
				Board:    map[string]uint32{},
				Position: []int{0, 4, 4},
				Expected: "[16 576460752303425536 1125899907104768 2199056809984 577588851267340304 2199056809984 1125899907104768 576460752303425536]",
			},
			{
				Name:     "movementEdge2",
				Board:    map[string]uint32{},
				Position: []int{4, 0, 4},
				Expected: "[9386671504487645697 43981136199680 19140298420781056 9367487224930664960 1 0 0 0]",
			},
			{
				Name: "stopAtNearestPiece",
				Board: map[string]uint32{
					"♗": 292,
					"♖": 365,
				},
				Position: []int{4, 4, 4},
				Expected: "[1 9367487224930664960 19140298420781056 43981136199680 9386671504487645697 8796764110848 1125899911299072 144115188075889152]",
			},
			{
				Name: "preventTakingFriends",
				Board: map[string]uint32{
					"♗": 292,
					"♖": 219,
				},
				Position: []int{4, 4, 4},
				Expected: "[0 9367487224930664448 19140298420518912 43981001981952 9386671504487645697 43981136199680 19140298420781056 9367487224930664960]",
			},
			{
				Name: "allowTakingEnemies",
				Board: map[string]uint32{
					"♗": 292,
					"♟": 221,
				},
				Position: []int{4, 4, 4},
				Expected: "[1 9367487224930632192 19140298416586752 43981136199680 9386671504487645697 43981136199680 19140298420781056 9367487224930664960]",
			},
		},
	},

	pieceTestCase{
		Piece: "♕",
		Scenarios: []movementScenario{
			{
				Name:     "movementCenter",
				Board:    map[string]uint32{},
				Position: []int{4, 4, 4},
				Expected: "[68719476737 9367487293650141696 19140367140257792 44049855676416 10544115227674579473 44049855676416 19140367140257792 9367487293650141696]",
			},
			{
				Name:     "movementCorner1",
				Board:    map[string]uint32{},
				Position: []int{0, 0, 0},
				Expected: "[9313761861428380670 513 262145 134217729 68719476737 35184372088833 18014398509481985 9223372036854775809]",
			},
			{
				Name:     "movementCorner2",
				Board:    map[string]uint32{},
				Position: []int{7, 7, 7},
				Expected: "[9223372036854775809 9223372036854776320 9223372036855037952 9223372036988993536 9223372105574252544 9223407221226864640 9241386435364257792 9205534180971414145]",
			},
			{
				Name:     "movementEdge1",
				Board:    map[string]uint32{},
				Position: []int{0, 4, 4},
				Expected: "[4294967312 576460756598392832 1125904202072064 2203351777280 649930110732142865 2203351777280 1125904202072064 576460756598392832]",
			},
			{
				Name:     "movementEdge2",
				Board:    map[string]uint32{},
				Position: []int{4, 0, 4},
				Expected: "[10544115227674579473 44049855676416 19140367140257792 9367487293650141696 68719476737 68719476736 68719476736 68719476736]",
			},
			{
				Name: "stopAtNearestPiece",
				Board: map[string]uint32{
					"♕": 292,
					"♖": 276,
				},
				Position: []int{4, 4, 4},
				Expected: "[68719476737 9367487293650141696 19140367140257792 44049855676416 10544115227673526785 44049855676416 19140367140257792 9367487293650141696]",
			},
			{
				Name: "preventTakingFriends",
				Board: map[string]uint32{
					"♕": 292,
					"♗": 164,
				},
				Position: []int{4, 4, 4},
				Expected: "[1 9367487224930664960 19140298420781056 44049855676416 10544115227674579473 44049855676416 19140367140257792 9367487293650141696]",
			},
			{
				Name: "allowTakingEnemies",
				Board: map[string]uint32{
					"♕": 292,
					"♟": 290,
				},
				Position: []int{4, 4, 4},
				Expected: "[68719476737 9367487293650141696 19140367140257792 44049855676416 10544115214789677585 44049855676416 19140367140257792 9367487293650141696]",
			},
		},
	},
	pieceTestCase{
		Piece: "♘",
		Scenarios: []movementScenario{
			{
				Name:     "movementCenter",
				Board:    map[string]uint32{},
				Position: []int{4, 4, 4},
				Expected: "[0 0 17764253171712 4503891686195200 11333767002587136 4503891686195200 17764253171712 0]",
			},
			{
				Name:     "movementCorner1",
				Board:    map[string]uint32{},
				Position: []int{0, 0, 0},
				Expected: "[132096 65540 258 0 0 0 0 0]",
			},
			{
				Name:     "movementCorner2",
				Board:    map[string]uint32{},
				Position: []int{7, 7, 7},
				Expected: "[0 0 0 0 0 4647714815446351872 2305983746702049280 9077567998918656]",
			},
			{
				Name:     "movementEdge1",
				Board:    map[string]uint32{},
				Position: []int{0, 4, 4},
				Expected: "[0 0 1108118339584 281492156645376 567348067172352 281492156645376 1108118339584 0]",
			},
			{
				Name:     "movementEdge2",
				Board:    map[string]uint32{},
				Position: []int{4, 0, 4},
				Expected: "[11333767002587136 4503891686195200 17764253171712 0 0 0 0 0]",
			},
			{
				Name: "stopAtNearestPiece",
				Board: map[string]uint32{
					"♘": 292,
					"♗": 228,
				},
				Position: []int{4, 4, 4},
				Expected: "[0 0 17764253171712 4503891686195200 11333767002587136 4503891686195200 17764253171712 0]",
			},
			{
				Name: "preventTakingFriends",
				Board: map[string]uint32{
					"♘": 292,
					"♗": 372,
				},
				Position: []int{4, 4, 4},
				Expected: "[0 0 17764253171712 4503891686195200 11333767002587136 292058824704 17764253171712 0]",
			},
			{
				Name: "allowTakingEnemies",
				Board: map[string]uint32{
					"♘": 292,
					"♟": 372,
				},
				Position: []int{4, 4, 4},
				Expected: "[0 0 17764253171712 4503891686195200 11333767002587136 4503891686195200 17764253171712 0]",
			},
		},
	},

	pieceTestCase{
		Piece: "♔",
		Scenarios: []movementScenario{
			{
				Name:     "movementCenter",
				Board:    map[string]uint32{},
				Position: []int{4, 4, 4},
				Expected: "[0 0 0 61814108848128 61745389371392 61814108848128 0 0]",
			},
			{
				Name:     "movementCorner1",
				Board:    map[string]uint32{},
				Position: []int{0, 0, 0},
				Expected: "[770 771 0 0 0 0 0 0]",
			},
			{
				Name:     "movementCorner2",
				Board:    map[string]uint32{},
				Position: []int{7, 7, 7},
				Expected: "[0 0 0 0 0 0 13889101250810609664 4665729213955833856]",
			},
			{
				Name:     "movementEdge1",
				Board:    map[string]uint32{},
				Position: []int{0, 4, 4},
				Expected: "[0 0 0 3311470116864 3307175149568 3311470116864 0 0]",
			},
			{
				Name:     "movementEdge2",
				Board:    map[string]uint32{},
				Position: []int{4, 0, 4},
				Expected: "[61745389371392 61814108848128 0 0 0 0 0 0]",
			},
			{
				Name: "stopAtNearestPiece",
				Board: map[string]uint32{
					"♔": 292,
					"♗": 356,
				},
				Position: []int{4, 4, 4},
				Expected: "[0 0 0 61814108848128 61745389371392 61745389371392 0 0]",
			},
			{
				Name: "preventTakingFriends",
				Board: map[string]uint32{
					"♔": 292,
					"♗": 356,
				},
				Position: []int{4, 4, 4},
				Expected: "[0 0 0 61814108848128 61745389371392 61745389371392 0 0]",
			},
			{
				Name: "allowTakingEnemies",
				Board: map[string]uint32{
					"♔": 292,
					"♟": 356,
				},
				Position: []int{4, 4, 4},
				Expected: "[0 0 0 61814108848128 61745389371392 61814108848128 0 0]",
			},
		},
	},

	pieceTestCase{
		Piece: "♙",
		Scenarios: []movementScenario{
			{
				Name:     "movementCenter",
				Board:    map[string]uint32{},
				Position: []int{4, 4, 4},
				Expected: "[0 0 0 17592186044416 17592186044416 17592186044416 0 0]",
			},
			{
				Name:     "movementCorner1",
				Board:    map[string]uint32{},
				Position: []int{0, 0, 0},
				Expected: "[256 256 0 0 0 0 0 0]",
			},
			{
				Name:     "movementCorner2",
				Board:    map[string]uint32{},
				Position: []int{7, 7, 7},
				Expected: "[0 0 0 0 0 0 0 0]",
			},
			{
				Name:     "movementEdge1",
				Board:    map[string]uint32{},
				Position: []int{0, 4, 4},
				Expected: "[0 0 0 1099511627776 1099511627776 1099511627776 0 0]",
			},
			{
				Name:     "movementEdge2",
				Board:    map[string]uint32{},
				Position: []int{4, 0, 4},
				Expected: "[17592186044416 17592186044416 0 0 0 0 0 0]",
			},
			{
				Name: "stopAtNearestPiece",
				Board: map[string]uint32{
					"♙": 292,
					"♟": 284,
				},
				Position: []int{4, 4, 4},
				Expected: "[0 0 0 268435456 0 268435456 0 0]",
			},
			{
				Name: "preventTakingFriends",
				Board: map[string]uint32{
					"♙": 292,
					"♗": 229,
				},
				Position: []int{4, 4, 4},
				Expected: "[0 0 0 268435456 268435456 268435456 0 0]",
			},
			{
				Name: "allowTakingEnemies",
				Board: map[string]uint32{
					"♙": 292,
					"♟": 219,
				},
				Position: []int{4, 4, 4},
				Expected: "[0 0 0 402653184 268435456 268435456 0 0]",
			},
		},
	},

	pieceTestCase{
		// Pawns are the only Piece whose movement depends on their team,
		// so it makes sense to test both teams.
		Piece: "♟",
		Scenarios: []movementScenario{
			{
				Name:     "movementCenter",
				Board:    map[string]uint32{},
				Position: []int{4, 4, 4},
				Expected: "[0 0 0 17592186044416 17592186044416 17592186044416 0 0]",
			},
			{
				Name:     "movementCorner1",
				Board:    map[string]uint32{},
				Position: []int{0, 0, 0},
				Expected: "[256 256 0 0 0 0 0 0]",
			},
			{
				Name:     "movementCorner2",
				Board:    map[string]uint32{},
				Position: []int{7, 7, 7},
				Expected: "[0 0 0 0 0 0 0 0]",
			},
			{
				Name:     "movementEdge1",
				Board:    map[string]uint32{},
				Position: []int{0, 4, 4},
				Expected: "[0 0 0 1099511627776 1099511627776 1099511627776 0 0]",
			},
			{
				Name:     "movementEdge2",
				Board:    map[string]uint32{},
				Position: []int{4, 0, 4},
				Expected: "[17592186044416 17592186044416 0 0 0 0 0 0]",
			},
			{
				Name: "stopAtNearestPiece",
				Board: map[string]uint32{
					"♙": 300,
					"♟": 292,
				},
				Position: []int{4, 4, 4},
				Expected: "[0 0 0 17592186044416 0 17592186044416 0 0]",
			},
			{
				Name: "preventTakingFriends",
				Board: map[string]uint32{
					"♟": 292,
					"♝": 301,
				},
				Position: []int{4, 4, 4},
				Expected: "[0 0 0 17592186044416 17592186044416 17592186044416 0 0]",
			},
			{
				Name: "allowTakingEnemies",
				Board: map[string]uint32{
					"♟": 292,
					"♗": 235,
				},
				Position: []int{4, 4, 4},
				Expected: "[0 0 0 26388279066624 17592186044416 17592186044416 0 0]",
			},
			{
				Name: "ValidDoubleMove",
				Board: map[string]uint32{
					"♙": 180,
				},
				Position: []int{4, 2, 6},
				Expected: "[0 17660905521152 17660905521152 17660905521152 0 0 0 0]",
			},
			{
				Name: "NONValidDoubleMove",
				Board: map[string]uint32{
					"♙": 172,
				},
				Position: []int{4, 2, 5},
				Expected: "[0 68719476736 68719476736 68719476736 0 0 0 0]",
			},
			{
				Name: "CannontDoubleThroughPiece",
				Board: map[string]uint32{
					"♙": 180,
					"♗": 172,
				},
				Position: []int{4, 2, 6},
				Expected: "[0 17660905521152 0 17660905521152 0 0 0 0]",
			},
		},
	},
}

// Runs a specific portion of a test case, loads the Board state associated with it, then compares its result.
// x,y,z represents the Position at which the Piece is located
func checkMovement(t *testing.T, Piece string, piecesToLoad map[string]uint32, x, y, z int, Expected string) {
	allLoadedPieces := new.GenerateSinglePiece(piecesToLoad)

	bs, _ := load.GenerateBoardState(allLoadedPieces, Piece)

	loc := bitutil.VecToUint(x, y, z)

	moves := genMoves.MoveMap[Piece](bs, loc, x, y, z)

	if got := fmt.Sprint(moves); got != Expected {
		t.Errorf("Expected %s, got %s", Expected, got)
		// t.Errorf("ALL PIECES : Expected %064b\n", bs.AllPieces)
		Expected2, _ := testutil.BitmapStringToBinary(Expected)
		t.Errorf("Human Readable Moves: Expected %s\n", Expected2)
		t.Errorf("Human Readable Moves: got %064b\n", moves)

	}
}

func AllMovementTest(t *testing.T, tc pieceTestCase) {
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
	// testutil.DumpExpectedMoves(PiecesUnderTest)
}

// The functions below were AI generated to speed up the initial validation of piece moves
// All outputs from the below functions were throughly checked with the 3D debugger to ensure accuracy
// While not used at runtime they are kept here for reference's sake

func DumpExpectedMoves(PiecesUnderTest []pieceTestCase) {
	for _, tc := range PiecesUnderTest {
		fmt.Printf("\n=== %s ===\n", tc.Piece)

		for _, scenario := range tc.Scenarios {
			// Copy the board so we don't modify the original test data.
			board := make(map[string]uint32, len(scenario.Board)+1)
			for piece, loc := range scenario.Board {
				board[piece] = loc
			}

			// Ensure the piece under test is present.
			board[tc.Piece] = bitutil.VecToUint(
				scenario.Position[0],
				scenario.Position[1],
				scenario.Position[2],
			)

			loaded := new.GenerateSinglePiece(board)
			bs, _ := load.GenerateBoardState(loaded, tc.Piece)

			loc := bitutil.VecToUint(
				scenario.Position[0],
				scenario.Position[1],
				scenario.Position[2],
			)

			moves := genMoves.MoveMap[tc.Piece](
				bs,
				loc,
				scenario.Position[0],
				scenario.Position[1],
				scenario.Position[2],
			)

			asStr, _ := testutil.BitmapStringToBinary(fmt.Sprint(moves))
			fmt.Printf("%-24s %q\n", scenario.Name+":", asStr)
		}
	}
}
