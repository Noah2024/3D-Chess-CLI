package ls_test

import (
	"3DChessCLI/cmd/ls"
	"bytes"
	"testing"
)

func TestLSCmd(t *testing.T) {
	//Create Greet Command
	cmd := ls.LsCommand()

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
