// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

package dataPlanes

import (
	"3DC/util/dataplane"

	"github.com/spf13/cobra"
)

func DataPlanes() *cobra.Command {
	return &cobra.Command{
		Use:   "planes",
		Short: "Regenerate data planes used for computation of moves ",
		Long:  "Re-executes algorithms used to generate data planes, DOES NOT affect data planes already in use. They come with the compiled binary",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			dataplane.GenerateAllPlanes()
			return nil
		},
	}
}
