// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

package load_test

import (
	"3DC/internal/game/load"
	"3DC/internal/game/new"
	"3DC/internal/game/save"
	"3DC/util/metadata"
	"fmt"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Run("LoadGame", func(t *testing.T) {
		dir := t.TempDir()
		location := filepath.Join(dir, "CurrentGame")

		allPieces := new.GenerateSinglePiece(map[string]uint32{
			"♚": 0,
			"♔": 511,
		})
		meta := metadata.MetaData{}

		//Was already tested, so we use without testing now
		saveErr := save.SaveGame(allPieces, meta, location)
		if saveErr != nil {
			t.Errorf("Could not save game because of error %s", saveErr)
		}

		loadedPieces, _, loadErr := load.LoadGame(location)
		if loadErr != nil {
			t.Errorf("Could not load game because of error %s", saveErr)
		}

		expected := "map[♔:[0 0 0 0 0 0 0 9223372036854775808] ♕:[0 0 0 0 0 0 0 0] ♖:[0 0 0 0 0 0 0 0] ♗:[0 0 0 0 0 0 0 0] ♘:[0 0 0 0 0 0 0 0] ♙:[0 0 0 0 0 0 0 0] ♚:[1 0 0 0 0 0 0 0] ♛:[0 0 0 0 0 0 0 0] ♜:[0 0 0 0 0 0 0 0] ♝:[0 0 0 0 0 0 0 0] ♞:[0 0 0 0 0 0 0 0] ♟:[0 0 0 0 0 0 0 0]]"

		if fmt.Sprint(loadedPieces) != expected {
			t.Errorf("Pieces Loaded Incorrectly, expected: %s but got %s", expected, fmt.Sprint(loadedPieces))
		}
	})

	t.Run("GenerateBoardState", func(t *testing.T) {
		allPieces := new.GenerateSinglePiece(map[string]uint32{
			"♚": 0,
			"♔": 511,
		})

		BoardState, err := load.GenerateBoardState(allPieces, "♚")
		if err != nil {
			t.Errorf("Could not generate board state because : '%s'", err)
		}

		friendAsStr := fmt.Sprint(BoardState.FriendPieces)
		expectedFriend := "[1 0 0 0 0 0 0 0]"
		if friendAsStr != expectedFriend {
			t.Errorf("Friend bitmap not right, expected: %s, but got %s", expectedFriend, friendAsStr)
		}

		enemyAsStr := fmt.Sprint(BoardState.EnemyPieces)
		expectedEnemy := "[0 0 0 0 0 0 0 9223372036854775808]"
		if enemyAsStr != expectedEnemy {
			t.Errorf("Enemy bitmap not right, expected: %s, but got %s", expectedEnemy, enemyAsStr)
		}

		allPiecesAsStr := fmt.Sprint(BoardState.AllPieces)
		expectedAllPieces := "[1 0 0 0 0 0 0 9223372036854775808]"
		if allPiecesAsStr != expectedAllPieces {
			t.Errorf("AllPieces bitmap not right, expected: %s, but got %s", expectedAllPieces, allPiecesAsStr)
		}

		if BoardState.PieceLoadError != nil {
			t.Errorf("PieceLoadError not right, expected: %v, but got %v", nil, BoardState.PieceLoadError)
		}

		friendKingAsStr := fmt.Sprint(BoardState.FriendKing)
		expectedFriendKing := "[1 0 0 0 0 0 0 0]"
		if friendKingAsStr != expectedFriendKing {
			t.Errorf("FriendKing bitmap not right, expected: %s, but got %s", expectedFriendKing, friendKingAsStr)
		}

		expectedFriendKingLoc := uint32(0)
		if BoardState.FriendKingLoc != expectedFriendKingLoc {
			t.Errorf("FriendKingLoc not right, expected: %d, but got %d", expectedFriendKingLoc, BoardState.FriendKingLoc)
		}

		expectedSwapPawn := false
		if BoardState.SwapPawn != expectedSwapPawn {
			t.Errorf("SwapPawn not right, expected: %v, but got %v", expectedSwapPawn, BoardState.SwapPawn)
		}

		expectedReferencePiece := "♚"
		if BoardState.ReferencePiece != expectedReferencePiece {
			t.Errorf("ReferencePiece not right, expected: %q, but got %q", expectedReferencePiece, BoardState.ReferencePiece)
		}

		expectedIndividualPieces := map[string]string{
			// "♙": blackPawn,
			// "♘": blackKnight,
			// "♗": blackBishop,
			// "♖": blackRook,
			// "♕": blackQueen,
			"♔": "[0 0 0 0 0 0 0 9223372036854775808]",
			// "♟": whitePawn,
			// "♞": whiteKnight,
			// "♝": whiteBishop,
			// "♜": whiteRook,
			// "♛": whiteQueen,
			"♚": "[1 0 0 0 0 0 0 0]",
		}

		for piece, expected := range expectedIndividualPieces {
			got := fmt.Sprint(BoardState.AllIndividualPieces[piece])
			if got != expected {
				t.Errorf("AllIndividualPieces[%q] not right, expected: %s, but got %s", piece, expected, got)
			}
		}
	})
}
