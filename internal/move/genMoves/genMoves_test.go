package genMoves_test

import (
	"3DC/internal/game/load"
	"3DC/internal/game/new"
	"3DC/internal/move/genMoves"
	"3DC/util/bitutil"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

type movementScenario struct {
	name string

	board    map[string]uint32
	position []int

	expected string
}

type pieceTestCase struct {
	piece     string
	scenarios []movementScenario
}

//As it stands many movement scenarious are exactly the same
//So there is room to create common scenarios which could be better

var PiecesUnderTest = []pieceTestCase{
	{
		piece: "♖",
		scenarios: []movementScenario{
			{
				name:     "movementCenter",
				board:    map[string]uint32{},
				position: []int{5, 5, 5},
				expected: "[68719476736 68719476736 68719476736 68719476736 1157443723186933776 68719476736 68719476736 68719476736]",
			},
			{
				name:     "movementCorner1",
				board:    map[string]uint32{},
				position: []int{1, 1, 1},
				expected: "[72340172838076926 1 1 1 1 1 1 1]",
			},
			{
				name:     "movementCorner2",
				board:    map[string]uint32{},
				position: []int{8, 8, 8},
				expected: "[9223372036854775808 9223372036854775808 9223372036854775808 9223372036854775808 9223372036854775808 9223372036854775808 9223372036854775808 9187484529235886208]",
			},
			{
				name:     "movementEdge1",
				board:    map[string]uint32{},
				position: []int{1, 5, 5},
				expected: "[4294967296 4294967296 4294967296 4294967296 72341259464802561 4294967296 4294967296 4294967296]",
			},
			{
				name:     "movementEdge2",
				board:    map[string]uint32{},
				position: []int{5, 1, 5},
				expected: "[1157443723186933776 68719476736 68719476736 68719476736 68719476736 68719476736 68719476736 68719476736]",
			},
			{
				name: "stopAtNearestPiece",
				board: map[string]uint32{
					"♖": 292,
					"♗": 276,
				},
				position: []int{5, 5, 5},
				expected: "[68719476736 68719476736 68719476736 68719476736 1157443723185881088 68719476736 68719476736 68719476736]",
			},
			{
				name: "preventTakingFriends",
				board: map[string]uint32{
					"♖": 292,
					"♗": 164,
				},
				position: []int{5, 5, 5},
				expected: "[0 0 0 68719476736 1157443723186933776 68719476736 68719476736 68719476736]",
			},
			{
				name: "allowTakingEnemies",
				board: map[string]uint32{
					"♖": 292,
					"♟": 290,
				},
				position: []int{5, 5, 5},
				expected: "[68719476736 68719476736 68719476736 68719476736 1157443710302031888 68719476736 68719476736 68719476736]",
			},
		},
	},
	pieceTestCase{
		piece: "♗",
		scenarios: []movementScenario{
			{
				name:     "movementCenter",
				board:    map[string]uint32{},
				position: []int{5, 5, 5},
				expected: "[1 9367487224930664960 19140298420781056 43981136199680 9386671504487645697 43981136199680 19140298420781056 9367487224930664960]",
			},
			{
				name:     "movementCorner1",
				board:    map[string]uint32{},
				position: []int{1, 1, 1},
				expected: "[9241421688590303744 512 262144 134217728 68719476736 35184372088832 18014398509481984 9223372036854775808]",
			},
			{
				name:     "movementCorner2",
				board:    map[string]uint32{},
				position: []int{8, 8, 8},
				expected: "[1 512 262144 134217728 68719476736 35184372088832 18014398509481984 18049651735527937]",
			},
			{
				name:     "movementEdge1",
				board:    map[string]uint32{},
				position: []int{1, 5, 5},
				expected: "[16 576460752303425536 1125899907104768 2199056809984 577588851267340304 2199056809984 1125899907104768 576460752303425536]",
			},
			{
				name:     "movementEdge2",
				board:    map[string]uint32{},
				position: []int{5, 1, 5},
				expected: "[9386671504487645697 43981136199680 19140298420781056 9367487224930664960 1 0 0 0]",
			},
			{
				name: "stopAtNearestPiece",
				board: map[string]uint32{
					"♗": 292,
					"♖": 365,
				},
				position: []int{5, 5, 5},
				expected: "[1 9367487224930664960 19140298420781056 43981136199680 9386671504487645697 8796764110848 1125899911299072 144115188075889152]",
			},
			{
				name: "preventTakingFriends",
				board: map[string]uint32{
					"♗": 292,
					"♖": 219,
				},
				position: []int{5, 5, 5},
				expected: "[0 9367487224930664448 19140298420518912 43981001981952 9386671504487645697 43981136199680 19140298420781056 9367487224930664960]",
			},
			{
				name: "allowTakingEnemies",
				board: map[string]uint32{
					"♗": 292,
					"♟": 221,
				},
				position: []int{5, 5, 5},
				expected: "[1 9367487224930632192 19140298416586752 43981136199680 9386671504487645697 43981136199680 19140298420781056 9367487224930664960]",
			},
		},
	},

	pieceTestCase{
		piece: "♕",
		scenarios: []movementScenario{
			{
				name:     "movementCenter",
				board:    map[string]uint32{},
				position: []int{5, 5, 5},
				expected: "[68719476737 9367487293650141696 19140367140257792 44049855676416 10544115227674579473 44049855676416 19140367140257792 9367487293650141696]",
			},
			{
				name:     "movementCorner1",
				board:    map[string]uint32{},
				position: []int{1, 1, 1},
				expected: "[9313761861428380670 513 262145 134217729 68719476737 35184372088833 18014398509481985 9223372036854775809]",
			},
			{
				name:     "movementCorner2",
				board:    map[string]uint32{},
				position: []int{8, 8, 8},
				expected: "[9223372036854775809 9223372036854776320 9223372036855037952 9223372036988993536 9223372105574252544 9223407221226864640 9241386435364257792 9205534180971414145]",
			},
			{
				name:     "movementEdge1",
				board:    map[string]uint32{},
				position: []int{1, 5, 5},
				expected: "[4294967312 576460756598392832 1125904202072064 2203351777280 649930110732142865 2203351777280 1125904202072064 576460756598392832]",
			},
			{
				name:     "movementEdge2",
				board:    map[string]uint32{},
				position: []int{5, 1, 5},
				expected: "[10544115227674579473 44049855676416 19140367140257792 9367487293650141696 68719476737 68719476736 68719476736 68719476736]",
			},
			{
				name: "stopAtNearestPiece",
				board: map[string]uint32{
					"♕": 292,
					"♖": 276,
				},
				position: []int{5, 5, 5},
				expected: "[68719476737 9367487293650141696 19140367140257792 44049855676416 10544115227673526785 44049855676416 19140367140257792 9367487293650141696]",
			},
			{
				name: "preventTakingFriends",
				board: map[string]uint32{
					"♕": 292,
					"♗": 164,
				},
				position: []int{5, 5, 5},
				expected: "[1 9367487224930664960 19140298420781056 44049855676416 10544115227674579473 44049855676416 19140367140257792 9367487293650141696]",
			},
			{
				name: "allowTakingEnemies",
				board: map[string]uint32{
					"♕": 292,
					"♟": 290,
				},
				position: []int{5, 5, 5},
				expected: "[68719476737 9367487293650141696 19140367140257792 44049855676416 10544115214789677585 44049855676416 19140367140257792 9367487293650141696]",
			},
		},
	},
	pieceTestCase{
		piece: "♘",
		scenarios: []movementScenario{
			{
				name:     "movementCenter",
				board:    map[string]uint32{},
				position: []int{5, 5, 5},
				expected: "[0 0 17764253171712 4503891686195200 11333767002587136 4503891686195200 17764253171712 0]",
			},
			{
				name:     "movementCorner1",
				board:    map[string]uint32{},
				position: []int{1, 1, 1},
				expected: "[132096 65540 258 0 0 0 0 0]",
			},
			{
				name:     "movementCorner2",
				board:    map[string]uint32{},
				position: []int{8, 8, 8},
				expected: "[0 0 0 0 0 4647714815446351872 2305983746702049280 9077567998918656]",
			},
			{
				name:     "movementEdge1",
				board:    map[string]uint32{},
				position: []int{1, 5, 5},
				expected: "[0 0 1108118339584 281492156645376 567348067172352 281492156645376 1108118339584 0]",
			},
			{
				name:     "movementEdge2",
				board:    map[string]uint32{},
				position: []int{5, 1, 5},
				expected: "[11333767002587136 4503891686195200 17764253171712 0 0 0 0 0]",
			},
			{
				name: "stopAtNearestPiece",
				board: map[string]uint32{
					"♘": 292,
					"♗": 228,
				},
				position: []int{5, 5, 5},
				expected: "[0 0 17764253171712 4503891686195200 11333767002587136 4503891686195200 17764253171712 0]",
			},
			{
				name: "preventTakingFriends",
				board: map[string]uint32{
					"♘": 292,
					"♗": 372,
				},
				position: []int{5, 5, 5},
				expected: "[0 0 17764253171712 4503891686195200 11333767002587136 292058824704 17764253171712 0]",
			},
			{
				name: "allowTakingEnemies",
				board: map[string]uint32{
					"♘": 292,
					"♟": 372,
				},
				position: []int{5, 5, 5},
				expected: "[0 0 17764253171712 4503891686195200 11333767002587136 4503891686195200 17764253171712 0]",
			},
		},
	},

	pieceTestCase{
		piece: "♔",
		scenarios: []movementScenario{
			{
				name:     "movementCenter",
				board:    map[string]uint32{},
				position: []int{5, 5, 5},
				expected: "[0 0 0 61814108848128 61745389371392 61814108848128 0 0]",
			},
			{
				name:     "movementCorner1",
				board:    map[string]uint32{},
				position: []int{1, 1, 1},
				expected: "[770 771 0 0 0 0 0 0]",
			},
			{
				name:     "movementCorner2",
				board:    map[string]uint32{},
				position: []int{8, 8, 8},
				expected: "[0 0 0 0 0 0 13889101250810609664 4665729213955833856]",
			},
			{
				name:     "movementEdge1",
				board:    map[string]uint32{},
				position: []int{1, 5, 5},
				expected: "[0 0 0 3311470116864 3307175149568 3311470116864 0 0]",
			},
			{
				name:     "movementEdge2",
				board:    map[string]uint32{},
				position: []int{5, 1, 5},
				expected: "[61745389371392 61814108848128 0 0 0 0 0 0]",
			},
			{
				name: "stopAtNearestPiece",
				board: map[string]uint32{
					"♔": 292,
					"♗": 356,
				},
				position: []int{5, 5, 5},
				expected: "[0 0 0 61814108848128 61745389371392 61745389371392 0 0]",
			},
			{
				name: "preventTakingFriends",
				board: map[string]uint32{
					"♔": 292,
					"♗": 356,
				},
				position: []int{5, 5, 5},
				expected: "[0 0 0 61814108848128 61745389371392 61745389371392 0 0]",
			},
			{
				name: "allowTakingEnemies",
				board: map[string]uint32{
					"♔": 292,
					"♟": 356,
				},
				position: []int{5, 5, 5},
				expected: "[0 0 0 61814108848128 61745389371392 61814108848128 0 0]",
			},
		},
	},

	pieceTestCase{
		piece: "♙",
		scenarios: []movementScenario{
			{
				name:     "movementCenter",
				board:    map[string]uint32{},
				position: []int{5, 5, 5},
				expected: "[0 0 0 17592186044416 17592186044416 17592186044416 0 0]",
			},
			{
				name:     "movementCorner1",
				board:    map[string]uint32{},
				position: []int{1, 1, 1},
				expected: "[256 256 0 0 0 0 0 0]",
			},
			{
				name:     "movementCorner2",
				board:    map[string]uint32{},
				position: []int{8, 8, 8},
				expected: "[0 0 0 0 0 0 0 0]",
			},
			{
				name:     "movementEdge1",
				board:    map[string]uint32{},
				position: []int{1, 5, 5},
				expected: "[0 0 0 1099511627776 1099511627776 1099511627776 0 0]",
			},
			{
				name:     "movementEdge2",
				board:    map[string]uint32{},
				position: []int{5, 1, 5},
				expected: "[17592186044416 17592186044416 0 0 0 0 0 0]",
			},
			{
				name: "stopAtNearestPiece",
				board: map[string]uint32{
					"♙": 292,
					"♟": 284,
				},
				position: []int{5, 5, 5},
				expected: "[0 0 0 268435456 0 268435456 0 0]",
			},
			{
				name: "preventTakingFriends",
				board: map[string]uint32{
					"♙": 292,
					"♗": 229,
				},
				position: []int{5, 5, 5},
				expected: "[0 0 0 268435456 268435456 268435456 0 0]",
			},
			{
				name: "allowTakingEnemies",
				board: map[string]uint32{
					"♙": 292,
					"♟": 219,
				},
				position: []int{5, 5, 5},
				expected: "[0 0 0 402653184 268435456 268435456 0 0]",
			},
		},
	},

	pieceTestCase{
		// Pawns are the only piece whose movement depends on their team,
		// so it makes sense to test both teams.
		piece: "♟",
		scenarios: []movementScenario{
			{
				name:     "movementCenter",
				board:    map[string]uint32{},
				position: []int{5, 5, 5},
				expected: "[0 0 0 17592186044416 17592186044416 17592186044416 0 0]",
			},
			{
				name:     "movementCorner1",
				board:    map[string]uint32{},
				position: []int{1, 1, 1},
				expected: "[256 256 0 0 0 0 0 0]",
			},
			{
				name:     "movementCorner2",
				board:    map[string]uint32{},
				position: []int{8, 8, 8},
				expected: "[0 0 0 0 0 0 0 0]",
			},
			{
				name:     "movementEdge1",
				board:    map[string]uint32{},
				position: []int{1, 5, 5},
				expected: "[0 0 0 1099511627776 1099511627776 1099511627776 0 0]",
			},
			{
				name:     "movementEdge2",
				board:    map[string]uint32{},
				position: []int{5, 1, 5},
				expected: "[17592186044416 17592186044416 0 0 0 0 0 0]",
			},
			{
				name: "stopAtNearestPiece",
				board: map[string]uint32{
					"♙": 300,
					"♟": 292,
				},
				position: []int{5, 5, 5},
				expected: "[0 0 0 17592186044416 0 17592186044416 0 0]",
			},
			{
				name: "preventTakingFriends",
				board: map[string]uint32{
					"♟": 292,
					"♝": 301,
				},
				position: []int{5, 5, 5},
				expected: "[0 0 0 17592186044416 17592186044416 17592186044416 0 0]",
			},
			{
				name: "allowTakingEnemies",
				board: map[string]uint32{
					"♟": 292,
					"♗": 235,
				},
				position: []int{5, 5, 5},
				expected: "[0 0 0 26388279066624 17592186044416 17592186044416 0 0]",
			},
		},
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
		// t.Errorf("ALL PIECES : expected %064b\n", bs.AllPieces)
		expected2, _ := BitmapStringToBinary(expected)
		t.Errorf("Human Readable Moves: expected %s\n", expected2)
		t.Errorf("Human Readable Moves: got %064b\n", moves)

	}
}

func AllMovementTest(t *testing.T, tc pieceTestCase) {
	for _, scenario := range tc.scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			checkMovement(
				t,
				tc.piece,
				scenario.board,
				scenario.position[0],
				scenario.position[1],
				scenario.position[2],
				scenario.expected,
			)
		})
	}
}

