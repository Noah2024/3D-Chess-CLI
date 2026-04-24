package board_test

import (
	"3DC/cmd/board"
	"bytes"
	"testing"
)

func TestLSCmd(t *testing.T) {
	//Create Greet Command
	cmd := board.Board()

	//Set up stdout to capure output
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)

	//Pass Arguments
	cmd.SetArgs([]string{"Printing Board now"})

	//Execute the command
	if err := cmd.Execute(); err != nil {
		t.Errorf("Unexpected error at ls %v", err)
	}

	expectedOutput := "Testing args Printing Board now"
	if expectedOutput != stdout.String() {
		t.Errorf("Expected output: %q, but got: %q", expectedOutput, stdout.String())
	}
	//Check output

}
