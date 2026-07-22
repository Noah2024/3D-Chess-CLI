// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

package debug

import (
	"3DC/cmd/debug/dataPlanes"
	"3DC/util/logger"

	"github.com/spf13/cobra"
)

func Debug() *cobra.Command {
	boardCMD := &cobra.Command{
		Use:   "debug",
		Short: "execute debug commans and functions",
		Long:  "Used to debug very specific debug commands and functions, DO NOT use uless you know whats happening",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("Calling Debug command")
			// fmt.Fprintf(cmd.OutOrStdout(), "Testing args %s\n", args[0])
		},
	}
	boardCMD.AddCommand(dataPlanes.DataPlanes())
	return boardCMD

}
