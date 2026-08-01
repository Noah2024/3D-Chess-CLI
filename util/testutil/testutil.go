// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// Basic utility functions for testing and debugging, such as converting short bitmap strings to long form.
// Some debug functions are more proply placed in differnet packages, thus this package does not contian all debug functions, but those that didn't fit into their own packages.
package testutil

import (
	"fmt"
	"strconv"
	"strings"
)

// Generated with AI below, just cause I needed a quick way to convert the short form of the bitmaps into a longer form which is interperted later for debug

// Turns the shortened form of the bitmaps into a longer form which is interperted later for debug
func BitmapStringToBinary(bitmap string) (string, error) {
	bitmap = strings.Trim(bitmap, "[]")

	fields := strings.Fields(bitmap)

	var out strings.Builder
	out.WriteByte('[')

	for i, field := range fields {
		value, err := strconv.ParseUint(field, 10, 64)
		if err != nil {
			return "", fmt.Errorf("failed to parse %q: %w", field, err)
		}

		if i != 0 {
			out.WriteByte(' ')
		}

		fmt.Fprintf(&out, "%064b", value)
	}

	out.WriteByte(']')

	return out.String(), nil
}
