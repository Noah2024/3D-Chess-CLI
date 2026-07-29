// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

package debug

import (
	"3DC/cmd/debug/dataPlanes"
	"3DC/cmd/debug/moves"
	"3DC/cmd/debug/uintvec"
	"3DC/cmd/debug/vecuint"
	"3DC/util/logger"

	"github.com/spf13/cobra"
)

func Debug() *cobra.Command {
	debugCMD := &cobra.Command{
		Use:   "debug",
		Short: "execute debug commans and functions",
		Long:  "Used to debug very specific debug commands and functions, DO NOT use uless you know whats happening",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("Calling Debug command")
			// fmt.Fprintf(cmd.OutOrStdout(), "Testing args %s\n", args[0])
		},
	}
	debugCMD.AddCommand(dataPlanes.DataPlanes())
	debugCMD.AddCommand(moves.DebugMoveGen())
	debugCMD.AddCommand(uintvec.UintTVec())
	debugCMD.AddCommand(vecuint.VecTUint())
	return debugCMD

}
