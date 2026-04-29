package list

import (
	"3DC/internal/game/list"

	"github.com/spf13/cobra"
)

func ListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "lists all games currently saved",
		Long:  "lists all games currently saved",
		RunE: func(cmd *cobra.Command, args []string) error {
			list.ListGames()
			return nil
		},
	}
}
