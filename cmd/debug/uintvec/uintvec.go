// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// uintvec Provides CLI interface for accessing the UintToVec function
package uintvec

import (
	"3DC/util/bitutil"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// UintTVec converts a given unsigned integer to a vector in game space using the UintToVec function from the bitutil package.
func UintTVec() *cobra.Command {
	return &cobra.Command{
		Use:   "uintvec",
		Short: "Turns a given unsigned integer to vector in game space ",
		Long:  "Directly runs the UintToVec command contained in bitutil",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			uLoc, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			fmt.Println(bitutil.UintToVec(uint32(uLoc)))
			return nil
		},
	}
}
