// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// Test cases for special moves, including esPessent and checking
package special_test

//BEFORE FINISHING I NEED TO CHAGE HOW THE ENPESSENT NUMBERS WORK
//CAN NOT BE LEFT AT ZERO (otherwise zero could always be enpessented)

import (
	"3DC/internal/game/load"
	"3DC/internal/game/new"
	"3DC/internal/move/genMoves"
	"3DC/internal/move/special"
	"3DC/util/bitutil"
	"3DC/util/metadata"
	"3DC/util/testutil"
	"fmt"
	"testing"
)

// Test cases for determining if a piece can succesfully add enPessent moves to their generated moves
type DetectMakeEnPessent struct {
	Name           string
	RefPiece       string
	RefLoc         uint32
	PiecesToLoad   map[string]uint32
	WhiteEnPessent uint8
	BlackEnPessent uint8
	Expected       string
}

// Test cases for determining if an En-Pessent move actually removes the pawn its attacking
type HandlRemovalEnPessentMove struct {
	Name     string
	RefPiece string
	From     uint32
	To       uint32

	PiecesToLoad map[string]uint32

	WhiteEnPessent uint8
	BlackEnPessent uint8

	ExpectedNumOfPieces uint8
}

// Test case for determing if a double move correctly changes metadata
type HandlUpdateEnPessent struct {
	Name     string
	RefPiece string
	From     uint32
	To       uint32

	PiecesToLoad map[string]uint32

	WhiteEnPessent         uint8
	BlackEnPessent         uint8
	ExpectedWhiteEnPessent uint8
	ExpectedBlackEnPessent uint8

	Expected string
}

// ♟
var AllDetectionEnPessentCases = []DetectMakeEnPessent{
	{
		Name:           "BlackEnPessentRightEdge",
		RefPiece:       "♙",
		RefLoc:         159,
		PiecesToLoad:   map[string]uint32{"♙": 159, "♟": 158},
		WhiteEnPessent: 9,
		BlackEnPessent: 6,
		Expected:       "[0 12582912 12582912 12582912 0 0 0 0]",
	},
	{
		Name:           "WhiteEnPessentRightEdge",
		RefPiece:       "♟",
		RefLoc:         167,
		PiecesToLoad:   map[string]uint32{"♙": 166, "♟": 167},
		WhiteEnPessent: 6,
		BlackEnPessent: 9,
		Expected:       "[0 211106232532992 211106232532992 211106232532992 0 0 0 0]",
	},
	//{"♙": 169, "♟": 168}
	{
		Name:           "WhiteEnPessentLeftEdge",
		RefPiece:       "♟",
		RefLoc:         161,
		PiecesToLoad:   map[string]uint32{"♙": 160, "♟": 161},
		WhiteEnPessent: 0,
		BlackEnPessent: 9,
		Expected:       "[0 3298534883328 3298534883328 3298534883328 0 0 0 0]",
	},
	{
		Name:           "BlackEnPessentLeftEdge",
		RefPiece:       "♙",
		RefLoc:         153,
		PiecesToLoad:   map[string]uint32{"♙": 153, "♟": 152},
		WhiteEnPessent: 9,
		BlackEnPessent: 0,
		Expected:       "[0 196608 196608 196608 0 0 0 0]",
	},
}

var AllHandleRemovalEnPessentCases = []HandlRemovalEnPessentMove{
	{
		Name:                "BlackEnPessentRightEdge",
		RefPiece:            "♙",
		From:                159,
		To:                  150,
		PiecesToLoad:        map[string]uint32{"♙": 159, "♟": 158},
		WhiteEnPessent:      0,
		BlackEnPessent:      6,
		ExpectedNumOfPieces: 1,
	},
	{
		Name:                "WhiteEnPessentRightEdge",
		RefPiece:            "♟",
		From:                167,
		To:                  174,
		PiecesToLoad:        map[string]uint32{"♙": 166, "♟": 167},
		WhiteEnPessent:      6,
		BlackEnPessent:      0,
		ExpectedNumOfPieces: 1,
	},
	//{"♙": 169, "♟": 168}
	{
		Name:                "WhiteEnPessentLeftEdge",
		RefPiece:            "♟",
		From:                161,
		To:                  168,
		PiecesToLoad:        map[string]uint32{"♙": 160, "♟": 161},
		WhiteEnPessent:      0,
		BlackEnPessent:      9,
		ExpectedNumOfPieces: 1,
	},
	{
		Name:                "BlackEnPessentLeftEdge",
		RefPiece:            "♙",
		From:                153,
		To:                  144,
		PiecesToLoad:        map[string]uint32{"♙": 153, "♟": 152},
		WhiteEnPessent:      9,
		BlackEnPessent:      0,
		ExpectedNumOfPieces: 1,
	},
}

