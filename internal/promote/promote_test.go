// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// Tests wether a pawn can be promoted
// Does NOT test the whole pipeline of AttemptPawnPromotion
package promote_test

import (
	"3DC/internal/game/new"
	"3DC/internal/promote"
	"testing"
)

type movementScenario struct {
	Name string

	Board map[string]uint32
	Team  bool

	Promotion string
	expected  bool
}

func checkPromote(scenario movementScenario, t *testing.T) {
	allLoadedPieces := new.GenerateSinglePiece(scenario.Board)
	_, canPromote := promote.CanPromotePawn(scenario.Team, allLoadedPieces)

	if scenario.expected != canPromote {
		t.Errorf("Failed promotion check, expected '%t', but got '%t'", scenario.expected, canPromote)
	}
}

var allScenarios = []movementScenario{
	{
		Name:      "WhiteValidPromotion",
		Board:     map[string]uint32{"♟": 511, "♚": 64, "♙": 0, "♔": 455},
		Team:      true,
		Promotion: "Queen",
		expected:  true,
	},
	{
		Name:      "BlackValidPromotion",
		Board:     map[string]uint32{"♟": 511, "♚": 64, "♙": 0, "♔": 455},
		Team:      false,
		Promotion: "Queen",
		expected:  true,
	},
	{
		Name:      "WhiteNONValidPromotion",
		Board:     map[string]uint32{"♟": 495, "♚": 64, "♙": 16, "♔": 455},
		Team:      true,
		Promotion: "Queen",
		expected:  false,
	},
	{
		Name:      "BlackNONValidPromotion",
		Board:     map[string]uint32{"♟": 495, "♚": 64, "♙": 16, "♔": 455},
		Team:      false,
		Promotion: "Queen",
		expected:  false,
	},
}

func TestPawnCanPromote(t *testing.T) {
	for _, scenario := range allScenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			checkPromote(scenario, t)
		})

	}
}

func TestSystemAttemptPawnPromotion(t *testing.T) {

}
