package root_test

import (
	"3DChessCLI/cmd/root"
	"bytes"
	"testing"
)

func TestRootCommand(t *testing.T) {
	//Create the root command
	cmd := root.RootCommand()

	//redirect stoud to a buffer to capture it
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)

	//set args

	//execute the root command w/ args
	// err := cmd.Execute()
	if err := cmd.Execute(); err != nil {
		t.Errorf("Unexpected error at root %v", err)
	}

	//check the output
	expectedOutput := "bruh"
	if stdout.String() != expectedOutput {
		t.Errorf("Expected output: %q, but got: %q", expectedOutput, stdout.String())
	}

}
