package board

import (
	"3DC/cmd/board/view"
	"3DC/util/logger"

	"github.com/spf13/cobra"
)

func Board() *cobra.Command {
	boardCMD := &cobra.Command{
		Use:   "board",
		Short: "view and manage the board",
		Long:  "view and manage the board",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("Calling Board command")
			// fmt.Fprintf(cmd.OutOrStdout(), "Testing args %s\n", args[0])
		},
	}
	boardCMD.AddCommand(view.View())
	return boardCMD

}
