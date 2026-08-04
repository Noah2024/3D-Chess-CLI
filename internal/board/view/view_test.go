// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

package view_test

import (
	"3DC/internal/board/view"
	"testing"

	"github.com/kelindar/bitmap"
)

//Most of these test cases are AI Generated

func TestView(t *testing.T) {
	t.Run("EmptyBoard", func(t *testing.T) {
		allPieces := make(map[string]bitmap.Bitmap)

		board := view.BuildLayer(allPieces, 0)

		for z := range board {
			for x := range board[z] {
				if board[z][x] != "" {
					t.Errorf("expected empty square at [%d][%d], got %q", z, x, board[z][x])
				}
			}
		}
	})
	t.Run("SinglePiece", func(t *testing.T) {
		allPieces := make(map[string]bitmap.Bitmap)

		var bm bitmap.Bitmap
		bm.Set(0) // Replace with a known index if 0 isn't A1 on layer A.

		allPieces["♟"] = bm

		board := view.BuildLayer(allPieces, 0)

		if board[0][0] != "♟" {
			t.Errorf("expected P at [0][0], got %q", board[0][0])
		}
	})

}

// // "♙": blackPawn,
// 			// "♘": blackKnight,
// 			// "♗": blackBishop,
// 			// "♖": blackRook,
// 			// "♕": blackQueen,
// 			"♔": "[0 0 0 0 0 0 0 9223372036854775808]",
// 			// "♟": whitePawn,
// 			// "♞": whiteKnight,
// 			// "♝": whiteBishop,
// 			// "♜": whiteRook,
// 			// "♛": whiteQueen,
// 			"♚": "[1 0 0 0 0 0 0 0]",
// 		}
