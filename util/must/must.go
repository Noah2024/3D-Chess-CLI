package must

import (
	"3DC/util/logger"
	"os"
)

// func Must(err error) {
// 	if err != nil {
// 		logger.Error(err.Error())
// 		// panic(err)
// 	}
// }

// Wraps function signature of (rtn, error) and handles error if any occured
func Must[T any](val T, err error) T {
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	return val
}