// Runs all movement tests for all pieces
func TestMovement(t *testing.T) {
	for _, testCase := range PiecesUnderTest {
		t.Run(testCase.piece, func(t *testing.T) {
			AllMovementTest(t, testCase)
		})
	}
	dumpExpectedMoves()
}

// The functions below were AI generated to speed up the initial validation of piece moves
// All outputs from the below functions were throughly checked with the 3D debugger to ensure accuracy
// While not used at runtime they are kept here for reference's sake

func dumpExpectedMoves() {
	for _, tc := range PiecesUnderTest {
		fmt.Printf("\n=== %s ===\n", tc.piece)

		for _, scenario := range tc.scenarios {
			// Copy the board so we don't modify the original test data.
			board := make(map[string]uint32, len(scenario.board)+1)
			for piece, loc := range scenario.board {
				board[piece] = loc
			}

			// Ensure the piece under test is present.
			board[tc.piece] = bitutil.VecToUint(
				scenario.position[0],
				scenario.position[1],
				scenario.position[2],
			)

			loaded := new.GenerateSinglePiece(board)
			bs, _ := load.GenerateBoardState(loaded, tc.piece)

			loc := bitutil.VecToUint(
				scenario.position[0],
				scenario.position[1],
				scenario.position[2],
			)

			moves := genMoves.MoveMap[tc.piece](
				bs,
				loc,
				scenario.position[0],
				scenario.position[1],
				scenario.position[2],
			)

			asStr, _ := BitmapStringToBinary(fmt.Sprint(moves))
			fmt.Printf("%-24s %q\n", scenario.name+":", asStr)
		}
	}
}

// Turns the shortened form of the bitmaps into a longer form which is interperted later for debug
func BitmapStringToBinary(bitmap string) (string, error) {
	bitmap = strings.Trim(bitmap, "[]")

	fields := strings.Fields(bitmap)

	var out strings.Builder
	out.WriteByte('[')

	for i, field := range fields {
		value, err := strconv.ParseUint(field, 10, 64)
		if err != nil {
			return "", fmt.Errorf("failed to parse %q: %w", field, err)
		}

		if i != 0 {
			out.WriteByte(' ')
		}

		fmt.Fprintf(&out, "%064b", value)
	}

	out.WriteByte(']')

	return out.String(), nil
}
