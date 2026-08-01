// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// vecuint Provides CLI interface for accessing the VecToUint function
package vecuint

import (
	"3DC/util/bitutil"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// VecTUint converts a given x,y, and z vector into unsigned integer in game space using the VecToUint function from the bitutil package.
func VecTUint() *cobra.Command {
	return &cobra.Command{
		Use:   "vecuint",
		Short: "Turns a given x,y, and z vector into unsigned integer in game sapce ",
		Long:  "Directly runs the VecToUint command contained in bitutil",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {

			x, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			y, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}

			z, err := strconv.Atoi(args[2])
			if err != nil {
				return err
			}

			fmt.Println(bitutil.VecToUint(x, y, z))
			return nil
		},
	}
}
