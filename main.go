package main

import (
	"3DC/cmd/root"
	"fmt"
	// "github.com/spf13/cobra"
)

func main() {
	fmt.Println("Main Run")

	rootCmd := root.RootCommand()
	rootCmd.Execute()
}
