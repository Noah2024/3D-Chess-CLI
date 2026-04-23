package must

import "3DChessCLI/util/logger"

func Must(err error) {
	if err != nil {
		logger.Error(err.Error())
		// panic(err)
	}
}
