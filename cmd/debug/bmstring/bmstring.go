// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

package bmstring

import (
	"3DC/util/testutil"
	"fmt"

	"github.com/spf13/cobra"
)

func Bmstring() *cobra.Command {
	return &cobra.Command{
		Use:   "bmstring",
		Short: "converts from short to long string representation of bitmap",
		Long:  "converts short string bitmap in form [0 0 0 0 0 0 0 0] to long form [00000000000...~512] ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			fmt.Println(testutil.BitmapStringToBinary(args[0]))
			return nil
		},
	}
}
