// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// Package move_test contains test cases for each piece movments, it relies on new and view commands to properly validate movement.
// New and delete commands must also pass test cases for this test to be accurate.
package move_test

import (
	"3DC/cmd/game/new"
	"3DC/cmd/move"
	"3DC/config"
	"3DC/util/logger"
	"bytes"
	"path/filepath"
	"testing"
)

//Yes moveTestCase could be an array, but I want it as a struct for readability

// Contains the information for a single move test.
// Currently the expected test case includes ascii escape codes as thats how outut is formatted.
// Be aware that changing this format without changing the expected output may break test cases.
type MoveTestCase struct {
	moveFrom string
	moveTo   string
	reason   string
	expected string
}

// Contains all test cases to test for the move command
// Executed sequentially, so be aware that the state of the board will change after each test case
var allTestCases = []MoveTestCase{
	// ============================================
	// Rook Test Cases
	// ============================================
	MoveTestCase{
		moveFrom: "a8C",
		moveTo:   "a8H",
		reason:   "general movment",
		expected: `[34mINFO: Piece Moved Successfully![0m` + "\n",
	},
	MoveTestCase{
		moveFrom: "a8H",
		moveTo:   "h8H",
		reason:   "general movment",
		expected: `[34mINFO: Piece Moved Successfully![0m` + "\n",
	},
	MoveTestCase{
		moveFrom: "h8H",
		moveTo:   "h1H",
		reason:   "general movment",
		expected: `[34mINFO: Piece Moved Successfully![0m` + "\n",
	},
	MoveTestCase{
		moveFrom: "h1H",
		moveTo:   "a1H",
		reason:   "general movment",
		expected: `[34mINFO: Piece Moved Successfully![0m` + "\n",
	},
	MoveTestCase{ //Ensuring we can't move thorugh pieces
		moveFrom: "a1H",
		moveTo:   "a1A",
		reason:   "friendly passthoough",
		expected: `[31mERROR: Piece ♖ cannot move in that way[0m`,
	},
	MoveTestCase{ //Ensuring we can take enemy pieces
		moveFrom: "a1H",
		moveTo:   "a1C",
		reason:   "taking enemy",
		expected: `[34mINFO: Piece Moved Successfully![0m` + "\n",
	},
	MoveTestCase{
		moveFrom: "a1C",
		moveTo:   "a1A",
		reason:   "general movment",
		expected: `[34mINFO: Piece Moved Successfully![0m` + "\n",
	},
	MoveTestCase{
		moveFrom: "a1A",
		moveTo:   "a7A",
		reason:   "general movment",
		expected: `[34mINFO: Piece Moved Successfully![0m` + "\n",
	},
	MoveTestCase{ //Testing Can't move through friendly pieces
		moveFrom: "a7A",
		moveTo:   "a7H",
		reason:   "friendly passthrough",
		expected: `[31mERROR: Piece ♖ cannot move in that way[0m`,
	},
	MoveTestCase{ //Testing Can't TAKE friendly pieces
		moveFrom: "a7A",
		moveTo:   "a7C",
		reason:   "friendly protection",
		expected: `[31mERROR: Piece ♖ cannot move in that way[0m`,
	},
	// ============================================
	// Bishop Test Cases
	// ============================================
	MoveTestCase{ //Testing Can't TAKE friendly pieces
		moveFrom: "f8C",
		moveTo:   "g7C",
		reason:   "friendly protection",
		expected: `[31mERROR: Piece ♗ cannot move in that way[0m`,
	},
	MoveTestCase{ //Testing Can't TAKE friendly pieces on the other side
		moveFrom: "f8C",
		moveTo:   "a7C",
		reason:   "friendly protection",
		expected: `[31mERROR: Piece ♗ cannot move in that way[0m`,
	},
	MoveTestCase{ //Testing Can't move thorugh friendly pieces
		moveFrom: "f8C",
		moveTo:   "c5C",
		reason:   "friendly protection",
		expected: `[31mERROR: Piece ♗ cannot move in that way[0m`,
	},
	MoveTestCase{ //Testing Can't TAKE friendly pieces
		moveFrom: "f8C",
		moveTo:   "a7C",
		reason:   "friendly protection",
		expected: `[31mERROR: Piece ♗ cannot move in that way[0m`,
	},
	MoveTestCase{
		moveFrom: "f8C",
		moveTo:   "d6E",
		reason:   "general movment",
		expected: `[34mINFO: Piece Moved Successfully![0m` + "\n",
	},
	MoveTestCase{ //////////////////////Will need to see if checking works here later
		moveFrom: "d6E",
		moveTo:   "g3E",
		reason:   "general movment",
		expected: `[34mINFO: Piece Moved Successfully![0m` + "\n",
	},
	MoveTestCase{
		moveFrom: "g3E",
		moveTo:   "e1C",
		reason:   "taking enemy",
		expected: `[34mINFO: Piece Moved Successfully![0m` + "\n",
	},
	MoveTestCase{
		moveFrom: "e1C",
		moveTo:   "g3C",
		reason:   "enemy movethrough",
		expected: `[31mERROR: Piece ♗ cannot move in that way[0m`,
	},
}

func TestMoveCommand(t *testing.T) {

	logger.Debug("--- Starting Move Testing  ---")

	//Create the move command
	newCMD := new.NewCommand()
	moveCMD := move.MoveCommand()

	//redirect stoud to a buffer to capture it
	var stdout bytes.Buffer
	moveCMD.SetOut(&stdout)

	//Change config variables for the sake of not messing with users active game
	config.CurrentGame = filepath.Join(config.DataDir, "TestGame")

	//Create new game to test move command (this is depdendent on new game working)
	if err := newCMD.Execute(); err != nil {
		t.Fatalf("new command failed: %v", err)
	}

	//Changing stdout for logger before executing the part we actually want to test, so that we can capture ONLY the output of the move command
	//And not the debug output of the new command
	logger.SetOutput(&stdout)

	//Loop all indivudal test cases and execute
	for _, testCase := range allTestCases {
		//set args
		moveCMD.SetArgs([]string{testCase.moveFrom, testCase.moveTo})

		// execute the move command w/ args
		if err := moveCMD.Execute(); err != nil {
			t.Errorf("Unexpected error moving from %s to %s at move %v", testCase.moveFrom, testCase.moveTo, err)
			break
		}

		//check the output
		if stdout.String() != testCase.expected {
			t.Errorf("Failed at move check FROM: '%s' TO: '%s'. Failed check type: '%s'", testCase.moveFrom, testCase.moveTo, testCase.reason)
			t.Errorf("Expected output: %q, but got: %q", testCase.expected, stdout.String())
			break
		}

		//Reset the buffer for the next test case
		stdout.Reset()

	}

	logger.Debug("--- Finished Move Testing  ---")

}
