// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// checking Contians all the business logic for determing the state of checked pieces
// Mostly used inside genMoves
package checking_test

import (
	"3DC/internal/game/load"
	"3DC/internal/game/new"
	"3DC/internal/move/checking"
	"3DC/util/bitutil"
	"fmt"
	"testing"
)

type CheckExpected struct {
	InCheck     bool
	InCheckMate bool

	AllowedKingMoves string
	SavingKingMoves  string
	ProtectingMoves  string

	KingDanger bool
}

type CheckingScenario struct {
	Name        string
	RefPiece    string
	RefPosition []int

	Board map[string]uint32

	Expected CheckExpected
}

//Every piece needs to be able to
//Check the enemy king
//Take their own king out of danger
//Not reveal an attack on their king
//Able to deliver checkmate situation

var AllCheckingTests = []CheckingScenario{
	CheckingScenario{
		Name:        "BasicCheck",
		RefPiece:    "♔",
		RefPosition: []int{4, 8, 8},
		Board:       map[string]uint32{"♔": 507, "♜": 451},
		Expected: CheckExpected{
			InCheck:          true,
			AllowedKingMoves: "[0 0 0 0 0 0 2025493932409880576 1446781380292771840]",
			SavingKingMoves:  "[8 8 8 8 8 8 8 2260630401190135]",
			ProtectingMoves:  "[]",
		},
	},
}

// Runs whole portion of checking system from the perspective of the given piece.
// If any part of the returned board state is not as expected it will error
func checkCheckingSystem(t *testing.T, piece string, piecesToLoad map[string]uint32, x, y, z int, expected CheckExpected) {
	allLoadedPieces := new.GenerateSinglePiece(piecesToLoad)

	bs, _ := load.GenerateBoardState(allLoadedPieces, piece)

	loc := bitutil.VecToUint(x, y, z)

	//Run main functions of the checking system
	inCheck, inCheckMate, allowedKingMoves, savingKingMoves := checking.IsKingInCheck(bs)
	protectingMoves, kingDanger := checking.KingInDanger(bs, loc)

	if inCheck != expected.InCheck {
		t.Errorf("Unexpected CheckState expected: %t \n BUT GOT: %t", inCheck, expected.InCheck)
	}
	if inCheckMate != expected.InCheckMate {
		t.Errorf("Unexpected CheckMate State expected: %t \n BUT GOT: %t", inCheckMate, expected.InCheckMate)

	}
	allowedAsStr := fmt.Sprint(allowedKingMoves)
	if allowedAsStr != expected.AllowedKingMoves {
		t.Errorf("Unexpected allowedKingMoves expected: %s \n BUT GOT: %s", allowedAsStr, expected.AllowedKingMoves)

	}
	savingAsString := fmt.Sprint(savingKingMoves)
	if savingAsString != expected.SavingKingMoves {
		t.Errorf("Unexpected savingKingMoves expected: %s \n BUT GOT: %s", savingAsString, expected.SavingKingMoves)

	}
	protectingAsStr := fmt.Sprint(protectingMoves)
	if protectingAsStr != expected.ProtectingMoves {
		t.Errorf("Unexpected protectingMoves expected: %s \n BUT GOT: %s", protectingAsStr, expected.ProtectingMoves)

	}
	if kingDanger != expected.KingDanger {
		t.Errorf("Unexpected kingDanger expected: %t \n BUT GOT: %t", kingDanger, expected.KingDanger)

	}

}

// Runs all movement tests for all pieces
func TestChecking(t *testing.T) {
	for _, testCase := range AllCheckingTests {
		t.Run(testCase.Name, func(t *testing.T) {
			checkCheckingSystem(
				t,
				testCase.RefPiece,
				testCase.Board,
				testCase.RefPosition[0],
				testCase.RefPosition[1],
				testCase.RefPosition[2],
				testCase.Expected,
			)
		})
	}
}
