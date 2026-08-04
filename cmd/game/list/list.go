// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// list Provides CLI interface for listing all saved games
package list

import (
	"3DC/internal/game/list"
	"os"

	"github.com/spf13/cobra"
)

// Lists all saved games from the CurrentGame folder (prints to stdout). Takes no arguments.
func ListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "lists all games currently saved",
		Long:  "lists all games currently saved",
		RunE: func(cmd *cobra.Command, args []string) error {
			list.ListGames(os.Stdout)
			return nil
		},
	}
}
