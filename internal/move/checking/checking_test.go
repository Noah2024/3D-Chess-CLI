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
	InCheck     bool //Obvious
	InCheckMate bool //Obious

	AllowedKingMoves string //All legal moves the king can make
	ProtectingMoves  string //All moves a pinned piece can make (from the perspective of the pref piece)
	SavingKingMoves  string //All moves which could take the king out of check
	//A blank [] represents that the pinningAlgorithm did not need to run for the given piece
	//Meanwhile [0 0 0 0 0 0 0 0] means that it did in some capacity

	KingDanger bool
}

type CheckingScenario struct {
	Name        string //Name of the given scenario
	RefPiece    string //Determines who's team were seeing is in check
	RefPosition []int  //Determines the position at which the ref piece is located
	//Ref position dosn't matter much until checking for pinning

	Board map[string]uint32 //State of the board for the check

	Expected CheckExpected //Exepcted outputs of the check
}

//Every piece needs to be able to
//Check the enemy king
//Take their own king out of danger
//Not reveal an attack on their king
//Able to deliver checkmate situation

var AllCheckingTests = []CheckingScenario{

	CheckingScenario{
		Name:        "No Check",
		RefPiece:    "♔",
		RefPosition: []int{4, 8, 8},
		Board:       map[string]uint32{"♔": 511, "♜": 292},
		Expected: CheckExpected{
			InCheck:          false,
			AllowedKingMoves: "[0 0 0 0 0 0 13889101250810609664 4665729213955833856]",
			SavingKingMoves:  "[0 0 0 0 0 0 0 0]", //Again no piece can save
			ProtectingMoves:  "[]",
		},
	},

	// ============================
	// Can deliver check
	// ============================

	CheckingScenario{
		Name:        "BasicRookCheck",
		RefPiece:    "♔",
		RefPosition: []int{4, 8, 8},
		Board:       map[string]uint32{"♔": 507, "♜": 451},
		Expected: CheckExpected{
			InCheck:          true,
			AllowedKingMoves: "[0 0 0 0 0 0 2025493932409880576 1446781380292771840]",
			SavingKingMoves:  "[0 0 0 0 0 0 0 0]", //No piece can protect this king
			ProtectingMoves:  "[]",
		},
	},

	CheckingScenario{
		Name:        "BasicBishopCheck",
		RefPiece:    "♔",
		RefPosition: []int{4, 8, 8},
		Board:       map[string]uint32{"♔": 511, "♝": 292},
		Expected: CheckExpected{
			InCheck:          true,
			AllowedKingMoves: "[0 0 0 0 0 0 13871086852301127680 4665729213955833856]",
			SavingKingMoves:  "[0 0 0 0 0 0 0 0]", //Again no piece can save
			ProtectingMoves:  "[]",
		},
	},
	CheckingScenario{
		Name:        "BasicKnightCheck",
		RefPiece:    "♔",
		RefPosition: []int{4, 8, 8},
		Board:       map[string]uint32{"♔": 511, "♞": 375},
		Expected: CheckExpected{
			InCheck:          true,
			AllowedKingMoves: "[0 0 0 0 0 0 13889101250810609664 4647714815446351872]",
			SavingKingMoves:  "[0 0 0 0 0 0 0 0]", //Again no piece can save
			ProtectingMoves:  "[]",
		},
	},
	CheckingScenario{
		Name:        "BasicPawnCheck",
		RefPiece:    "♔",
		RefPosition: []int{4, 8, 8},
		Board:       map[string]uint32{"♔": 63, "♟": 118},
		Expected: CheckExpected{
			InCheck:          true,
			AllowedKingMoves: "[4665729213955833856 4647714815446351872 0 0 0 0 0 0]",
			SavingKingMoves:  "[0 0 0 0 0 0 0 0]", //Again no piece can save
			ProtectingMoves:  "[]",
		},
	},

	// ============================
	// Can be pinned correctly
	// ============================

	CheckingScenario{
		Name:        "RookSaving",
		RefPiece:    "♖",
		RefPosition: []int{4, 8, 8},
		Board:       map[string]uint32{"♔": 507, "♜": 451, "♖": 27},
		Expected: CheckExpected{
			InCheck:          true,
			AllowedKingMoves: "[0 0 0 0 0 0 2025493932409880576 1446781380292771840]",
			SavingKingMoves:  "[0 0 0 0 0 0 0 134217728]", //Now the black rook should be able to protect its king
			ProtectingMoves:  "[0 0 0 0 0 0 0 0]",
		},
	},

	CheckingScenario{
		Name:        "BishopSaving",
		RefPiece:    "♗",
		RefPosition: []int{7, 6, 5},
		Board:       map[string]uint32{"♔": 511, "♝": 292, "♗": 358},
		Expected: CheckExpected{
			InCheck:          true,
			AllowedKingMoves: "[0 0 0 0 0 0 13871086852301127680 4665729213955833856]",
			SavingKingMoves:  "[0 0 0 0 0 35184372088832 0 0]",
			ProtectingMoves:  "[]",
		},
	},
	CheckingScenario{
		Name:        "KnightSaving",
		RefPiece:    "♘",
		RefPosition: []int{6, 6, 6},
		Board:       map[string]uint32{"♔": 511, "♞": 375, "♘": 365},
		Expected: CheckExpected{
			InCheck:          true,
			AllowedKingMoves: "[0 0 0 0 0 0 13889101250810609664 4647714815446351872]",
			SavingKingMoves:  "[0 0 0 0 0 36028797018963968 0 0]",
			ProtectingMoves:  "[]",
		},
	},
	CheckingScenario{
		Name:        "PawnSaving",
		RefPiece:    "♔",
		RefPosition: []int{4, 8, 8},
		Board:       map[string]uint32{"♔": 480, "♜": 487, "♙": 491},
		Expected: CheckExpected{
			InCheck:          true,
			AllowedKingMoves: "[0 0 0 0 0 0 3311470116864 3298585214976]",
			SavingKingMoves:  "[0 0 0 0 0 0 0 34359738368]",
			ProtectingMoves:  "[]",
		},
	},
	CheckingScenario{ //Cause the pawn has to be fucking different all the damn time
		Name:        "PawnSavingAsAttack",
		RefPiece:    "♔",
		RefPosition: []int{4, 8, 8},
		Board:       map[string]uint32{"♔": 480, "♜": 487, "♙": 494},
		Expected: CheckExpected{
			InCheck:          true,
			AllowedKingMoves: "[0 0 0 0 0 0 3311470116864 3298585214976]",
			SavingKingMoves:  "[0 0 0 0 0 0 0 824633720832]",
			ProtectingMoves:  "[]",
		},
	},

	// // ============================
	// // Can be pinned king out of check
	// // ============================

	// CheckingScenario{
	// 	Name:        "RookPinned",
	// 	RefPiece:    "♖",
	// 	RefPosition: []int{5, 8, 5},
	// 	Board:       map[string]uint32{"♔": 507, "♜": 451, "♖": 484},
	// 	Expected: CheckExpected{
	// 		InCheck:          true,
	// 		AllowedKingMoves: "[0 0 0 0 0 0 2025493932409880576 1446781380292771840]",
	// 		SavingKingMoves:  "[0 0 0 0 0 0 0 134217728]", //Now the black rook should be able to protect its king
	// 		ProtectingMoves:  "[0 0 0 0 0 0 0 0]",
	// 	},
	// },

	// CheckingScenario{
	// 	Name:        "BishopSaving",
	// 	RefPiece:    "♗",
	// 	RefPosition: []int{7, 6, 5},
	// 	Board:       map[string]uint32{"♔": 511, "♝": 292, "♗": 358},
	// 	Expected: CheckExpected{
	// 		InCheck:          true,
	// 		AllowedKingMoves: "[0 0 0 0 0 0 13871086852301127680 4665729213955833856]",
	// 		SavingKingMoves:  "[0 0 0 0 0 35184372088832 0 0]",
	// 		ProtectingMoves:  "[]",
	// 	},
	// },
	// CheckingScenario{
	// 	Name:        "KnightSaving",
	// 	RefPiece:    "♘",
	// 	RefPosition: []int{6, 6, 6},
	// 	Board:       map[string]uint32{"♔": 511, "♞": 375, "♘": 365},
	// 	Expected: CheckExpected{
	// 		InCheck:          true,
	// 		AllowedKingMoves: "[0 0 0 0 0 0 13889101250810609664 4647714815446351872]",
	// 		SavingKingMoves:  "[0 0 0 0 0 36028797018963968 0 0]",
	// 		ProtectingMoves:  "[]",
	// 	},
	// },
	// CheckingScenario{
	// 	Name:        "PawnSaving",
	// 	RefPiece:    "♔",
	// 	RefPosition: []int{4, 8, 8},
	// 	Board:       map[string]uint32{"♔": 480, "♜": 487, "♙": 491},
	// 	Expected: CheckExpected{
	// 		InCheck:          true,
	// 		AllowedKingMoves: "[4665729213955833856 4647714815446351872 0 0 0 0 0 0]",
	// 		SavingKingMoves:  "[0 0 0 0 0 36028797018963968 0 0]",
	// 		ProtectingMoves:  "[]",
	// 	},
	// },
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
		t.Errorf("Unexpected CheckState expected: %t \n BUT GOT: %t", expected.InCheck, inCheck)
	}
	if inCheckMate != expected.InCheckMate {
		t.Errorf("Unexpected CheckMate State expected: %t \n BUT GOT: %t", expected.InCheckMate, inCheckMate)
	}
	allowedAsStr := fmt.Sprint(allowedKingMoves)
	if allowedAsStr != expected.AllowedKingMoves {
		t.Errorf("Unexpected allowedKingMoves expected: %s \n BUT GOT: %s", expected.AllowedKingMoves, allowedAsStr)
	}
	savingAsString := fmt.Sprint(savingKingMoves)
	if savingAsString != expected.SavingKingMoves {
		t.Errorf("Unexpected savingKingMoves expected: %s \n BUT GOT: %s", expected.SavingKingMoves, savingAsString)
	}
	protectingAsStr := fmt.Sprint(protectingMoves)
	if protectingAsStr != expected.ProtectingMoves {
		t.Errorf("Unexpected protectingMoves expected: %s \n BUT GOT: %s", expected.ProtectingMoves, protectingAsStr)
	}
	if kingDanger != expected.KingDanger {
		t.Errorf("Unexpected kingDanger expected: %t \n BUT GOT: %t", expected.KingDanger, kingDanger)
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
