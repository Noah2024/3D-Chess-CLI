// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// Test package with validates the existance and useage of the board root command
package board_test

import (
	"3DC/cmd/board"
	"3DC/util/logger"
	"bytes"
	"testing"
)

// Tests the board command by executating it and setting standard output to a buffer and then checking the output of the command against the expected output
func TestLSCmd(t *testing.T) {
	//Create Greet Command
	cmd := board.Board()

	//Set up stdout to capure output
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)

	logger.SetOutput(&stdout)

	//Pass Arguments
	// cmd.SetArgs([]string{"Printing Board now"})

	//Execute the command
	if err := cmd.Execute(); err != nil {
		t.Errorf("Unexpected error at ls %v", err)
	}

	expectedOutput := `[34mINFO: Calling Board command[0m` + "\n"
	if expectedOutput != stdout.String() {
		t.Errorf("Expected output: %q, but got: %q", expectedOutput, stdout.String())
	}
	//Check output

}
