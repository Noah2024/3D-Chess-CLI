// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// board Provides CLI interface for accessing the board function
package board

import (
	"3DC/cmd/board/view"
	"3DC/util/logger"

	"github.com/spf13/cobra"
)

// Board returns the root command for the board subcommand, which allows users to view and manage the game board.
func Board() *cobra.Command {
	boardCMD := &cobra.Command{
		Use:   "board",
		Short: "view and manage the board",
		Long:  "board contains one subcommand which allows users to view the game board",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("Calling Board command")
			// fmt.Fprintf(cmd.OutOrStdout(), "Testing args %s\n", args[0])
		},
	}
	boardCMD.AddCommand(view.View())
	return boardCMD

}
