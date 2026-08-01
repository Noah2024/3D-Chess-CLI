// Imports the root command and executes it. This is the entry point of the application.
package main

import (
	"3DC/cmd/root"
)

// Executes the root command, which in turn executes the appropriate subcommand based on user input.
func main() {
	rootCmd := root.RootCommand()
	rootCmd.Execute()
}
