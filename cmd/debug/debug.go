// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// debug Provides CLI interface for accessing the debug functions
package debug

import (
	"3DC/cmd/debug/bmstring"
	"3DC/cmd/debug/dataPlanes"
	"3DC/cmd/debug/moves"
	"3DC/cmd/debug/uintvec"
	"3DC/cmd/debug/vecuint"
	"3DC/util/logger"

	"github.com/spf13/cobra"
)

// Debug returns the root command for the debug subcommand, which allows users to execute various debug commands and functions.
func Debug() *cobra.Command {
	debugCMD := &cobra.Command{
		Use:   "debug",
		Short: "execute debug commans and functions",
		Long:  "Used to debug very specific debug commands and functions, DO NOT use uless you know whats happening",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Debug("Calling Debug command")
		},
	}
	debugCMD.AddCommand(dataPlanes.DataPlanes())
	debugCMD.AddCommand(moves.DebugMoveGen())
	debugCMD.AddCommand(uintvec.UintTVec())
	debugCMD.AddCommand(vecuint.VecTUint())
	debugCMD.AddCommand(bmstring.Bmstring())
	return debugCMD

}
