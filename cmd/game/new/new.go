// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

package new

import (
	"3DC/internal/game/new"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "new",
		Short: "Create a new game",
		Long:  "Creates a new game. Overwriting previous game stored in CurretGame folder.",
		RunE: func(cmd *cobra.Command, args []string) error {
			new.NewCommand()
			return nil
		},
	}
}
