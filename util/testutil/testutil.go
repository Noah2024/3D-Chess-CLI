package testutil

import (
	"3DC/internal/game/load"
	"3DC/internal/game/new"
	"3DC/internal/move/genMoves"
	"3DC/util/bitutil"
	"fmt"
	"strconv"
	"strings"
)

type MovementScenario struct {
	Name string

	Board    map[string]uint32
	Position []int

	Expected string
}

type PieceTestCase struct {
	Piece     string
	Scenarios []MovementScenario
}

// The functions below were AI generated to speed up the initial validation of piece moves
// All outputs from the below functions were throughly checked with the 3D debugger to ensure accuracy
// While not used at runtime they are kept here for reference's sake

func DumpExpectedMoves(PiecesUnderTest []PieceTestCase) {
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

			asStr, _ := BitmapStringToBinary(fmt.Sprint(moves))
			fmt.Printf("%-24s %q\n", scenario.Name+":", asStr)
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
