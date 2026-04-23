package main

import (
	"3DChessCLI/cmd/root"
	util "3DChessCLI/util/must"
	"fmt"
	// "github.com/spf13/cobra"
)

func main() {
	fmt.Println("Main Run")

	rootCmd := root.RootCommand()
	util.Must((rootCmd.Execute()))
}
