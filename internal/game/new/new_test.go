// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

package new_test

import (
	"3DC/internal/game/new"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {

	t.Run("DefaultStartState", func(t *testing.T) {
		allPieces := new.DefaultStartState()
		// fmt.Println(allPieceLoc)
		expectedStart := "map[♔:[0 0 576460752303423488 0 0 0 0 0] ♕:[0 0 1152921504606846976 0 0 0 0 0] ♖:[0 0 9295429630892703744 0 0 0 0 0] ♗:[0 0 2594073385365405696 0 0 0 0 0] ♘:[0 0 4755801206503243776 0 0 0 0 0] ♙:[0 0 71776119061217280 0 0 0 0 0] ♚:[0 0 16 0 0 0 0 0] ♛:[0 0 8 0 0 0 0 0] ♜:[0 0 129 0 0 0 0 0] ♝:[0 0 36 0 0 0 0 0] ♞:[0 0 66 0 0 0 0 0] ♟:[0 0 65280 0 0 0 0 0]]"
		if expectedStart != fmt.Sprint(allPieces) {
			t.Errorf("Default start state was initalized correctly")
			t.Errorf("Expected '%s'\n", expectedStart)
			t.Errorf("But got '%s'\n", fmt.Sprint(allPieces))
		}
	})

	t.Run("UniqueStartState", func(t *testing.T) {
		allPieces := new.GenerateSinglePiece(map[string]uint32{
			"♚": 0,
			"♔": 511,
		})
		expectedStart := "map[♔:[0 0 0 0 0 0 0 9223372036854775808] ♕:[0 0 0 0 0 0 0 0] ♖:[0 0 0 0 0 0 0 0] ♗:[0 0 0 0 0 0 0 0] ♘:[0 0 0 0 0 0 0 0] ♙:[0 0 0 0 0 0 0 0] ♚:[1 0 0 0 0 0 0 0] ♛:[0 0 0 0 0 0 0 0] ♜:[0 0 0 0 0 0 0 0] ♝:[0 0 0 0 0 0 0 0] ♞:[0 0 0 0 0 0 0 0] ♟:[0 0 0 0 0 0 0 0]]"
		if expectedStart != fmt.Sprint(allPieces) {
			t.Errorf("Default start state was initalized correctly")
			t.Errorf("Expected '%s'\n", expectedStart)
			t.Errorf("But got '%s'\n", fmt.Sprint(allPieces))
		}

	})
}
