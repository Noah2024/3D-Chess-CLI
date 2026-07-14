// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3
// Copyright 2026 Your Company Name
package board_test

import (
	"3DC/cmd/board"
	"3DC/util/logger"
	"bytes"
	"testing"
)

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
