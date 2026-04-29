package main

import (
	"3DC/cmd/root"
)

func main() {
	rootCmd := root.RootCommand()
	rootCmd.Execute()
}
