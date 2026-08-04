// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// functional and system tests for move package
package move_test

import (
	"3DC/internal/game/new"
	"3DC/internal/move"
	"fmt"
	"testing"
)

func TestFunctionalMove(t *testing.T) {
	t.Run("ParseLocBoundry#1", func(t *testing.T) {
		expX, expY, expZ := 0, 0, 0
		uloc, gotX, gotY, gotZ := move.ParseLoc("a1A")
		if gotX != expX || gotY != expY || gotZ != expZ {
			t.Errorf("Expected (%d, %d, %d) but got (%d, %d, %d) ", expX, expY, expZ, gotZ, gotY, gotZ)
		}
		if uloc != 0 {
			t.Errorf("Expected uloc of %d but got %d", 0, uloc)
		}
	})
	t.Run("ParseLocBoundry#2", func(t *testing.T) {
		expX, expY, expZ := 7, 7, 7
		uloc, gotX, gotY, gotZ := move.ParseLoc("h8H")
		if gotX != expX || gotY != expY || gotZ != expZ {
			t.Errorf("Expected (%d, %d, %d) but got (%d, %d, %d) ", expX, expY, expZ, gotZ, gotY, gotZ)
		}
		if uloc != 511 {
			t.Errorf("Expected uloc of %d but got %d", 511, uloc)
		}
	})

	t.Run("ParseLoc", func(t *testing.T) {
		expX, expY, expZ := 3, 5, 6
		uloc, gotX, gotY, gotZ := move.ParseLoc("d7F")
		if gotX != expX || gotY != expY || gotZ != expZ {
			t.Errorf("Expected (%d, %d, %d) but got (%d, %d, %d) ", expX, expY, expZ, gotX, gotY, gotZ)
		}
		if uloc != 371 {
			t.Errorf("Expected uloc of %d but got %d", 371, uloc)
		}
	})

	t.Run("ParseLocOutsideRange", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic")
			}
		}()
		fmt.Println("!!! A Formatted Error Below Here is Expected!!!")
		move.ParseLoc("d7f")
	})

	t.Run("PieceType", func(t *testing.T) {
		allLoadedPieces := new.GenerateSinglePiece(map[string]uint32{"♟": 292})
		vis, bm := move.PieceType(allLoadedPieces, 292)
		fmt.Println(bm)
		if vis != "♟" {
			t.Errorf("Expected '♟' but got '%s'", vis)
		}

		if fmt.Sprint(bm) != "[0 0 0 0 68719476736 0 0 0]" {
			t.Errorf("Expected '[0 0 0 0 68719476736 0 0 0]' but got '%s'", fmt.Sprint(bm))
		}
	})

	t.Run("NoPieceType", func(t *testing.T) {
		allLoadedPieces := new.GenerateSinglePiece(map[string]uint32{"♟": 292})
		vis, bm := move.PieceType(allLoadedPieces, 310)
		fmt.Println(bm)
		if vis != "" {
			t.Errorf("Expected '' but got '%s'", vis)
		}

		if fmt.Sprint(bm) != "[]" {
			t.Errorf("Expected '[]' but got '%s'", fmt.Sprint(bm))
		}
	})
}

func TestSystemMove(t *testing.T) {

}
