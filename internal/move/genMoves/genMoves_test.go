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

//I may still seperate out the pawn test cases from the rest
//There are quite alot of things that the pawn does that no one else does
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
		movementEdge1:   "[4294967296 4294967296 4294967296 4294967296 72341259464802561 4294967296 4294967296 4294967296]",
		movementEdge2:   "[1157443723186933776 68719476736 68719476736 68719476736 68719476736 68719476736 68719476736 68719476736]",

		//The combination of each of these will check all 3 cardinal directions
		//to ensure a piece cannont move through, but more importantly that it
		nearestBoardState:  map[string]uint32{"♖": 292, "♗": 276},
		stopAtNearestPiece: "[68719476736 68719476736 68719476736 68719476736 1157443723185881088 68719476736 68719476736 68719476736]",

		friendBoardState:     map[string]uint32{"♖": 292, "♗": 164},
		preventTakingFriends: "[0 0 0 68719476736 1157443723186933776 68719476736 68719476736 68719476736]",

		enemyBoardState:    map[string]uint32{"♖": 292, "♟": 290},
		allowTakingEnemies: "[68719476736 68719476736 68719476736 68719476736 1157443710302031888 68719476736 68719476736 68719476736]",
	},
	pieceTestCase{
		piece: "♗",

		// Uses an empty board to ensure there is nothing that could interfere with movement
		movementCenter:  "[1 9367487224930664960 19140298420781056 43981136199680 9386671504487645697 43981136199680 19140298420781056 9367487224930664960]",
		movementCorner1: "[9241421688590303744 512 262144 134217728 68719476736 35184372088832 18014398509481984 9223372036854775808]",
		movementCorner2: "[1 512 262144 134217728 68719476736 35184372088832 18014398509481984 18049651735527937]",
		movementEdge1:   "[16 576460752303425536 1125899907104768 2199056809984 577588851267340304 2199056809984 1125899907104768 576460752303425536]",
		movementEdge2:   "[9386671504487645697 43981136199680 19140298420781056 9367487224930664960 1 0 0 0]",

		nearestBoardState:  map[string]uint32{"♗": 292, "♖": 365},
		stopAtNearestPiece: "[1 9367487224930664960 19140298420781056 43981136199680 9386671504487645697 8796764110848 1125899911299072 144115188075889152]",

		friendBoardState:     map[string]uint32{"♗": 292, "♖": 219},
		preventTakingFriends: "[0 9367487224930664448 19140298420518912 43981001981952 9386671504487645697 43981136199680 19140298420781056 9367487224930664960]",

		enemyBoardState:    map[string]uint32{"♗": 292, "♟": 221},
		allowTakingEnemies: "[1 9367487224930632192 19140298416586752 43981136199680 9386671504487645697 43981136199680 19140298420781056 9367487224930664960]",
	},

	pieceTestCase{
		piece: "♕",

		// Uses an empty board to ensure there is nothing that could interfere with movement
		movementCenter:  "[68719476737 9367487293650141696 19140367140257792 44049855676416 10544115227674579473 44049855676416 19140367140257792 9367487293650141696]",
		movementCorner1: "[9313761861428380670 513 262145 134217729 68719476737 35184372088833 18014398509481985 9223372036854775809]",
		movementCorner2: "[9223372036854775809 9223372036854776320 9223372036855037952 9223372036988993536 9223372105574252544 9223407221226864640 9241386435364257792 9205534180971414145]",
		movementEdge1:   "[4294967312 576460756598392832 1125904202072064 2203351777280 649930110732142865 2203351777280 1125904202072064 576460756598392832]",
		movementEdge2:   "[10544115227674579473 44049855676416 19140367140257792 9367487293650141696 68719476737 68719476736 68719476736 68719476736]",

		nearestBoardState:  map[string]uint32{"♕": 292, "♖": 276},
		stopAtNearestPiece: "[68719476737 9367487293650141696 19140367140257792 44049855676416 10544115227673526785 44049855676416 19140367140257792 9367487293650141696]",

		friendBoardState:     map[string]uint32{"♕": 292, "♗": 164},
		preventTakingFriends: "[1 9367487224930664960 19140298420781056 44049855676416 10544115227674579473 44049855676416 19140367140257792 9367487293650141696]",

		enemyBoardState:    map[string]uint32{"♕": 292, "♟": 290},
		allowTakingEnemies: "[68719476737 9367487293650141696 19140367140257792 44049855676416 10544115214789677585 44049855676416 19140367140257792 9367487293650141696]",
	},

	pieceTestCase{
		piece: "♘",

		// Uses an empty board to ensure there is nothing that could interfere with movement
		movementCenter:  "[0 0 17764253171712 4503891686195200 11333767002587136 4503891686195200 17764253171712 0]",
		movementCorner1: "[132096 65540 258 0 0 0 0 0]",
		movementCorner2: "[0 0 0 0 0 4647714815446351872 2305983746702049280 9077567998918656]",
		movementEdge1:   "[0 0 1108118339584 281492156645376 567348067172352 281492156645376 1108118339584 0]",
		movementEdge2:   "[11333767002587136 4503891686195200 17764253171712 0 0 0 0 0]",

		//On Knight this tests the knights ability to hop over other pieces
		nearestBoardState:  map[string]uint32{"♘": 292, "♗": 228},
		stopAtNearestPiece: "[0 0 17764253171712 4503891686195200 11333767002587136 4503891686195200 17764253171712 0]",

		friendBoardState:     map[string]uint32{"♘": 292, "♗": 372},
		preventTakingFriends: "[0 0 17764253171712 4503891686195200 11333767002587136 292058824704 17764253171712 0]",

		enemyBoardState:    map[string]uint32{"♘": 292, "♟": 372},
		allowTakingEnemies: "[0 0 17764253171712 4503891686195200 11333767002587136 4503891686195200 17764253171712 0]",
	},

	pieceTestCase{
		piece: "♔",

		// Uses an empty board to ensure there is nothing that could interfere with movement
		movementCenter:  "[0 0 0 61814108848128 61745389371392 61814108848128 0 0]",
		movementCorner1: "[770 771 0 0 0 0 0 0]",
		movementCorner2: "[0 0 0 0 0 0 13889101250810609664 4665729213955833856]",
		movementEdge1:   "[0 0 0 3311470116864 3307175149568 3311470116864 0 0]",
		movementEdge2:   "[61745389371392 61814108848128 0 0 0 0 0 0]",

		nearestBoardState:  map[string]uint32{"♔": 292, "♗": 356},
		stopAtNearestPiece: "[0 0 0 61814108848128 61745389371392 61745389371392 0 0]",

		friendBoardState:     map[string]uint32{"♔": 292, "♗": 356},
		preventTakingFriends: "[0 0 0 61814108848128 61745389371392 61745389371392 0 0]",

		enemyBoardState:    map[string]uint32{"♔": 292, "♟": 356},
		allowTakingEnemies: "[0 0 0 61814108848128 61745389371392 61814108848128 0 0]",
	},

	// NOTE TO SELF: Some of the pawn checks (such as corners) will need to change once pawn promotion gets introducted
	//IF we want to check for it, otherwise these can stay as is
	pieceTestCase{
		piece: "♙",

		// Uses an empty board to ensure there is nothing that could interfere with movement
		movementCenter:  "[0 0 0 17592186044416 17592186044416 17592186044416 0 0]",
		movementCorner1: "[256 256 0 0 0 0 0 0]",
		movementCorner2: "[0 0 0 0 0 0 0 0]",
		movementEdge1:   "[0 0 0 1099511627776 1099511627776 1099511627776 0 0]",
		movementEdge2:   "[17592186044416 17592186044416 0 0 0 0 0 0]",

		//For pawn this checks to ensure the pawn cannont attack forwards
		nearestBoardState:  map[string]uint32{"♙": 292, "♟": 284},
		stopAtNearestPiece: "[0 0 0 268435456 0 268435456 0 0]",

		friendBoardState:     map[string]uint32{"♙": 292, "♗": 229},
		preventTakingFriends: "[0 0 0 268435456 268435456 268435456 0 0]",

		enemyBoardState:    map[string]uint32{"♙": 292, "♟": 219},
		allowTakingEnemies: "[0 0 0 402653184 268435456 268435456 0 0]",
	},
	pieceTestCase{
		//Pawns are the only piece who's movement depends on their team
		//So it makes sense to test both teams for this case
		piece: "♟",

		// Uses an empty board to ensure there is nothing that could interfere with movement
		movementCenter:  "[0 0 0 17592186044416 17592186044416 17592186044416 0 0]",
		movementCorner1: "[256 256 0 0 0 0 0 0]",
		movementCorner2: "[0 0 0 0 0 0 0 0]",
		movementEdge1:   "[0 0 0 1099511627776 1099511627776 1099511627776 0 0]",
		movementEdge2:   "[17592186044416 17592186044416 0 0 0 0 0 0]",

		//For pawn this checks to ensure the pawn cannont attack forwards
		nearestBoardState:  map[string]uint32{"♙": 300, "♟": 292},
		stopAtNearestPiece: "[0 0 0 17592186044416 0 17592186044416 0 0]",

		friendBoardState:     map[string]uint32{"♟": 292, "♝": 301},
		preventTakingFriends: "[0 0 0 17592186044416 17592186044416 17592186044416 0 0]",

		enemyBoardState:    map[string]uint32{"♟": 292, "♗": 235},
		allowTakingEnemies: "[0 0 0 26388279066624 17592186044416 17592186044416 0 0]",
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

	t.Run("movementEdge1", func(t *testing.T) {
		checkMovement(t, testCase.piece, map[string]uint32{}, 1, 5, 5, testCase.movementEdge1)
	})

	t.Run("movementEdge2", func(t *testing.T) {
		checkMovement(t, testCase.piece, map[string]uint32{}, 5, 1, 5, testCase.movementEdge2)
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
	dumpExpectedMoves()
}

// The functions below were AI generated to speed up the initial validation of piece moves
// All outputs from the below functions were throughly checked with the 3D debugger to ensure accuracy
// While not used at runtime they are kept here for reference's sake

func dumpExpectedMoves() {
	for _, tc := range PiecesUnderTest {
		fmt.Printf("\n=== %s ===\n", tc.piece)

		dump := func(name string, board map[string]uint32, x, y, z int) {
			loaded := new.GenerateSinglePiece(board)
			bs, _ := load.GenerateBoardState(loaded, tc.piece)

			loc := bitutil.VecToUint(x, y, z)
			moves := genMoves.MoveMap[tc.piece](bs, loc, x, y, z)

			// fmt.Printf("%-24s %q\n", name+":", fmt.Sprint(moves))
			asStr, _ := BitmapStringToBinary(fmt.Sprint(moves))
			fmt.Printf("%-24s %q\n", name+":", asStr)

		}

		// Empty board movement
		dump("movementCenter", map[string]uint32{
			tc.piece: bitutil.VecToUint(5, 5, 5),
		}, 5, 5, 5)

		dump("movementCorner1", map[string]uint32{
			tc.piece: bitutil.VecToUint(1, 1, 1),
		}, 1, 1, 1)

		dump("movementCorner2", map[string]uint32{
			tc.piece: bitutil.VecToUint(8, 8, 8),
		}, 8, 8, 8)

		dump("movementEdge1", map[string]uint32{
			tc.piece: bitutil.VecToUint(1, 5, 5),
		}, 1, 5, 5)

		dump("movementEdge2", map[string]uint32{
			tc.piece: bitutil.VecToUint(5, 1, 5),
		}, 5, 1, 5)

		// Blocker tests
		dump("stopAtNearestPiece", tc.nearestBoardState, 5, 5, 5)
		dump("preventTakingFriends", tc.friendBoardState, 5, 5, 5)
		dump("allowTakingEnemies", tc.enemyBoardState, 5, 5, 5)
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
