package view

import (
	"3DC/config"
	"3DC/internal/board/view"
	"3DC/util/logger"
	"3DC/util/must"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func View() *cobra.Command {
	ViewCommand := &cobra.Command{
		Use:   "view",
		Short: "View a given vertical slice of baord",
		Long:  "Takes one optional integer argument noting what layer to display (0-7). No argument passed will show every layer",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("Calling Board command")
			//Validating view command input
			if len(args) == 1 {
				layerNum := must.Must(strconv.Atoi(args[0]))
				numOfLayers := int((config.BoardSize / config.LayerSize)) - 1
				if (layerNum > 0) && (layerNum <= numOfLayers) {
					view.ViewLayer(layerNum)
				} else {
					logger.Error(fmt.Sprintf("Layer %d does not exist; provide a number between (0-%d)", layerNum, numOfLayers))
				}
			} else {
				view.ViewAllLayers()
			}
			// fmt.Fprintf(cmd.OutOrStdout(), "Testing args %s\n", args[0])
		},
	}
	return ViewCommand
}