var AllHandleUpdateEnPessentCases = []HandlUpdateEnPessent{
	{
		Name:                   "BlackEnPessentRightEdge",
		RefPiece:               "♙",
		From:                   159,
		To:                     150,
		PiecesToLoad:           map[string]uint32{"♙": 159, "♟": 158},
		WhiteEnPessent:         9,
		BlackEnPessent:         6,
		ExpectedWhiteEnPessent: 9,
		ExpectedBlackEnPessent: 9,
	},
	{
		Name:                   "WhiteEnPessentRightEdge",
		RefPiece:               "♟",
		From:                   167,
		To:                     174,
		PiecesToLoad:           map[string]uint32{"♙": 166, "♟": 167},
		WhiteEnPessent:         6,
		BlackEnPessent:         9,
		ExpectedWhiteEnPessent: 9,
		ExpectedBlackEnPessent: 9,
	},
	//{"♙": 169, "♟": 168}
	{
		Name:                   "WhiteEnPessentLeftEdge",
		RefPiece:               "♟",
		From:                   161,
		To:                     168,
		PiecesToLoad:           map[string]uint32{"♙": 160, "♟": 161},
		WhiteEnPessent:         0,
		BlackEnPessent:         9,
		ExpectedWhiteEnPessent: 9,
		ExpectedBlackEnPessent: 9,
	},
	{
		Name:                   "BlackEnPessentLeftEdge",
		RefPiece:               "♙",
		From:                   153,
		To:                     144,
		PiecesToLoad:           map[string]uint32{"♙": 153, "♟": 152},
		WhiteEnPessent:         9,
		BlackEnPessent:         0,
		ExpectedWhiteEnPessent: 9,
		ExpectedBlackEnPessent: 9,
	},
}

func TestDetectEnPessent(t *testing.T) {

	for _, tc := range AllDetectionEnPessentCases {
		t.Run(tc.Name, func(t *testing.T) {
			allLoadedPieces := new.GenerateSinglePiece(tc.PiecesToLoad)

			bs, _ := load.GenerateBoardState(allLoadedPieces, tc.RefPiece)

			bs.Meta = metadata.CreateDefaultMetaData()

			bs.Meta.WhiteEnPessent = tc.WhiteEnPessent
			bs.Meta.BlackEnPessent = tc.BlackEnPessent

			x, y, z := bitutil.UintToVec(tc.RefLoc)

			moves := genMoves.MoveMap[tc.RefPiece](bs, tc.RefLoc, x, y, z)

			if got := fmt.Sprint(moves); got != tc.Expected {
				t.Errorf("Expected %s, got %s", tc.Expected, got)
				// t.Errorf("ALL PIECES : Expected %064b\n", bs.AllPieces)
				Expected2, _ := testutil.BitmapStringToBinary(tc.Expected)
				t.Errorf("Human Readable Moves: Expected %s\n", Expected2)
				t.Errorf("Human Readable Moves: got %064b\n", moves)

			}
		})
	}
}

func TestHandleEnPessent(t *testing.T) {

	for _, tc := range AllHandleRemovalEnPessentCases {
		t.Run(tc.Name, func(t *testing.T) {
			allLoadedPieces := new.GenerateSinglePiece(tc.PiecesToLoad)

			bs, _ := load.GenerateBoardState(allLoadedPieces, tc.RefPiece)

			bs.Meta = metadata.CreateDefaultMetaData()

			bs.Meta.WhiteEnPessent = tc.WhiteEnPessent
			bs.Meta.BlackEnPessent = tc.BlackEnPessent

			tX, _, _ := bitutil.UintToVec(tc.To)

			// moves := genMoves.MoveMap[tc.RefPiece](bs, tc.RefLoc, x, y, z)
			// fmt.Println(bs.AllIndividualPieces)
			special.UpdateEnPessent(&bs, tc.From, tc.To, tX)

			rbs, _ := load.GenerateBoardState(allLoadedPieces, tc.RefPiece)

			actualNumberOfPieces := rbs.AllPieces.Count()

			if actualNumberOfPieces != int(tc.ExpectedNumOfPieces) {
				t.Errorf("En-Pessent did not take piece Expected %d pieces, but got %d pieces", tc.ExpectedNumOfPieces, actualNumberOfPieces)
			}
		})
	}

	for _, tc := range AllHandleUpdateEnPessentCases {
		t.Run(tc.Name, func(t *testing.T) {
			allLoadedPieces := new.GenerateSinglePiece(tc.PiecesToLoad)

			bs, _ := load.GenerateBoardState(allLoadedPieces, tc.RefPiece)

			bs.Meta = metadata.CreateDefaultMetaData()

			bs.Meta.WhiteEnPessent = tc.WhiteEnPessent
			bs.Meta.BlackEnPessent = tc.BlackEnPessent

			tX, _, _ := bitutil.UintToVec(tc.To)

			special.UpdateEnPessent(&bs, tc.From, tc.To, tX)

			if (bs.Meta.WhiteEnPessent != tc.ExpectedWhiteEnPessent) || (bs.Meta.BlackEnPessent != tc.ExpectedBlackEnPessent) {
				t.Errorf("En-Pessent did not correctly update Expected White:%d Black: %d, but got White: %d  Black: %d",
					tc.ExpectedWhiteEnPessent, tc.ExpectedBlackEnPessent, bs.Meta.WhiteEnPessent, bs.Meta.BlackEnPessent)
			}
		})
	}
}

func TestSpecialMoves(t *testing.T) {
	TestDetectEnPessent(t)
	TestHandleEnPessent(t)
}
