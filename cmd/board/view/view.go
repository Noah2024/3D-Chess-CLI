package view

import (
	"3DC/util/logger"

	view "3DC/internal/board/view"

	"github.com/spf13/cobra"
)

func View() *cobra.Command {
	ViewCommand := &cobra.Command{
		Use:   "view",
		Short: "View the current board",
		Long:  "View the current board",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("Viewing Board")
			view.ViewBoard()
		},
	}
	return ViewCommand
}
