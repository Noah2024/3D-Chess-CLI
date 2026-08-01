// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// must package provides a utility function for handling errors in Go.
//
//	The Must function takes a value and an error as input, and if the error is not nil, it logs the error message and terminates the program. This is useful for simplifying error handling in situations where an error is unexpected or unrecoverable.
package must

import (
	"3DC/util/logger"
	"os"
)

// The Must function takes a value and an error as input, and if the error is not nil, it logs the error message and terminates the program.
// This is useful for simplifying error handling in situations where an error is unexpected or unrecoverable.
func Must[T any](val T, err error) T {
	if err != nil {
		logger.Error(err.Error())
		panic(err)
		os.Exit(1)
	}
	return val
}
