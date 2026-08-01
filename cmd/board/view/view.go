// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// view Provides CLI interface for accessing the view board function
// It preforms basic input validation on behalf of the view function and then calls the view function with the appropriate arguments
package view

import (
	"3DC/config"
	"3DC/internal/board/view"
	"3DC/util/logger"
	"fmt"

	"github.com/spf13/cobra"
)

func View() *cobra.Command {
	ViewCommand := &cobra.Command{
		Use:   "view",
		Short: "View given vertical slice/s of the game board",
		Long:  "Takes one optional integer argument noting what layer to display (0-7). No argument passed will show every layer",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger.Debug("Calling Board command")

			if len(args) == 1 {
				asRune := []rune(args[0])

				if len(asRune) != 1 {
					logger.Error(fmt.Sprintf("Optional argument takes only one character"))
				}
				layerNum := int(asRune[0] - 'A') //must.Must(strconv.Atoi(args[0]))

				numOfLayers := int((config.BoardSize / config.LayerSize)) - 1
				if (layerNum >= 0) && (layerNum <= numOfLayers) {
					view.ViewLayer(layerNum, true)
				} else {
					logger.Error(fmt.Sprintf("Layer %d does not exist; provide a number between (0-%d)", layerNum, numOfLayers))
				}
			} else {
				view.ViewAllLayers()
			}
		},
	}
	return ViewCommand
}
