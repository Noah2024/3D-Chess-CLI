// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// new Provides CLI interface for creating a new game
package new

import (
	"3DC/internal/game/new"

	"github.com/spf13/cobra"
)

// NewCommand returns a cobra command that creates a new game and saves it to the CurrentGame folder.
// It overwrites any existing game in the folder without prompting for confirmation.
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
