package main

import (
	"3DChessCLI/cmd/3DChessCLI/root"
	must "3DChessCLI/util"
	"fmt"
	// "github.com/spf13/cobra"
)

func main(){
	fmt.Println("Main Run")
	rootCmd := root.RootCommand()
	must.Must((rootCmd.Execute()))
}